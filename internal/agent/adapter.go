package agent

import (
	"fmt"
	"os"
	"path/filepath"
)

// Adapter represents an AI coding agent that can consume skills.
type Adapter struct {
	Name        string   `json:"name"`
	DisplayName string   `json:"display_name"`
	ProjectDir  string   `json:"project_dir"`
	GlobalDir   string   `json:"global_dir"`
	DetectPaths []string `json:"detect_paths"`
	IsBuiltin   bool     `json:"is_builtin"`
}

var builtins = []Adapter{
	{
		Name:        "claude",
		DisplayName: "Claude Code",
		ProjectDir:  ".claude/skills",
		GlobalDir:   ".claude/skills",
		DetectPaths: []string{".claude"},
		IsBuiltin:   true,
	},
	{
		Name:        "cursor",
		DisplayName: "Cursor",
		ProjectDir:  ".cursor/skills",
		GlobalDir:   ".cursor/skills",
		DetectPaths: []string{".cursor"},
		IsBuiltin:   true,
	},
	{
		Name:        "codex",
		DisplayName: "Codex",
		ProjectDir:  ".codex/skills",
		GlobalDir:   ".codex/skills",
		DetectPaths: []string{".codex"},
		IsBuiltin:   true,
	},
}

// Builtin returns a copy of all built-in agent adapters.
func Builtin() []Adapter {
	result := make([]Adapter, len(builtins))
	copy(result, builtins)
	return result
}

// Find looks up a built-in agent adapter by name.
func Find(name string) (Adapter, bool) {
	for _, a := range builtins {
		if a.Name == name {
			return a, true
		}
	}
	return Adapter{}, false
}

// InstallPath returns the filesystem path where skills should be installed
// for the given adapter. If global is true, the path is under the user's
// home directory; otherwise it is relative to cwd.
func InstallPath(a Adapter, global bool, cwd string) (string, error) {
	if global {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("get home dir: %w", err)
		}
		return filepath.Join(home, a.GlobalDir), nil
	}
	return filepath.Join(cwd, a.ProjectDir), nil
}
