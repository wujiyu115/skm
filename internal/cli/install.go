package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/fatih/color"
	"github.com/google/uuid"
	"github.com/spf13/cobra"

	"github.com/ejoy/skm/internal/agent"
	"github.com/ejoy/skm/internal/fetcher"
	"github.com/ejoy/skm/internal/skill"
	skmsync "github.com/ejoy/skm/internal/sync"
)

func newInstallCmd() *cobra.Command {
	var (
		global   bool
		agents   []string
		yes      bool
		listOnly bool
	)

	cmd := &cobra.Command{
		Use:   "install <source>",
		Short: "Install a skill from GitHub URL, shorthand, or local path",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := NewConfig()
			if err != nil {
				return err
			}
			defer cfg.Close()

			return runInstall(cfg, args[0], global, agents, yes, listOnly)
		},
	}

	cmd.Flags().BoolVarP(&global, "global", "g", false, "Install to global skills directory")
	cmd.Flags().StringSliceVarP(&agents, "agent", "a", nil, "Target agent(s)")
	cmd.Flags().BoolVarP(&yes, "yes", "y", false, "Skip confirmation prompts")
	cmd.Flags().BoolVarP(&listOnly, "list", "l", false, "List discovered skills without installing")

	return cmd
}

func runInstall(cfg *Config, source string, global bool, agentNames []string, yes, listOnly bool) error {
	src := fetcher.Parse(source)

	var baseDir string
	if src.Type == "local" {
		absPath, err := filepath.Abs(src.Subpath)
		if err != nil {
			return fmt.Errorf("resolve path: %w", err)
		}
		baseDir = absPath
	} else {
		fmt.Printf("Cloning %s...\n", src.GitURL())
		var err error
		baseDir, err = fetcher.CloneToTemp(src)
		if err != nil {
			return fmt.Errorf("clone: %w", err)
		}
		defer fetcher.CleanupTemp(baseDir)
	}

	skills, err := skill.Discover(baseDir, src.Subpath)
	if err != nil {
		return fmt.Errorf("discover skills: %w", err)
	}
	if len(skills) == 0 {
		return fmt.Errorf("no skills found in %s", source)
	}

	if listOnly {
		fmt.Printf("Found %d skill(s):\n", len(skills))
		for _, s := range skills {
			fmt.Printf("  %s — %s\n", color.CyanString(s.Name), s.Description)
		}
		return nil
	}

	targetAgents := agent.Resolve(agentNames)
	if len(targetAgents) == 0 {
		return fmt.Errorf("no agents specified or detected. Use -a <agent> or install an agent first")
	}

	cwd, _ := os.Getwd()
	green := color.New(color.FgGreen)

	for _, sk := range skills {
		centralPath := filepath.Join(cfg.SkillsDir, sk.Name)

		if _, err := os.Stat(centralPath); err == nil {
			fmt.Printf("  %s already in library, updating...\n", sk.Name)
			os.RemoveAll(centralPath)
		}

		if err := fetcher.CopyLocal(sk.CentralPath, centralPath); err != nil {
			color.Red("  ✗ Failed to copy %s: %v", sk.Name, err)
			continue
		}

		hash, _ := skill.HashDir(centralPath)
		sk.ID = uuid.New().String()
		sk.CentralPath = centralPath
		sk.ContentHash = hash
		if src.Type == "local" {
			sk.SourceType = "local"
			sk.SourceRef = source
		} else {
			sk.SourceType = "git"
			sk.SourceRef = src.GitURL()
		}

		cfg.Store.DeleteSkillByName(sk.Name)
		if err := cfg.Store.InsertSkill(sk); err != nil {
			color.Red("  ✗ Failed to save %s: %v", sk.Name, err)
			continue
		}

		for _, ag := range targetAgents {
			targetDir, err := agent.InstallPath(ag, global, cwd)
			if err != nil {
				color.Red("  ✗ %s → %s: %v", sk.Name, ag.DisplayName, err)
				continue
			}
			targetPath := filepath.Join(targetDir, sk.Name)

			if err := skmsync.SyncSkill(centralPath, targetPath, "symlink"); err != nil {
				color.Red("  ✗ %s → %s: %v", sk.Name, ag.DisplayName, err)
				continue
			}

			cfg.Store.UpsertTarget(sk.ID, ag.Name, targetPath, "symlink", hash)
			green.Printf("  ✓ %s → %s (%s)\n", sk.Name, ag.DisplayName, targetPath)
		}
	}

	cfg.WriteMetadata()
	return nil
}

