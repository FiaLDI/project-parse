package githubactions

import (
	"context"
	"strings"

	"github.com/FiaLDI/project-parse/internal/domain"
	"github.com/FiaLDI/project-parse/internal/plugins/support"
	"github.com/FiaLDI/project-parse/internal/ports"
)

// Plugin analyzes GitHub Actions workflows.
type Plugin struct {
	support.Base
}

// New creates the GitHub Actions plugin.
func New(cache ports.FileCache) *Plugin {
	return &Plugin{Base: support.Base{PluginName: "githubactions", PluginPriority: 15, Cache: cache}}
}

func (p *Plugin) Supports(pctx domain.ProjectContext) bool {
	return support.HasMarker(pctx, ".github/workflows") || support.HasPathPrefix(pctx, ".github/workflows/")
}

func (p *Plugin) Analyze(ctx context.Context, pctx domain.ProjectContext) (domain.PluginResult, error) {
	_ = ctx
	var res domain.PluginResult
	var paths []string
	for _, f := range pctx.Files.All() {
		if strings.HasPrefix(f.RelPath, ".github/workflows/") {
			ext := strings.ToLower(f.Ext)
			if ext == ".yml" || ext == ".yaml" {
				paths = append(paths, f.RelPath)
				res.Evidence = append(res.Evidence, support.Evidence(f.RelPath, "workflow"))
			}
		}
	}
	if len(paths) == 0 {
		return res, nil
	}
	res.Findings = append(res.Findings, support.Finding(
		"ci.github_actions", domain.CategoryCI, 0.92,
		"GitHub Actions", "Detected workflow files",
		paths, map[string]any{"count": len(paths)},
	))
	return res, nil
}
