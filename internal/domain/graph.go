package domain

// NodeKind classifies an architecture graph node.
type NodeKind string

const (
	NodeKindService   NodeKind = "service"
	NodeKindLayer     NodeKind = "layer"
	NodeKindModule    NodeKind = "module"
	NodeKindDatastore NodeKind = "datastore"
	NodeKindInfra     NodeKind = "infra"
	NodeKindExternal  NodeKind = "external"
)

// EdgeKind classifies a relationship between nodes.
type EdgeKind string

const (
	EdgeKindDependsOn EdgeKind = "depends_on"
	EdgeKindContains  EdgeKind = "contains"
	EdgeKindDeploys   EdgeKind = "deploys"
	EdgeKindReads     EdgeKind = "reads"
	EdgeKindWrites    EdgeKind = "writes"
)

// Node is a vertex in the architecture graph.
type Node struct {
	ID         string         `json:"id"`
	Label      string         `json:"label"`
	Kind       NodeKind       `json:"kind"`
	Attributes map[string]any `json:"attributes,omitempty"`
}

// Edge is a directed relationship between nodes.
type Edge struct {
	From       string         `json:"from"`
	To         string         `json:"to"`
	Kind       EdgeKind       `json:"kind"`
	Attributes map[string]any `json:"attributes,omitempty"`
}

// ArchitectureGraph is a structural view derived from a Report.
type ArchitectureGraph struct {
	Nodes []Node `json:"nodes,omitempty"`
	Edges []Edge `json:"edges,omitempty"`
}

// RenderOptions controls output rendering.
type RenderOptions struct {
	IncludeEvidence bool              `json:"include_evidence"`
	Title           string            `json:"title,omitempty"`
	Extra           map[string]string `json:"extra,omitempty"`
}

// RenderDocument is the payload passed to format renderers.
type RenderDocument struct {
	Report  Report             `json:"report"`
	Graph   *ArchitectureGraph `json:"graph,omitempty"`
	Options RenderOptions      `json:"options"`
}
