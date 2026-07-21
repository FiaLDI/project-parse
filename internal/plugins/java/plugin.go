package java

import (
	"context"
	"regexp"

	"github.com/FiaLDI/project-parse/internal/domain"
	"github.com/FiaLDI/project-parse/internal/plugins/support"
	"github.com/FiaLDI/project-parse/internal/ports"
)

var (
	mavenArtifact = regexp.MustCompile(`<artifactId>([^<]+)</artifactId>`)
	gradlePlugin  = regexp.MustCompile(`id\s+['"]([^'"]+)['"]`)
)

// Plugin analyzes Java projects.
type Plugin struct {
	support.Base
}

// New creates the Java plugin.
func New(cache ports.FileCache) *Plugin {
	return &Plugin{Base: support.Base{PluginName: "java", PluginPriority: 20, Cache: cache}}
}

func (p *Plugin) Supports(pctx domain.ProjectContext) bool {
	return support.HasMarker(pctx, "pom.xml", "build.gradle") ||
		len(support.FilesByName(pctx, "pom.xml")) > 0 ||
		len(support.FilesByNamePrefix(pctx, "build.gradle")) > 0
}

func (p *Plugin) Analyze(ctx context.Context, pctx domain.ProjectContext) (domain.PluginResult, error) {
	_ = ctx
	var res domain.PluginResult

	for _, f := range support.FilesByName(pctx, "pom.xml") {
		res.Evidence = append(res.Evidence, support.Evidence(f.RelPath, "Maven POM"))
		attrs := map[string]any{"build": "maven"}
		if data, err := support.ReadMeta(p.Cache, f); err == nil {
			if m := mavenArtifact.FindSubmatch(data); len(m) == 2 {
				attrs["artifact"] = string(m[1])
			}
		} else {
			return res, err
		}
		res.Findings = append(res.Findings, support.Finding(
			"java.maven", domain.CategoryLanguage, 0.92,
			"Java (Maven)", "Detected pom.xml",
			[]string{f.RelPath}, attrs,
		))
	}

	for _, f := range support.FilesByNamePrefix(pctx, "build.gradle") {
		res.Evidence = append(res.Evidence, support.Evidence(f.RelPath, "Gradle build"))
		attrs := map[string]any{"build": "gradle"}
		if data, err := support.ReadMeta(p.Cache, f); err == nil {
			if m := gradlePlugin.FindSubmatch(data); len(m) == 2 {
				attrs["plugin"] = string(m[1])
			}
		} else {
			return res, err
		}
		res.Findings = append(res.Findings, support.Finding(
			"java.gradle", domain.CategoryLanguage, 0.92,
			"Java (Gradle)", "Detected Gradle build file",
			[]string{f.RelPath}, attrs,
		))
	}
	return res, nil
}
