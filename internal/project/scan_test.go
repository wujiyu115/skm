package project

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ejoy/skm/internal/agent"
)

const testSkillMD = "---\nname: test-skill\ndescription: A test\n---\n# Test"

func writeSkillMD(t *testing.T, dir, name, content string) {
	t.Helper()
	skillDir := filepath.Join(dir, name)
	if err := os.MkdirAll(skillDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(skillDir, "SKILL.md"), []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func testAgents() []agent.Adapter {
	return []agent.Adapter{
		{
			Name:        "claude",
			DisplayName: "Claude Code",
			ProjectDir:  ".claude/skills",
		},
	}
}

func TestScanSkills_FindsEnabledAndDisabled(t *testing.T) {
	tmp := t.TempDir()

	writeSkillMD(t, filepath.Join(tmp, ".claude/skills"), "s1",
		"---\nname: s1\ndescription: Skill one\n---\n# S1")
	writeSkillMD(t, filepath.Join(tmp, ".claude/skills-disabled"), "s2",
		"---\nname: s2\ndescription: Skill two\n---\n# S2")

	skills, err := ScanSkills(tmp, testAgents())
	if err != nil {
		t.Fatal(err)
	}

	if len(skills) != 2 {
		t.Fatalf("expected 2 skills, got %d", len(skills))
	}

	// Sorted by name: s1 then s2.
	if skills[0].SkillName != "s1" || skills[0].Enabled != true {
		t.Errorf("expected s1 enabled, got %+v", skills[0])
	}
	if skills[1].SkillName != "s2" || skills[1].Enabled != false {
		t.Errorf("expected s2 disabled, got %+v", skills[1])
	}
}

func TestScanSkills_EmptyProject(t *testing.T) {
	tmp := t.TempDir()

	skills, err := ScanSkills(tmp, testAgents())
	if err != nil {
		t.Fatal(err)
	}

	if len(skills) != 0 {
		t.Fatalf("expected 0 skills, got %d", len(skills))
	}
}

func TestScanSkills_SkipsNonSkillDirs(t *testing.T) {
	tmp := t.TempDir()

	// Create a subdir without SKILL.md.
	noSkillDir := filepath.Join(tmp, ".claude/skills/not-a-skill")
	if err := os.MkdirAll(noSkillDir, 0o755); err != nil {
		t.Fatal(err)
	}
	// Create a regular file in there.
	if err := os.WriteFile(filepath.Join(noSkillDir, "README.md"), []byte("hello"), 0o644); err != nil {
		t.Fatal(err)
	}

	// Also add one real skill.
	writeSkillMD(t, filepath.Join(tmp, ".claude/skills"), "real-skill", testSkillMD)

	skills, err := ScanSkills(tmp, testAgents())
	if err != nil {
		t.Fatal(err)
	}

	if len(skills) != 1 {
		t.Fatalf("expected 1 skill, got %d", len(skills))
	}
	if skills[0].SkillName != "test-skill" {
		t.Errorf("expected test-skill, got %s", skills[0].SkillName)
	}
}

func TestDisabledDir(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{".claude/skills", ".claude/skills-disabled"},
		{".cursor/skills", ".cursor/skills-disabled"},
		{"skills", "skills-disabled"},
	}

	for _, tt := range tests {
		got := disabledDir(tt.input)
		if got != tt.want {
			t.Errorf("disabledDir(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestToggleSkill(t *testing.T) {
	tmp := t.TempDir()

	writeSkillMD(t, filepath.Join(tmp, ".claude/skills"), "my-skill", testSkillMD)

	// Toggle to disabled.
	err := ToggleSkill(tmp, ".claude/skills", "my-skill", false)
	if err != nil {
		t.Fatal(err)
	}

	// Verify it moved.
	disabledPath := filepath.Join(tmp, ".claude/skills-disabled/my-skill/SKILL.md")
	if _, err := os.Stat(disabledPath); os.IsNotExist(err) {
		t.Error("skill was not moved to disabled dir")
	}

	enabledPath := filepath.Join(tmp, ".claude/skills/my-skill")
	if _, err := os.Stat(enabledPath); !os.IsNotExist(err) {
		t.Error("skill still exists in enabled dir")
	}
}

func TestToggleSkill_SourceNotFound(t *testing.T) {
	tmp := t.TempDir()

	err := ToggleSkill(tmp, ".claude/skills", "nonexistent", false)
	if err == nil {
		t.Fatal("expected error for missing source")
	}
}

func TestRemoveSkill(t *testing.T) {
	tmp := t.TempDir()

	skillDir := filepath.Join(tmp, "my-skill")
	writeSkillMD(t, tmp, "my-skill", testSkillMD)

	err := RemoveSkill(skillDir)
	if err != nil {
		t.Fatal(err)
	}

	if _, err := os.Stat(skillDir); !os.IsNotExist(err) {
		t.Error("skill dir still exists after removal")
	}
}

func TestAddFromCentral(t *testing.T) {
	tmp := t.TempDir()

	// Create a central skill with multiple files.
	centralDir := filepath.Join(tmp, "central/my-skill")
	if err := os.MkdirAll(centralDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(centralDir, "SKILL.md"), []byte(testSkillMD), 0o644); err != nil {
		t.Fatal(err)
	}
	// Add a subdirectory with a file.
	subDir := filepath.Join(centralDir, "templates")
	if err := os.MkdirAll(subDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(subDir, "tmpl.txt"), []byte("template content"), 0o644); err != nil {
		t.Fatal(err)
	}

	projectDir := filepath.Join(tmp, "project")
	if err := os.MkdirAll(projectDir, 0o755); err != nil {
		t.Fatal(err)
	}

	err := AddFromCentral(centralDir, projectDir, ".claude/skills", "my-skill")
	if err != nil {
		t.Fatal(err)
	}

	// Verify SKILL.md was copied.
	targetSkill := filepath.Join(projectDir, ".claude/skills/my-skill/SKILL.md")
	data, err := os.ReadFile(targetSkill)
	if err != nil {
		t.Fatalf("SKILL.md not copied: %v", err)
	}
	if string(data) != testSkillMD {
		t.Errorf("SKILL.md content mismatch: got %q", string(data))
	}

	// Verify subdirectory file was copied.
	tmplFile := filepath.Join(projectDir, ".claude/skills/my-skill/templates/tmpl.txt")
	data, err = os.ReadFile(tmplFile)
	if err != nil {
		t.Fatalf("template file not copied: %v", err)
	}
	if string(data) != "template content" {
		t.Errorf("template content mismatch: got %q", string(data))
	}
}

func TestAddFromCentral_AlreadyExists(t *testing.T) {
	tmp := t.TempDir()

	centralDir := filepath.Join(tmp, "central/my-skill")
	writeSkillMD(t, filepath.Join(tmp, "central"), "my-skill", testSkillMD)

	projectDir := filepath.Join(tmp, "project")
	// Pre-create the target.
	target := filepath.Join(projectDir, ".claude/skills/my-skill")
	if err := os.MkdirAll(target, 0o755); err != nil {
		t.Fatal(err)
	}

	err := AddFromCentral(centralDir, projectDir, ".claude/skills", "my-skill")
	if err == nil {
		t.Fatal("expected error when target already exists")
	}
}
