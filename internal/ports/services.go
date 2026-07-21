package ports

import (
	"context"
	"time"

	"github.com/FiaLDI/project-parse/internal/domain"
)

// Analyzer runs enabled plugins against a project context.
type Analyzer interface {
	Run(ctx context.Context, pctx domain.ProjectContext, plugins []Plugin) ([]domain.PluginResult, error)
}

// Aggregator merges plugin results into a single Report.
type Aggregator interface {
	Aggregate(results []domain.PluginResult) (domain.Report, error)
}

// GraphBuilder derives an architecture graph from a report.
type GraphBuilder interface {
	Build(report domain.Report) (domain.ArchitectureGraph, error)
}

// Renderer serializes a render document into a specific format.
type Renderer interface {
	Format() string
	Render(ctx context.Context, doc domain.RenderDocument) ([]byte, error)
}

// Clock abstracts time for testability.
type Clock interface {
	Now() time.Time
}
