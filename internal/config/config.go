package config

import (
	"fmt"
	"os"
	"runtime"

	"gopkg.in/yaml.v3"
)

// Config is the top-level application configuration.
type Config struct {
	Scan    ScanConfig    `yaml:"scan"`
	Plugins PluginsConfig `yaml:"plugins"`
	Report  ReportConfig  `yaml:"report"`
	Graph   GraphConfig   `yaml:"graph"`
	Log     LogConfig     `yaml:"log"`
}

// ScanConfig controls filesystem traversal and parallelism.
type ScanConfig struct {
	Root           string   `yaml:"root"`
	Jobs           int      `yaml:"jobs"`
	MaxFileBytes   int64    `yaml:"max_file_bytes"`
	FollowSymlinks bool     `yaml:"follow_symlinks"`
	Include        []string `yaml:"include"`
	Exclude        []string `yaml:"exclude"`
}

// PluginsConfig selects which plugins run.
type PluginsConfig struct {
	Enabled  []string `yaml:"enabled"`
	Disabled []string `yaml:"disabled"`
}

// ReportConfig controls report generation.
type ReportConfig struct {
	Formats         []string `yaml:"formats"`
	OutputDir       string   `yaml:"output_dir"`
	IncludeEvidence bool     `yaml:"include_evidence"`
}

// GraphConfig controls architecture graph output.
type GraphConfig struct {
	Enabled bool   `yaml:"enabled"`
	Format  string `yaml:"format"`
}

// LogConfig controls slog setup.
type LogConfig struct {
	Level  string `yaml:"level"`
	Format string `yaml:"format"`
}

// Default returns a sensible baseline configuration.
func Default() Config {
	return Config{
		Scan: ScanConfig{
			Root:           ".",
			Jobs:           0,
			MaxFileBytes:   2 * 1024 * 1024,
			FollowSymlinks: false,
			Exclude: []string{
				"**/node_modules/**",
				"**/.git/**",
				"**/vendor/**",
				"**/dist/**",
				"**/build/**",
				"**/.next/**",
				"**/target/**",
				"**/__pycache__/**",
			},
		},
		Plugins: PluginsConfig{
			Enabled: []string{
				"node",
				"python",
				"golang",
				"rust",
				"java",
				"docker",
				"git",
				"githubactions",
				"kubernetes",
				"database",
				"architecture",
				"documentation",
			},
		},
		Report: ReportConfig{
			Formats:         []string{"json", "markdown"},
			OutputDir:       "./.project-parser",
			IncludeEvidence: true,
		},
		Graph: GraphConfig{
			Enabled: true,
			Format:  "svg",
		},
		Log: LogConfig{
			Level:  "info",
			Format: "text",
		},
	}
}

// Load reads YAML from path and merges it onto defaults.
func Load(path string) (Config, error) {
	cfg := Default()
	if path == "" {
		return cfg, nil
	}
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return cfg, nil
		}
		return Config{}, fmt.Errorf("read config: %w", err)
	}
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return Config{}, fmt.Errorf("parse config: %w", err)
	}
	cfg.ApplyDefaults()
	if err := cfg.Validate(); err != nil {
		return Config{}, err
	}
	return cfg, nil
}

// ApplyDefaults fills zero-value fields that must not stay empty.
func (c *Config) ApplyDefaults() {
	def := Default()
	if c.Scan.Root == "" {
		c.Scan.Root = def.Scan.Root
	}
	if c.Scan.MaxFileBytes == 0 {
		c.Scan.MaxFileBytes = def.Scan.MaxFileBytes
	}
	if c.Report.OutputDir == "" {
		c.Report.OutputDir = def.Report.OutputDir
	}
	if len(c.Report.Formats) == 0 {
		c.Report.Formats = def.Report.Formats
	}
	if c.Graph.Format == "" {
		c.Graph.Format = def.Graph.Format
	}
	if c.Log.Level == "" {
		c.Log.Level = def.Log.Level
	}
	if c.Log.Format == "" {
		c.Log.Format = def.Log.Format
	}
}

// Validate checks configuration invariants.
func (c Config) Validate() error {
	if c.Scan.Jobs < 0 {
		return fmt.Errorf("scan.jobs must be >= 0")
	}
	if c.Scan.MaxFileBytes < 0 {
		return fmt.Errorf("scan.max_file_bytes must be >= 0")
	}
	switch c.Log.Level {
	case "debug", "info", "warn", "error":
	default:
		return fmt.Errorf("log.level must be one of: debug, info, warn, error")
	}
	switch c.Log.Format {
	case "text", "json":
	default:
		return fmt.Errorf("log.format must be one of: text, json")
	}
	return nil
}

// EffectiveJobs returns the worker count, defaulting to NumCPU when Jobs is 0.
func (c Config) EffectiveJobs() int {
	if c.Scan.Jobs > 0 {
		return c.Scan.Jobs
	}
	n := runtime.NumCPU()
	if n < 1 {
		return 1
	}
	return n
}
