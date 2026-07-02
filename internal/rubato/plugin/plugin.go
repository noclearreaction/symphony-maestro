package plugin

import (
	"context"
	"fmt"
)

// Plugin is the contract all rubato plugins must implement.
type Plugin interface {
	Name() string
	// Execute runs the plugin with the given static args and returns its output.
	// Args may be nil when none were declared in the anchor.
	Execute(ctx context.Context, args map[string]any) (string, error)
}

// Registry resolves and executes declared plugins by name.
// Plugins execute fresh on every call — no session reuse.
type Registry struct {
	plugins map[string]Plugin
}

// NewRegistry creates a Registry pre-populated with the given plugins.
func NewRegistry(plugins ...Plugin) *Registry {
	r := &Registry{plugins: make(map[string]Plugin, len(plugins))}
	for _, p := range plugins {
		r.plugins[p.Name()] = p
	}
	return r
}

// Execute runs each name in declared, in order, and returns their outputs keyed by name.
// Returns an error immediately for unknown plugin names or execution failures (fail-fast).
func (r *Registry) Execute(ctx context.Context, declared []string, args map[string]map[string]any) (map[string]string, error) {
	out := make(map[string]string, len(declared))
	for _, name := range declared {
		p, ok := r.plugins[name]
		if !ok {
			return nil, fmt.Errorf("unknown plugin: %q", name)
		}
		result, err := p.Execute(ctx, args[name])
		if err != nil {
			return nil, fmt.Errorf("plugin %q: %w", name, err)
		}
		out[name] = result
	}
	return out, nil
}
