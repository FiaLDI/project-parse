package analyzer

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/FiaLDI/project-parse/internal/domain"
	"github.com/FiaLDI/project-parse/internal/ports"
)

type stubPlugin struct {
	name     string
	priority int
	support  bool
	delay    time.Duration
	err      error
}

func (s stubPlugin) Name() string  { return s.name }
func (s stubPlugin) Priority() int { return s.priority }
func (s stubPlugin) Supports(domain.ProjectContext) bool {
	return s.support
}
func (s stubPlugin) Analyze(ctx context.Context, _ domain.ProjectContext) (domain.PluginResult, error) {
	if s.delay > 0 {
		select {
		case <-time.After(s.delay):
		case <-ctx.Done():
			return domain.PluginResult{}, ctx.Err()
		}
	}
	if s.err != nil {
		return domain.PluginResult{}, s.err
	}
	return domain.PluginResult{
		Findings: []domain.Finding{{ID: s.name + ".ok", Category: domain.CategoryLanguage, Confidence: 1, Title: s.name}},
	}, nil
}

func TestRunParallelAndSoftFail(t *testing.T) {
	a := New(Options{Jobs: 2})
	pctx := domain.ProjectContext{Files: domain.NewFileIndex()}
	plugins := []ports.Plugin{
		stubPlugin{name: "fast", priority: 10, support: true},
		stubPlugin{name: "slow", priority: 5, support: true, delay: 20 * time.Millisecond},
		stubPlugin{name: "skip", priority: 1, support: false},
		stubPlugin{name: "fail", priority: 1, support: true, err: errors.New("boom")},
	}

	results, err := a.Run(context.Background(), pctx, plugins)
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 3 {
		t.Fatalf("got %d results, want 3 supported plugins", len(results))
	}
	var failed bool
	for _, r := range results {
		if r.PluginName == "fail" && len(r.Errors) == 0 {
			t.Fatal("expected fail plugin to record error")
		}
		if r.PluginName == "fail" {
			failed = true
		}
	}
	if !failed {
		t.Fatal("missing fail plugin result")
	}
}
