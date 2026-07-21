package scanner

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/FiaLDI/project-parse/internal/domain"
	"github.com/FiaLDI/project-parse/internal/ports"
)

// Scanner walks a project tree and builds a ProjectContext file index.
type Scanner struct{}

// New creates a filesystem scanner.
func New() *Scanner {
	return &Scanner{}
}

// Scan walks root according to opts and returns an indexed ProjectContext.
func (s *Scanner) Scan(ctx context.Context, root string, opts ports.ScanOptions) (domain.ProjectContext, error) {
	if root == "" {
		root = "."
	}
	absRoot, err := filepath.Abs(root)
	if err != nil {
		return domain.ProjectContext{}, fmt.Errorf("resolve root: %w", err)
	}
	info, err := os.Stat(absRoot)
	if err != nil {
		return domain.ProjectContext{}, fmt.Errorf("stat root: %w", err)
	}
	if !info.IsDir() {
		return domain.ProjectContext{}, fmt.Errorf("root is not a directory: %s", absRoot)
	}

	exclude := newGlobMatcher(opts.Exclude)
	include := newGlobMatcher(opts.Include)
	useInclude := len(opts.Include) > 0

	idx := domain.NewFileIndex()
	markerSet := map[string]struct{}{}
	started := time.Now()

	err = filepath.WalkDir(absRoot, func(path string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			if d != nil && d.IsDir() {
				return fs.SkipDir
			}
			return nil
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		rel, relErr := filepath.Rel(absRoot, path)
		if relErr != nil {
			return nil
		}
		if rel == "." {
			return nil
		}
		relSlash := filepath.ToSlash(rel)

		if d.IsDir() {
			if !opts.FollowSymlinks && isSymlink(d) {
				return fs.SkipDir
			}
			if shouldSkipDir(relSlash, exclude) {
				return fs.SkipDir
			}
			return nil
		}

		if !opts.FollowSymlinks && isSymlink(d) {
			return nil
		}
		if exclude.Match(relSlash) {
			return nil
		}
		if useInclude && !include.Match(relSlash) {
			return nil
		}

		fi, infoErr := d.Info()
		if infoErr != nil {
			return nil
		}

		idx.Add(domain.FileMeta{
			Path:    path,
			RelPath: relSlash,
			Size:    fi.Size(),
			ModTime: fi.ModTime(),
			Ext:     strings.ToLower(filepath.Ext(d.Name())),
			Name:    d.Name(),
			IsDir:   false,
		})
		recordMarkers(markerSet, relSlash, d.Name())
		return nil
	})
	if err != nil {
		return domain.ProjectContext{}, err
	}

	markers := make([]string, 0, len(markerSet))
	for m := range markerSet {
		markers = append(markers, m)
	}
	sort.Strings(markers)
	idx.SetMarkers(markers)

	return domain.ProjectContext{
		Root:      absRoot,
		Files:     idx,
		StartedAt: started,
	}, nil
}

func isSymlink(d fs.DirEntry) bool {
	return d.Type()&fs.ModeSymlink != 0
}

func shouldSkipDir(relSlash string, exclude *globMatcher) bool {
	if exclude == nil || len(exclude.patterns) == 0 {
		return false
	}
	if exclude.Match(relSlash) {
		return true
	}
	for _, pattern := range exclude.patterns {
		if dirExcludedByPattern(pattern, relSlash) {
			return true
		}
	}
	return false
}

// dirExcludedByPattern reports whether walking into relDir is pointless because
// an exclude rule covers the whole directory tree.
func dirExcludedByPattern(pattern, relDir string) bool {
	pattern = strings.ReplaceAll(pattern, "\\", "/")
	relDir = strings.Trim(relDir, "/")

	trimmed := strings.TrimSuffix(pattern, "/")
	if trimmed == relDir {
		return true
	}

	if strings.HasSuffix(pattern, "/**") {
		prefix := strings.TrimSuffix(pattern, "/**")
		prefix = strings.TrimSuffix(prefix, "/")
		if prefix == relDir {
			return true
		}
		if strings.HasPrefix(prefix, "**/") {
			suffix := strings.TrimPrefix(prefix, "**/")
			if relDir == suffix || strings.HasSuffix(relDir, "/"+suffix) {
				return true
			}
		}
	}
	return false
}
