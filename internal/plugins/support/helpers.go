package support

import (
	"strings"

	"github.com/FiaLDI/project-parse/internal/domain"
	"github.com/FiaLDI/project-parse/internal/ports"
)

// Base holds shared plugin metadata and file cache access.
type Base struct {
	PluginName     string
	PluginPriority int
	Cache          ports.FileCache
}

func (b Base) Name() string     { return b.PluginName }
func (b Base) Priority() int    { return b.PluginPriority }

// HasMarker reports whether any scan marker basename matches.
func HasMarker(pctx domain.ProjectContext, markers ...string) bool {
	for _, m := range markers {
		if pctx.MarkerExists(m) {
			return true
		}
	}
	return false
}

// FilesByName returns indexed files with the given basename.
func FilesByName(pctx domain.ProjectContext, name string) []domain.FileMeta {
	if pctx.Files == nil {
		return nil
	}
	return pctx.Files.ByName(name)
}

// FilesByNamePrefix returns files whose basename starts with prefix.
func FilesByNamePrefix(pctx domain.ProjectContext, prefix string) []domain.FileMeta {
	if pctx.Files == nil {
		return nil
	}
	var out []domain.FileMeta
	for _, f := range pctx.Files.All() {
		if strings.HasPrefix(f.Name, prefix) {
			out = append(out, f)
		}
	}
	return out
}

// HasDir reports whether any indexed path starts with or equals dir prefix.
func HasDir(pctx domain.ProjectContext, dir string) bool {
	dir = strings.TrimSuffix(strings.TrimPrefix(dir, "/"), "/")
	if pctx.Files == nil {
		return false
	}
	for _, f := range pctx.Files.All() {
		p := strings.TrimPrefix(f.RelPath, "./")
		if p == dir || strings.HasPrefix(p, dir+"/") {
			return true
		}
	}
	return false
}

// HasPathPrefix reports whether any file path has the given prefix.
func HasPathPrefix(pctx domain.ProjectContext, prefix string) bool {
	prefix = strings.TrimPrefix(prefix, "./")
	if pctx.Files == nil {
		return false
	}
	for _, f := range pctx.Files.All() {
		p := strings.TrimPrefix(f.RelPath, "./")
		if strings.HasPrefix(p, prefix) {
			return true
		}
	}
	return false
}

// ReadMeta reads file contents through the plugin cache.
func ReadMeta(cache ports.FileCache, meta domain.FileMeta) ([]byte, error) {
	if cache == nil {
		return nil, nil
	}
	return cache.Read(meta.Path)
}

// Finding creates a finding with optional attributes.
func Finding(id string, cat domain.Category, conf float64, title, summary string, paths []string, attrs map[string]any) domain.Finding {
	return domain.Finding{
		ID:           id,
		Category:     cat,
		Confidence:   conf,
		Title:        title,
		Summary:      summary,
		Attributes:   attrs,
		RelatedPaths: paths,
	}
}

// Evidence creates evidence for a source path.
func Evidence(path, reason string) domain.Evidence {
	return domain.Evidence{Path: path, Reason: reason}
}
