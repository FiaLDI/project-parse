package app

import (
	"github.com/FiaLDI/project-parse/internal/config"
	"github.com/FiaLDI/project-parse/internal/ports"
)

func scanOptionsFromConfig(cfg config.Config) ports.ScanOptions {
	return ports.ScanOptions{
		Jobs:           cfg.EffectiveJobs(),
		MaxFileBytes:   cfg.Scan.MaxFileBytes,
		FollowSymlinks: cfg.Scan.FollowSymlinks,
		Include:        append([]string(nil), cfg.Scan.Include...),
		Exclude:        append([]string(nil), cfg.Scan.Exclude...),
	}
}
