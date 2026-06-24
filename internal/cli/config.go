package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/ejoy/skm/internal/agent"
	"github.com/ejoy/skm/internal/store"
)

type Config struct {
	HomeDir   string
	DBPath    string
	SkillsDir string
	CacheDir  string
	Store     *store.Store
}

func NewConfig() (*Config, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("get home dir: %w", err)
	}

	skmDir := filepath.Join(home, ".skm")
	cfg := &Config{
		HomeDir:   home,
		DBPath:    filepath.Join(skmDir, "skm.db"),
		SkillsDir: filepath.Join(skmDir, "skills"),
		CacheDir:  filepath.Join(skmDir, "cache"),
	}

	os.MkdirAll(cfg.SkillsDir, 0755)
	os.MkdirAll(cfg.CacheDir, 0755)

	s, err := store.New(cfg.DBPath)
	if err != nil {
		return nil, err
	}
	cfg.Store = s

	s.SeedBuiltinAgents(agent.Builtin())

	return cfg, nil
}

func (c *Config) Close() {
	if c.Store != nil {
		c.Store.Close()
	}
}
