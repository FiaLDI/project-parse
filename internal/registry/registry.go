package registry

import (
	"sort"
	"strings"
	"sync"

	"github.com/FiaLDI/project-parse/internal/ports"
)

// Registry is an in-memory plugin registry.
type Registry struct {
	mu      sync.RWMutex
	plugins map[string]ports.Plugin
}

// New creates an empty registry.
func New() *Registry {
	return &Registry{
		plugins: make(map[string]ports.Plugin),
	}
}

// Register adds or replaces a plugin by name.
func (r *Registry) Register(p ports.Plugin) {
	if p == nil {
		return
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	r.plugins[p.Name()] = p
}

// Enabled returns plugins allowed by policy, sorted by priority descending.
func (r *Registry) Enabled(enabled, disabled []string) []ports.Plugin {
	r.mu.RLock()
	defer r.mu.RUnlock()

	disabledSet := toSet(disabled)
	enabledSet := toSet(enabled)

	out := make([]ports.Plugin, 0, len(r.plugins))
	for name, p := range r.plugins {
		if disabledSet[name] {
			continue
		}
		if len(enabledSet) > 0 && !enabledSet[name] {
			continue
		}
		out = append(out, p)
	}
	sort.SliceStable(out, func(i, j int) bool {
		if out[i].Priority() == out[j].Priority() {
			return out[i].Name() < out[j].Name()
		}
		return out[i].Priority() > out[j].Priority()
	})
	return out
}

// List returns metadata for all registered plugins with enabled flags.
func (r *Registry) List(enabled, disabled []string) []ports.PluginMeta {
	active := r.Enabled(enabled, disabled)
	activeSet := make(map[string]struct{}, len(active))
	for _, p := range active {
		activeSet[p.Name()] = struct{}{}
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	out := make([]ports.PluginMeta, 0, len(r.plugins))
	for _, p := range r.plugins {
		_, on := activeSet[p.Name()]
		out = append(out, ports.PluginMeta{
			Name:     p.Name(),
			Priority: p.Priority(),
			Enabled:  on,
		})
	}
	sort.SliceStable(out, func(i, j int) bool {
		if out[i].Priority == out[j].Priority {
			return out[i].Name < out[j].Name
		}
		return out[i].Priority > out[j].Priority
	})
	return out
}

func toSet(items []string) map[string]bool {
	set := make(map[string]bool, len(items))
	for _, item := range items {
		item = strings.TrimSpace(item)
		if item == "" {
			continue
		}
		set[item] = true
	}
	return set
}
