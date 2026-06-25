package agent

import (
	"os"
	"path/filepath"
	"testing"
)

func TestBuiltin_Returns52Agents(t *testing.T) {
	agents := Builtin()
	if len(agents) < 50 {
		t.Fatalf("expected 50+ built-in agents, got %d", len(agents))
	}
	if len(agents) != 52 {
		t.Fatalf("expected exactly 52 built-in agents, got %d", len(agents))
	}
}

func TestBuiltin_NoDuplicateNames(t *testing.T) {
	agents := Builtin()
	seen := map[string]bool{}
	for _, a := range agents {
		if seen[a.Name] {
			t.Errorf("duplicate agent name: %s", a.Name)
		}
		seen[a.Name] = true
	}
}

func TestBuiltin_AllFieldsNonEmpty(t *testing.T) {
	for _, a := range Builtin() {
		if a.Name == "" {
			t.Error("agent has empty Name")
		}
		if a.DisplayName == "" {
			t.Errorf("agent %s has empty DisplayName", a.Name)
		}
		if a.ProjectDir == "" {
			t.Errorf("agent %s has empty ProjectDir", a.Name)
		}
		if a.GlobalDir == "" {
			t.Errorf("agent %s has empty GlobalDir", a.Name)
		}
		if len(a.DetectPaths) == 0 {
			t.Errorf("agent %s has empty DetectPaths", a.Name)
		}
		if !a.IsBuiltin {
			t.Errorf("agent %s has IsBuiltin=false", a.Name)
		}
	}
}

func TestBuiltin_CategoryIsValid(t *testing.T) {
	for _, a := range Builtin() {
		if a.Category != "coding" && a.Category != "lobster" {
			t.Errorf("agent %s has invalid Category=%q (want coding or lobster)", a.Name, a.Category)
		}
	}
}

func TestBuiltin_CategoryCounts(t *testing.T) {
	coding, lobster := 0, 0
	for _, a := range Builtin() {
		switch a.Category {
		case "coding":
			coding++
		case "lobster":
			lobster++
		}
	}
	if coding != 46 {
		t.Errorf("expected 46 coding agents, got %d", coding)
	}
	if lobster != 6 {
		t.Errorf("expected 6 lobster agents, got %d", lobster)
	}
}

func TestFind_KeyAgents(t *testing.T) {
	tests := []struct {
		name        string
		displayName string
		category    string
	}{
		{"claude", "Claude Code", "coding"},
		{"cursor", "Cursor", "coding"},
		{"codex", "Codex", "coding"},
		{"windsurf", "Windsurf", "coding"},
		{"gemini_cli", "Gemini CLI", "coding"},
		{"github_copilot", "GitHub Copilot", "coding"},
		{"amp", "Amp", "coding"},
		{"roo_code", "Roo Code", "coding"},
		{"kilo_code", "Kilo Code", "coding"},
		{"cline", "Cline", "coding"},
		{"trae", "Trae", "coding"},
		{"continue_dev", "Continue", "coding"},
		{"aider", "Aider", "coding"},
		{"kiro", "Kiro", "coding"},
		{"openclaw", "OpenClaw", "lobster"},
		{"hermes", "Hermes", "lobster"},
		{"workbuddy", "WorkBuddy", "lobster"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a, ok := Find(tt.name)
			if !ok {
				t.Fatalf("Find(%s) returned false", tt.name)
			}
			if a.DisplayName != tt.displayName {
				t.Errorf("DisplayName = %q, want %q", a.DisplayName, tt.displayName)
			}
			if a.Category != tt.category {
				t.Errorf("Category = %q, want %q", a.Category, tt.category)
			}
		})
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
