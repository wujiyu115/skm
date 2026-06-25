package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/ejoy/skm/internal/agent"
)

func newUnsyncCmd() *cobra.Command {
	var (
		agentNames []string
		global     bool
	)

	cmd := &cobra.Command{
		Use:   "unsync <skill>",
		Short: "Unsync a skill from specific agents",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := NewConfig()
			if err != nil {
				return err
			}
			defer cfg.Close()

			sk, err := cfg.Store.GetSkill(args[0])
			if err != nil {
				return fmt.Errorf("skill %q: %w", args[0], err)
			}

			targetAgents := agent.Resolve(agentNames)
			if len(targetAgents) == 0 {
				return fmt.Errorf("no agents specified or detected")
			}

			cwd, _ := os.Getwd()

			for _, ag := range targetAgents {
				targetDir, err := agent.InstallPath(ag, global, cwd)
				if err != nil {
					color.Red("  ✗ %s: %v", ag.DisplayName, err)
					continue
				}
				targetPath := filepath.Join(targetDir, sk.Name)

				if err := os.Remove(targetPath); err != nil && !os.IsNotExist(err) {
					color.Red("  ✗ %s: %v", ag.DisplayName, err)
					continue
				}

				_ = cfg.Store.DeleteTarget(sk.ID, ag.Name)
				color.Green("  ✓ Unsynced %s from %s", sk.Name, ag.DisplayName)
			}
			return nil
		},
	}

	cmd.Flags().StringSliceVarP(&agentNames, "agent", "a", nil, "Target agent(s) to unsync from")
	cmd.Flags().BoolVarP(&global, "global", "g", false, "Unsync from global directories")
	cmd.MarkFlagRequired("agent")
	return cmd
}
