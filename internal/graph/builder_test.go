package graph

import (
	"testing"

	"github.com/FiaLDI/project-parse/internal/domain"
)

func TestBuildGraph(t *testing.T) {
	b := New()
	report := domain.Report{
		Meta: domain.ReportMeta{Root: "/tmp/demo-app"},
		Stack: []domain.TechComponent{
			{Name: "Go module", Kind: "golang", Confidence: 0.95},
		},
		Infrastructure: []domain.TechComponent{
			{Name: "Docker", Kind: "infra", Confidence: 0.9},
		},
		Databases: []domain.TechComponent{
			{Name: "Prisma schema", Kind: "database", Confidence: 0.9},
		},
		Architecture: domain.ArchAssessment{
			Patterns: []string{"Clean Architecture"},
		},
		Documentation: domain.DocsAssessment{HasReadme: true},
	}

	g, err := b.Build(report)
	if err != nil {
		t.Fatal(err)
	}
	if len(g.Nodes) < 5 {
		t.Fatalf("nodes=%d want >=5: %+v", len(g.Nodes), g.Nodes)
	}
	if len(g.Edges) < 4 {
		t.Fatalf("edges=%d want >=4", len(g.Edges))
	}
}

func TestDedupeGraph(t *testing.T) {
	g := dedupeGraph(domain.ArchitectureGraph{
		Nodes: []domain.Node{
			{ID: "project", Label: "p"},
			{ID: "project", Label: "dup"},
			{ID: "stack-go", Label: "go"},
		},
		Edges: []domain.Edge{
			{From: "project", To: "stack-go", Kind: domain.EdgeKindDependsOn},
			{From: "project", To: "stack-go", Kind: domain.EdgeKindDependsOn},
		},
	})
	if len(g.Nodes) != 2 {
		t.Fatalf("nodes=%d", len(g.Nodes))
	}
	if len(g.Edges) != 1 {
		t.Fatalf("edges=%d", len(g.Edges))
	}
}
