package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/fatih/color"
	"github.com/google/uuid"
	"github.com/spf13/cobra"

	"github.com/ejoy/skm/internal/agent"
	"github.com/ejoy/skm/internal/project"
)

func newProjectCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "project",
		Short: "Manage project workspaces",
	}

	cmd.AddCommand(
		newProjectAddCmd(),
		newProjectListCmd(),
		newProjectRemoveCmd(),
		newProjectScanCmd(),
	)

	return cmd
}

func newProjectAddCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "add <path>",
		Short: "Register a project directory",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := NewConfig()
			if err != nil {
				return err
			}
			defer cfg.Close()

			return projectAdd(cfg, args[0])
		},
	}
}

func projectAdd(cfg *Config, path string) error {
	id := uuid.New().String()
	name := filepath.Base(path)

	absPath, err := filepath.Abs(path)
	if err != nil {
		return fmt.Errorf("resolve path: %w", err)
	}

	info, err := os.Stat(absPath)
	if err != nil {
		return fmt.Errorf("path does not exist: %s", absPath)
	}
	if !info.IsDir() {
		return fmt.Errorf("path is not a directory: %s", absPath)
	}

	name = filepath.Base(absPath)
	if err := cfg.Store.InsertProject(id, name, absPath); err != nil {
		return fmt.Errorf("register project: %w", err)
	}

	if err := cfg.Store.InsertAuditLog("project_create", name, fmt.Sprintf("path=%s", absPath)); err != nil {
		fmt.Fprintf(os.Stderr, "warning: audit log: %v\n", err)
	}

	color.Green("Registered project %q (%s)", name, id)
	return nil
}

func newProjectListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List registered projects",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := NewConfig()
			if err != nil {
				return err
			}
			defer cfg.Close()

			return projectList(cfg)
		},
	}
}

func projectList(cfg *Config) error {
	projects, err := cfg.Store.ListProjects()
	if err != nil {
		return err
	}
	if len(projects) == 0 {
		fmt.Println("No projects. Use 'skm project add <path>' to register one.")
		return nil
	}
	for _, p := range projects {
		fmt.Printf("  %s  %s  %s\n", color.CyanString(p.ID[:8]), color.WhiteString(p.Name), color.HiBlackString(p.Path))
	}
	return nil
}

func newProjectRemoveCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "remove <id>",
		Short: "Unregister a project (DB only, does not delete files)",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := NewConfig()
			if err != nil {
				return err
			}
			defer cfg.Close()

			return projectRemove(cfg, args[0])
		},
	}
}

func projectRemove(cfg *Config, id string) error {
	proj, _ := cfg.Store.GetProjectByID(id)

	if err := cfg.Store.DeleteProject(id); err != nil {
		return fmt.Errorf("delete project: %w", err)
	}

	target := id
	if proj != nil {
		target = proj.Name
	}
	if err := cfg.Store.InsertAuditLog("project_delete", target, ""); err != nil {
		fmt.Fprintf(os.Stderr, "warning: audit log: %v\n", err)
	}

	color.Yellow("Removed project %q", target)
	return nil
}

func newProjectScanCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "scan <id>",
		Short: "Scan and display skills installed in a project",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := NewConfig()
			if err != nil {
				return err
			}
			defer cfg.Close()

			return projectScan(cfg, args[0])
		},
	}
}

func projectScan(cfg *Config, id string) error {
	proj, err := cfg.Store.GetProjectByID(id)
	if err != nil {
		return fmt.Errorf("project not found: %w", err)
	}

	adapters := agent.Builtin()
	skills, err := project.ScanSkills(proj.Path, adapters)
	if err != nil {
		return fmt.Errorf("scan skills: %w", err)
	}

	if len(skills) == 0 {
		fmt.Printf("No skills found in project %q (%s)\n", proj.Name, proj.Path)
		return nil
	}

	fmt.Printf("Skills in project %q (%s):\n\n", proj.Name, proj.Path)
	fmt.Printf("  %-15s %-25s %-8s %s\n",
		color.HiBlackString("AGENT"),
		color.HiBlackString("SKILL"),
		color.HiBlackString("STATUS"),
		color.HiBlackString("PATH"))

	for _, sk := range skills {
		status := color.GreenString("enabled")
		if !sk.Enabled {
			status = color.YellowString("disabled")
		}
		fmt.Printf("  %-15s %-25s %-8s %s\n",
			sk.AgentDisplay, sk.SkillName, status, color.HiBlackString(sk.SkillPath))
	}

	return nil
}
