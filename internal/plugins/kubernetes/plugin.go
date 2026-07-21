package kubernetes

import (
	"context"
	"strings"

	"github.com/FiaLDI/project-parse/internal/domain"
	"github.com/FiaLDI/project-parse/internal/plugins/support"
	"github.com/FiaLDI/project-parse/internal/ports"
)

var k8sKinds = []string{
	"apiVersion:", "kind: Deployment", "kind: Service", "kind: Ingress",
	"kind: StatefulSet", "kind: ConfigMap", "kind: Secret", "kind: Pod",
}

// Plugin analyzes Kubernetes manifests and Helm charts.
type Plugin struct {
	support.Base
}

// New creates the Kubernetes plugin.
func New(cache ports.FileCache) *Plugin {
	return &Plugin{Base: support.Base{PluginName: "kubernetes", PluginPriority: 15, Cache: cache}}
}

func (p *Plugin) Supports(pctx domain.ProjectContext) bool {
	if support.HasMarker(pctx, "helm/") || support.HasDir(pctx, "helm") {
		return true
	}
	for _, f := range pctx.Files.All() {
		ext := strings.ToLower(f.Ext)
		if ext != ".yaml" && ext != ".yml" {
			continue
		}
		if looksLikeK8s(f.RelPath) {
			return true
		}
	}
	return false
}

func (p *Plugin) Analyze(ctx context.Context, pctx domain.ProjectContext) (domain.PluginResult, error) {
	_ = ctx
	var res domain.PluginResult

	if support.HasDir(pctx, "helm") || support.HasMarker(pctx, "helm/") {
		res.Findings = append(res.Findings, support.Finding(
			"k8s.helm", domain.CategoryInfra, 0.88,
			"Helm charts", "Detected helm/ directory",
			[]string{"helm/"}, nil,
		))
	}

	var manifestPaths []string
	for _, f := range pctx.Files.All() {
		ext := strings.ToLower(f.Ext)
		if ext != ".yaml" && ext != ".yml" {
			continue
		}
		data, err := support.ReadMeta(p.Cache, f)
		if err != nil {
			return res, err
		}
		body := string(data)
		if !containsK8sSignal(body) {
			continue
		}
		manifestPaths = append(manifestPaths, f.RelPath)
		res.Evidence = append(res.Evidence, support.Evidence(f.RelPath, "k8s manifest"))
	}
	if len(manifestPaths) > 0 {
		res.Findings = append(res.Findings, support.Finding(
			"k8s.manifests", domain.CategoryInfra, 0.85,
			"Kubernetes manifests", "Detected YAML files with Kubernetes resources",
			manifestPaths, map[string]any{"count": len(manifestPaths)},
		))
	}
	return res, nil
}

func looksLikeK8s(path string) bool {
	lower := strings.ToLower(path)
	return strings.Contains(lower, "/k8s/") || strings.Contains(lower, "/kubernetes/") ||
		strings.Contains(lower, "/deploy/") || strings.Contains(lower, "/manifests/")
}

func containsK8sSignal(body string) bool {
	if !strings.Contains(body, "apiVersion:") {
		return false
	}
	for _, k := range k8sKinds[1:] {
		if strings.Contains(body, k) {
			return true
		}
	}
	return strings.Contains(body, "kind:")
}
