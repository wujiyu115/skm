package cli

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/ejoy/skm/internal/agent"
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
