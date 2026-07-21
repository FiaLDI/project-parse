package output

import (
	"context"
	"strings"
	"testing"

	"github.com/FiaLDI/project-parse/internal/domain"
)

func TestSVGRenderer(t *testing.T) {
	doc := domain.RenderDocument{
		Graph: &domain.ArchitectureGraph{
			Nodes: []domain.Node{
				{ID: "project", Label: "demo", Kind: domain.NodeKindService},
				{ID: "stack-go", Label: "Go", Kind: domain.NodeKindModule},
			},
			Edges: []domain.Edge{
				{From: "project", To: "stack-go", Kind: domain.EdgeKindDependsOn},
			},
		},
		Options: domain.RenderOptions{Title: "Demo Graph"},
	}
	data, err := NewSVG().Render(context.Background(), doc)
	if err != nil {
		t.Fatal(err)
	}
	body := string(data)
	for _, want := range []string{"<svg", "Demo Graph", "demo", "depends_on"} {
		if !strings.Contains(body, want) {
			t.Fatalf("missing %q in svg output", want)
		}
	}
}

func TestGraphJSONRenderer(t *testing.T) {
	doc := domain.RenderDocument{
		Graph: &domain.ArchitectureGraph{
			Nodes: []domain.Node{{ID: "project", Label: "demo", Kind: domain.NodeKindService}},
		},
	}
	data, err := NewGraphJSON().Render(context.Background(), doc)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(data), `"id": "project"`) {
		t.Fatalf("unexpected: %s", data)
	}
}

func TestPDFRendererNotImplemented(t *testing.T) {
	_, err := NewPDF().Render(context.Background(), domain.RenderDocument{})
	if err != ErrNotImplemented {
		t.Fatalf("got %v want ErrNotImplemented", err)
	}
}
