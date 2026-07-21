package report

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/FiaLDI/project-parse/internal/domain"
	"github.com/FiaLDI/project-parse/internal/version"
)

// Aggregator merges plugin results into a structured Report.
type Aggregator struct{}

// New creates a report aggregator.
func New() *Aggregator {
	return &Aggregator{}
}

// Aggregate builds a Report from plugin results.
func (a *Aggregator) Aggregate(results []domain.PluginResult) (domain.Report, error) {
	report := domain.Report{
		Meta: domain.ReportMeta{
			GeneratedAt: time.Now().UTC(),
			ToolVersion: version.String(),
		},
		Documentation: domain.DocsAssessment{},
	}

	findings := dedupeFindings(collectFindings(results))
	report.Findings = findings

	for _, r := range results {
		for _, e := range r.Errors {
			report.Warnings = append(report.Warnings, fmt.Sprintf("[%s] %s", r.PluginName, e))
		}
	}

	report.Stack = buildStack(findings)
	report.Infrastructure = buildInfra(findings)
	report.Databases = buildDatabases(findings)
	report.Architecture = buildArchitecture(findings)
	report.Documentation = buildDocumentation(findings)

	var total time.Duration
	for _, r := range results {
		total += r.Duration
	}
	report.Meta.Duration = total
	return report, nil
}

func collectFindings(results []domain.PluginResult) []domain.Finding {
	var out []domain.Finding
	for _, r := range results {
		out = append(out, r.Findings...)
	}
	return out
}

func dedupeFindings(in []domain.Finding) []domain.Finding {
	best := map[string]domain.Finding{}
	for _, f := range in {
		if f.ID == "" {
			continue
		}
		cur, ok := best[f.ID]
		if !ok || f.Confidence > cur.Confidence {
			best[f.ID] = f
		}
	}
	out := make([]domain.Finding, 0, len(best))
	for _, f := range best {
		out = append(out, f)
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].Category == out[j].Category {
			return out[i].ID < out[j].ID
		}
		return out[i].Category < out[j].Category
	})
	return out
}

func buildStack(findings []domain.Finding) []domain.TechComponent {
	var stack []domain.TechComponent
	for _, f := range findings {
		if f.Category != domain.CategoryLanguage && f.Category != domain.CategoryFramework && f.Category != domain.CategoryTooling {
			continue
		}
		if !strings.Contains(f.ID, ".") {
			continue
		}
		parts := strings.SplitN(f.ID, ".", 2)
		kind := parts[0]
		if f.Category == domain.CategoryFramework {
			kind = "framework"
		}
		stack = append(stack, domain.TechComponent{
			Name:       f.Title,
			Kind:       kind,
			Confidence: f.Confidence,
			Attributes: f.Attributes,
			Paths:      f.RelatedPaths,
		})
	}
	sort.Slice(stack, func(i, j int) bool { return stack[i].Name < stack[j].Name })
	return stack
}

func buildInfra(findings []domain.Finding) []domain.TechComponent {
	return componentsForCategories(findings, domain.CategoryInfra, domain.CategoryCI)
}

func buildDatabases(findings []domain.Finding) []domain.TechComponent {
	return componentsForCategories(findings, domain.CategoryDatabase)
}

func componentsForCategories(findings []domain.Finding, cats ...domain.Category) []domain.TechComponent {
	set := map[domain.Category]bool{}
	for _, c := range cats {
		set[c] = true
	}
	var out []domain.TechComponent
	for _, f := range findings {
		if !set[f.Category] {
			continue
		}
		out = append(out, domain.TechComponent{
			Name:       f.Title,
			Kind:       string(f.Category),
			Confidence: f.Confidence,
			Attributes: f.Attributes,
			Paths:      f.RelatedPaths,
		})
	}
	return out
}

func buildArchitecture(findings []domain.Finding) domain.ArchAssessment {
	var patterns []string
	var conf float64
	for _, f := range findings {
		if f.Category != domain.CategoryArch || f.ID == "arch.unknown" {
			continue
		}
		patterns = append(patterns, f.Title)
		if f.Confidence > conf {
			conf = f.Confidence
		}
	}
	sort.Strings(patterns)
	summary := "No known architecture pattern detected"
	if len(patterns) > 0 {
		summary = "Detected: " + strings.Join(patterns, ", ")
	}
	return domain.ArchAssessment{
		Patterns:   patterns,
		Confidence: conf,
		Summary:    summary,
	}
}

func buildDocumentation(findings []domain.Finding) domain.DocsAssessment {
	doc := domain.DocsAssessment{}
	for _, f := range findings {
		if f.Category != domain.CategoryDocs {
			continue
		}
		switch f.ID {
		case "docs.readme":
			doc.HasReadme = true
			doc.Paths = append(doc.Paths, f.RelatedPaths...)
		case "docs.directory":
			doc.HasDocsDir = true
			doc.Paths = append(doc.Paths, f.RelatedPaths...)
		}
	}
	switch {
	case doc.HasReadme && doc.HasDocsDir:
		doc.Summary = "README and docs/ directory present"
	case doc.HasReadme:
		doc.Summary = "README present"
	case doc.HasDocsDir:
		doc.Summary = "docs/ directory present"
	default:
		doc.Summary = "Limited documentation signals"
	}
	return doc
}
