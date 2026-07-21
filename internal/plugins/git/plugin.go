package git

import (
	"context"
	"os"
	"path/filepath"

	"github.com/FiaLDI/project-parse/internal/domain"
	"github.com/FiaLDI/project-parse/internal/plugins/support"
	"github.com/FiaLDI/project-parse/internal/ports"
)

// Plugin analyzes Git repository signals.
type Plugin struct {
	support.Base
}

// New creates the Git plugin.
func New(cache ports.FileCache) *Plugin {
	return &Plugin{Base: support.Base{PluginName: "git", PluginPriority: 10, Cache: cache}}
}

func (p *Plugin) Supports(_ domain.ProjectContext) bool {
	return true
}

func (p *Plugin) Analyze(ctx context.Context, pctx domain.ProjectContext) (domain.PluginResult, error) {
	_ = ctx
	var res domain.PluginResult

	gitPath := filepath.Join(pctx.Root, ".git")
	if info, err := os.Stat(gitPath); err == nil && (info.IsDir() || info.Mode()&os.ModeSymlink != 0) {
		res.Findings = append(res.Findings, support.Finding(
			"git.repository", domain.CategoryVCS, 0.99,
			"Git repository", "Detected .git directory",
			[]string{".git"}, nil,
		))
	}

	for _, f := range support.FilesByName(pctx, ".gitignore") {
		res.Evidence = append(res.Evidence, support.Evidence(f.RelPath, "gitignore"))
		res.Findings = append(res.Findings, support.Finding(
			"git.gitignore", domain.CategoryVCS, 0.85,
			"Git ignore rules", "Detected .gitignore",
			[]string{f.RelPath}, nil,
		))
	}
	return res, nil
}
