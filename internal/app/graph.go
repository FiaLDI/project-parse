package app

import (
	"context"
	"fmt"

	"github.com/FiaLDI/project-parse/internal/domain"
	"github.com/FiaLDI/project-parse/internal/output"
)

// GraphOptions controls graph generation.
type GraphOptions struct {
	Root   string
	Format string
	OutDir string
}

// GraphArtifact is a rendered architecture graph output.
type GraphArtifact struct {
	Graph  domain.ArchitectureGraph
	Format string
	Path   string
	Bytes  []byte
}

// Graph builds an architecture graph and renders it.
func (a *App) Graph(ctx context.Context, opts GraphOptions) (GraphArtifact, error) {
	if a.scanner == nil || a.registry == nil || a.analyzer == nil || a.agg == nil || a.graph == nil {
		return GraphArtifact{}, fmt.Errorf("%w: graph pipeline", ErrDependencyMissing)
	}

	root := opts.Root
	if root == "" {
		root = a.cfg.Scan.Root
	}
	format := normalizeGraphFormat(opts.Format)
	if format == "" {
		format = normalizeGraphFormat(a.cfg.Graph.Format)
	}

	outDir := opts.OutDir
	if outDir == "" {
		outDir = a.cfg.Report.OutputDir
	}

	pctx, err := a.scanner.Scan(ctx, root, scanOptionsFromConfig(a.cfg))
	if err != nil {
		return GraphArtifact{}, err
	}

	plugins := a.registry.Enabled(a.cfg.Plugins.Enabled, a.cfg.Plugins.Disabled)
	results, err := a.analyzer.Run(ctx, pctx, plugins)
	if err != nil {
		return GraphArtifact{}, err
	}

	report, err := a.agg.Aggregate(results)
	if err != nil {
		return GraphArtifact{}, err
	}
	report.Meta.Root = pctx.Root
	report.Meta.FileCount = pctx.FileCount()

	g, err := a.graph.Build(report)
	if err != nil {
		return GraphArtifact{}, err
	}

	r, ok := a.renderers[format]
	if !ok || r == nil {
		return GraphArtifact{}, fmt.Errorf("%w: renderer %q", ErrDependencyMissing, format)
	}

	doc := domain.RenderDocument{
		Report: report,
		Graph:  &g,
		Options: domain.RenderOptions{
			IncludeEvidence: a.cfg.Report.IncludeEvidence,
			Title:           "Architecture Graph",
		},
	}
	data, err := r.Render(ctx, doc)
	if err != nil {
		return GraphArtifact{}, err
	}

	path, err := output.WriteFile(outDir, format, data)
	if err != nil {
		return GraphArtifact{}, fmt.Errorf("write %s: %w", format, err)
	}

	return GraphArtifact{
		Graph:  g,
		Format: format,
		Path:   path,
		Bytes:  data,
	}, nil
}

func normalizeGraphFormat(format string) string {
	switch format {
	case "json":
		return "graph-json"
	default:
		return format
	}
}
