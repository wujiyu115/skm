package cli

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

func newAuditCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "audit",
		Short: "View and manage audit log",
	}

	cmd.AddCommand(
		newAuditListCmd(),
		newAuditPruneCmd(),
	)

	return cmd
}

func newAuditListCmd() *cobra.Command {
	var limit int

	cmd := &cobra.Command{
		Use:   "list",
		Short: "Show recent audit log entries",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := NewConfig()
			if err != nil {
				return err
			}
			defer cfg.Close()

			entries, err := cfg.Store.ListAuditLog(limit)
			if err != nil {
				return err
			}
			if len(entries) == 0 {
				fmt.Println("No audit entries.")
				return nil
			}

			fmt.Printf("  %-20s %-15s %-20s %s\n",
				color.HiBlackString("TIME"),
				color.HiBlackString("ACTION"),
				color.HiBlackString("TARGET"),
				color.HiBlackString("DETAIL"))
			for _, e := range entries {
				fmt.Printf("  %-20s %-15s %-20s %s\n",
					e.CreatedAt, color.CyanString(e.Action), e.Target, e.Detail)
			}
			return nil
		},
	}

	cmd.Flags().IntVarP(&limit, "limit", "n", 50, "Max entries to show")
	return cmd
}

func newAuditPruneCmd() *cobra.Command {
	var keep int

	cmd := &cobra.Command{
		Use:   "prune",
		Short: "Prune old audit log entries",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := NewConfig()
			if err != nil {
				return err
			}
			defer cfg.Close()

			if err := cfg.Store.PruneAuditLog(keep); err != nil {
				return err
			}
			color.Green("Pruned audit log (keeping latest %d entries)", keep)
			return nil
		},
	}

	cmd.Flags().IntVar(&keep, "keep", 200, "Number of recent entries to keep")
	return cmd
}
