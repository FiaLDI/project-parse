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
			art, err := a.Graph(cmd.Context(), app.GraphOptions{
				Root:   root,
				Format: format,
				OutDir: outDir,
			})
			if err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "nodes: %d\n", len(art.Graph.Nodes))
			fmt.Fprintf(cmd.OutOrStdout(), "edges: %d\n", len(art.Graph.Edges))
			fmt.Fprintf(cmd.OutOrStdout(), "rendered %s (%d bytes) → %s\n", art.Format, len(art.Bytes), art.Path)
			return nil
		},
	}
	cmd.Flags().StringVar(&format, "format", "", "graph format: svg, graph-json (alias: json)")
	cmd.Flags().StringVar(&outDir, "out", "", "output directory")
	return cmd
}
