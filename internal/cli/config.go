package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/ejoy/skm/internal/agent"
	"github.com/ejoy/skm/internal/logger"
	"github.com/ejoy/skm/internal/store"
)

type Config struct {
	HomeDir   string
	DBPath    string
	SkillsDir string
	CacheDir  string
	MetaDir   string
	Store     *store.Store
}

func NewConfig() (*Config, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("get home dir: %w", err)
	}

	skmDir := filepath.Join(home, ".skm")
	logDir := filepath.Join(skmDir, "logs")
	cfg := &Config{
		HomeDir:   home,
		DBPath:    filepath.Join(skmDir, "skm.db"),
		SkillsDir: filepath.Join(skmDir, "skills"),
		CacheDir:  filepath.Join(skmDir, "cache"),
		MetaDir:   filepath.Join(skmDir, "metadata"),
	}

	os.MkdirAll(cfg.SkillsDir, 0755)
	os.MkdirAll(cfg.CacheDir, 0755)

	if err := logger.Setup(logger.Config{Dir: logDir, MaxSizeMB: 10, MaxBackups: 5, MaxAgeDays: 30}); err != nil {
		fmt.Fprintf(os.Stderr, "warning: setup logger: %v\n", err)
	}

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

// WriteMetadata writes JSON metadata mirror files. Errors are logged but
// not treated as fatal since metadata is a convenience mirror of the DB.
func (c *Config) WriteMetadata() {
	if err := c.Store.WriteMetadata(c.MetaDir); err != nil {
		fmt.Fprintf(os.Stderr, "warning: write metadata: %v\n", err)
	}
}
