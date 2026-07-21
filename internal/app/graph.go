package app

import (
	"context"
	"fmt"

	"github.com/FiaLDI/project-parse/internal/domain"
)

// GraphOptions controls graph generation.
type GraphOptions struct {
	Root   string
	Format string
	OutDir string
}

// Graph builds an architecture graph and renders it.
func (a *App) Graph(ctx context.Context, opts GraphOptions) (domain.ArchitectureGraph, []byte, error) {
	if a.scanner == nil || a.registry == nil || a.analyzer == nil || a.agg == nil || a.graph == nil {
		return domain.ArchitectureGraph{}, nil, fmt.Errorf("%w: graph pipeline", ErrDependencyMissing)
	}

	root := opts.Root
	if root == "" {
		root = a.cfg.Scan.Root
	}
	format := opts.Format
	if format == "" {
		format = a.cfg.Graph.Format
	}

	pctx, err := a.scanner.Scan(ctx, root, scanOptionsFromConfig(a.cfg))
	if err != nil {
		return domain.ArchitectureGraph{}, nil, err
	}

	plugins := a.registry.Enabled(a.cfg.Plugins.Enabled, a.cfg.Plugins.Disabled)
	results, err := a.analyzer.Run(ctx, pctx, plugins)
	if err != nil {
		return domain.ArchitectureGraph{}, nil, err
	}

	report, err := a.agg.Aggregate(results)
	if err != nil {
		return domain.ArchitectureGraph{}, nil, err
	}

	g, err := a.graph.Build(report)
	if err != nil {
		return domain.ArchitectureGraph{}, nil, err
	}

	r, ok := a.renderers[format]
	if !ok || r == nil {
		return g, nil, fmt.Errorf("%w: renderer %q", ErrDependencyMissing, format)
	}

	doc := domain.RenderDocument{
		Report: report,
		Graph:  &g,
		Options: domain.RenderOptions{
			IncludeEvidence: a.cfg.Report.IncludeEvidence,
			Title:           "project-parser graph",
		},
	}
	data, err := r.Render(ctx, doc)
	if err != nil {
		return g, nil, err
	}
	return g, data, nil
}
