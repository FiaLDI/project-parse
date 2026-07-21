package ports

import (
	"context"

	"github.com/FiaLDI/project-parse/internal/domain"
)

// Plugin analyzes a project and returns structured findings.
// Implementations must be safe for concurrent Analyze calls on distinct contexts.
type Plugin interface {
	Name() string
	Priority() int
	Supports(ctx domain.ProjectContext) bool
	Analyze(ctx context.Context, pctx domain.ProjectContext) (domain.PluginResult, error)
}

// PluginMeta is a lightweight descriptor for listing plugins.
type PluginMeta struct {
	Name     string `json:"name"`
	Priority int    `json:"priority"`
	Enabled  bool   `json:"enabled"`
}
