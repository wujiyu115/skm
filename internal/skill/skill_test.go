package skill

import (
	"os"
	"path/filepath"
	"testing"
)

func writeSkillMD(t *testing.T, dir, content string) {
	t.Helper()
	if err := os.MkdirAll(dir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "SKILL.md"), []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
}

func TestParseSkillMD_FullFrontmatter(t *testing.T) {
	dir := t.TempDir()
	writeSkillMD(t, dir, `---
name: test-skill
description: A test skill
metadata:
  type: coding
  tags: [go, testing]
---

# Test Skill Content
`)

	s, err := ParseSkillMD(dir)
	if err != nil {
		t.Fatalf("ParseSkillMD() error: %v", err)
	}
	if s.Name != "test-skill" {
		t.Errorf("Name = %q, want %q", s.Name, "test-skill")
	}
	if s.Description != "A test skill" {
		t.Errorf("Description = %q, want %q", s.Description, "A test skill")
	}
}

func TestParseSkillMD_NoFrontmatter(t *testing.T) {
	dir := t.TempDir()
	writeSkillMD(t, dir, "# Just markdown\n\nSome content here.")

	s, err := ParseSkillMD(dir)
	if err != nil {
		t.Fatalf("ParseSkillMD() error: %v", err)
	}
	if s.Name != filepath.Base(dir) {
		t.Errorf("Name = %q, want dir basename %q", s.Name, filepath.Base(dir))
	}
}

func TestParseSkillMD_NameFallbackToDir(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "my-cool-skill")
	writeSkillMD(t, dir, `---
description: No name field
---
Content`)

	s, err := ParseSkillMD(dir)
	if err != nil {
		t.Fatalf("ParseSkillMD() error: %v", err)
	}
	if s.Name != "my-cool-skill" {
		t.Errorf("Name = %q, want %q", s.Name, "my-cool-skill")
	}
}

func TestParseSkillMD_NoFile(t *testing.T) {
	dir := t.TempDir()
	_, err := ParseSkillMD(dir)
	if err == nil {
		t.Fatal("expected error for missing SKILL.md")
	}
}

func TestParseSkillMD_CaseInsensitive(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "skill.md"), []byte("---\nname: lower\n---\n"), 0644); err != nil {
		t.Fatal(err)
	}

	s, err := ParseSkillMD(dir)
	if err != nil {
		t.Fatalf("ParseSkillMD() error: %v", err)
	}
	if s.Name != "lower" {
		t.Errorf("Name = %q, want %q", s.Name, "lower")
	}
}
