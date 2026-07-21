package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/FiaLDI/project-parse/internal/app"
)

func newGraphCmd() *cobra.Command {
	var (
		format string
		outDir string
	)
	cmd := &cobra.Command{
		Use:   "graph [path]",
		Short: "Build an architecture graph for a project",
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
			g, data, err := a.Graph(cmd.Context(), app.GraphOptions{
				Root:   root,
				Format: format,
				OutDir: outDir,
			})
			if err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "nodes: %d\nedges: %d\nrendered: %d bytes\n", len(g.Nodes), len(g.Edges), len(data))
			return nil
		},
	}
	cmd.Flags().StringVar(&format, "format", "", "graph format (default: svg)")
	cmd.Flags().StringVar(&outDir, "out", "", "output directory")
	return cmd
}
