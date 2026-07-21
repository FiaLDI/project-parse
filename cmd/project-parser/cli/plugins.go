package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newPluginsCmd() *cobra.Command {
	var showAll bool
	cmd := &cobra.Command{
		Use:   "plugins",
		Short: "List registered analysis plugins",
		RunE: func(cmd *cobra.Command, args []string) error {
			a, err := requireApp()
			if err != nil {
				return err
			}
			list, err := a.ListPlugins()
			if err != nil {
				return err
			}
			if len(list) == 0 {
				fmt.Fprintln(cmd.OutOrStdout(), "no plugins registered yet")
				return nil
			}
			for _, p := range list {
				if !showAll && !p.Enabled {
					continue
				}
				state := "disabled"
				if p.Enabled {
					state = "enabled"
				}
				fmt.Fprintf(cmd.OutOrStdout(), "%-16s priority=%-3d %s\n", p.Name, p.Priority, state)
			}
			return nil
		},
	}
	cmd.Flags().BoolVar(&showAll, "all", false, "include disabled plugins")
	return cmd
}
