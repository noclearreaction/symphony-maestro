package mutate

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sort"
	"strings"

	"github.com/noclearreaction/symphony-maestro/internal/rubato/anchor"
	"github.com/noclearreaction/symphony-maestro/internal/rubato/plugin"
)

const guidanceVersion = "1" // retained for any future version-tagged guidance

// rawMsg is a type alias for a JSON object whose fields are all preserved as raw JSON.
// This lets us parse, modify, and re-serialise messages without dropping unknown fields
// (e.g. name, tool_call_id, cache_control, etc.).
type rawMsg = map[string]json.RawMessage

// Injector applies runtime plugin injection to OpenAI-compatible chat completion
// request bodies. Each call to Apply is independent with no session state.
type Injector struct {
	registry *plugin.Registry
}

// NewInjector returns an Injector backed by the given registry.
func NewInjector(registry *plugin.Registry) *Injector {
	return &Injector{registry: registry}
}

// Apply inspects body for a rubato anchor in messages[0].content.
//   - No anchor → body returned unchanged.
//   - Anchor found → execute declared plugins, mutate messages, return new body.
//   - Malformed anchor, unknown plugin, invalid structure → return error (caller rejects request).
func (inj *Injector) Apply(ctx context.Context, body []byte) ([]byte, error) {
	// Parse the full request into a map so all unknown top-level fields are preserved.
	var req rawMsg
	if err := json.Unmarshal(body, &req); err != nil {
		return nil, fmt.Errorf("mutate: invalid JSON: %w", err)
	}

	messagesRaw, ok := req["messages"]
	if !ok {
		return body, nil
	}

	var msgs []json.RawMessage
	if err := json.Unmarshal(messagesRaw, &msgs); err != nil {
		return nil, fmt.Errorf("mutate: invalid messages array: %w", err)
	}
	if len(msgs) == 0 {
		return body, nil
	}

	// Parse messages[0] preserving all fields.
	var first rawMsg
	if err := json.Unmarshal(msgs[0], &first); err != nil {
		return nil, fmt.Errorf("mutate: messages[0]: %w", err)
	}
	firstContent, ok := first["content"]
	if !ok {
		return nil, fmt.Errorf("mutate: messages[0]: missing content field")
	}

	// Extract combined text from messages[0] to search for the anchor.
	firstText, err := textFrom(firstContent)
	if err != nil {
		return nil, fmt.Errorf("mutate: messages[0].content: %w", err)
	}

	block, err := anchor.Find(firstText)
	if err != nil {
		return nil, err // malformed anchor — caller rejects
	}
	if block == nil {
		log.Printf("rubato: no anchor in system message; pass-through")
		return body, nil // no anchor — pass through unchanged
	}

	// Extract plugin names early for logging.
	names := make([]string, len(block.Plugins))
	for i, d := range block.Plugins {
		names[i] = d.Plugin
	}
	log.Printf("rubato: anchor detected plugins=%v max_age=%d", names, block.MaxAge())

	// Execute declared plugins.
	outputs, err := inj.registry.Execute(ctx, block.Plugins)
	if err != nil {
		log.Printf("rubato: plugin execution failed: %v", err)
		return nil, err
	}
	log.Printf("rubato: plugins executed successfully count=%d", len(outputs))

	// Need at least 2 messages: messages[0] (system) and messages[-1] (user turn).
	last := len(msgs) - 1
	if last == 0 {
		return nil, fmt.Errorf("mutate: need at least 2 messages for injection")
	}

	// On-change injection: determine which plugins need re-injection.
	maxAge := block.MaxAge()
	var injectNames []string
	if maxAge == 0 {
		// max_age 0: always inject all plugins unconditionally.
		injectNames = names
	} else {
		prior := scanPriorState(msgs, maxAge)
		for _, name := range names {
			if outputs[name] != prior[name] {
				injectNames = append(injectNames, name)
			}
		}
	}

	if len(injectNames) == 0 {
		log.Printf("rubato: all plugins stable; no state block injected")
	} else {
		log.Printf("rubato: injecting plugins=%v stable=%v", injectNames, stableNames(names, injectNames))
	}

	// Mutate messages[0]: inject guidance when absent (idempotent).
	guidance := buildGuidance(names)
	if !strings.Contains(firstText, "```rubato:guidance") {
		first["content"], err = appendToContent(firstContent, "\n\n"+guidance)
		if err != nil {
			return nil, fmt.Errorf("mutate: messages[0] guidance: %w", err)
		}
		msgs[0], err = json.Marshal(first)
		if err != nil {
			return nil, fmt.Errorf("mutate: re-marshal messages[0]: %w", err)
		}
	}

	// Only prepend a state block when at least one plugin needs injection.
	if len(injectNames) > 0 {
		var lastMsg rawMsg
		if err := json.Unmarshal(msgs[last], &lastMsg); err != nil {
			return nil, fmt.Errorf("mutate: messages[-1]: %w", err)
		}
		lastContent, ok := lastMsg["content"]
		if !ok {
			return nil, fmt.Errorf("mutate: messages[-1]: missing content field")
		}
		injectOutputs := make(map[string]string, len(injectNames))
		for _, n := range injectNames {
			injectOutputs[n] = outputs[n]
		}
		stateBlock := buildStateBlock(injectNames, injectOutputs)
		lastMsg["content"], err = prependToContent(lastContent, stateBlock+"\n\n")
		if err != nil {
			return nil, fmt.Errorf("mutate: messages[-1].content: %w", err)
		}
		msgs[last], err = json.Marshal(lastMsg)
		if err != nil {
			return nil, fmt.Errorf("mutate: re-marshal messages[-1]: %w", err)
		}
	}

	// Re-pack the mutated messages back into the request map.
	req["messages"], err = json.Marshal(msgs)
	if err != nil {
		return nil, fmt.Errorf("mutate: re-marshal messages: %w", err)
	}

	return marshalNoEscape(req)
}

// marshalNoEscape encodes v as JSON without HTML-escaping < > & so that
// the forwarded request body remains human-readable and easily inspectable.
func marshalNoEscape(v any) ([]byte, error) {
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(false)
	if err := enc.Encode(v); err != nil {
		return nil, err
	}
	// json.Encoder.Encode appends a newline; trim it.
	return bytes.TrimRight(buf.Bytes(), "\n"), nil
}

// buildStateBlock assembles the runtime-state fenced block.
// Plugins appear in declaration order.
func buildStateBlock(declared []string, outputs map[string]string) string {
	var sb strings.Builder
	sb.WriteString("```rubato:state\n")
	for _, name := range declared {
		sb.WriteString("[" + name + "]\n")
		sb.WriteString(outputs[name])
		sb.WriteString("\n")
	}
	sb.WriteString("```")
	return sb.String()
}

// buildGuidance assembles a deterministic guidance fenced block for the declared plugin set.
// Plugin names are sorted so that identical plugin sets always produce identical bytes.
func buildGuidance(declared []string) string {
	sorted := make([]string, len(declared))
	copy(sorted, declared)
	sort.Strings(sorted)

	var sb strings.Builder
	sb.WriteString("```rubato:guidance\n")
	sb.WriteString("Runtime context is injected per request by the rubato proxy.\n")
	sb.WriteString("Active plugins:\n")
	for _, name := range sorted {
		sb.WriteString("- " + name + ": " + pluginDesc(name) + "\n")
	}
	sb.WriteString("```")
	return sb.String()
}

// pluginDesc returns a fixed one-line description for a known plugin.
func pluginDesc(name string) string {
	switch name {
	case "git_status":
		return "current git repository status (branch, ahead/behind, staged, unstaged, untracked)"
	default:
		return "runtime plugin output"
	}
}

// stableNames returns the elements of all that are not in inject.
func stableNames(all, inject []string) []string {
	injectSet := make(map[string]struct{}, len(inject))
	for _, n := range inject {
		injectSet[n] = struct{}{}
	}
	var stable []string
	for _, n := range all {
		if _, found := injectSet[n]; !found {
			stable = append(stable, n)
		}
	}
	return stable
}

// textFrom extracts combined plain text from a content field.
// content may be a JSON string or a JSON array of content-part objects.
func textFrom(content json.RawMessage) (string, error) {
	// Try plain string first.
	var s string
	if err := json.Unmarshal(content, &s); err == nil {
		return s, nil
	}
	// Try array of content parts.
	var parts []struct {
		Type string `json:"type"`
		Text string `json:"text"`
	}
	if err := json.Unmarshal(content, &parts); err != nil {
		return "", fmt.Errorf("unexpected content format (not a string or array)")
	}
	var sb strings.Builder
	for _, p := range parts {
		if p.Type == "text" {
			sb.WriteString(p.Text)
		}
	}
	return sb.String(), nil
}

// prependToContent prepends prefix to a string or array content field.
// For array content, a new text-type part is inserted at the front.
func prependToContent(content json.RawMessage, prefix string) (json.RawMessage, error) {
	var s string
	if err := json.Unmarshal(content, &s); err == nil {
		return json.Marshal(prefix + s)
	}
	var parts []json.RawMessage
	if err := json.Unmarshal(content, &parts); err != nil {
		return nil, fmt.Errorf("unexpected content format for prepend")
	}
	prefixPart, _ := json.Marshal(map[string]string{"type": "text", "text": prefix})
	result := make([]json.RawMessage, 0, 1+len(parts))
	result = append(result, prefixPart)
	result = append(result, parts...)
	return json.Marshal(result)
}

// appendToContent appends suffix to a string or array content field.
// For array content, a new text-type part is added at the end.
func appendToContent(content json.RawMessage, suffix string) (json.RawMessage, error) {
	var s string
	if err := json.Unmarshal(content, &s); err == nil {
		return json.Marshal(s + suffix)
	}
	var parts []json.RawMessage
	if err := json.Unmarshal(content, &parts); err != nil {
		return nil, fmt.Errorf("unexpected content format for append")
	}
	suffixPart, _ := json.Marshal(map[string]string{"type": "text", "text": suffix})
	result := make([]json.RawMessage, 0, len(parts)+1)
	result = append(result, parts...)
	result = append(result, suffixPart)
	return json.Marshal(result)
}

// scanPriorState scans backward through msgs[0..len-2] (up to maxAge positions)
// for rubato:state blocks, returning the most recent known output per plugin name.
// msgs[-1] (the current user turn being mutated) is excluded from the scan.
func scanPriorState(msgs []json.RawMessage, maxAge int) map[string]string {
	result := make(map[string]string)
	last := len(msgs) - 2 // index of msgs[-2]: exclude current user turn
	if last < 0 {
		return result
	}
	// Scan at most maxAge messages, from newest (last) to oldest.
	for i := last; i >= 0 && i > last-maxAge; i-- {
		text := textFromMsg(msgs[i])
		if text == "" {
			continue
		}
		for name, output := range parseStateBlock(text) {
			if _, seen := result[name]; !seen {
				result[name] = output // most recent wins
			}
		}
	}
	return result
}

// textFromMsg extracts combined plain text from a message JSON object.
// Returns "" on any error (best-effort for backward scanning).
func textFromMsg(msg json.RawMessage) string {
	var m map[string]json.RawMessage
	if err := json.Unmarshal(msg, &m); err != nil {
		return ""
	}
	content, ok := m["content"]
	if !ok {
		return ""
	}
	text, _ := textFrom(content)
	return text
}

// parseStateBlock parses a rubato:state fenced block from text, returning
// per-plugin output keyed by plugin name. Returns an empty map when no state
// block is present.
func parseStateBlock(text string) map[string]string {
	result := make(map[string]string)
	lines := strings.Split(text, "\n")
	inBlock := false
	currentPlugin := ""
	var currentLines []string

	for _, line := range lines {
		if !inBlock {
			if line == "```rubato:state" {
				inBlock = true
			}
			continue
		}
		// Close fence ends the block.
		if line == "```" {
			if currentPlugin != "" {
				result[currentPlugin] = strings.TrimRight(strings.Join(currentLines, "\n"), "\n")
			}
			break
		}
		// Section header: [plugin_name]
		if len(line) >= 2 && line[0] == '[' && line[len(line)-1] == ']' {
			if currentPlugin != "" {
				result[currentPlugin] = strings.TrimRight(strings.Join(currentLines, "\n"), "\n")
			}
			currentPlugin = line[1 : len(line)-1]
			currentLines = nil
			continue
		}
		if currentPlugin != "" {
			currentLines = append(currentLines, line)
		}
	}
	return result
}
