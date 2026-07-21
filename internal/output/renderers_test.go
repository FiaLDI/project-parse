package output

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/FiaLDI/project-parse/internal/domain"
)

func TestJSONRenderer(t *testing.T) {
	doc := sampleDoc()
	data, err := NewJSON().Render(context.Background(), doc)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(data), `"root": "/tmp/demo"`) {
		t.Fatalf("unexpected json: %s", data)
	}
}

func TestMarkdownRenderer(t *testing.T) {
	data, err := NewMarkdown().Render(context.Background(), sampleDoc())
	if err != nil {
		t.Fatal(err)
	}
	body := string(data)
	for _, want := range []string{"# Project Report", "## Technology Stack", "Go module"} {
		if !strings.Contains(body, want) {
			t.Fatalf("missing %q in:\n%s", want, body)
		}
	}
}

func TestHTMLRenderer(t *testing.T) {
	data, err := NewHTML().Render(context.Background(), sampleDoc())
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(data), "<html>") || !strings.Contains(string(data), "Go module") {
		t.Fatalf("unexpected html: %s", data)
	}
}

func sampleDoc() domain.RenderDocument {
	return domain.RenderDocument{
		Report: domain.Report{
			Meta: domain.ReportMeta{
				Root:        "/tmp/demo",
				GeneratedAt: time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC),
				FileCount:   10,
				ToolVersion: "test",
			},
			Stack: []domain.TechComponent{{Name: "Go module", Kind: "golang", Confidence: 0.95}},
			Architecture: domain.ArchAssessment{Summary: "Detected: Clean Architecture"},
			Documentation: domain.DocsAssessment{Summary: "README present", HasReadme: true},
		},
		Options: domain.RenderOptions{Title: "Project Report"},
	}
}
