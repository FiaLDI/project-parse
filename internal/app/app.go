package app

import (
	"log/slog"

	"github.com/FiaLDI/project-parse/internal/config"
	"github.com/FiaLDI/project-parse/internal/ports"
)

// App orchestrates use-cases. CLI must not contain business logic.
type App struct {
	cfg      config.Config
	log      *slog.Logger
	scanner  ports.Scanner
	cache    ports.FileCache
	registry ports.Registry
	analyzer ports.Analyzer
	agg      ports.Aggregator
	graph    ports.GraphBuilder
	renderers map[string]ports.Renderer
}

// Deps holds optional collaborators injected at construction time.
type Deps struct {
	Scanner   ports.Scanner
	Cache     ports.FileCache
	Registry  ports.Registry
	Analyzer  ports.Analyzer
	Agg       ports.Aggregator
	Graph     ports.GraphBuilder
	Renderers []ports.Renderer
}

// New constructs an App. Missing deps are allowed in stage 0; use-cases
// return ErrNotImplemented until later stages wire real implementations.
func New(cfg config.Config, log *slog.Logger, deps Deps) *App {
	if log == nil {
		log = slog.Default()
	}
	a := &App{
		cfg:       cfg,
		log:       log,
		scanner:   deps.Scanner,
		cache:     deps.Cache,
		registry:  deps.Registry,
		analyzer:  deps.Analyzer,
		agg:       deps.Agg,
		graph:     deps.Graph,
		renderers: make(map[string]ports.Renderer),
	}
	for _, r := range deps.Renderers {
		if r != nil {
			a.renderers[r.Format()] = r
		}
	}
	return a
}

// Config returns the active configuration snapshot.
func (a *App) Config() config.Config {
	return a.cfg
}
