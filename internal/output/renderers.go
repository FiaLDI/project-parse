package output

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"html"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/FiaLDI/project-parse/internal/domain"
)

// JSONRenderer renders reports as JSON.
type JSONRenderer struct{}

func NewJSON() *JSONRenderer { return &JSONRenderer{} }

func (r *JSONRenderer) Format() string { return "json" }

func (r *JSONRenderer) Render(_ context.Context, doc domain.RenderDocument) ([]byte, error) {
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetIndent("", "  ")
	if err := enc.Encode(doc.Report); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// MarkdownRenderer renders reports as Markdown.
type MarkdownRenderer struct{}

func NewMarkdown() *MarkdownRenderer { return &MarkdownRenderer{} }

func (r *MarkdownRenderer) Format() string { return "markdown" }

func (r *MarkdownRenderer) Render(_ context.Context, doc domain.RenderDocument) ([]byte, error) {
	rep := doc.Report
	var b strings.Builder
	title := doc.Options.Title
	if title == "" {
		title = "Project Report"
	}
	fmt.Fprintf(&b, "# %s\n\n", title)
	fmt.Fprintf(&b, "- **Root:** `%s`\n", rep.Meta.Root)
	fmt.Fprintf(&b, "- **Generated:** %s\n", rep.Meta.GeneratedAt.Format(time.RFC3339))
	fmt.Fprintf(&b, "- **Files scanned:** %d\n", rep.Meta.FileCount)
	fmt.Fprintf(&b, "- **Tool:** %s\n\n", rep.Meta.ToolVersion)

	if len(rep.Stack) > 0 {
		b.WriteString("## Technology Stack\n\n")
		for _, c := range rep.Stack {
			fmt.Fprintf(&b, "- **%s** (%s, confidence %.0f%%)\n", c.Name, c.Kind, c.Confidence*100)
		}
		b.WriteString("\n")
	}

	if rep.Architecture.Summary != "" {
		b.WriteString("## Architecture\n\n")
		fmt.Fprintf(&b, "%s\n\n", rep.Architecture.Summary)
	}

	if len(rep.Infrastructure) > 0 {
		b.WriteString("## Infrastructure\n\n")
		for _, c := range rep.Infrastructure {
			fmt.Fprintf(&b, "- **%s** (confidence %.0f%%)\n", c.Name, c.Confidence*100)
		}
		b.WriteString("\n")
	}

	if len(rep.Databases) > 0 {
		b.WriteString("## Databases\n\n")
		for _, c := range rep.Databases {
			fmt.Fprintf(&b, "- **%s**\n", c.Name)
		}
		b.WriteString("\n")
	}

	b.WriteString("## Documentation\n\n")
	fmt.Fprintf(&b, "%s\n\n", rep.Documentation.Summary)

	if len(rep.Findings) > 0 {
		b.WriteString("## Findings\n\n")
		for _, f := range rep.Findings {
			fmt.Fprintf(&b, "### %s\n\n", f.Title)
			fmt.Fprintf(&b, "- **ID:** `%s`\n", f.ID)
			fmt.Fprintf(&b, "- **Category:** %s\n", f.Category)
			fmt.Fprintf(&b, "- **Confidence:** %.0f%%\n", f.Confidence*100)
			if f.Summary != "" {
				fmt.Fprintf(&b, "- **Summary:** %s\n", f.Summary)
			}
			b.WriteString("\n")
		}
	}

	if len(rep.Warnings) > 0 {
		b.WriteString("## Warnings\n\n")
		for _, w := range rep.Warnings {
			fmt.Fprintf(&b, "- %s\n", w)
		}
	}
	return []byte(b.String()), nil
}

// HTMLRenderer renders reports as HTML.
type HTMLRenderer struct{}

func NewHTML() *HTMLRenderer { return &HTMLRenderer{} }

func (r *HTMLRenderer) Format() string { return "html" }

func (r *HTMLRenderer) Render(_ context.Context, doc domain.RenderDocument) ([]byte, error) {
	rep := doc.Report
	title := doc.Options.Title
	if title == "" {
		title = "Project Report"
	}
	var b strings.Builder
	b.WriteString("<!DOCTYPE html><html><head><meta charset=\"utf-8\">")
	fmt.Fprintf(&b, "<title>%s</title>", html.EscapeString(title))
	b.WriteString(`<style>
body{font-family:system-ui,sans-serif;max-width:960px;margin:2rem auto;padding:0 1rem;line-height:1.5}
h1,h2{border-bottom:1px solid #ddd;padding-bottom:.3rem}
table{border-collapse:collapse;width:100%;margin:1rem 0}
th,td{border:1px solid #ddd;padding:.5rem;text-align:left}
.warn{color:#a40}
</style></head><body>`)
	fmt.Fprintf(&b, "<h1>%s</h1>", html.EscapeString(title))
	fmt.Fprintf(&b, "<p><strong>Root:</strong> <code>%s</code><br>", html.EscapeString(rep.Meta.Root))
	fmt.Fprintf(&b, "<strong>Generated:</strong> %s<br>", html.EscapeString(rep.Meta.GeneratedAt.Format(time.RFC3339)))
	fmt.Fprintf(&b, "<strong>Files:</strong> %d</p>", rep.Meta.FileCount)

	writeSection := func(name string, items []domain.TechComponent) {
		if len(items) == 0 {
			return
		}
		fmt.Fprintf(&b, "<h2>%s</h2><table><tr><th>Name</th><th>Kind</th><th>Confidence</th></tr>", html.EscapeString(name))
		for _, c := range items {
			fmt.Fprintf(&b, "<tr><td>%s</td><td>%s</td><td>%.0f%%</td></tr>",
				html.EscapeString(c.Name), html.EscapeString(c.Kind), c.Confidence*100)
		}
		b.WriteString("</table>")
	}
	writeSection("Technology Stack", rep.Stack)
	if rep.Architecture.Summary != "" {
		fmt.Fprintf(&b, "<h2>Architecture</h2><p>%s</p>", html.EscapeString(rep.Architecture.Summary))
	}
	writeSection("Infrastructure", rep.Infrastructure)
	writeSection("Databases", rep.Databases)

	fmt.Fprintf(&b, "<h2>Documentation</h2><p>%s</p>", html.EscapeString(rep.Documentation.Summary))

	if len(rep.Findings) > 0 {
		b.WriteString("<h2>Findings</h2><table><tr><th>Title</th><th>Category</th><th>Confidence</th></tr>")
		for _, f := range rep.Findings {
			fmt.Fprintf(&b, "<tr><td>%s</td><td>%s</td><td>%.0f%%</td></tr>",
				html.EscapeString(f.Title), html.EscapeString(string(f.Category)), f.Confidence*100)
		}
		b.WriteString("</table>")
	}
	if len(rep.Warnings) > 0 {
		b.WriteString("<h2>Warnings</h2><ul>")
		for _, w := range rep.Warnings {
			fmt.Fprintf(&b, "<li class=\"warn\">%s</li>", html.EscapeString(w))
		}
		b.WriteString("</ul>")
	}
	b.WriteString("</body></html>")
	return []byte(b.String()), nil
}

// WriteFile writes rendered bytes to outDir using the standard filename for format.
func WriteFile(outDir, format string, data []byte) (string, error) {
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		return "", err
	}
	name := filenameForFormat(format)
	path := filepath.Join(outDir, name)
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return "", err
	}
	return path, nil
}

func filenameForFormat(format string) string {
	switch format {
	case "json":
		return "report.json"
	case "markdown", "md":
		return "report.md"
	case "html":
		return "report.html"
	case "svg":
		return "graph.svg"
	case "graph-json":
		return "graph.json"
	case "pdf":
		return "report.pdf"
	default:
		return "report." + format
	}
}
