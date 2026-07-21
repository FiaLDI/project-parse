package graph

import (
	"path/filepath"
	"sort"
	"strings"

	"github.com/FiaLDI/project-parse/internal/domain"
)

// Builder derives an architecture graph from an aggregated report.
type Builder struct{}

// New creates a graph builder.
func New() *Builder {
	return &Builder{}
}

// Build constructs nodes and edges representing project structure.
func (b *Builder) Build(report domain.Report) (domain.ArchitectureGraph, error) {
	g := domain.ArchitectureGraph{}
	projectID := "project"
	label := filepath.Base(report.Meta.Root)
	if label == "" || label == "." {
		label = "project"
	}
	g.Nodes = append(g.Nodes, domain.Node{
		ID:    projectID,
		Label: label,
		Kind:  domain.NodeKindService,
		Attributes: map[string]any{
			"root": report.Meta.Root,
		},
	})

	for _, c := range report.Stack {
		id := nodeID("stack", c.Kind, c.Name)
		g.Nodes = append(g.Nodes, domain.Node{
			ID:         id,
			Label:      c.Name,
			Kind:       nodeKindForStack(c.Kind),
			Attributes: attrs(c.Attributes, c.Version),
		})
		g.Edges = append(g.Edges, domain.Edge{
			From: projectID,
			To:   id,
			Kind: domain.EdgeKindDependsOn,
		})
	}

	for _, c := range report.Infrastructure {
		id := nodeID("infra", c.Kind, c.Name)
		g.Nodes = append(g.Nodes, domain.Node{
			ID:         id,
			Label:      c.Name,
			Kind:       domain.NodeKindInfra,
			Attributes: attrs(c.Attributes, c.Version),
		})
		g.Edges = append(g.Edges, domain.Edge{
			From: projectID,
			To:   id,
			Kind: domain.EdgeKindDeploys,
		})
	}

	for _, c := range report.Databases {
		id := nodeID("db", c.Kind, c.Name)
		g.Nodes = append(g.Nodes, domain.Node{
			ID:         id,
			Label:      c.Name,
			Kind:       domain.NodeKindDatastore,
			Attributes: attrs(c.Attributes, c.Version),
		})
		g.Edges = append(g.Edges, domain.Edge{
			From: projectID,
			To:   id,
			Kind: domain.EdgeKindReads,
		})
	}

	for _, pattern := range report.Architecture.Patterns {
		id := nodeID("layer", pattern, pattern)
		g.Nodes = append(g.Nodes, domain.Node{
			ID:    id,
			Label: pattern,
			Kind:  domain.NodeKindLayer,
			Attributes: map[string]any{
				"pattern": pattern,
			},
		})
		g.Edges = append(g.Edges, domain.Edge{
			From: projectID,
			To:   id,
			Kind: domain.EdgeKindContains,
		})
	}

	if report.Documentation.HasReadme || report.Documentation.HasDocsDir {
		id := "docs"
		g.Nodes = append(g.Nodes, domain.Node{
			ID:    id,
			Label: "Documentation",
			Kind:  domain.NodeKindExternal,
		})
		g.Edges = append(g.Edges, domain.Edge{
			From: projectID,
			To:   id,
			Kind: domain.EdgeKindDependsOn,
		})
	}

	g = dedupeGraph(g)
	return g, nil
}

func nodeKindForStack(kind string) domain.NodeKind {
	switch kind {
	case "framework":
		return domain.NodeKindModule
	default:
		return domain.NodeKindModule
	}
}

func nodeID(parts ...string) string {
	slug := strings.Join(parts, "-")
	slug = strings.Map(func(r rune) rune {
		switch {
		case r >= 'a' && r <= 'z', r >= 'A' && r <= 'Z', r >= '0' && r <= '9':
			return r
		case r == ' ', r == '_', r == '.', r == '/':
			return '-'
		default:
			return '-'
		}
	}, slug)
	for strings.Contains(slug, "--") {
		slug = strings.ReplaceAll(slug, "--", "-")
	}
	return strings.Trim(slug, "-")
}

func attrs(in map[string]any, version string) map[string]any {
	out := map[string]any{}
	for k, v := range in {
		out[k] = v
	}
	if version != "" {
		out["version"] = version
	}
	if len(out) == 0 {
		return nil
	}
	return out
}

func dedupeGraph(g domain.ArchitectureGraph) domain.ArchitectureGraph {
	nodeIndex := map[string]int{}
	var nodes []domain.Node
	for _, n := range g.Nodes {
		if _, ok := nodeIndex[n.ID]; ok {
			continue
		}
		nodeIndex[n.ID] = len(nodes)
		nodes = append(nodes, n)
	}

	type edgeKey struct{ from, to, kind string }
	seen := map[edgeKey]struct{}{}
	var edges []domain.Edge
	for _, e := range g.Edges {
		if _, ok := nodeIndex[e.From]; !ok {
			continue
		}
		if _, ok := nodeIndex[e.To]; !ok {
			continue
		}
		k := edgeKey{e.From, e.To, string(e.Kind)}
		if _, ok := seen[k]; ok {
			continue
		}
		seen[k] = struct{}{}
		edges = append(edges, e)
	}

	sort.Slice(nodes, func(i, j int) bool { return nodes[i].ID < nodes[j].ID })
	sort.Slice(edges, func(i, j int) bool {
		if edges[i].From == edges[j].From {
			if edges[i].To == edges[j].To {
				return edges[i].Kind < edges[j].Kind
			}
			return edges[i].To < edges[j].To
		}
		return edges[i].From < edges[j].From
	})

	return domain.ArchitectureGraph{Nodes: nodes, Edges: edges}
}
