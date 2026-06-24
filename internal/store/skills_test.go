package store

import (
	"path/filepath"
	"testing"

	"github.com/ejoy/skm/internal/agent"
	"github.com/ejoy/skm/internal/skill"
)

func testStore(t *testing.T) *Store {
	t.Helper()
	s, err := New(filepath.Join(t.TempDir(), "test.db"))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { s.Close() })
	return s
}

func TestInsertAndGetSkill(t *testing.T) {
	s := testStore(t)

	sk := &skill.Skill{
		ID:          "test-id-1",
		Name:        "my-skill",
		Description: "A test skill",
		SourceType:  "git",
		SourceRef:   "https://github.com/owner/repo",
		CentralPath: "/home/user/.skm/skills/my-skill",
		ContentHash: "abc123",
		Enabled:     true,
	}

	if err := s.InsertSkill(sk); err != nil {
		t.Fatalf("InsertSkill() error: %v", err)
	}

	got, err := s.GetSkill("my-skill")
	if err != nil {
		t.Fatalf("GetSkill() error: %v", err)
	}
	if got.Name != "my-skill" {
		t.Errorf("Name = %q", got.Name)
	}
	if got.Description != "A test skill" {
		t.Errorf("Description = %q", got.Description)
	}
	if got.ID != "test-id-1" {
		t.Errorf("ID = %q, want %q", got.ID, "test-id-1")
	}
	if got.SourceType != "git" {
		t.Errorf("SourceType = %q", got.SourceType)
	}
	if got.ContentHash != "abc123" {
		t.Errorf("ContentHash = %q", got.ContentHash)
	}
	if !got.Enabled {
		t.Error("Enabled = false, want true")
	}
}

func TestInsertSkill_DuplicateName(t *testing.T) {
	s := testStore(t)
	sk := &skill.Skill{ID: "1", Name: "dup", SourceType: "local", CentralPath: "/a", Enabled: true}
	if err := s.InsertSkill(sk); err != nil {
		t.Fatal(err)
	}
	sk2 := &skill.Skill{ID: "2", Name: "dup", SourceType: "local", CentralPath: "/b", Enabled: true}
	if err := s.InsertSkill(sk2); err == nil {
		t.Fatal("expected error on duplicate name insert")
	}
}

func TestGetSkill_NotFound(t *testing.T) {
	s := testStore(t)
	_, err := s.GetSkill("nonexistent")
	if err == nil {
		t.Fatal("expected error for nonexistent skill")
	}
}

func TestGetSkillByID(t *testing.T) {
	s := testStore(t)
	sk := &skill.Skill{ID: "byid-1", Name: "byid-skill", SourceType: "local", CentralPath: "/byid", Enabled: true}
	s.InsertSkill(sk)

	got, err := s.GetSkillByID("byid-1")
	if err != nil {
		t.Fatalf("GetSkillByID() error: %v", err)
	}
	if got.Name != "byid-skill" {
		t.Errorf("Name = %q, want %q", got.Name, "byid-skill")
	}
}

func TestListSkills(t *testing.T) {
	s := testStore(t)

	s.InsertSkill(&skill.Skill{ID: "1", Name: "a", SourceType: "local", CentralPath: "/a", Enabled: true})
	s.InsertSkill(&skill.Skill{ID: "2", Name: "b", SourceType: "local", CentralPath: "/b", Enabled: true})

	skills, err := s.ListSkills()
	if err != nil {
		t.Fatalf("ListSkills() error: %v", err)
	}
	if len(skills) != 2 {
		t.Fatalf("expected 2 skills, got %d", len(skills))
	}
	// Should be ordered by name
	if skills[0].Name != "a" || skills[1].Name != "b" {
		t.Errorf("unexpected order: %q, %q", skills[0].Name, skills[1].Name)
	}
}

func TestListSkills_Empty(t *testing.T) {
	s := testStore(t)
	skills, err := s.ListSkills()
	if err != nil {
		t.Fatalf("ListSkills() error: %v", err)
	}
	if len(skills) != 0 {
		t.Fatalf("expected 0 skills, got %d", len(skills))
	}
}

func TestDeleteSkill(t *testing.T) {
	s := testStore(t)
	s.InsertSkill(&skill.Skill{ID: "del-1", Name: "gone", SourceType: "local", CentralPath: "/g", Enabled: true})

	if err := s.DeleteSkill("del-1"); err != nil {
		t.Fatalf("DeleteSkill() error: %v", err)
	}

	_, err := s.GetSkill("gone")
	if err == nil {
		t.Fatal("expected error after delete")
	}
}

func TestDeleteSkill_NotFound(t *testing.T) {
	s := testStore(t)
	err := s.DeleteSkill("nonexistent-id")
	if err == nil {
		t.Fatal("expected error deleting nonexistent skill")
	}
}

func TestDeleteSkillByName(t *testing.T) {
	s := testStore(t)
	s.InsertSkill(&skill.Skill{ID: "dbn-1", Name: "delete-by-name", SourceType: "local", CentralPath: "/dbn", Enabled: true})

	if err := s.DeleteSkillByName("delete-by-name"); err != nil {
		t.Fatalf("DeleteSkillByName() error: %v", err)
	}

	_, err := s.GetSkill("delete-by-name")
	if err == nil {
		t.Fatal("expected error after delete by name")
	}
}

func TestUpdateSkillHash(t *testing.T) {
	s := testStore(t)
	s.InsertSkill(&skill.Skill{ID: "hash-1", Name: "hashed", SourceType: "local", CentralPath: "/h", ContentHash: "old", Enabled: true})

	if err := s.UpdateSkillHash("hash-1", "newhash"); err != nil {
		t.Fatalf("UpdateSkillHash() error: %v", err)
	}

	got, err := s.GetSkill("hashed")
	if err != nil {
		t.Fatal(err)
	}
	if got.ContentHash != "newhash" {
		t.Errorf("ContentHash = %q, want %q", got.ContentHash, "newhash")
	}
}

func TestUpsertAndListTargets(t *testing.T) {
	s := testStore(t)
	s.InsertSkill(&skill.Skill{ID: "sk1", Name: "targeted", SourceType: "local", CentralPath: "/t", Enabled: true})

	if err := s.UpsertTarget("sk1", "claude", "/home/.claude/skills/targeted", "symlink", "hash1"); err != nil {
		t.Fatalf("UpsertTarget() error: %v", err)
	}

	targets, err := s.ListTargets("sk1")
	if err != nil {
		t.Fatalf("ListTargets() error: %v", err)
	}
	if len(targets) != 1 {
		t.Fatalf("expected 1 target, got %d", len(targets))
	}
	if targets[0].Agent != "claude" {
		t.Errorf("Agent = %q", targets[0].Agent)
	}
	if targets[0].TargetPath != "/home/.claude/skills/targeted" {
		t.Errorf("TargetPath = %q", targets[0].TargetPath)
	}
	if targets[0].Mode != "symlink" {
		t.Errorf("Mode = %q", targets[0].Mode)
	}
}

func TestUpsertTarget_Update(t *testing.T) {
	s := testStore(t)
	s.InsertSkill(&skill.Skill{ID: "sk2", Name: "upsert-target", SourceType: "local", CentralPath: "/u", Enabled: true})

	s.UpsertTarget("sk2", "claude", "/path1", "symlink", "h1")
	s.UpsertTarget("sk2", "claude", "/path2", "copy", "h2")

	targets, err := s.ListTargets("sk2")
	if err != nil {
		t.Fatal(err)
	}
	if len(targets) != 1 {
		t.Fatalf("expected 1 target after upsert, got %d", len(targets))
	}
	if targets[0].TargetPath != "/path2" {
		t.Errorf("TargetPath = %q, want /path2", targets[0].TargetPath)
	}
	if targets[0].Mode != "copy" {
		t.Errorf("Mode = %q, want copy", targets[0].Mode)
	}
}

func TestDeleteTarget(t *testing.T) {
	s := testStore(t)
	s.InsertSkill(&skill.Skill{ID: "sk3", Name: "del-target", SourceType: "local", CentralPath: "/dt", Enabled: true})
	s.UpsertTarget("sk3", "claude", "/some/path", "symlink", "h1")

	if err := s.DeleteTarget("sk3", "claude"); err != nil {
		t.Fatalf("DeleteTarget() error: %v", err)
	}

	targets, err := s.ListTargets("sk3")
	if err != nil {
		t.Fatal(err)
	}
	if len(targets) != 0 {
		t.Fatalf("expected 0 targets after delete, got %d", len(targets))
	}
}

func TestSeedBuiltinAgents(t *testing.T) {
	s := testStore(t)

	agents := []agent.Adapter{
		{Name: "claude", DisplayName: "Claude Code", ProjectDir: ".claude/skills", GlobalDir: ".claude/skills", DetectPaths: []string{".claude"}, IsBuiltin: true},
		{Name: "cursor", DisplayName: "Cursor", ProjectDir: ".cursor/skills", GlobalDir: ".cursor/skills", DetectPaths: []string{".cursor"}, IsBuiltin: true},
	}

	if err := s.SeedBuiltinAgents(agents); err != nil {
		t.Fatalf("SeedBuiltinAgents() error: %v", err)
	}

	listed, err := s.ListAgents()
	if err != nil {
		t.Fatalf("ListAgents() error: %v", err)
	}
	if len(listed) != 2 {
		t.Fatalf("expected 2 agents, got %d", len(listed))
	}
}

func TestSeedBuiltinAgents_Idempotent(t *testing.T) {
	s := testStore(t)

	agents := []agent.Adapter{
		{Name: "claude", DisplayName: "Claude Code", ProjectDir: ".claude/skills", GlobalDir: ".claude/skills", DetectPaths: []string{".claude"}, IsBuiltin: true},
	}

	s.SeedBuiltinAgents(agents)
	s.SeedBuiltinAgents(agents) // second call should not error

	listed, err := s.ListAgents()
	if err != nil {
		t.Fatal(err)
	}
	if len(listed) != 1 {
		t.Fatalf("expected 1 agent after double seed, got %d", len(listed))
	}
}

func TestInsertAndDeleteAgent(t *testing.T) {
	s := testStore(t)

	a := agent.Adapter{
		Name:        "custom",
		DisplayName: "Custom Agent",
		ProjectDir:  ".custom/skills",
		GlobalDir:   ".custom/skills",
		DetectPaths: []string{".custom"},
		IsBuiltin:   false,
	}

	if err := s.InsertAgent(a); err != nil {
		t.Fatalf("InsertAgent() error: %v", err)
	}

	listed, err := s.ListAgents()
	if err != nil {
		t.Fatal(err)
	}
	if len(listed) != 1 {
		t.Fatalf("expected 1 agent, got %d", len(listed))
	}
	if listed[0].Name != "custom" {
		t.Errorf("Name = %q, want custom", listed[0].Name)
	}

	if err := s.DeleteAgent("custom"); err != nil {
		t.Fatalf("DeleteAgent() error: %v", err)
	}

	listed, _ = s.ListAgents()
	if len(listed) != 0 {
		t.Fatalf("expected 0 agents after delete, got %d", len(listed))
	}
}

func TestDeleteAgent_BuiltinProtected(t *testing.T) {
	s := testStore(t)

	agents := []agent.Adapter{
		{Name: "claude", DisplayName: "Claude Code", ProjectDir: ".claude/skills", GlobalDir: ".claude/skills", DetectPaths: []string{".claude"}, IsBuiltin: true},
	}
	s.SeedBuiltinAgents(agents)

	// Delete should not remove builtin agents
	s.DeleteAgent("claude")

	listed, _ := s.ListAgents()
	if len(listed) != 1 {
		t.Fatalf("expected builtin agent to survive delete, got %d", len(listed))
	}
}

func TestListTargets_Empty(t *testing.T) {
	s := testStore(t)
	targets, err := s.ListTargets("nonexistent")
	if err != nil {
		t.Fatalf("ListTargets() error: %v", err)
	}
	if len(targets) != 0 {
		t.Fatalf("expected 0 targets, got %d", len(targets))
	}
}

func TestDeleteSkill_CascadesTargets(t *testing.T) {
	s := testStore(t)
	s.InsertSkill(&skill.Skill{ID: "cascade-1", Name: "cascade", SourceType: "local", CentralPath: "/c", Enabled: true})
	s.UpsertTarget("cascade-1", "claude", "/path", "symlink", "h1")

	s.DeleteSkill("cascade-1")

	targets, err := s.ListTargets("cascade-1")
	if err != nil {
		t.Fatal(err)
	}
	if len(targets) != 0 {
		t.Fatalf("expected targets to be cascaded deleted, got %d", len(targets))
	}
}
