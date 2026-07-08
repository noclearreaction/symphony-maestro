package plugin

import (
	"context"
	"fmt"

	"github.com/noclearreaction/symphony-maestro/internal/rubato/anchor"
)

// Plugin is the contract all rubato plugins must implement.
type Plugin interface {
	Name() string
	// Execute runs the plugin with the given options and returns its output.
	// Options may be empty when none were declared in the anchor descriptor.
	Execute(ctx context.Context, options []anchor.Option) (string, error)
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

// Execute runs each descriptor in order and returns outputs keyed by plugin name.
// Returns an error immediately for unknown plugin names or execution failures (fail-fast).
func (r *Registry) Execute(ctx context.Context, descriptors []anchor.PluginDescriptor) (map[string]string, error) {
	out := make(map[string]string, len(descriptors))
	for _, d := range descriptors {
		p, ok := r.plugins[d.Plugin]
		if !ok {
			return nil, fmt.Errorf("unknown plugin: %q", d.Plugin)
		}
		result, err := p.Execute(ctx, d.Options)
		if err != nil {
			return nil, fmt.Errorf("plugin %q: %w", d.Plugin, err)
		}
		out[d.Plugin] = result
	}
	return out, nil
}
