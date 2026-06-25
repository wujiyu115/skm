package cli

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/ejoy/skm/internal/skill"
)

func newSearchCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "search <query>",
		Short: "Search skills by name, description, or tag",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := NewConfig()
			if err != nil {
				return err
			}
			defer cfg.Close()

			results, err := searchSkills(cfg, args[0])
			if err != nil {
				return err
			}

			if len(results) == 0 {
				fmt.Printf("No skills matching %q\n", args[0])
				return nil
			}

			printSearchTable(results)
			return nil
		},
	}
}

type searchResult struct {
	Skill *skill.Skill
	Tags  []string
}

func searchSkills(cfg *Config, query string) ([]searchResult, error) {
	skills, err := cfg.Store.ListSkills()
	if err != nil {
		return nil, err
	}

	q := strings.ToLower(query)
	var results []searchResult

	for _, sk := range skills {
		tags, _ := cfg.Store.ListSkillTags(sk.ID)

		if matchesQuery(sk, tags, q) {
			results = append(results, searchResult{Skill: sk, Tags: tags})
		}
	}

	return results, nil
}

func matchesQuery(sk *skill.Skill, tags []string, query string) bool {
	if strings.Contains(strings.ToLower(sk.Name), query) {
		return true
	}
	if strings.Contains(strings.ToLower(sk.Description), query) {
		return true
	}
	for _, tag := range tags {
		if strings.Contains(strings.ToLower(tag), query) {
			return true
		}
	}
	return false
}

func printSearchTable(results []searchResult) {
	// Calculate column widths.
	nameW, descW, tagW := 4, 11, 4 // header lengths
	for _, r := range results {
		if len(r.Skill.Name) > nameW {
			nameW = len(r.Skill.Name)
		}
		desc := truncate(r.Skill.Description, 50)
		if len(desc) > descW {
			descW = len(desc)
		}
		tagStr := strings.Join(r.Tags, ", ")
		if len(tagStr) > tagW {
			tagW = len(tagStr)
		}
	}

	// Header.
	hdr := fmt.Sprintf("%-*s  %-*s  %-*s  %s", nameW, "Name", descW, "Description", tagW, "Tags", "Enabled")
	fmt.Println(hdr)
	fmt.Println(strings.Repeat("-", len(hdr)))

	// Rows.
	for _, r := range results {
		enabled := color.GreenString("yes")
		if !r.Skill.Enabled {
			enabled = color.RedString("no")
		}
		fmt.Printf("%-*s  %-*s  %-*s  %s\n",
			nameW, r.Skill.Name,
			descW, truncate(r.Skill.Description, 50),
			tagW, strings.Join(r.Tags, ", "),
			enabled,
		)
	}

	fmt.Printf("\n%d result(s)\n", len(results))
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max-3] + "..."
}
