package docker

import (
	"context"
	"strings"

	"github.com/FiaLDI/project-parse/internal/domain"
	"github.com/FiaLDI/project-parse/internal/plugins/support"
	"github.com/FiaLDI/project-parse/internal/ports"
)

// Plugin analyzes Docker and Compose files.
type Plugin struct {
	support.Base
}

// New creates the Docker plugin.
func New(cache ports.FileCache) *Plugin {
	return &Plugin{Base: support.Base{PluginName: "docker", PluginPriority: 15, Cache: cache}}
}

func (p *Plugin) Supports(pctx domain.ProjectContext) bool {
	if support.HasMarker(pctx, "Dockerfile", "docker-compose.yml", "docker-compose.yaml", "compose.yaml", "compose.yml") {
		return true
	}
	return len(support.FilesByNamePrefix(pctx, "Dockerfile")) > 0
}

func (p *Plugin) Analyze(ctx context.Context, pctx domain.ProjectContext) (domain.PluginResult, error) {
	_ = ctx
	var res domain.PluginResult

	for _, f := range pctx.Files.All() {
		name := strings.ToLower(f.Name)
		switch {
		case name == "dockerfile" || strings.HasPrefix(name, "dockerfile."):
			res.Evidence = append(res.Evidence, support.Evidence(f.RelPath, "Dockerfile"))
			res.Findings = append(res.Findings, support.Finding(
				"docker.containerfile", domain.CategoryInfra, 0.9,
				"Docker", "Detected Dockerfile",
				[]string{f.RelPath}, map[string]any{"kind": "dockerfile"},
			))
		case name == "docker-compose.yml", name == "docker-compose.yaml", name == "compose.yaml", name == "compose.yml":
			res.Evidence = append(res.Evidence, support.Evidence(f.RelPath, "Compose file"))
			res.Findings = append(res.Findings, support.Finding(
				"docker.compose", domain.CategoryInfra, 0.9,
				"Docker Compose", "Detected compose manifest",
				[]string{f.RelPath}, map[string]any{"kind": "compose"},
			))
		}
	}
	return res, nil
}
