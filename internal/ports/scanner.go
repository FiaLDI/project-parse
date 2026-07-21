package ports

import (
	"context"
	"io/fs"

	"github.com/FiaLDI/project-parse/internal/domain"
)

// ScanOptions configures filesystem traversal.
type ScanOptions struct {
	Jobs            int
	MaxFileBytes    int64
	FollowSymlinks  bool
	Include         []string
	Exclude         []string
}

// Scanner walks a project tree and builds a ProjectContext index.
type Scanner interface {
	Scan(ctx context.Context, root string, opts ScanOptions) (domain.ProjectContext, error)
}

// FileCache provides read-once access to file contents and metadata.
type FileCache interface {
	Read(path string) ([]byte, error)
	Stat(path string) (fs.FileInfo, error)
	Exists(path string) bool
}

// IgnoreMatcher reports whether a relative path should be skipped.
type IgnoreMatcher interface {
	Match(relPath string) bool
}
