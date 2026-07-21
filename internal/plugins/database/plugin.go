package database

import (
	"context"

	"github.com/FiaLDI/project-parse/internal/domain"
	"github.com/FiaLDI/project-parse/internal/plugins/support"
	"github.com/FiaLDI/project-parse/internal/ports"
)

// Plugin analyzes database-related artifacts.
type Plugin struct {
	support.Base
}

// New creates the Database plugin.
func New(cache ports.FileCache) *Plugin {
	return &Plugin{Base: support.Base{PluginName: "database", PluginPriority: 18, Cache: cache}}
}

func (p *Plugin) Supports(pctx domain.ProjectContext) bool {
	return support.HasMarker(pctx, "schema.prisma", "sql/", "migrations/") ||
		support.HasDir(pctx, "sql") || support.HasDir(pctx, "migrations") ||
		support.HasPathPrefix(pctx, "prisma/")
}

func (p *Plugin) Analyze(ctx context.Context, pctx domain.ProjectContext) (domain.PluginResult, error) {
	_ = ctx
	var res domain.PluginResult

	for _, f := range support.FilesByName(pctx, "schema.prisma") {
		res.Evidence = append(res.Evidence, support.Evidence(f.RelPath, "Prisma schema"))
		res.Findings = append(res.Findings, support.Finding(
			"db.prisma", domain.CategoryDatabase, 0.9,
			"Prisma schema", "Detected schema.prisma",
			[]string{f.RelPath}, map[string]any{"engine": "prisma"},
		))
	}

	if support.HasDir(pctx, "sql") || support.HasMarker(pctx, "sql/") {
		res.Findings = append(res.Findings, support.Finding(
			"db.sql", domain.CategoryDatabase, 0.82,
			"SQL scripts", "Detected sql/ directory",
			[]string{"sql/"}, nil,
		))
	}

	if support.HasDir(pctx, "migrations") || support.HasMarker(pctx, "migrations/") {
		res.Findings = append(res.Findings, support.Finding(
			"db.migrations", domain.CategoryDatabase, 0.82,
			"Database migrations", "Detected migrations directory",
			[]string{"migrations/"}, nil,
		))
	}
	return res, nil
}
