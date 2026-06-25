package cli

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

func newTagCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tag",
		Short: "Manage skill tags",
	}

	cmd.AddCommand(
		newTagListCmd(),
		newTagAddCmd(),
		newTagRemoveCmd(),
		newTagRenameCmd(),
		newTagDeleteCmd(),
	)

	return cmd
}

func newTagListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List all tags",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := NewConfig()
			if err != nil {
				return err
			}
			defer cfg.Close()

			tags, err := cfg.Store.ListTags()
			if err != nil {
				return err
			}
			if len(tags) == 0 {
				fmt.Println("No tags. Use 'skm tag add <skill> <tag>' to add one.")
				return nil
			}
			for _, tag := range tags {
				fmt.Printf("  %s\n", color.CyanString(tag))
			}
			return nil
		},
	}
}

func newTagAddCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "add <skill> <tag>",
		Short: "Add a tag to a skill",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := NewConfig()
			if err != nil {
				return err
			}
			defer cfg.Close()

			return tagAdd(cfg, args[0], args[1])
		},
	}
}

func tagAdd(cfg *Config, skillName, tag string) error {
	sk, err := cfg.Store.GetSkill(skillName)
	if err != nil {
		return fmt.Errorf("skill %q: %w", skillName, err)
	}

	if err := cfg.Store.AddTag(sk.ID, tag); err != nil {
		return fmt.Errorf("add tag: %w", err)
	}

	color.Green("Added tag %q to %s", tag, sk.Name)
	return nil
}

func newTagRemoveCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "remove <skill> <tag>",
		Short: "Remove a tag from a skill",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := NewConfig()
			if err != nil {
				return err
			}
			defer cfg.Close()

			return tagRemove(cfg, args[0], args[1])
		},
	}
}

func tagRemove(cfg *Config, skillName, tag string) error {
	sk, err := cfg.Store.GetSkill(skillName)
	if err != nil {
		return fmt.Errorf("skill %q: %w", skillName, err)
	}

	if err := cfg.Store.RemoveTag(sk.ID, tag); err != nil {
		return fmt.Errorf("remove tag: %w", err)
	}

	color.Yellow("Removed tag %q from %s", tag, sk.Name)
	return nil
}

func newTagRenameCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "rename <old> <new>",
		Short: "Rename a tag globally",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := NewConfig()
			if err != nil {
				return err
			}
			defer cfg.Close()

			return tagRename(cfg, args[0], args[1])
		},
	}
}

func tagRename(cfg *Config, oldTag, newTag string) error {
	if err := cfg.Store.RenameTag(oldTag, newTag); err != nil {
		return fmt.Errorf("rename tag: %w", err)
	}

	color.Green("Renamed tag %q to %q", oldTag, newTag)
	return nil
}

func newTagDeleteCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "delete <tag>",
		Short: "Delete a tag from all skills",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := NewConfig()
			if err != nil {
				return err
			}
			defer cfg.Close()

			if err := cfg.Store.DeleteTag(args[0]); err != nil {
				return fmt.Errorf("delete tag: %w", err)
			}
			color.Green("Deleted tag %q from all skills", args[0])
			return nil
		},
	}
}
