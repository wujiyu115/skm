package cli

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	skmsync "github.com/ejoy/skm/internal/sync"
)

func newRemoveCmd() *cobra.Command {
	var yes bool

	cmd := &cobra.Command{
		Use:   "remove <skill>",
		Short: "Remove an installed skill",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := NewConfig()
			if err != nil {
				return err
			}
			defer cfg.Close()

			sk, err := cfg.Store.GetSkill(args[0])
			if err != nil {
				return fmt.Errorf("skill %q not found", args[0])
			}

			if !yes {
				fmt.Printf("Remove skill %q and unsync from all agents? [y/N] ", sk.Name)
				var confirm string
				fmt.Scanln(&confirm)
				if confirm != "y" && confirm != "Y" {
					fmt.Println("Cancelled.")
					return nil
				}
			}

			targets, _ := cfg.Store.ListTargets(sk.ID)
			for _, t := range targets {
				if err := skmsync.Unsync(t.TargetPath); err != nil {
					color.Yellow("  Warning: unsync %s failed: %v", t.TargetPath, err)
				}
			}

			os.RemoveAll(sk.CentralPath)

			if err := cfg.Store.DeleteSkill(sk.ID); err != nil {
				return err
			}

			color.Green("✓ Removed %s", sk.Name)
			return nil
		},
	}

	cmd.Flags().BoolVarP(&yes, "yes", "y", false, "Skip confirmation")
	return cmd
}
