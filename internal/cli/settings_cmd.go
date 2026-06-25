package cli

import (
	"fmt"
	"sort"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

func newConfigCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Manage SKM settings",
	}

	cmd.AddCommand(
		newConfigListCmd(),
		newConfigGetCmd(),
		newConfigSetCmd(),
	)

	return cmd
}

func newConfigListCmd() *cobra.Command {
	var jsonOutput bool

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all settings",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := NewConfig()
			if err != nil {
				return err
			}
			defer cfg.Close()

			settings, err := cfg.Store.ListSettings()
			if err != nil {
				return err
			}

			// Add config-derived values so the list is complete.
			settings["skills_dir"] = cfg.SkillsDir
			settings["cache_dir"] = cfg.CacheDir

			if jsonOutput {
				return printJSON(settings)
			}

			if len(settings) == 0 {
				fmt.Println("No settings configured.")
				return nil
			}

			// Sort keys for deterministic output.
			keys := make([]string, 0, len(settings))
			for k := range settings {
				keys = append(keys, k)
			}
			sort.Strings(keys)

			fmt.Println("Settings:")
			for _, k := range keys {
				fmt.Printf("  %-28s %s\n", color.CyanString(k), settings[k])
			}
			return nil
		},
	}

	cmd.Flags().BoolVar(&jsonOutput, "json", false, "Output as JSON")
	return cmd
}

func newConfigGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get <key>",
		Short: "Get a single setting value",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := NewConfig()
			if err != nil {
				return err
			}
			defer cfg.Close()

			val, err := cfg.Store.GetSetting(args[0])
			if err != nil {
				return fmt.Errorf("setting %q: %w", args[0], err)
			}

			fmt.Println(val)
			return nil
		},
	}
}

func newConfigSetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "set <key> <value>",
		Short: "Set a single setting value",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := NewConfig()
			if err != nil {
				return err
			}
			defer cfg.Close()

			if err := cfg.Store.SetSetting(args[0], args[1]); err != nil {
				return fmt.Errorf("set %q: %w", args[0], err)
			}

			color.Green("Set %s = %s", args[0], args[1])
			return nil
		},
	}
}
