package analyzer

import (
	"context"
	"log/slog"
	"runtime"
	"sync"
	"time"

	"github.com/FiaLDI/project-parse/internal/domain"
	"github.com/FiaLDI/project-parse/internal/ports"
)

// Options configures the plugin analyzer worker pool.
type Options struct {
	Jobs int
	Log  *slog.Logger
}

// Analyzer runs plugins concurrently with bounded parallelism.
type Analyzer struct {
	jobs int
	log  *slog.Logger
}

// New creates an analyzer.
func New(opts Options) *Analyzer {
	jobs := opts.Jobs
	if jobs <= 0 {
		jobs = runtime.NumCPU()
		if jobs < 1 {
			jobs = 1
		}
	}
	log := opts.Log
	if log == nil {
		log = slog.Default()
	}
	return &Analyzer{jobs: jobs, log: log}
}

// Run executes supported plugins in parallel. Individual plugin failures are
// recorded in PluginResult.Errors and do not fail the whole run.
func (a *Analyzer) Run(ctx context.Context, pctx domain.ProjectContext, plugins []ports.Plugin) ([]domain.PluginResult, error) {
	var supported []ports.Plugin
	for _, p := range plugins {
		if p == nil {
			continue
		}
		if p.Supports(pctx) {
			supported = append(supported, p)
		}
	}
	if len(supported) == 0 {
		return nil, nil
	}

	results := make([]domain.PluginResult, len(supported))
	sem := make(chan struct{}, a.jobs)
	var wg sync.WaitGroup

	for i, plugin := range supported {
		wg.Add(1)
		go func(i int, p ports.Plugin) {
			defer wg.Done()
			select {
			case <-ctx.Done():
				results[i] = domain.PluginResult{
					PluginName: p.Name(),
					Priority:   p.Priority(),
					Errors:     []string{ctx.Err().Error()},
				}
				return
			default:
			}

			sem <- struct{}{}
			defer func() { <-sem }()

			start := time.Now()
			res, err := p.Analyze(ctx, pctx)
			res.PluginName = p.Name()
			res.Priority = p.Priority()
			res.Duration = time.Since(start)
			if err != nil {
				a.log.Warn("plugin analyze failed", "plugin", p.Name(), "err", err)
				res.Errors = append(res.Errors, err.Error())
			}
			results[i] = res
		}(i, plugin)
	}

	wg.Wait()
	if ctx.Err() != nil {
		return results, ctx.Err()
	}
	return results, nil
}
