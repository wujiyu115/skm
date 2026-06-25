package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/ejoy/skm/internal/agent"
	"github.com/ejoy/skm/internal/project"
	skmsync "github.com/ejoy/skm/internal/sync"
)

func newAgentCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "agent",
		Short: "Manage agent adapters",
	}

	cmd.AddCommand(
		newAgentListCmd(),
		newAgentAddCmd(),
		newAgentRemoveCmd(),
		newAgentSkillsCmd(),
		newAgentAddSkillCmd(),
		newAgentRemoveSkillCmd(),
	)
	return cmd
}

func newAgentListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List supported agents",
		RunE: func(cmd *cobra.Command, args []string) error {
			detected := agent.Detect()
			detectedMap := map[string]bool{}
			for _, a := range detected {
				detectedMap[a.Name] = true
			}

			for _, a := range agent.Builtin() {
				status := color.RedString("not detected")
				if detectedMap[a.Name] {
					status = color.GreenString("detected")
				}
				fmt.Printf("  %s (%s)  [%s]\n", color.CyanString(a.Name), a.DisplayName, status)
				fmt.Printf("    project: %s\n    global:  ~/%s\n", a.ProjectDir, a.GlobalDir)
			}
			return nil
		},
	}
}

func newAgentAddCmd() *cobra.Command {
	var (
		projectDir string
		globalDir  string
		detectPath string
	)

	cmd := &cobra.Command{
		Use:   "add <name>",
		Short: "Add a custom agent",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := NewConfig()
			if err != nil {
				return err
			}
			defer cfg.Close()

			a := agent.Adapter{
				Name:        args[0],
				DisplayName: args[0],
				ProjectDir:  projectDir,
				GlobalDir:   globalDir,
				DetectPaths: []string{detectPath},
			}

			if err := cfg.Store.InsertAgent(a); err != nil {
				return fmt.Errorf("add agent: %w", err)
			}
			color.Green("Added agent %q", args[0])
			return nil
		},
	}

	cmd.Flags().StringVar(&projectDir, "project-dir", "", "Project-level skills directory")
	cmd.Flags().StringVar(&globalDir, "global-dir", "", "Global skills directory")
	cmd.Flags().StringVar(&detectPath, "detect", "", "Detection path")
	cmd.MarkFlagRequired("project-dir")
	cmd.MarkFlagRequired("global-dir")
	return cmd
}

func newAgentRemoveCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "remove <name>",
		Short: "Remove a custom agent",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := NewConfig()
			if err != nil {
				return err
			}
			defer cfg.Close()

			if err := cfg.Store.DeleteAgent(args[0]); err != nil {
				return err
			}
			color.Green("Removed agent %q", args[0])
			return nil
		},
	}
}

func newAgentSkillsCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "skills <name>",
		Short: "List skills in an agent's global directory",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			adapters := agent.Builtin()
			var ag *agent.Adapter
			for i := range adapters {
				if adapters[i].Name == args[0] {
					ag = &adapters[i]
					break
				}
			}
			if ag == nil {
				return fmt.Errorf("agent %q not found", args[0])
			}

			home, _ := os.UserHomeDir()
			skills, err := project.ScanSkills(home, []agent.Adapter{*ag})
			if err != nil {
				return err
			}

			if len(skills) == 0 {
				fmt.Printf("No skills in %s global directory.\n", ag.DisplayName)
				return nil
			}

			for _, sk := range skills {
				status := color.GreenString("enabled")
				if !sk.Enabled {
					status = color.YellowString("disabled")
				}
				fmt.Printf("  %-25s %s\n", color.CyanString(sk.SkillName), status)
			}
			return nil
		},
	}
}

func newAgentAddSkillCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "add-skill <agent> <skill>",
		Short: "Add a skill from central library to agent's global directory",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			agentName, skillName := args[0], args[1]

			cfg, err := NewConfig()
			if err != nil {
				return err
			}
			defer cfg.Close()

			sk, err := cfg.Store.GetSkill(skillName)
			if err != nil {
				return fmt.Errorf("skill %q: %w", skillName, err)
			}

			adapters := agent.Builtin()
			var ag *agent.Adapter
			for i := range adapters {
				if adapters[i].Name == agentName {
					ag = &adapters[i]
					break
				}
			}
			if ag == nil {
				return fmt.Errorf("agent %q not found", agentName)
			}

			home, _ := os.UserHomeDir()
			targetDir := filepath.Join(home, ag.GlobalDir)
			targetPath := filepath.Join(targetDir, sk.Name)

			if err := os.MkdirAll(targetDir, 0o755); err != nil {
				return fmt.Errorf("create directory: %w", err)
			}

			if err := skmsync.SyncSkill(sk.CentralPath, targetPath, "symlink"); err != nil {
				return fmt.Errorf("sync: %w", err)
			}

			_ = cfg.Store.UpsertTarget(sk.ID, ag.Name, targetPath, "symlink", sk.ContentHash)
			color.Green("✓ Added %s to %s", sk.Name, ag.DisplayName)
			return nil
		},
	}
}

func newAgentRemoveSkillCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "remove-skill <agent> <skill>",
		Short: "Remove a skill from agent's global directory",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			agentName, skillName := args[0], args[1]

			adapters := agent.Builtin()
			var ag *agent.Adapter
			for i := range adapters {
				if adapters[i].Name == agentName {
					ag = &adapters[i]
					break
				}
			}
			if ag == nil {
				return fmt.Errorf("agent %q not found", agentName)
			}

			home, _ := os.UserHomeDir()
			targetPath := filepath.Join(home, ag.GlobalDir, skillName)

			if err := os.RemoveAll(targetPath); err != nil {
				return fmt.Errorf("remove: %w", err)
			}

			color.Green("✓ Removed %s from %s", skillName, ag.DisplayName)
			return nil
		},
	}
}
