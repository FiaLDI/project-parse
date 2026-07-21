package cli

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
)

func newScanCmd() *cobra.Command {
	var format string
	cmd := &cobra.Command{
		Use:   "scan [path]",
		Short: "Scan a project tree and print an index summary",
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
			res, err := a.Scan(cmd.Context(), root)
			if err != nil {
				return err
			}
			switch format {
			case "json":
				enc := json.NewEncoder(cmd.OutOrStdout())
				enc.SetIndent("", "  ")
				return enc.Encode(res)
			default:
				fmt.Fprintf(cmd.OutOrStdout(), "root: %s\n", res.Root)
				fmt.Fprintf(cmd.OutOrStdout(), "files: %d\n", res.FileCount)
				fmt.Fprintf(cmd.OutOrStdout(), "jobs: %d\n", res.Jobs)
				if len(res.Markers) > 0 {
					fmt.Fprintf(cmd.OutOrStdout(), "markers: %v\n", res.Markers)
				}
				return nil
			}
		},
	}
	cmd.Flags().StringVar(&format, "format", "summary", "output format: summary|json")
	return cmd
}
