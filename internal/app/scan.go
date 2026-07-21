package app

import (
	"context"
	"fmt"
)

// ScanResult is a lightweight scan summary for the scan command.
type ScanResult struct {
	Root      string   `json:"root"`
	FileCount int      `json:"file_count"`
	Markers   []string `json:"markers,omitempty"`
	Jobs      int      `json:"jobs"`
}

// Scan walks the project and returns an index summary.
func (a *App) Scan(ctx context.Context, root string) (ScanResult, error) {
	if a.scanner == nil {
		return ScanResult{}, fmt.Errorf("%w: scanner", ErrDependencyMissing)
	}
	if root == "" {
		root = a.cfg.Scan.Root
	}
	opts := scanOptionsFromConfig(a.cfg)
	pctx, err := a.scanner.Scan(ctx, root, opts)
	if err != nil {
		return ScanResult{}, err
	}
	markers := []string(nil)
	if pctx.Files != nil {
		markers = pctx.Files.Markers()
	}
	return ScanResult{
		Root:      pctx.Root,
		FileCount: pctx.FileCount(),
		Markers:   markers,
		Jobs:      a.cfg.EffectiveJobs(),
	}, nil
}
