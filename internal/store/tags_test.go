package store

import (
	"path/filepath"
	"testing"

	"github.com/ejoy/skm/internal/skill"
)

func newTestStoreWithSkill(t *testing.T) (*Store, *skill.Skill) {
	t.Helper()
	dir := t.TempDir()
	s, err := New(filepath.Join(dir, "test.db"))
	if err != nil {
		t.Fatalf("New() error: %v", err)
	}
	sk := &skill.Skill{
		ID:          "sk-1",
		Name:        "test-skill",
		Description: "A test skill",
		SourceType:  "local",
		SourceRef:   "/tmp/skill",
		CentralPath: "/central/test-skill",
		ContentHash: "abc123",
		Enabled:     true,
	}
	if err := s.InsertSkill(sk); err != nil {
		t.Fatalf("InsertSkill() error: %v", err)
	}
	return s, sk
}

func TestAddTag(t *testing.T) {
	s, sk := newTestStoreWithSkill(t)
	defer s.Close()

	if err := s.AddTag(sk.ID, "networking"); err != nil {
		t.Fatalf("AddTag() error: %v", err)
	}

	tags, err := s.ListSkillTags(sk.ID)
	if err != nil {
		t.Fatalf("ListSkillTags() error: %v", err)
	}
	if len(tags) != 1 || tags[0] != "networking" {
		t.Fatalf("expected [networking], got %v", tags)
	}
}

func TestAddTag_Duplicate(t *testing.T) {
	s, sk := newTestStoreWithSkill(t)
	defer s.Close()

	if err := s.AddTag(sk.ID, "networking"); err != nil {
		t.Fatalf("AddTag() error: %v", err)
	}
	// Adding same tag again should not error
	if err := s.AddTag(sk.ID, "networking"); err != nil {
		t.Fatalf("AddTag() duplicate error: %v", err)
	}

	tags, _ := s.ListSkillTags(sk.ID)
	if len(tags) != 1 {
		t.Fatalf("expected 1 tag after duplicate add, got %d", len(tags))
	}
}

func TestRemoveTag(t *testing.T) {
	s, sk := newTestStoreWithSkill(t)
	defer s.Close()

	s.AddTag(sk.ID, "networking")
	s.AddTag(sk.ID, "security")

	if err := s.RemoveTag(sk.ID, "networking"); err != nil {
		t.Fatalf("RemoveTag() error: %v", err)
	}

	tags, _ := s.ListSkillTags(sk.ID)
	if len(tags) != 1 || tags[0] != "security" {
		t.Fatalf("expected [security], got %v", tags)
	}
}

func TestSetSkillTags(t *testing.T) {
	s, sk := newTestStoreWithSkill(t)
	defer s.Close()

	s.AddTag(sk.ID, "old-tag")

	if err := s.SetSkillTags(sk.ID, []string{"alpha", "beta", "gamma"}); err != nil {
		t.Fatalf("SetSkillTags() error: %v", err)
	}

	tags, _ := s.ListSkillTags(sk.ID)
	if len(tags) != 3 {
		t.Fatalf("expected 3 tags, got %d: %v", len(tags), tags)
	}
	expected := []string{"alpha", "beta", "gamma"}
	for i, want := range expected {
		if tags[i] != want {
			t.Fatalf("tag[%d] = %q, want %q", i, tags[i], want)
		}
	}
}

func TestSetSkillTags_Empty(t *testing.T) {
	s, sk := newTestStoreWithSkill(t)
	defer s.Close()

	s.AddTag(sk.ID, "old-tag")

	if err := s.SetSkillTags(sk.ID, []string{}); err != nil {
		t.Fatalf("SetSkillTags() error: %v", err)
	}

	tags, _ := s.ListSkillTags(sk.ID)
	if len(tags) != 0 {
		t.Fatalf("expected 0 tags, got %d", len(tags))
	}
}

func TestListTags(t *testing.T) {
	s, sk := newTestStoreWithSkill(t)
	defer s.Close()

	// Insert a second skill
	sk2 := &skill.Skill{
		ID: "sk-2", Name: "skill-2", SourceType: "local",
		CentralPath: "/central/skill-2", Enabled: true,
	}
	s.InsertSkill(sk2)

	s.AddTag(sk.ID, "networking")
	s.AddTag(sk.ID, "security")
	s.AddTag(sk2.ID, "security")
	s.AddTag(sk2.ID, "automation")

	tags, err := s.ListTags()
	if err != nil {
		t.Fatalf("ListTags() error: %v", err)
	}
	// Should return distinct tags sorted: automation, networking, security
	if len(tags) != 3 {
		t.Fatalf("expected 3 distinct tags, got %d: %v", len(tags), tags)
	}
	expected := []string{"automation", "networking", "security"}
	for i, want := range expected {
		if tags[i] != want {
			t.Fatalf("tag[%d] = %q, want %q", i, tags[i], want)
		}
	}
}

func TestListSkillTags_Empty(t *testing.T) {
	s, sk := newTestStoreWithSkill(t)
	defer s.Close()

	tags, err := s.ListSkillTags(sk.ID)
	if err != nil {
		t.Fatalf("ListSkillTags() error: %v", err)
	}
	if len(tags) != 0 {
		t.Fatalf("expected 0 tags for new skill, got %d", len(tags))
	}
}

func TestRenameTag(t *testing.T) {
	s, sk := newTestStoreWithSkill(t)
	defer s.Close()

	sk2 := &skill.Skill{
		ID: "sk-2", Name: "skill-2", SourceType: "local",
		CentralPath: "/central/skill-2", Enabled: true,
	}
	s.InsertSkill(sk2)

	s.AddTag(sk.ID, "old-name")
	s.AddTag(sk2.ID, "old-name")

	if err := s.RenameTag("old-name", "new-name"); err != nil {
		t.Fatalf("RenameTag() error: %v", err)
	}

	// Both skills should now have "new-name"
	tags1, _ := s.ListSkillTags(sk.ID)
	tags2, _ := s.ListSkillTags(sk2.ID)

	if len(tags1) != 1 || tags1[0] != "new-name" {
		t.Fatalf("skill 1 tags: expected [new-name], got %v", tags1)
	}
	if len(tags2) != 1 || tags2[0] != "new-name" {
		t.Fatalf("skill 2 tags: expected [new-name], got %v", tags2)
	}
}

func TestRenameTag_ConflictBothTags(t *testing.T) {
	s, sk := newTestStoreWithSkill(t)
	defer s.Close()

	// Skill has both "tagA" and "tagB".
	s.AddTag(sk.ID, "tagA")
	s.AddTag(sk.ID, "tagB")

	// Renaming "tagA" -> "tagB" should succeed (not fail with UNIQUE constraint).
	if err := s.RenameTag("tagA", "tagB"); err != nil {
		t.Fatalf("RenameTag() with conflict error: %v", err)
	}

	tags, err := s.ListSkillTags(sk.ID)
	if err != nil {
		t.Fatalf("ListSkillTags() error: %v", err)
	}
	if len(tags) != 1 || tags[0] != "tagB" {
		t.Fatalf("expected [tagB], got %v", tags)
	}
}

func TestDeleteTag(t *testing.T) {
	s, sk := newTestStoreWithSkill(t)
	defer s.Close()

	sk2 := &skill.Skill{
		ID: "sk-2", Name: "skill-2", SourceType: "local",
		CentralPath: "/central/skill-2", Enabled: true,
	}
	s.InsertSkill(sk2)

	s.AddTag(sk.ID, "to-delete")
	s.AddTag(sk.ID, "keep")
	s.AddTag(sk2.ID, "to-delete")

	if err := s.DeleteTag("to-delete"); err != nil {
		t.Fatalf("DeleteTag() error: %v", err)
	}

	tags1, _ := s.ListSkillTags(sk.ID)
	tags2, _ := s.ListSkillTags(sk2.ID)

	if len(tags1) != 1 || tags1[0] != "keep" {
		t.Fatalf("skill 1 tags: expected [keep], got %v", tags1)
	}
	if len(tags2) != 0 {
		t.Fatalf("skill 2 tags: expected [], got %v", tags2)
	}
}

func TestTagsCascadeOnSkillDelete(t *testing.T) {
	s, sk := newTestStoreWithSkill(t)
	defer s.Close()

	s.AddTag(sk.ID, "tag-1")
	s.AddTag(sk.ID, "tag-2")

	if err := s.DeleteSkill(sk.ID); err != nil {
		t.Fatalf("DeleteSkill() error: %v", err)
	}

	// Tags should be cascade-deleted
	tags, _ := s.ListSkillTags(sk.ID)
	if len(tags) != 0 {
		t.Fatalf("expected tags to be cascade-deleted, got %v", tags)
	}
}
