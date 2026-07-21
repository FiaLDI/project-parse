package python

import (
	"context"
	"regexp"
	"strings"

	"github.com/FiaLDI/project-parse/internal/domain"
	"github.com/FiaLDI/project-parse/internal/plugins/support"
	"github.com/FiaLDI/project-parse/internal/ports"
)

var pyProjectName = regexp.MustCompile(`(?m)^name\s*=\s*"([^"]+)"`)

// Plugin analyzes Python projects.
type Plugin struct {
	support.Base
}

// New creates the Python plugin.
func New(cache ports.FileCache) *Plugin {
	return &Plugin{Base: support.Base{PluginName: "python", PluginPriority: 20, Cache: cache}}
}

func (p *Plugin) Supports(pctx domain.ProjectContext) bool {
	return support.HasMarker(pctx, "requirements.txt", "pyproject.toml", "poetry.lock", "uv.lock")
}

func (p *Plugin) Analyze(ctx context.Context, pctx domain.ProjectContext) (domain.PluginResult, error) {
	_ = ctx
	var res domain.PluginResult
	res.Findings = append(res.Findings, support.Finding(
		"python.runtime", domain.CategoryLanguage, 0.9,
		"Python project", "Detected Python dependency manifest",
		nil, nil,
	))

	switch {
	case support.HasMarker(pctx, "uv.lock"):
		res.Findings = append(res.Findings, toolFinding("uv", "uv.lock"))
	case support.HasMarker(pctx, "poetry.lock"):
		res.Findings = append(res.Findings, toolFinding("poetry", "poetry.lock"))
	case support.HasMarker(pctx, "pyproject.toml"):
		res.Findings = append(res.Findings, toolFinding("pyproject", "pyproject.toml"))
	case support.HasMarker(pctx, "requirements.txt"):
		res.Findings = append(res.Findings, toolFinding("pip", "requirements.txt"))
	}

	for _, f := range support.FilesByName(pctx, "pyproject.toml") {
		res.Evidence = append(res.Evidence, support.Evidence(f.RelPath, "pyproject.toml"))
		data, err := support.ReadMeta(p.Cache, f)
		if err != nil {
			return res, err
		}
		if m := pyProjectName.FindSubmatch(data); len(m) == 2 {
			res.Findings = append(res.Findings, support.Finding(
				"python.project", domain.CategoryLanguage, 0.85,
				"Python package: "+string(m[1]), "Parsed name from pyproject.toml",
				[]string{f.RelPath}, map[string]any{"name": string(m[1])},
			))
		}
		body := string(data)
		if strings.Contains(body, "[tool.django") || strings.Contains(body, "django>=") {
			res.Findings = append(res.Findings, frameworkFinding("django", f.RelPath))
		}
		if strings.Contains(body, "fastapi") {
			res.Findings = append(res.Findings, frameworkFinding("fastapi", f.RelPath))
		}
	}
	return res, nil
}

func toolFinding(tool, marker string) domain.Finding {
	return support.Finding(
		"python.tool."+tool, domain.CategoryTooling, 0.88,
		"Python tooling: "+tool, "Detected "+marker,
		[]string{marker}, map[string]any{"tool": tool},
	)
}

func frameworkFinding(name, path string) domain.Finding {
	return support.Finding(
		"python.framework."+name, domain.CategoryFramework, 0.8,
		name, "Detected in pyproject.toml",
		[]string{path}, map[string]any{"framework": name},
	)
}
