package cli

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/ejoy/skm/internal/fetcher"
	"github.com/ejoy/skm/internal/skill"
)

func newUpdateCmd() *cobra.Command {
	var all bool

	cmd := &cobra.Command{
		Use:   "update [skill]",
		Short: "Update git-sourced skills",
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

			var toUpdate []*skill.Skill
			if all {
				for _, sk := range skills {
					if sk.SourceType == "git" {
						toUpdate = append(toUpdate, sk)
					}
				}
			} else if len(args) > 0 {
				sk, err := cfg.Store.GetSkill(args[0])
				if err != nil {
					return err
				}
				if sk.SourceType != "git" {
					return fmt.Errorf("skill %q is not git-sourced", sk.Name)
				}
				toUpdate = append(toUpdate, sk)
			} else {
				return fmt.Errorf("specify a skill name or use --all")
			}

			if len(toUpdate) == 0 {
				fmt.Println("No git-sourced skills to update.")
				return nil
			}

			for _, sk := range toUpdate {
				src := fetcher.Parse(sk.SourceRef)
				fmt.Printf("Updating %s from %s...\n", sk.Name, sk.SourceRef)

				tmpDir, err := fetcher.CloneToTemp(src)
				if err != nil {
					color.Red("  ✗ %s: %v", sk.Name, err)
					continue
				}

				discovered, err := skill.Discover(tmpDir, src.Subpath)
				fetcher.CleanupTemp(tmpDir)
				if err != nil || len(discovered) == 0 {
					color.Red("  ✗ %s: no skills found", sk.Name)
					continue
				}

				var found *skill.Skill
				for _, d := range discovered {
					if d.Name == sk.Name {
						found = d
						break
					}
				}
				if found == nil {
					found = discovered[0]
				}

				newHash, _ := skill.HashDir(found.CentralPath)
				if newHash == sk.ContentHash {
					fmt.Printf("  %s already up to date\n", sk.Name)
					continue
				}

				if err := fetcher.CopyLocal(found.CentralPath, sk.CentralPath); err != nil {
					color.Red("  ✗ %s: copy failed: %v", sk.Name, err)
					continue
				}

				cfg.Store.UpdateSkillHash(sk.ID, newHash)
				color.Green("  ✓ %s updated", sk.Name)
			}

			cfg.WriteMetadata()
			return nil
		},
	}

	cmd.Flags().BoolVar(&all, "all", false, "Update all git-sourced skills")
	return cmd
}
