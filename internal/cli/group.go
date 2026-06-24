package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/fatih/color"
	"github.com/google/uuid"
	"github.com/spf13/cobra"

	"github.com/ejoy/skm/internal/agent"
	skmsync "github.com/ejoy/skm/internal/sync"
)

func newGroupCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "group",
		Short: "Manage skill groups",
	}

	cmd.AddCommand(
		newGroupCreateCmd(),
		newGroupListCmd(),
		newGroupShowCmd(),
		newGroupAddCmd(),
		newGroupRemoveCmd(),
		newGroupInstallCmd(),
		newGroupDeleteCmd(),
	)

	return cmd
}

func newGroupCreateCmd() *cobra.Command {
	var desc string

	cmd := &cobra.Command{
		Use:   "create <name>",
		Short: "Create a skill group",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := NewConfig()
			if err != nil {
				return err
			}
			defer cfg.Close()

			id := uuid.New().String()
			if err := cfg.Store.InsertGroup(id, args[0], desc); err != nil {
				return fmt.Errorf("create group: %w", err)
			}
			cfg.WriteMetadata()
			color.Green("✓ Created group %q", args[0])
			return nil
		},
	}

	cmd.Flags().StringVarP(&desc, "description", "d", "", "Group description")
	return cmd
}

func newGroupListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List all groups",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := NewConfig()
			if err != nil {
				return err
			}
			defer cfg.Close()

			groups, err := cfg.Store.ListGroups()
			if err != nil {
				return err
			}
			if len(groups) == 0 {
				fmt.Println("No groups. Use 'skm group create <name>' to create one.")
				return nil
			}
			for _, g := range groups {
				skills, _ := cfg.Store.ListGroupSkills(g.ID)
				fmt.Printf("  %s  (%d skills)", color.CyanString(g.Name), len(skills))
				if g.Description != "" {
					fmt.Printf("  — %s", g.Description)
				}
				fmt.Println()
			}
			return nil
		},
	}
}

func newGroupShowCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "show <name>",
		Short: "Show group with skills",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := NewConfig()
			if err != nil {
				return err
			}
			defer cfg.Close()

			g, err := cfg.Store.GetGroup(args[0])
			if err != nil {
				return err
			}
			skills, _ := cfg.Store.ListGroupSkills(g.ID)

			color.Cyan("Group: %s\n", g.Name)
			if g.Description != "" {
				fmt.Printf("  %s\n", g.Description)
			}
			fmt.Printf("  Skills (%d):\n", len(skills))
			for _, sk := range skills {
				fmt.Printf("    %s — %s\n", color.CyanString(sk.Name), sk.Description)
			}
			return nil
		},
	}
}

func newGroupAddCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "add <group> <skill> [skill...]",
		Short: "Add skills to a group",
		Args:  cobra.MinimumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := NewConfig()
			if err != nil {
				return err
			}
			defer cfg.Close()

			g, err := cfg.Store.GetGroup(args[0])
			if err != nil {
				return err
			}

			for i, name := range args[1:] {
				sk, err := cfg.Store.GetSkill(name)
				if err != nil {
					color.Red("  ✗ Skill %q not found", name)
					continue
				}
				if err := cfg.Store.AddSkillToGroup(g.ID, sk.ID, i); err != nil {
					color.Red("  ✗ Failed to add %s: %v", name, err)
					continue
				}
				color.Green("  ✓ Added %s to %s", name, g.Name)
			}
			cfg.WriteMetadata()
			return nil
		},
	}
}

func newGroupRemoveCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "remove <group> <skill>",
		Short: "Remove a skill from a group",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := NewConfig()
			if err != nil {
				return err
			}
			defer cfg.Close()

			g, err := cfg.Store.GetGroup(args[0])
			if err != nil {
				return err
			}
			sk, err := cfg.Store.GetSkill(args[1])
			if err != nil {
				return err
			}
			if err := cfg.Store.RemoveSkillFromGroup(g.ID, sk.ID); err != nil {
				return err
			}
			cfg.WriteMetadata()
			color.Green("✓ Removed %s from %s", sk.Name, g.Name)
			return nil
		},
	}
}

func newGroupInstallCmd() *cobra.Command {
	var (
		agentNames []string
		global     bool
	)

	cmd := &cobra.Command{
		Use:   "install <group>",
		Short: "Install all skills in a group to agents",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := NewConfig()
			if err != nil {
				return err
			}
			defer cfg.Close()

			g, err := cfg.Store.GetGroup(args[0])
			if err != nil {
				return err
			}

			skills, err := cfg.Store.ListGroupSkills(g.ID)
			if err != nil {
				return err
			}
			if len(skills) == 0 {
				return fmt.Errorf("group %q has no skills", g.Name)
			}

			targetAgents := agent.Resolve(agentNames)
			if len(targetAgents) == 0 {
				return fmt.Errorf("no agents specified or detected")
			}

			cwd, _ := os.Getwd()

			for _, sk := range skills {
				for _, ag := range targetAgents {
					targetDir, err := agent.InstallPath(ag, global, cwd)
					if err != nil {
						color.Red("  ✗ %s → %s: %v", sk.Name, ag.DisplayName, err)
						continue
					}
					targetPath := filepath.Join(targetDir, sk.Name)

					if skmsync.IsCurrent(sk.CentralPath, targetPath, "symlink") {
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
				}
			}
			return nil
		},
	}

	cmd.Flags().StringSliceVarP(&agentNames, "agent", "a", nil, "Target agent(s)")
	cmd.Flags().BoolVarP(&global, "global", "g", false, "Install to global directories")
	return cmd
}

func newGroupDeleteCmd() *cobra.Command {
	var yes bool

	cmd := &cobra.Command{
		Use:   "delete <name>",
		Short: "Delete a group",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := NewConfig()
			if err != nil {
				return err
			}
			defer cfg.Close()

			g, err := cfg.Store.GetGroup(args[0])
			if err != nil {
				return err
			}

			if !yes {
				fmt.Printf("Delete group %q? Skills are not removed. [y/N] ", g.Name)
				var confirm string
				fmt.Scanln(&confirm)
				if confirm != "y" && confirm != "Y" {
					fmt.Println("Cancelled.")
					return nil
				}
			}

			if err := cfg.Store.DeleteGroup(g.ID); err != nil {
				return err
			}
			cfg.WriteMetadata()
			color.Green("✓ Deleted group %q", g.Name)
			return nil
		},
	}

	cmd.Flags().BoolVarP(&yes, "yes", "y", false, "Skip confirmation")
	return cmd
}
