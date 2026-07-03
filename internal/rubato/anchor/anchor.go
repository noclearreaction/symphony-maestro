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
//	{"plugins":[{"plugin":"git_status"}]}
//	```
const (
	anchorOpen  = "```rubato:anchor\n"
	anchorClose = "\n```"
)

// Option is a name/setting pair used for both per-plugin options and top-level
// rubato config. Setting is optional; nil represents a flag-style option with
// no value.
type Option struct {
	Name    string `json:"name"`
	Setting any    `json:"setting,omitempty"`
}

// PluginDescriptor declares a plugin and its per-plugin options.
type PluginDescriptor struct {
	Plugin  string   `json:"plugin"`
	Options []Option `json:"options,omitempty"`
}

// Block holds the parsed contents of a rubato runtime anchor declaration.
type Block struct {
	// Plugins lists declared plugin descriptors in declaration order.
	Plugins []PluginDescriptor
	// Options holds top-level rubato-level config options.
	Options []Option
}

// MaxAge scans Options for "max_age" and returns its value as int.
// Defaults to 100 when not found. Returns 0 when explicitly set to 0.
func (b *Block) MaxAge() int {
	return IntOption(b.Options, "max_age", 100)
}

// IntOption scans opts for an option named name and returns its setting as int.
// Returns def when the name is not found or the setting cannot be interpreted as int.
// Handles JSON float64 coercion (Go decodes JSON numbers to float64 when target is any).
func IntOption(opts []Option, name string, def int) int {
	for _, o := range opts {
		if o.Name == name {
			switch v := o.Setting.(type) {
			case int:
				return v
			case float64:
				return int(v)
			case int64:
				return int(v)
			}
			return def
		}
	}
	return def
}

// StringOption scans opts for an option named name and returns its string setting.
// Returns ("", false) when not found or setting is not a string.
func StringOption(opts []Option, name string) (string, bool) {
	for _, o := range opts {
		if o.Name == name {
			s, ok := o.Setting.(string)
			return s, ok
		}
	}
	return "", false
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

	var envelope struct {
		Plugins []PluginDescriptor `json:"plugins"`
		Options []Option           `json:"options"`
	}
	if err := json.Unmarshal([]byte(raw), &envelope); err != nil {
		return nil, fmt.Errorf("malformed anchor: %w", err)
	}
	if len(envelope.Plugins) == 0 {
		return nil, fmt.Errorf("malformed anchor: no plugins declared")
	}

	// Normalise nil slices to empty slices for consistency.
	for i := range envelope.Plugins {
		if envelope.Plugins[i].Options == nil {
			envelope.Plugins[i].Options = []Option{}
		}
	}
	if envelope.Options == nil {
		envelope.Options = []Option{}
	}

	return &Block{Plugins: envelope.Plugins, Options: envelope.Options}, nil
}
