package architecture

import (
	"context"

	"github.com/FiaLDI/project-parse/internal/domain"
	"github.com/FiaLDI/project-parse/internal/plugins/support"
	"github.com/FiaLDI/project-parse/internal/ports"
)

type pattern struct {
	name       string
	id         string
	confidence float64
	match      func(domain.ProjectContext) bool
}

// Plugin detects common architecture patterns from directory layout.
type Plugin struct {
	support.Base
}

// New creates the Architecture plugin.
func New(cache ports.FileCache) *Plugin {
	return &Plugin{Base: support.Base{PluginName: "architecture", PluginPriority: 12, Cache: cache}}
}

func (p *Plugin) Supports(_ domain.ProjectContext) bool {
	return true
}

func (p *Plugin) Analyze(ctx context.Context, pctx domain.ProjectContext) (domain.PluginResult, error) {
	_ = ctx
	var res domain.PluginResult

	patterns := []pattern{
		{
			name: "Feature-Sliced Design", id: "arch.fsd", confidence: 0.75,
			match: func(pc domain.ProjectContext) bool {
				hits := 0
				for _, d := range []string{"shared", "entities", "features", "widgets", "pages", "app"} {
					if support.HasDir(pc, "src/"+d) || support.HasDir(pc, d) {
						hits++
					}
				}
				return hits >= 3
			},
		},
		{
			name: "Clean Architecture", id: "arch.clean", confidence: 0.78,
			match: func(pc domain.ProjectContext) bool {
				return (support.HasDir(pc, "internal/domain") || support.HasDir(pc, "domain")) &&
					(support.HasDir(pc, "internal/usecase") || support.HasDir(pc, "usecase") || support.HasDir(pc, "application"))
			},
		},
		{
			name: "Domain-Driven Design", id: "arch.ddd", confidence: 0.76,
			match: func(pc domain.ProjectContext) bool {
				return support.HasDir(pc, "domain") && (support.HasDir(pc, "infrastructure") || support.HasDir(pc, "application"))
			},
		},
		{
			name: "Hexagonal Architecture", id: "arch.hexagonal", confidence: 0.77,
			match: func(pc domain.ProjectContext) bool {
				return support.HasDir(pc, "ports") && support.HasDir(pc, "adapters")
			},
		},
	}

	var matched []string
	for _, pat := range patterns {
		if pat.match(pctx) {
			matched = append(matched, pat.name)
			res.Findings = append(res.Findings, support.Finding(
				pat.id, domain.CategoryArch, pat.confidence,
				pat.name, "Detected characteristic directory layout",
				nil, map[string]any{"pattern": pat.name},
			))
		}
	}

	if len(matched) == 0 {
		res.Findings = append(res.Findings, support.Finding(
			"arch.unknown", domain.CategoryArch, 0.4,
			"No known pattern", "No FSD/Clean/DDD/Hexagonal layout detected",
			nil, nil,
		))
	}
	return res, nil
}
