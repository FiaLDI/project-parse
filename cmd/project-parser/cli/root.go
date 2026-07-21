package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/FiaLDI/project-parse/internal/app"
	"github.com/FiaLDI/project-parse/internal/analyzer"
	"github.com/FiaLDI/project-parse/internal/cache"
	"github.com/FiaLDI/project-parse/internal/config"
	"github.com/FiaLDI/project-parse/internal/graph"
	"github.com/FiaLDI/project-parse/internal/logger"
	"github.com/FiaLDI/project-parse/internal/output"
	"github.com/FiaLDI/project-parse/internal/plugins"
	"github.com/FiaLDI/project-parse/internal/ports"
	"github.com/FiaLDI/project-parse/internal/registry"
	"github.com/FiaLDI/project-parse/internal/report"
	"github.com/FiaLDI/project-parse/internal/scanner"
	"github.com/FiaLDI/project-parse/internal/version"
)

type rootOptions struct {
	configPath string
	logLevel   string
	logFormat  string
	jobs       int
}

var (
	opts rootOptions
	application *app.App
)

// Execute runs the root Cobra command.
func Execute() error {
	root := newRootCmd()
	if err := root.Execute(); err != nil {
		printErr(err)
		return err
	}
	return nil
}

func newRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "project-parser",
		Short:         "Analyze project stack, architecture, and infrastructure",
		SilenceUsage:  true,
		SilenceErrors: true,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return bootstrap()
		},
	}

	cmd.PersistentFlags().StringVar(&opts.configPath, "config", "parser.yaml", "path to YAML config")
	cmd.PersistentFlags().StringVar(&opts.logLevel, "log-level", "", "override log level (debug|info|warn|error)")
	cmd.PersistentFlags().StringVar(&opts.logFormat, "log-format", "", "override log format (text|json)")
	cmd.PersistentFlags().IntVar(&opts.jobs, "jobs", -1, "worker pool size (0=NumCPU, -1=use config)")

	cmd.AddCommand(newConfigCmd())
	cmd.AddCommand(newScanCmd())
	cmd.AddCommand(newReportCmd())
	cmd.AddCommand(newGraphCmd())
	cmd.AddCommand(newPluginsCmd())
	cmd.AddCommand(newVersionCmd())
	cmd.AddCommand(newDoctorCmd())

	return cmd
}

func bootstrap() error {
	cfg, err := config.Load(opts.configPath)
	if err != nil {
		return err
	}
	if opts.logLevel != "" {
		cfg.Log.Level = opts.logLevel
	}
	if opts.logFormat != "" {
		cfg.Log.Format = opts.logFormat
	}
	if opts.jobs >= 0 {
		cfg.Scan.Jobs = opts.jobs
	}
	if err := cfg.Validate(); err != nil {
		return err
	}

	log, err := logger.Setup(cfg.Log.Level, cfg.Log.Format)
	if err != nil {
		return err
	}

	reg := registry.New()
	fileCache := cache.New(cache.Options{MaxFileBytes: cfg.Scan.MaxFileBytes})
	plugins.RegisterAll(reg, fileCache)

	fsScanner := scanner.New()
	pluginAnalyzer := analyzer.New(analyzer.Options{Jobs: cfg.EffectiveJobs(), Log: log})

	application = app.New(cfg, log, app.Deps{
		Scanner:  fsScanner,
		Cache:    fileCache,
		Registry: reg,
		Analyzer: pluginAnalyzer,
		Agg:      report.New(),
		Graph:    graph.New(),
		Renderers: []ports.Renderer{
			output.NewJSON(),
			output.NewMarkdown(),
			output.NewHTML(),
			output.NewSVG(),
			output.NewGraphJSON(),
			output.NewPDF(),
		},
	})
	return nil
}

func requireApp() (*app.App, error) {
	if application == nil {
		return nil, fmt.Errorf("application is not initialized")
	}
	return application, nil
}

func printErr(err error) {
	fmt.Fprintln(os.Stderr, "error:", err)
}

func newVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print version information",
		RunE: func(cmd *cobra.Command, args []string) error {
			a, err := requireApp()
			if err != nil {
				// version must work even if bootstrap partially fails
				fmt.Println(version.String())
				return nil
			}
			fmt.Println(a.Version())
			return nil
		},
	}
}
