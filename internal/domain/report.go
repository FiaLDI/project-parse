package domain

import "time"

// ReportMeta holds generation metadata for a report.
type ReportMeta struct {
	Root        string        `json:"root"`
	GeneratedAt time.Time     `json:"generated_at"`
	Duration    time.Duration `json:"duration"`
	ToolVersion string        `json:"tool_version"`
	FileCount   int           `json:"file_count"`
}

// TechComponent describes a detected technology or service.
type TechComponent struct {
	Name       string         `json:"name"`
	Kind       string         `json:"kind"`
	Version    string         `json:"version,omitempty"`
	Confidence float64        `json:"confidence"`
	Attributes map[string]any `json:"attributes,omitempty"`
	Paths      []string       `json:"paths,omitempty"`
}

// ArchAssessment summarizes detected architectural patterns.
type ArchAssessment struct {
	Patterns   []string       `json:"patterns,omitempty"`
	Confidence float64        `json:"confidence"`
	Summary    string         `json:"summary,omitempty"`
	Attributes map[string]any `json:"attributes,omitempty"`
}

// DocsAssessment summarizes documentation coverage signals.
type DocsAssessment struct {
	HasReadme  bool     `json:"has_readme"`
	HasDocsDir bool     `json:"has_docs_dir"`
	Paths      []string `json:"paths,omitempty"`
	Summary    string   `json:"summary,omitempty"`
}

// Report is the aggregated analysis result.
type Report struct {
	Meta           ReportMeta       `json:"meta"`
	Stack          []TechComponent  `json:"stack,omitempty"`
	Architecture   ArchAssessment   `json:"architecture"`
	Infrastructure []TechComponent  `json:"infrastructure,omitempty"`
	Databases      []TechComponent  `json:"databases,omitempty"`
	Documentation  DocsAssessment   `json:"documentation"`
	Findings       []Finding        `json:"findings,omitempty"`
	Warnings       []string         `json:"warnings,omitempty"`
}
