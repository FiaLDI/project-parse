package golang

import (
	"bufio"
	"context"
	"strings"

	"github.com/FiaLDI/project-parse/internal/domain"
	"github.com/FiaLDI/project-parse/internal/plugins/support"
	"github.com/FiaLDI/project-parse/internal/ports"
)

// Plugin analyzes Go projects.
type Plugin struct {
	support.Base
}

// New creates the Go plugin.
func New(cache ports.FileCache) *Plugin {
	return &Plugin{Base: support.Base{PluginName: "golang", PluginPriority: 20, Cache: cache}}
}

func (p *Plugin) Supports(pctx domain.ProjectContext) bool {
	return support.HasMarker(pctx, "go.mod") || len(support.FilesByName(pctx, "go.mod")) > 0
}

func (p *Plugin) Analyze(ctx context.Context, pctx domain.ProjectContext) (domain.PluginResult, error) {
	_ = ctx
	var res domain.PluginResult
	files := support.FilesByName(pctx, "go.mod")
	if len(files) == 0 {
		return res, nil
	}
	mod := files[0]
	for _, f := range files {
		if f.RelPath == "go.mod" {
			mod = f
			break
		}
	}
	res.Evidence = append(res.Evidence, support.Evidence(mod.RelPath, "go module"))

	data, err := support.ReadMeta(p.Cache, mod)
	if err != nil {
		return res, err
	}
	module, goVer := parseGoMod(string(data))
	attrs := map[string]any{}
	if module != "" {
		attrs["module"] = module
	}
	if goVer != "" {
		attrs["go_version"] = goVer
	}
	res.Findings = append(res.Findings, support.Finding(
		"golang.module", domain.CategoryLanguage, 0.95,
		"Go module", "Detected go.mod",
		[]string{mod.RelPath}, attrs,
	))
	return res, nil
}

func parseGoMod(body string) (module, goVer string) {
	sc := bufio.NewScanner(strings.NewReader(body))
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if strings.HasPrefix(line, "module ") {
			module = strings.TrimSpace(strings.TrimPrefix(line, "module "))
		}
		if strings.HasPrefix(line, "go ") {
			goVer = strings.TrimSpace(strings.TrimPrefix(line, "go "))
		}
	}
	return module, goVer
}
