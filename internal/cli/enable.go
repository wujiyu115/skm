package cli

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/ejoy/skm/internal/logger"
)

func newEnableCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "enable <skill>",
		Short: "Enable an installed skill",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := NewConfig()
			if err != nil {
				return err
			}
			defer cfg.Close()

			return setSkillEnabled(cfg, args[0], true)
		},
	}
}

func newDisableCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "disable <skill>",
		Short: "Disable an installed skill",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := NewConfig()
			if err != nil {
				return err
			}
			defer cfg.Close()

			return setSkillEnabled(cfg, args[0], false)
		},
	}
}

func setSkillEnabled(cfg *Config, name string, enabled bool) error {
	sk, err := cfg.Store.GetSkill(name)
	if err != nil {
		return fmt.Errorf("skill %q not found", name)
	}

	if err := cfg.Store.SetSkillEnabled(sk.ID, enabled); err != nil {
		return fmt.Errorf("update skill: %w", err)
	}

	action := "enable"
	if !enabled {
		action = "disable"
	}

	if err := cfg.Store.InsertAuditLog(action, sk.Name, ""); err != nil {
		logger.Warn("audit log failed", "err", err)
	}

	if enabled {
		color.Green("Enabled %s", sk.Name)
	} else {
		color.Yellow("Disabled %s", sk.Name)
	}
	return nil
}
