package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/ejoy/skm/internal/agent"
	skmsync "github.com/ejoy/skm/internal/sync"
)

func newBatchCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "batch",
		Short: "Batch operations on multiple skills",
	}

	cmd.AddCommand(
		newBatchDeleteCmd(),
		newBatchEnableCmd(),
		newBatchDisableCmd(),
		newBatchTagCmd(),
		newBatchSyncCmd(),
	)

	return cmd
}

func newBatchDeleteCmd() *cobra.Command {
	var yes bool

	cmd := &cobra.Command{
		Use:   "delete <skill> [skill...]",
		Short: "Delete multiple skills",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if !yes {
				fmt.Printf("Delete %d skill(s)? [y/N] ", len(args))
				var confirm string
				fmt.Scanln(&confirm)
				if confirm != "y" && confirm != "Y" {
					fmt.Println("Cancelled.")
					return nil
				}
			}

			cfg, err := NewConfig()
			if err != nil {
				return err
			}
			defer cfg.Close()

			for _, name := range args {
				if err := cfg.Store.DeleteSkillByName(name); err != nil {
					color.Red("  ✗ %s: %v", name, err)
					continue
				}
				color.Green("  ✓ Deleted %s", name)
			}
			return nil
		},
	}

	cmd.Flags().BoolVarP(&yes, "yes", "y", false, "Skip confirmation")
	return cmd
}

func newBatchEnableCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "enable <skill> [skill...]",
		Short: "Enable multiple skills",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := NewConfig()
			if err != nil {
				return err
			}
			defer cfg.Close()

			for _, name := range args {
				sk, err := cfg.Store.GetSkill(name)
				if err != nil {
					color.Red("  ✗ %s: not found", name)
					continue
				}
				if err := cfg.Store.SetSkillEnabled(sk.ID, true); err != nil {
					color.Red("  ✗ %s: %v", name, err)
					continue
				}
				color.Green("  ✓ Enabled %s", name)
			}
			return nil
		},
	}
}

func newBatchDisableCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "disable <skill> [skill...]",
		Short: "Disable multiple skills",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := NewConfig()
			if err != nil {
				return err
			}
			defer cfg.Close()

			for _, name := range args {
				sk, err := cfg.Store.GetSkill(name)
				if err != nil {
					color.Red("  ✗ %s: not found", name)
					continue
				}
				if err := cfg.Store.SetSkillEnabled(sk.ID, false); err != nil {
					color.Red("  ✗ %s: %v", name, err)
					continue
				}
				color.Green("  ✓ Disabled %s", name)
			}
			return nil
		},
	}
}

func newBatchTagCmd() *cobra.Command {
	var (
		addTag    string
		removeTag string
	)

	cmd := &cobra.Command{
		Use:   "tag <skill> [skill...]",
		Short: "Add or remove a tag on multiple skills",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if addTag == "" && removeTag == "" {
				return fmt.Errorf("specify --add or --remove")
			}

			cfg, err := NewConfig()
			if err != nil {
				return err
			}
			defer cfg.Close()

			for _, name := range args {
				sk, err := cfg.Store.GetSkill(name)
				if err != nil {
					color.Red("  ✗ %s: not found", name)
					continue
				}
				if addTag != "" {
					if err := cfg.Store.AddTag(sk.ID, addTag); err != nil {
						color.Red("  ✗ %s: %v", name, err)
						continue
					}
					color.Green("  ✓ %s: +%s", name, addTag)
				}
				if removeTag != "" {
					if err := cfg.Store.RemoveTag(sk.ID, removeTag); err != nil {
						color.Red("  ✗ %s: %v", name, err)
						continue
					}
					color.Green("  ✓ %s: -%s", name, removeTag)
				}
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&addTag, "add", "", "Tag to add")
	cmd.Flags().StringVar(&removeTag, "remove", "", "Tag to remove")
	return cmd
}

func newBatchSyncCmd() *cobra.Command {
	var (
		agentNames []string
		global     bool
	)

	cmd := &cobra.Command{
		Use:   "sync <skill> [skill...]",
		Short: "Sync specific skills to agents",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := NewConfig()
			if err != nil {
				return err
			}
			defer cfg.Close()

			targetAgents := agent.Resolve(agentNames)
			if len(targetAgents) == 0 {
				return fmt.Errorf("no agents specified or detected")
			}

			cwd, _ := os.Getwd()

			for _, name := range args {
				sk, err := cfg.Store.GetSkill(name)
				if err != nil {
					color.Red("  ✗ %s: not found", name)
					continue
				}
				for _, ag := range targetAgents {
					targetDir, err := agent.InstallPath(ag, global, cwd)
					if err != nil {
						color.Red("  ✗ %s → %s: %v", name, ag.DisplayName, err)
						continue
					}
					targetPath := filepath.Join(targetDir, sk.Name)

					if err := skmsync.SyncSkill(sk.CentralPath, targetPath, "symlink"); err != nil {
						color.Red("  ✗ %s → %s: %v", name, ag.DisplayName, err)
						continue
					}
					_ = cfg.Store.UpsertTarget(sk.ID, ag.Name, targetPath, "symlink", sk.ContentHash)
					color.Green("  ✓ %s → %s", name, ag.DisplayName)
				}
			}
			return nil
		},
	}

	cmd.Flags().StringSliceVarP(&agentNames, "agent", "a", nil, "Target agent(s)")
	cmd.Flags().BoolVarP(&global, "global", "g", false, "Sync to global directories")
	return cmd
}
