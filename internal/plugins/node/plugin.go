package node

import (
	"context"
	"encoding/json"

	"github.com/FiaLDI/project-parse/internal/domain"
	"github.com/FiaLDI/project-parse/internal/plugins/support"
	"github.com/FiaLDI/project-parse/internal/ports"
)

// Plugin analyzes Node.js / JavaScript projects.
type Plugin struct {
	support.Base
}

// New creates the Node.js plugin.
func New(cache ports.FileCache) *Plugin {
	return &Plugin{Base: support.Base{PluginName: "node", PluginPriority: 20, Cache: cache}}
}

func (p *Plugin) Supports(pctx domain.ProjectContext) bool {
	return support.HasMarker(pctx, "package.json") || len(support.FilesByName(pctx, "package.json")) > 0
}

func (p *Plugin) Analyze(ctx context.Context, pctx domain.ProjectContext) (domain.PluginResult, error) {
	_ = ctx
	var res domain.PluginResult
	pkgs := support.FilesByName(pctx, "package.json")
	if len(pkgs) == 0 {
		return res, nil
	}

	rootPkg := pkgs[0]
	for _, pkg := range pkgs {
		if pkg.RelPath == "package.json" {
			rootPkg = pkg
			break
		}
	}

	res.Evidence = append(res.Evidence, support.Evidence(rootPkg.RelPath, "package manifest"))
	res.Findings = append(res.Findings, support.Finding(
		"node.runtime", domain.CategoryLanguage, 0.95,
		"Node.js project", "Detected package.json manifest",
		[]string{rootPkg.RelPath}, nil,
	))

	pm := detectPackageManager(pctx)
	if pm != "" {
		res.Findings = append(res.Findings, support.Finding(
			"node.package_manager."+pm, domain.CategoryTooling, 0.9,
			"Package manager: "+pm, "Detected "+pm+" lockfile or config",
			nil, map[string]any{"manager": pm},
		))
	}

	if support.HasMarker(pctx, "tsconfig.json") || len(support.FilesByName(pctx, "tsconfig.json")) > 0 {
		res.Findings = append(res.Findings, support.Finding(
			"node.typescript", domain.CategoryTooling, 0.9,
			"TypeScript", "Detected tsconfig.json",
			pathsForName(pctx, "tsconfig.json"), nil,
		))
	}

	for _, prefix := range []struct {
		id, label string
	}{
		{"node.framework.next", "Next.js"},
		{"node.framework.vite", "Vite"},
		{"node.framework.webpack", "Webpack"},
		{"node.framework.nest", "NestJS"},
	} {
		if marker, ok := frameworkMarker(prefix.id); ok && support.HasMarker(pctx, marker) {
			res.Findings = append(res.Findings, support.Finding(
				prefix.id, domain.CategoryFramework, 0.85,
				prefix.label, "Detected "+marker,
				nil, map[string]any{"marker": marker},
			))
		}
	}

	if support.HasPathPrefix(pctx, "prisma/") || support.HasMarker(pctx, "schema.prisma") {
		res.Findings = append(res.Findings, support.Finding(
			"node.orm.prisma", domain.CategoryTooling, 0.85,
			"Prisma ORM", "Detected Prisma schema",
			prismaPaths(pctx), nil,
		))
	}

	data, err := support.ReadMeta(p.Cache, rootPkg)
	if err != nil || len(data) == 0 {
		return res, err
	}

	var manifest struct {
		Name         string            `json:"name"`
		Dependencies map[string]string `json:"dependencies"`
		DevDeps      map[string]string `json:"devDependencies"`
	}
	if json.Unmarshal(data, &manifest) == nil {
		deps := mergeDeps(manifest.Dependencies, manifest.DevDeps)
		for pkg, id := range map[string]string{
			"react": "node.framework.react", "vue": "node.framework.vue",
			"express": "node.framework.express", "fastify": "node.framework.fastify",
		} {
			if _, ok := deps[pkg]; ok {
				res.Findings = append(res.Findings, support.Finding(
					id, domain.CategoryFramework, 0.8,
					pkg, "Listed in package.json dependencies",
					[]string{rootPkg.RelPath}, map[string]any{"package": pkg},
				))
			}
		}
		if manifest.Name != "" {
			for i := range res.Findings {
				if res.Findings[i].ID == "node.runtime" {
					if res.Findings[i].Attributes == nil {
						res.Findings[i].Attributes = map[string]any{}
					}
					res.Findings[i].Attributes["name"] = manifest.Name
					break
				}
			}
		}
	}
	return res, nil
}

func detectPackageManager(pctx domain.ProjectContext) string {
	switch {
	case support.HasMarker(pctx, "pnpm-lock.yaml"):
		return "pnpm"
	case support.HasMarker(pctx, "yarn.lock"):
		return "yarn"
	case support.HasMarker(pctx, "package-lock.json"):
		return "npm"
	default:
		return ""
	}
}

func frameworkMarker(id string) (string, bool) {
	switch id {
	case "node.framework.next":
		return "next.config.*", true
	case "node.framework.vite":
		return "vite.config.*", true
	case "node.framework.webpack":
		return "webpack.config.*", true
	case "node.framework.nest":
		return "nest-cli.json", true
	default:
		return "", false
	}
}

func pathsForName(pctx domain.ProjectContext, name string) []string {
	var paths []string
	for _, f := range support.FilesByName(pctx, name) {
		paths = append(paths, f.RelPath)
	}
	return paths
}

func prismaPaths(pctx domain.ProjectContext) []string {
	var paths []string
	for _, f := range support.FilesByName(pctx, "schema.prisma") {
		paths = append(paths, f.RelPath)
	}
	return paths
}

func mergeDeps(a, b map[string]string) map[string]string {
	out := map[string]string{}
	for k, v := range a {
		out[k] = v
	}
	for k, v := range b {
		out[k] = v
	}
	return out
}
