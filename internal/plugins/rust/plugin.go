package rust

import (
	"context"
	"regexp"

	"github.com/FiaLDI/project-parse/internal/domain"
	"github.com/FiaLDI/project-parse/internal/plugins/support"
	"github.com/FiaLDI/project-parse/internal/ports"
)

var cargoName = regexp.MustCompile(`(?m)^name\s*=\s*"([^"]+)"`)

// Plugin analyzes Rust projects.
type Plugin struct {
	support.Base
}

// New creates the Rust plugin.
func New(cache ports.FileCache) *Plugin {
	return &Plugin{Base: support.Base{PluginName: "rust", PluginPriority: 20, Cache: cache}}
}

func (p *Plugin) Supports(pctx domain.ProjectContext) bool {
	return support.HasMarker(pctx, "Cargo.toml") || len(support.FilesByName(pctx, "Cargo.toml")) > 0
}

func (p *Plugin) Analyze(ctx context.Context, pctx domain.ProjectContext) (domain.PluginResult, error) {
	_ = ctx
	var res domain.PluginResult
	files := support.FilesByName(pctx, "Cargo.toml")
	if len(files) == 0 {
		return res, nil
	}
	cargo := files[0]
	res.Evidence = append(res.Evidence, support.Evidence(cargo.RelPath, "Cargo manifest"))
	attrs := map[string]any{}
	if data, err := support.ReadMeta(p.Cache, cargo); err == nil {
		if m := cargoName.FindSubmatch(data); len(m) == 2 {
			attrs["name"] = string(m[1])
		}
	} else {
		return res, err
	}
	res.Findings = append(res.Findings, support.Finding(
		"rust.cargo", domain.CategoryLanguage, 0.95,
		"Rust project", "Detected Cargo.toml",
		[]string{cargo.RelPath}, attrs,
	))
	return res, nil
}
