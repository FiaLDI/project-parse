package report

import (
	"testing"

	"github.com/FiaLDI/project-parse/internal/domain"
)

func TestAggregateDedupeAndSections(t *testing.T) {
	agg := New()
	results := []domain.PluginResult{
		{
			PluginName: "golang",
			Findings: []domain.Finding{
				{ID: "golang.module", Category: domain.CategoryLanguage, Confidence: 0.95, Title: "Go module"},
			},
		},
		{
			PluginName: "docker",
			Findings: []domain.Finding{
				{ID: "docker.containerfile", Category: domain.CategoryInfra, Confidence: 0.9, Title: "Docker"},
			},
		},
		{
			PluginName: "architecture",
			Findings: []domain.Finding{
				{ID: "arch.clean", Category: domain.CategoryArch, Confidence: 0.78, Title: "Clean Architecture"},
			},
		},
		{
			PluginName: "fail",
			Errors:     []string{"boom"},
		},
	}

	report, err := agg.Aggregate(results)
	if err != nil {
		t.Fatal(err)
	}
	if len(report.Findings) != 3 {
		t.Fatalf("findings=%d want 3", len(report.Findings))
	}
	if len(report.Stack) == 0 {
		t.Fatal("expected stack entries")
	}
	if len(report.Infrastructure) != 1 {
		t.Fatalf("infra=%d want 1", len(report.Infrastructure))
	}
	if len(report.Architecture.Patterns) != 1 {
		t.Fatalf("patterns=%v", report.Architecture.Patterns)
	}
	if len(report.Warnings) != 1 {
		t.Fatalf("warnings=%v", report.Warnings)
	}
}

func TestDedupeKeepsHigherConfidence(t *testing.T) {
	out := dedupeFindings([]domain.Finding{
		{ID: "x", Category: domain.CategoryLanguage, Confidence: 0.5, Title: "low"},
		{ID: "x", Category: domain.CategoryLanguage, Confidence: 0.9, Title: "high"},
	})
	if len(out) != 1 || out[0].Title != "high" {
		t.Fatalf("got %+v", out)
	}
}
