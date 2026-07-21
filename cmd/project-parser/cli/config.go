package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/FiaLDI/project-parse/internal/app"
)

func newConfigCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Manage project-parser configuration",
	}
	cmd.AddCommand(newConfigInitCmd())
	return cmd
}

func newConfigInitCmd() *cobra.Command {
	var (
		force bool
		out   string
	)
	cmd := &cobra.Command{
		Use:   "init [path]",
		Short: "Create a default parser.yaml configuration file",
		Long:  "Writes a parser.yaml with default scan, plugin, report, graph, and log settings.",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			path := out
			if path == "" && len(args) > 0 {
				path = args[0]
			}
			written, err := app.InitConfig(app.ConfigInitOptions{
				Path:  path,
				Force: force,
			})
			if err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "created %s\n", written)
			return nil
		},
	}
	cmd.Flags().BoolVar(&force, "force", false, "overwrite existing configuration file")
	cmd.Flags().StringVar(&out, "out", "", "output path (default: parser.yaml)")
	return cmd
}
