package agent

import (
	"os"
	"path/filepath"
	"testing"
)

func TestBuiltin_ReturnsThreeAgents(t *testing.T) {
	agents := Builtin()
	if len(agents) != 3 {
		t.Fatalf("expected 3 built-in agents, got %d", len(agents))
	}

	names := map[string]bool{}
	for _, a := range agents {
		names[a.Name] = true
	}
	for _, want := range []string{"claude", "cursor", "codex"} {
		if !names[want] {
			t.Errorf("missing built-in agent: %s", want)
		}
	}
}

func TestFind_ExistingAgent(t *testing.T) {
	a, ok := Find("claude")
	if !ok {
		t.Fatal("Find(claude) returned false")
	}
	if a.DisplayName != "Claude Code" {
		t.Errorf("DisplayName = %q, want %q", a.DisplayName, "Claude Code")
	}
	if a.ProjectDir != ".claude/skills" {
		t.Errorf("ProjectDir = %q, want %q", a.ProjectDir, ".claude/skills")
	}
}

func TestFind_NonExistentAgent(t *testing.T) {
	_, ok := Find("nonexistent")
	if ok {
		t.Fatal("Find(nonexistent) should return false")
	}
}

func TestInstallPath_Global(t *testing.T) {
	a, _ := Find("claude")
	p, err := InstallPath(a, true, "/some/project")
	if err != nil {
		t.Fatalf("InstallPath() error: %v", err)
	}
	home, _ := os.UserHomeDir()
	want := filepath.Join(home, ".claude", "skills")
	if p != want {
		t.Errorf("global path = %q, want %q", p, want)
	}
}

func TestInstallPath_Project(t *testing.T) {
	a, _ := Find("claude")
	p, err := InstallPath(a, false, "/some/project")
	if err != nil {
		t.Fatalf("InstallPath() error: %v", err)
	}
	want := filepath.Join("/some/project", ".claude", "skills")
	if p != want {
		t.Errorf("project path = %q, want %q", p, want)
	}
}

func TestDetect_FindsInstalledAgents(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	os.MkdirAll(filepath.Join(home, ".claude"), 0755)

	detected := Detect()
	found := false
	for _, a := range detected {
		if a.Name == "claude" {
			found = true
		}
	}
	if !found {
		t.Error("expected claude to be detected when ~/.claude exists")
	}
}
