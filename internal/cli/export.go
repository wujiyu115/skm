package cli

import (
	"encoding/json"
	"fmt"
	"runtime"

	"github.com/spf13/cobra"

	"github.com/ejoy/skm/internal/agent"
)

func newExportCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "export",
		Short: "Export all skills as JSON",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := NewConfig()
			if err != nil {
				return err
			}
			defer cfg.Close()

			return exportSkills(cfg)
		},
	}
}

type exportEntry struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	SourceType  string   `json:"source_type"`
	SourceRef   string   `json:"source_ref"`
	Enabled     bool     `json:"enabled"`
	Tags        []string `json:"tags"`
}

func exportSkills(cfg *Config) error {
	skills, err := cfg.Store.ListSkills()
	if err != nil {
		return err
	}

	entries := make([]exportEntry, 0, len(skills))
	for _, sk := range skills {
		tags, _ := cfg.Store.ListSkillTags(sk.ID)
		if tags == nil {
			tags = []string{}
		}
		entries = append(entries, exportEntry{
			Name:        sk.Name,
			Description: sk.Description,
			SourceType:  sk.SourceType,
			SourceRef:   sk.SourceRef,
			Enabled:     sk.Enabled,
			Tags:        tags,
		})
	}

	data, err := json.MarshalIndent(entries, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(data))
	return nil
}

func newInfoCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "info",
		Short: "Show diagnostics information",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := NewConfig()
			if err != nil {
				return err
			}
			defer cfg.Close()

			return showInfo(cfg)
		},
	}
}

func showInfo(cfg *Config) error {
	skills, _ := cfg.Store.ListSkills()
	groups, _ := cfg.Store.ListGroups()
	agents := agent.Builtin()

	fmt.Printf("SKM Diagnostics\n\n")
	fmt.Printf("  Version:     %s\n", Version)
	fmt.Printf("  OS/Arch:     %s/%s\n", runtime.GOOS, runtime.GOARCH)
	fmt.Printf("  Go version:  %s\n", runtime.Version())
	fmt.Printf("  Skills dir:  %s\n", cfg.SkillsDir)
	fmt.Printf("  Cache dir:   %s\n", cfg.CacheDir)
	fmt.Printf("  DB path:     %s\n", cfg.DBPath)
	fmt.Printf("  Skills:      %d\n", len(skills))
	fmt.Printf("  Agents:      %d\n", len(agents))
	fmt.Printf("  Groups:      %d\n", len(groups))

	return nil
}
