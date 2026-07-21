package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newDoctorCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "doctor",
		Short: "Diagnose local environment and configuration",
		RunE: func(cmd *cobra.Command, args []string) error {
			a, err := requireApp()
			if err != nil {
				return err
			}
			res, err := a.Doctor(cmd.Context())
			if err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "project-parser %s\n", res.Version)
			fmt.Fprintf(cmd.OutOrStdout(), "go: %s\n\n", res.Go)
			for _, c := range res.Checks {
				mark := "ok"
				if !c.OK {
					if c.Critical {
						mark = "FAIL"
					} else {
						mark = "warn"
					}
				}
				fmt.Fprintf(cmd.OutOrStdout(), "[%s] %-12s %s\n", mark, c.Name, c.Message)
			}
			if res.HasCriticalFailures() {
				return fmt.Errorf("critical doctor check(s) failed")
			}
			return nil
		},
	}
}
