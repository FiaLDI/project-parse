package registry

import (
	"context"
	"testing"

	"github.com/FiaLDI/project-parse/internal/domain"
	"github.com/FiaLDI/project-parse/internal/ports"
)

type stubPlugin struct {
	name     string
	priority int
}

func (s stubPlugin) Name() string     { return s.name }
func (s stubPlugin) Priority() int    { return s.priority }
func (s stubPlugin) Supports(domain.ProjectContext) bool { return true }
func (s stubPlugin) Analyze(context.Context, domain.ProjectContext) (domain.PluginResult, error) {
	return domain.PluginResult{PluginName: s.name}, nil
}

func TestEnabledFilteringAndPriority(t *testing.T) {
	reg := New()
	reg.Register(stubPlugin{name: "node", priority: 10})
	reg.Register(stubPlugin{name: "python", priority: 20})
	reg.Register(stubPlugin{name: "docker", priority: 5})

	tests := []struct {
		name     string
		enabled  []string
		disabled []string
		want     []string
	}{
		{
			name:    "all_when_enabled_empty",
			want:    []string{"python", "node", "docker"},
		},
		{
			name:    "explicit_enabled",
			enabled: []string{"node", "docker"},
			want:    []string{"node", "docker"},
		},
		{
			name:     "disabled_wins",
			enabled:  []string{"node", "python"},
			disabled: []string{"python"},
			want:     []string{"node"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := reg.Enabled(tt.enabled, tt.disabled)
			if len(got) != len(tt.want) {
				t.Fatalf("len=%d want %d", len(got), len(tt.want))
			}
			for i := range tt.want {
				if got[i].Name() != tt.want[i] {
					t.Fatalf("idx %d: got %s want %s", i, got[i].Name(), tt.want[i])
				}
			}
		})
	}
}

func TestListMarksEnabled(t *testing.T) {
	reg := New()
	reg.Register(stubPlugin{name: "git", priority: 1})
	reg.Register(stubPlugin{name: "node", priority: 2})

	list := reg.List([]string{"git"}, nil)
	byName := map[string]ports.PluginMeta{}
	for _, m := range list {
		byName[m.Name] = m
	}
	if !byName["git"].Enabled {
		t.Fatal("git should be enabled")
	}
	if byName["node"].Enabled {
		t.Fatal("node should be disabled")
	}
}
