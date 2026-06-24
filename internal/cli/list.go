package cli

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

func newListCmd() *cobra.Command {
	var agentFilter string

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List installed skills",
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

			if len(skills) == 0 {
				fmt.Println("No skills installed. Use 'skm install <source>' to add one.")
				return nil
			}

			fmt.Printf("Installed skills (%d):\n\n", len(skills))
			for _, s := range skills {
				targets, _ := cfg.Store.ListTargets(s.ID)

				agentList := ""
				for _, t := range targets {
					if agentFilter != "" && t.Agent != agentFilter {
						continue
					}
					if agentList != "" {
						agentList += ", "
					}
					agentList += t.Agent
				}

				if agentFilter != "" && agentList == "" {
					continue
				}

				name := color.CyanString(s.Name)
				if agentList != "" {
					fmt.Printf("  %s  [%s]\n", name, color.GreenString(agentList))
				} else {
					fmt.Printf("  %s\n", name)
				}
				if s.Description != "" {
					fmt.Printf("    %s\n", s.Description)
				}
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&agentFilter, "agent", "a", "", "Filter by agent")
	return cmd
}
