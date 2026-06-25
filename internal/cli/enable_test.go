package cli

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/ejoy/skm/internal/skill"
	"github.com/ejoy/skm/internal/store"
)

func testConfig(t *testing.T) *Config {
	t.Helper()
	s, err := store.New(filepath.Join(t.TempDir(), "test.db"))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { s.Close() })
	return &Config{Store: s}
}

func insertTestSkill(t *testing.T, s *store.Store, name string) {
	t.Helper()
	sk := &skill.Skill{
		ID:          "test-" + name,
		Name:        name,
		Description: "test skill " + name,
		SourceType:  "local",
		SourceRef:   "/tmp/" + name,
		CentralPath: "/tmp/skills/" + name,
		ContentHash: "abc123",
		Enabled:     true,
	}
	if err := s.InsertSkill(sk); err != nil {
		t.Fatalf("InsertSkill(%q): %v", name, err)
	}
}

func TestSetSkillEnabled_Enable(t *testing.T) {
	cfg := testConfig(t)
	insertTestSkill(t, cfg.Store, "my-skill")

	// First disable it via the store, then re-enable through setSkillEnabled.
	if err := cfg.Store.SetSkillEnabled("test-my-skill", false); err != nil {
		t.Fatal(err)
	}

	if err := setSkillEnabled(cfg, "my-skill", true); err != nil {
		t.Fatalf("setSkillEnabled(enable) error: %v", err)
	}

	sk, err := cfg.Store.GetSkill("my-skill")
	if err != nil {
		t.Fatal(err)
	}
	if !sk.Enabled {
		t.Error("expected skill to be enabled after setSkillEnabled(true)")
	}
}

func TestSetSkillEnabled_Disable(t *testing.T) {
	cfg := testConfig(t)
	insertTestSkill(t, cfg.Store, "my-skill")

	if err := setSkillEnabled(cfg, "my-skill", false); err != nil {
		t.Fatalf("setSkillEnabled(disable) error: %v", err)
	}

	sk, err := cfg.Store.GetSkill("my-skill")
	if err != nil {
		t.Fatal(err)
	}
	if sk.Enabled {
		t.Error("expected skill to be disabled after setSkillEnabled(false)")
	}
}

func TestSetSkillEnabled_NotFound(t *testing.T) {
	cfg := testConfig(t)

	err := setSkillEnabled(cfg, "nonexistent", true)
	if err == nil {
		t.Fatal("expected error for nonexistent skill")
	}
	if !strings.Contains(err.Error(), "nonexistent") {
		t.Errorf("error should mention skill name, got: %v", err)
	}
}

func TestSetSkillEnabled_AuditLog(t *testing.T) {
	cfg := testConfig(t)
	insertTestSkill(t, cfg.Store, "audited")

	if err := setSkillEnabled(cfg, "audited", false); err != nil {
		t.Fatal(err)
	}

	entries, err := cfg.Store.ListAuditLog(10)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) == 0 {
		t.Fatal("expected at least one audit log entry")
	}
	if entries[0].Action != "disable" {
		t.Errorf("audit action = %q, want %q", entries[0].Action, "disable")
	}
	if entries[0].Target != "audited" {
		t.Errorf("audit target = %q, want %q", entries[0].Target, "audited")
	}
}
