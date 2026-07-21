package docplugin

import (
	"context"
	"strings"

	"github.com/FiaLDI/project-parse/internal/domain"
	"github.com/FiaLDI/project-parse/internal/plugins/support"
	"github.com/FiaLDI/project-parse/internal/ports"
)

// Plugin analyzes project documentation coverage.
type Plugin struct {
	support.Base
}

// New creates the Documentation plugin.
func New(cache ports.FileCache) *Plugin {
	return &Plugin{Base: support.Base{PluginName: "documentation", PluginPriority: 8, Cache: cache}}
}

func (p *Plugin) Supports(_ domain.ProjectContext) bool {
	return true
}

func (p *Plugin) Analyze(ctx context.Context, pctx domain.ProjectContext) (domain.PluginResult, error) {
	_ = ctx
	var res domain.PluginResult
	var readmePaths []string

	for _, f := range pctx.Files.All() {
		upper := strings.ToUpper(f.Name)
		if strings.HasPrefix(upper, "README") {
			readmePaths = append(readmePaths, f.RelPath)
			res.Evidence = append(res.Evidence, support.Evidence(f.RelPath, "readme"))
		}
	}

	if len(readmePaths) > 0 {
		res.Findings = append(res.Findings, support.Finding(
			"docs.readme", domain.CategoryDocs, 0.9,
			"README", "Detected README file",
			readmePaths, nil,
		))
	}

	if support.HasDir(pctx, "docs") {
		res.Findings = append(res.Findings, support.Finding(
			"docs.directory", domain.CategoryDocs, 0.85,
			"Documentation directory", "Detected docs/ folder",
			[]string{"docs/"}, nil,
		))
	}

	topLevelMD := 0
	for _, f := range pctx.Files.All() {
		if f.Ext == ".md" && !strings.Contains(f.RelPath, "/") {
			topLevelMD++
		}
	}
	if topLevelMD > 0 {
		res.Findings = append(res.Findings, support.Finding(
			"docs.markdown", domain.CategoryDocs, 0.7,
			"Markdown docs", "Detected top-level markdown files",
			nil, map[string]any{"count": topLevelMD},
		))
	}
	return res, nil
}
