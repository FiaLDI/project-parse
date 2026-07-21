package output

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"html"
	"math"
	"sort"
	"strings"

	"github.com/FiaLDI/project-parse/internal/domain"
)

const (
	svgNodeW = 160.0
	svgNodeH = 44.0
	svgGapX  = 24.0
	svgGapY  = 72.0
	svgPad   = 32.0
)

// SVGRenderer renders an architecture graph as SVG.
type SVGRenderer struct{}

func NewSVG() *SVGRenderer { return &SVGRenderer{} }

func (r *SVGRenderer) Format() string { return "svg" }

func (r *SVGRenderer) Render(_ context.Context, doc domain.RenderDocument) ([]byte, error) {
	if doc.Graph == nil || len(doc.Graph.Nodes) == 0 {
		return []byte(`<svg xmlns="http://www.w3.org/2000/svg" width="320" height="120"><text x="16" y="40" font-family="system-ui,sans-serif" font-size="14">No graph data</text></svg>`), nil
	}
	pos := layoutNodes(*doc.Graph)
	width, height := canvasSize(pos)
	title := doc.Options.Title
	if title == "" {
		title = "Architecture Graph"
	}

	var b strings.Builder
	fmt.Fprintf(&b, `<svg xmlns="http://www.w3.org/2000/svg" width="%.0f" height="%.0f" viewBox="0 0 %.0f %.0f">`, width, height, width, height)
	b.WriteString(`<defs><marker id="arrow" viewBox="0 0 10 10" refX="9" refY="5" markerWidth="6" markerHeight="6" orient="auto-start-reverse"><path d="M 0 0 L 10 5 L 0 10 z" fill="#64748b"/></marker></defs>`)
	fmt.Fprintf(&b, `<text x="%.0f" y="24" font-family="system-ui,sans-serif" font-size="16" font-weight="600" fill="#0f172a">%s</text>`, svgPad, html.EscapeString(title))

	nodeByID := map[string]domain.Node{}
	for _, n := range doc.Graph.Nodes {
		nodeByID[n.ID] = n
	}
	for _, e := range doc.Graph.Edges {
		from, ok1 := pos[e.From]
		to, ok2 := pos[e.To]
		if !ok1 || !ok2 {
			continue
		}
		x1, y1 := from.x+svgNodeW/2, from.y+svgNodeH
		x2, y2 := to.x+svgNodeW/2, to.y
		fmt.Fprintf(&b, `<line x1="%.1f" y1="%.1f" x2="%.1f" y2="%.1f" stroke="#94a3b8" stroke-width="1.5" marker-end="url(#arrow)"/>`, x1, y1, x2, y2)
		midX, midY := (x1+x2)/2, (y1+y2)/2
		fmt.Fprintf(&b, `<text x="%.1f" y="%.1f" font-family="system-ui,sans-serif" font-size="10" fill="#64748b" text-anchor="middle">%s</text>`, midX, midY-4, html.EscapeString(string(e.Kind)))
	}

	for _, n := range doc.Graph.Nodes {
		p, ok := pos[n.ID]
		if !ok {
			continue
		}
		fill, stroke := colorsForKind(n.Kind)
		fmt.Fprintf(&b, `<rect x="%.1f" y="%.1f" width="%.0f" height="%.0f" rx="8" fill="%s" stroke="%s" stroke-width="1.5"/>`, p.x, p.y, svgNodeW, svgNodeH, fill, stroke)
		fmt.Fprintf(&b, `<text x="%.1f" y="%.1f" font-family="system-ui,sans-serif" font-size="12" font-weight="600" fill="#0f172a" text-anchor="middle">%s</text>`, p.x+svgNodeW/2, p.y+svgNodeH/2+4, html.EscapeString(truncate(n.Label, 20)))
		fmt.Fprintf(&b, `<text x="%.1f" y="%.1f" font-family="system-ui,sans-serif" font-size="10" fill="#475569" text-anchor="middle">%s</text>`, p.x+svgNodeW/2, p.y+svgNodeH/2+18, html.EscapeString(string(n.Kind)))
	}

	b.WriteString("</svg>")
	return []byte(b.String()), nil
}

type point struct {
	x, y float64
}

func layoutNodes(g domain.ArchitectureGraph) map[string]point {
	tiers := map[domain.NodeKind][]domain.Node{}
	order := []domain.NodeKind{
		domain.NodeKindService,
		domain.NodeKindModule,
		domain.NodeKindLayer,
		domain.NodeKindInfra,
		domain.NodeKindDatastore,
		domain.NodeKindExternal,
	}
	for _, n := range g.Nodes {
		tiers[n.Kind] = append(tiers[n.Kind], n)
	}
	for k := range tiers {
		sort.Slice(tiers[k], func(i, j int) bool { return tiers[k][i].ID < tiers[k][j].ID })
	}

	pos := map[string]point{}
	row := 0
	for _, kind := range order {
		nodes := tiers[kind]
		if len(nodes) == 0 {
			continue
		}
		rowWidth := float64(len(nodes))*svgNodeW + float64(len(nodes)-1)*svgGapX
		startX := svgPad
		if rowWidth < 480 {
			startX = svgPad + (480-rowWidth)/2
		}
		y := svgPad + 28 + float64(row)*svgGapY
		for i, n := range nodes {
			x := startX + float64(i)*(svgNodeW+svgGapX)
			pos[n.ID] = point{x: x, y: y}
		}
		row++
	}
	return pos
}

func canvasSize(pos map[string]point) (float64, float64) {
	maxX, maxY := 320.0, 160.0
	for _, p := range pos {
		maxX = math.Max(maxX, p.x+svgNodeW+svgPad)
		maxY = math.Max(maxY, p.y+svgNodeH+svgPad)
	}
	return maxX, maxY
}

func colorsForKind(kind domain.NodeKind) (fill, stroke string) {
	switch kind {
	case domain.NodeKindService:
		return "#dbeafe", "#2563eb"
	case domain.NodeKindModule:
		return "#dcfce7", "#16a34a"
	case domain.NodeKindLayer:
		return "#fef3c7", "#d97706"
	case domain.NodeKindInfra:
		return "#ede9fe", "#7c3aed"
	case domain.NodeKindDatastore:
		return "#fce7f3", "#db2777"
	default:
		return "#f1f5f9", "#64748b"
	}
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n-1] + "…"
}

// GraphJSONRenderer renders only the architecture graph as JSON.
type GraphJSONRenderer struct{}

func NewGraphJSON() *GraphJSONRenderer { return &GraphJSONRenderer{} }

func (r *GraphJSONRenderer) Format() string { return "graph-json" }

func (r *GraphJSONRenderer) Render(_ context.Context, doc domain.RenderDocument) ([]byte, error) {
	if doc.Graph == nil {
		doc.Graph = &domain.ArchitectureGraph{}
	}
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetIndent("", "  ")
	if err := enc.Encode(doc.Graph); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
