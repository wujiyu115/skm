package skill

import (
	"os"
	"path/filepath"
	"testing"
)

func makeSkill(t *testing.T, base, name string) {
	t.Helper()
	dir := filepath.Join(base, name)
	os.MkdirAll(dir, 0755)
	os.WriteFile(filepath.Join(dir, "SKILL.md"), []byte("---\nname: "+name+"\n---\n# "+name), 0644)
}

func TestDiscover_DirectSkill(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "SKILL.md"), []byte("---\nname: direct\n---\n"), 0644)

	skills, err := Discover(dir, "")
	if err != nil {
		t.Fatalf("Discover() error: %v", err)
	}
	if len(skills) != 1 {
		t.Fatalf("expected 1 skill, got %d", len(skills))
	}
	if skills[0].Name != "direct" {
		t.Errorf("Name = %q, want direct", skills[0].Name)
	}
}

func TestDiscover_SkillsSubdir(t *testing.T) {
	dir := t.TempDir()
	skillsDir := filepath.Join(dir, "skills")
	makeSkill(t, skillsDir, "skill-a")
	makeSkill(t, skillsDir, "skill-b")

	skills, err := Discover(dir, "")
	if err != nil {
		t.Fatalf("Discover() error: %v", err)
	}
	if len(skills) != 2 {
		t.Fatalf("expected 2 skills, got %d", len(skills))
	}
}

func TestDiscover_WithSubpath(t *testing.T) {
	dir := t.TempDir()
	makeSkill(t, filepath.Join(dir, "skills"), "target")
	makeSkill(t, filepath.Join(dir, "skills"), "other")

	skills, err := Discover(dir, "skills/target")
	if err != nil {
		t.Fatalf("Discover() error: %v", err)
	}
	if len(skills) != 1 {
		t.Fatalf("expected 1 skill, got %d", len(skills))
	}
	if skills[0].Name != "target" {
		t.Errorf("Name = %q, want target", skills[0].Name)
	}
}

func TestDiscover_RecursiveSearch(t *testing.T) {
	dir := t.TempDir()
	makeSkill(t, filepath.Join(dir, "deep", "nested"), "hidden-skill")

	skills, err := Discover(dir, "")
	if err != nil {
		t.Fatalf("Discover() error: %v", err)
	}
	if len(skills) != 1 {
		t.Fatalf("expected 1 skill from recursive search, got %d", len(skills))
	}
}

func TestDiscover_DeduplicateByName(t *testing.T) {
	dir := t.TempDir()
	makeSkill(t, filepath.Join(dir, "skills"), "same-name")
	makeSkill(t, filepath.Join(dir, ".claude", "skills"), "same-name")

	skills, err := Discover(dir, "")
	if err != nil {
		t.Fatalf("Discover() error: %v", err)
	}
	if len(skills) != 1 {
		t.Fatalf("expected 1 deduplicated skill, got %d", len(skills))
	}
}
