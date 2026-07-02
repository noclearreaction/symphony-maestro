package anchor

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Anchor blocks use markdown fenced code blocks so that backtick characters
// appear in JSON strings without triggering HTML-escape in json.Marshal.
//
// Format in system message content:
//
//	```rubato:anchor
//	{"plugins":["git_status"]}
//	```
const (
	anchorOpen  = "```rubato:anchor\n"
	anchorClose = "\n```"
)

// Block holds the parsed contents of a rubato runtime anchor declaration.
type Block struct {
	// Plugins lists declared plugin names in declaration order.
	Plugins []string
	// Args holds per-plugin static args, keyed by plugin name.
	Args map[string]map[string]any
}

// Find searches content for a rubato anchor block.
//
// Returns nil, nil when no anchor is present — the request should be forwarded unchanged.
// Returns nil, error when an anchor tag is present but its content is malformed.
func Find(content string) (*Block, error) {
	start := strings.Index(content, anchorOpen)
	if start == -1 {
		return nil, nil // no anchor — bypass injection
	}
	bodyStart := start + len(anchorOpen)
	end := strings.Index(content[bodyStart:], anchorClose)
	if end == -1 {
		return nil, fmt.Errorf("malformed anchor: missing closing fence")
	}
	raw := strings.TrimSpace(content[bodyStart : bodyStart+end])
	if raw == "" {
		return nil, fmt.Errorf("malformed anchor: empty body")
	}

	// Parse the plugins list.
	var envelope struct {
		Plugins []string `json:"plugins"`
	}
	if err := json.Unmarshal([]byte(raw), &envelope); err != nil {
		return nil, fmt.Errorf("malformed anchor: %w", err)
	}
	if len(envelope.Plugins) == 0 {
		return nil, fmt.Errorf("malformed anchor: no plugins declared")
	}

	// Extract per-plugin args from top-level keys matching declared plugin names.
	var all map[string]json.RawMessage
	_ = json.Unmarshal([]byte(raw), &all)

	args := make(map[string]map[string]any, len(envelope.Plugins))
	for _, name := range envelope.Plugins {
		raw, ok := all[name]
		if !ok {
			continue
		}
		var pluginArgs map[string]any
		if err := json.Unmarshal(raw, &pluginArgs); err != nil {
			return nil, fmt.Errorf("malformed anchor: args for plugin %q: %w", name, err)
		}
		args[name] = pluginArgs
	}

	return &Block{Plugins: envelope.Plugins, Args: args}, nil
}
