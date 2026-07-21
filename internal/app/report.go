package app

import (
	"context"
	"fmt"

	"github.com/FiaLDI/project-parse/internal/domain"
)

// ReportOptions controls report generation for a single invocation.
type ReportOptions struct {
	Root    string
	Formats []string
	OutDir  string
}

// ReportArtifact is a rendered report payload.
type ReportArtifact struct {
	Format string
	Path   string
	Bytes  []byte
}

// Report runs scan → analyze → aggregate → render.
func (a *App) Report(ctx context.Context, opts ReportOptions) ([]ReportArtifact, error) {
	if a.scanner == nil || a.registry == nil || a.analyzer == nil || a.agg == nil {
		return nil, fmt.Errorf("%w: report pipeline", ErrDependencyMissing)
	}

	root := opts.Root
	if root == "" {
		root = a.cfg.Scan.Root
	}
	formats := opts.Formats
	if len(formats) == 0 {
		formats = a.cfg.Report.Formats
	}
	outDir := opts.OutDir
	if outDir == "" {
		outDir = a.cfg.Report.OutputDir
	}

	pctx, err := a.scanner.Scan(ctx, root, scanOptionsFromConfig(a.cfg))
	if err != nil {
		return nil, err
	}

	plugins := a.registry.Enabled(a.cfg.Plugins.Enabled, a.cfg.Plugins.Disabled)
	results, err := a.analyzer.Run(ctx, pctx, plugins)
	if err != nil {
		return nil, err
	}

	report, err := a.agg.Aggregate(results)
	if err != nil {
		return nil, err
	}
	report.Meta.Root = pctx.Root
	report.Meta.FileCount = pctx.FileCount()

	var graph *domain.ArchitectureGraph
	if a.cfg.Graph.Enabled && a.graph != nil {
		g, gerr := a.graph.Build(report)
		if gerr != nil {
			a.log.Warn("graph build failed", "err", gerr)
		} else {
			graph = &g
		}
	}

	doc := domain.RenderDocument{
		Report: report,
		Graph:  graph,
		Options: domain.RenderOptions{
			IncludeEvidence: a.cfg.Report.IncludeEvidence,
			Title:           "project-parser report",
		},
	}

	artifacts := make([]ReportArtifact, 0, len(formats))
	for _, format := range formats {
		r, ok := a.renderers[format]
		if !ok || r == nil {
			return nil, fmt.Errorf("%w: renderer %q", ErrDependencyMissing, format)
		}
		data, rerr := r.Render(ctx, doc)
		if rerr != nil {
			return nil, fmt.Errorf("render %s: %w", format, rerr)
		}
		artifacts = append(artifacts, ReportArtifact{
			Format: format,
			Path:   outDir,
			Bytes:  data,
		})
	}
	return artifacts, nil
}
