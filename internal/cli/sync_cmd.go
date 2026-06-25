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

func newSyncCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sync",
		Short: "Sync skills to agent directories",
	}

	cmd.AddCommand(newSyncStatusCmd())

	// Default action: run sync (keep backward-compatible)
	syncRun := newSyncRunCmd()
	cmd.RunE = syncRun.RunE
	cmd.Flags().AddFlagSet(syncRun.Flags())

	return cmd
}

func newSyncStatusCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "Show sync status summary",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := NewConfig()
			if err != nil {
				return err
			}
			defer cfg.Close()

			skills, err := cfg.Store.ListSkills()
			if err != nil {
				return err
			}

			total, synced, stale := 0, 0, 0
			for _, sk := range skills {
				if !sk.Enabled {
					continue
				}
				total++
				targets, _ := cfg.Store.ListTargets(sk.ID)
				if len(targets) > 0 {
					synced++
					for _, t := range targets {
						if t.SourceHash != sk.ContentHash {
							stale++
							break
						}
					}
				}
			}

			fmt.Printf("  Total enabled: %d\n", total)
			fmt.Printf("  Synced:        %s\n", color.GreenString("%d", synced))
			fmt.Printf("  Stale:         %s\n", color.YellowString("%d", stale))
			fmt.Printf("  Not synced:    %s\n", color.RedString("%d", total-synced))
			return nil
		},
	}
}

func newSyncRunCmd() *cobra.Command {
	var (
		agentNames []string
		dryRun     bool
		global     bool
	)

	cmd := &cobra.Command{
		Use:   "sync",
		Short: "Sync skills to agent directories",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := NewConfig()
			if err != nil {
				return err
			}
			defer cfg.Close()

			skills, err := cfg.Store.ListSkills()
			if err != nil {
				return err
			}

			targetAgents := agent.Resolve(agentNames)
			if len(targetAgents) == 0 {
				return fmt.Errorf("no agents specified or detected")
			}

			cwd, _ := os.Getwd()
			synced, skipped := 0, 0

			for _, sk := range skills {
				if !sk.Enabled {
					continue
				}
				for _, ag := range targetAgents {
					targetDir, err := agent.InstallPath(ag, global, cwd)
					if err != nil {
						color.Red("  ✗ %s → %s: %v", sk.Name, ag.DisplayName, err)
						continue
					}
					targetPath := filepath.Join(targetDir, sk.Name)

					if skmsync.IsCurrent(sk.CentralPath, targetPath, "symlink") {
						skipped++
						continue
					}

					if dryRun {
						fmt.Printf("  [dry-run] %s → %s (%s)\n", sk.Name, ag.DisplayName, targetPath)
						synced++
						continue
					}

					if err := skmsync.SyncSkill(sk.CentralPath, targetPath, "symlink"); err != nil {
						color.Red("  ✗ %s → %s: %v", sk.Name, ag.DisplayName, err)
						continue
					}

					if err := cfg.Store.UpsertTarget(sk.ID, ag.Name, targetPath, "symlink", sk.ContentHash); err != nil {
						color.Red("  ✗ %s → %s: failed to save target: %v", sk.Name, ag.DisplayName, err)
						continue
					}
					color.Green("  ✓ %s → %s", sk.Name, ag.DisplayName)
					synced++
				}
			}

			fmt.Printf("\nSynced: %d, Up-to-date: %d\n", synced, skipped)
			return nil
		},
	}

	cmd.Flags().StringSliceVarP(&agentNames, "agent", "a", nil, "Target agent(s)")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "Preview sync operations")
	cmd.Flags().BoolVarP(&global, "global", "g", false, "Sync to global directories")

	return cmd
}
