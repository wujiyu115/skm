package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

func newShowCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "show <skill>",
		Short: "Show skill details and SKILL.md content",
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

			color.Cyan("Skill: %s\n", sk.Name)
			fmt.Printf("  Description:  %s\n", sk.Description)
			fmt.Printf("  Source:        %s (%s)\n", sk.SourceRef, sk.SourceType)
			fmt.Printf("  Path:          %s\n", sk.CentralPath)
			fmt.Printf("  Hash:          %s\n", sk.ContentHash)

			targets, _ := cfg.Store.ListTargets(sk.ID)
			if len(targets) > 0 {
				fmt.Println("  Synced to:")
				for _, t := range targets {
					fmt.Printf("    %s → %s (%s)\n", color.GreenString(t.Agent), t.TargetPath, t.Mode)
				}
			}

			mdPath := filepath.Join(sk.CentralPath, "SKILL.md")
			content, err := os.ReadFile(mdPath)
			if err == nil {
				fmt.Printf("\n--- SKILL.md ---\n%s\n", string(content))
			}

			return nil
		},
	}
}
