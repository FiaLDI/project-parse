package cli

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/FiaLDI/project-parse/internal/app"
)

func newReportCmd() *cobra.Command {
	var (
		formats string
		outDir  string
	)
	cmd := &cobra.Command{
		Use:   "report [path]",
		Short: "Analyze a project and generate a structured report",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			a, err := requireApp()
			if err != nil {
				return err
			}
			root := "."
			if len(args) == 1 {
				root = args[0]
			}
			var formatList []string
			if formats != "" {
				formatList = splitCSV(formats)
			}
			artifacts, err := a.Report(cmd.Context(), app.ReportOptions{
				Root:    root,
				Formats: formatList,
				OutDir:  outDir,
			})
			if err != nil {
				return err
			}
			for _, art := range artifacts {
				fmt.Fprintf(cmd.OutOrStdout(), "rendered %s (%d bytes) → %s\n", art.Format, len(art.Bytes), art.Path)
			}
			return nil
		},
	}
	cmd.Flags().StringVar(&formats, "format", "", "comma-separated formats: json,markdown,html,svg")
	cmd.Flags().StringVar(&outDir, "out", "", "output directory")
	return cmd
}

func splitCSV(s string) []string {
	parts := strings.Split(s, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}
