package domain

import "time"

// Finding is a single analytical fact produced by a plugin.
type Finding struct {
	ID           string         `json:"id"`
	Category     Category       `json:"category"`
	Confidence   float64        `json:"confidence"`
	Title        string         `json:"title"`
	Summary      string         `json:"summary"`
	Attributes   map[string]any `json:"attributes,omitempty"`
	Tags         []string       `json:"tags,omitempty"`
	RelatedPaths []string       `json:"related_paths,omitempty"`
}

// Evidence points at a source file that justified a finding.
type Evidence struct {
	Path   string `json:"path"`
	Reason string `json:"reason,omitempty"`
}

// PluginResult is the outcome of a single plugin Analyze call.
type PluginResult struct {
	PluginName string        `json:"plugin_name"`
	Priority   int           `json:"priority"`
	Findings   []Finding     `json:"findings,omitempty"`
	Evidence   []Evidence    `json:"evidence,omitempty"`
	Errors     []string      `json:"errors,omitempty"`
	Duration   time.Duration `json:"duration"`
}
