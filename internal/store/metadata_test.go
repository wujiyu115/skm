package store

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/ejoy/skm/internal/skill"
)

func TestWriteMetadata_CreatesFiles(t *testing.T) {
	s := testStore(t)
	s.InsertSkill(&skill.Skill{ID: "m1", Name: "meta-test", Description: "test", SourceType: "git", SourceRef: "https://github.com/x/y", CentralPath: "/m", ContentHash: "abc", Enabled: true})
	s.InsertGroup("g1", "test-group", "desc")

	metaDir := filepath.Join(t.TempDir(), ".metadata")
	if err := s.WriteMetadata(metaDir); err != nil {
		t.Fatalf("WriteMetadata() error: %v", err)
	}

	skillFile := filepath.Join(metaDir, "skills", "meta-test.json")
	if _, err := os.Stat(skillFile); err != nil {
		t.Fatalf("skill metadata not written: %v", err)
	}

	var data map[string]interface{}
	raw, _ := os.ReadFile(skillFile)
	json.Unmarshal(raw, &data)
	if data["name"] != "meta-test" {
		t.Errorf("name = %v", data["name"])
	}

	groupFile := filepath.Join(metaDir, "groups", "test-group.json")
	if _, err := os.Stat(groupFile); err != nil {
		t.Fatalf("group metadata not written: %v", err)
	}
}

func TestWriteMetadata_IncludesGroupSkills(t *testing.T) {
	s := testStore(t)
	s.InsertSkill(&skill.Skill{ID: "s1", Name: "skill-a", Description: "A", SourceType: "local", CentralPath: "/a", Enabled: true})
	s.InsertSkill(&skill.Skill{ID: "s2", Name: "skill-b", Description: "B", SourceType: "local", CentralPath: "/b", Enabled: true})
	s.InsertGroup("g1", "my-group", "a group")
	s.AddSkillToGroup("g1", "s1", 0)
	s.AddSkillToGroup("g1", "s2", 1)

	metaDir := filepath.Join(t.TempDir(), ".metadata")
	if err := s.WriteMetadata(metaDir); err != nil {
		t.Fatalf("WriteMetadata() error: %v", err)
	}

	groupFile := filepath.Join(metaDir, "groups", "my-group.json")
	raw, err := os.ReadFile(groupFile)
	if err != nil {
		t.Fatalf("read group file: %v", err)
	}

	var data map[string]interface{}
	json.Unmarshal(raw, &data)
	skills, ok := data["skills"].([]interface{})
	if !ok {
		t.Fatalf("skills field not an array: %T", data["skills"])
	}
	if len(skills) != 2 {
		t.Errorf("expected 2 skills in group, got %d", len(skills))
	}
}

func TestWriteMetadata_EmptyStore(t *testing.T) {
	s := testStore(t)

	metaDir := filepath.Join(t.TempDir(), ".metadata")
	if err := s.WriteMetadata(metaDir); err != nil {
		t.Fatalf("WriteMetadata() error: %v", err)
	}

	// Directories should still be created
	if _, err := os.Stat(filepath.Join(metaDir, "skills")); err != nil {
		t.Fatalf("skills dir not created: %v", err)
	}
	if _, err := os.Stat(filepath.Join(metaDir, "groups")); err != nil {
		t.Fatalf("groups dir not created: %v", err)
	}
}
