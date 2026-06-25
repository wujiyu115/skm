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
	Category    string   `json:"category"` // "coding" or "lobster"
}

var builtins = []Adapter{
	// ── Coding agents (46) ──────────────────────────────────────────────
	{
		Name:        "claude",
		DisplayName: "Claude Code",
		ProjectDir:  ".claude/skills",
		GlobalDir:   ".claude/skills",
		DetectPaths: []string{".claude"},
		IsBuiltin:   true,
		Category:    "coding",
	},
	{
		Name:        "cursor",
		DisplayName: "Cursor",
		ProjectDir:  ".cursor/skills",
		GlobalDir:   ".cursor/skills",
		DetectPaths: []string{".cursor"},
		IsBuiltin:   true,
		Category:    "coding",
	},
	{
		Name:        "codex",
		DisplayName: "Codex",
		ProjectDir:  ".codex/skills",
		GlobalDir:   ".codex/skills",
		DetectPaths: []string{".codex"},
		IsBuiltin:   true,
		Category:    "coding",
	},
	{
		Name:        "windsurf",
		DisplayName: "Windsurf",
		ProjectDir:  ".windsurf/skills",
		GlobalDir:   ".windsurf/skills",
		DetectPaths: []string{".windsurf"},
		IsBuiltin:   true,
		Category:    "coding",
	},
	{
		Name:        "cline",
		DisplayName: "Cline",
		ProjectDir:  ".cline/skills",
		GlobalDir:   ".cline/skills",
		DetectPaths: []string{".cline"},
		IsBuiltin:   true,
		Category:    "coding",
	},
	{
		Name:        "github_copilot",
		DisplayName: "GitHub Copilot",
		ProjectDir:  ".github-copilot/skills",
		GlobalDir:   ".github-copilot/skills",
		DetectPaths: []string{".github-copilot"},
		IsBuiltin:   true,
		Category:    "coding",
	},
	{
		Name:        "gemini_cli",
		DisplayName: "Gemini CLI",
		ProjectDir:  ".gemini/skills",
		GlobalDir:   ".gemini/skills",
		DetectPaths: []string{".gemini"},
		IsBuiltin:   true,
		Category:    "coding",
	},
	{
		Name:        "amp",
		DisplayName: "Amp",
		ProjectDir:  ".amp/skills",
		GlobalDir:   ".amp/skills",
		DetectPaths: []string{".amp"},
		IsBuiltin:   true,
		Category:    "coding",
	},
	{
		Name:        "roo_code",
		DisplayName: "Roo Code",
		ProjectDir:  ".roo-code/skills",
		GlobalDir:   ".roo-code/skills",
		DetectPaths: []string{".roo-code"},
		IsBuiltin:   true,
		Category:    "coding",
	},
	{
		Name:        "kilo_code",
		DisplayName: "Kilo Code",
		ProjectDir:  ".kilo-code/skills",
		GlobalDir:   ".kilo-code/skills",
		DetectPaths: []string{".kilo-code"},
		IsBuiltin:   true,
		Category:    "coding",
	},
	{
		Name:        "trae",
		DisplayName: "Trae",
		ProjectDir:  ".trae/skills",
		GlobalDir:   ".trae/skills",
		DetectPaths: []string{".trae"},
		IsBuiltin:   true,
		Category:    "coding",
	},
	{
		Name:        "continue_dev",
		DisplayName: "Continue",
		ProjectDir:  ".continue/skills",
		GlobalDir:   ".continue/skills",
		DetectPaths: []string{".continue"},
		IsBuiltin:   true,
		Category:    "coding",
	},
	{
		Name:        "aider",
		DisplayName: "Aider",
		ProjectDir:  ".aider/skills",
		GlobalDir:   ".aider/skills",
		DetectPaths: []string{".aider"},
		IsBuiltin:   true,
		Category:    "coding",
	},
	{
		Name:        "omp_agent",
		DisplayName: "OMP Agent",
		ProjectDir:  ".omp_agent/skills",
		GlobalDir:   ".omp_agent/skills",
		DetectPaths: []string{".omp_agent"},
		IsBuiltin:   true,
		Category:    "coding",
	},
	{
		Name:        "grok",
		DisplayName: "Grok",
		ProjectDir:  ".grok/skills",
		GlobalDir:   ".grok/skills",
		DetectPaths: []string{".grok"},
		IsBuiltin:   true,
		Category:    "coding",
	},
	{
		Name:        "opencode",
		DisplayName: "OpenCode",
		ProjectDir:  ".opencode/skills",
		GlobalDir:   ".opencode/skills",
		DetectPaths: []string{".opencode"},
		IsBuiltin:   true,
		Category:    "coding",
	},
	{
		Name:        "antigravity",
		DisplayName: "Antigravity",
		ProjectDir:  ".antigravity/skills",
		GlobalDir:   ".antigravity/skills",
		DetectPaths: []string{".antigravity"},
		IsBuiltin:   true,
		Category:    "coding",
	},
	{
		Name:        "goose",
		DisplayName: "Goose",
		ProjectDir:  ".goose/skills",
		GlobalDir:   ".goose/skills",
		DetectPaths: []string{".goose"},
		IsBuiltin:   true,
		Category:    "coding",
	},
	{
		Name:        "droid",
		DisplayName: "Droid",
		ProjectDir:  ".droid/skills",
		GlobalDir:   ".droid/skills",
		DetectPaths: []string{".droid"},
		IsBuiltin:   true,
		Category:    "coding",
	},
	{
		Name:        "deepagents",
		DisplayName: "DeepAgents",
		ProjectDir:  ".deepagents/skills",
		GlobalDir:   ".deepagents/skills",
		DetectPaths: []string{".deepagents"},
		IsBuiltin:   true,
		Category:    "coding",
	},
	{
		Name:        "firebender",
		DisplayName: "Firebender",
		ProjectDir:  ".firebender/skills",
		GlobalDir:   ".firebender/skills",
		DetectPaths: []string{".firebender"},
		IsBuiltin:   true,
		Category:    "coding",
	},
	{
		Name:        "kimi",
		DisplayName: "Kimi",
		ProjectDir:  ".kimi/skills",
		GlobalDir:   ".kimi/skills",
		DetectPaths: []string{".kimi"},
		IsBuiltin:   true,
		Category:    "coding",
	},
	{
		Name:        "replit",
		DisplayName: "Replit",
		ProjectDir:  ".replit/skills",
		GlobalDir:   ".replit/skills",
		DetectPaths: []string{".replit"},
		IsBuiltin:   true,
		Category:    "coding",
	},
	{
		Name:        "warp",
		DisplayName: "Warp",
		ProjectDir:  ".warp/skills",
		GlobalDir:   ".warp/skills",
		DetectPaths: []string{".warp"},
		IsBuiltin:   true,
		Category:    "coding",
	},
	{
		Name:        "augment",
		DisplayName: "Augment",
		ProjectDir:  ".augment/skills",
		GlobalDir:   ".augment/skills",
		DetectPaths: []string{".augment"},
		IsBuiltin:   true,
		Category:    "coding",
	},
	{
		Name:        "bob",
		DisplayName: "Bob",
		ProjectDir:  ".bob/skills",
		GlobalDir:   ".bob/skills",
		DetectPaths: []string{".bob"},
		IsBuiltin:   true,
		Category:    "coding",
	},
	{
		Name:        "codebuddy",
		DisplayName: "CodeBuddy",
		ProjectDir:  ".codebuddy/skills",
		GlobalDir:   ".codebuddy/skills",
		DetectPaths: []string{".codebuddy"},
		IsBuiltin:   true,
		Category:    "coding",
	},
	{
		Name:        "command_code",
		DisplayName: "Command Code",
		ProjectDir:  ".command_code/skills",
		GlobalDir:   ".command_code/skills",
		DetectPaths: []string{".command_code"},
		IsBuiltin:   true,
		Category:    "coding",
	},
	{
		Name:        "cortex",
		DisplayName: "Cortex",
		ProjectDir:  ".cortex/skills",
		GlobalDir:   ".cortex/skills",
		DetectPaths: []string{".cortex"},
		IsBuiltin:   true,
		Category:    "coding",
	},
	{
		Name:        "crush",
		DisplayName: "Crush",
		ProjectDir:  ".crush/skills",
		GlobalDir:   ".crush/skills",
		DetectPaths: []string{".crush"},
		IsBuiltin:   true,
		Category:    "coding",
	},
	{
		Name:        "iflow",
		DisplayName: "iFlow",
		ProjectDir:  ".iflow/skills",
		GlobalDir:   ".iflow/skills",
		DetectPaths: []string{".iflow"},
		IsBuiltin:   true,
		Category:    "coding",
	},
	{
		Name:        "junie",
		DisplayName: "Junie",
		ProjectDir:  ".junie/skills",
		GlobalDir:   ".junie/skills",
		DetectPaths: []string{".junie"},
		IsBuiltin:   true,
		Category:    "coding",
	},
	{
		Name:        "kiro",
		DisplayName: "Kiro",
		ProjectDir:  ".kiro/skills",
		GlobalDir:   ".kiro/skills",
		DetectPaths: []string{".kiro"},
		IsBuiltin:   true,
		Category:    "coding",
	},
	{
		Name:        "kode",
		DisplayName: "Kode",
		ProjectDir:  ".kode/skills",
		GlobalDir:   ".kode/skills",
		DetectPaths: []string{".kode"},
		IsBuiltin:   true,
		Category:    "coding",
	},
	{
		Name:        "mcpjam",
		DisplayName: "MCPJam",
		ProjectDir:  ".mcpjam/skills",
		GlobalDir:   ".mcpjam/skills",
		DetectPaths: []string{".mcpjam"},
		IsBuiltin:   true,
		Category:    "coding",
	},
	{
		Name:        "mistral_vibe",
		DisplayName: "Mistral Vibe",
		ProjectDir:  ".mistral_vibe/skills",
		GlobalDir:   ".mistral_vibe/skills",
		DetectPaths: []string{".mistral_vibe"},
		IsBuiltin:   true,
		Category:    "coding",
	},
	{
		Name:        "mux",
		DisplayName: "Mux",
		ProjectDir:  ".mux/skills",
		GlobalDir:   ".mux/skills",
		DetectPaths: []string{".mux"},
		IsBuiltin:   true,
		Category:    "coding",
	},
	{
		Name:        "neovate",
		DisplayName: "Neovate",
		ProjectDir:  ".neovate/skills",
		GlobalDir:   ".neovate/skills",
		DetectPaths: []string{".neovate"},
		IsBuiltin:   true,
		Category:    "coding",
	},
	{
		Name:        "openhands",
		DisplayName: "OpenHands",
		ProjectDir:  ".openhands/skills",
		GlobalDir:   ".openhands/skills",
		DetectPaths: []string{".openhands"},
		IsBuiltin:   true,
		Category:    "coding",
	},
	{
		Name:        "pi",
		DisplayName: "Pi",
		ProjectDir:  ".pi/skills",
		GlobalDir:   ".pi/skills",
		DetectPaths: []string{".pi"},
		IsBuiltin:   true,
		Category:    "coding",
	},
	{
		Name:        "pochi",
		DisplayName: "Pochi",
		ProjectDir:  ".pochi/skills",
		GlobalDir:   ".pochi/skills",
		DetectPaths: []string{".pochi"},
		IsBuiltin:   true,
		Category:    "coding",
	},
	{
		Name:        "qoder",
		DisplayName: "Qoder",
		ProjectDir:  ".qoder/skills",
		GlobalDir:   ".qoder/skills",
		DetectPaths: []string{".qoder"},
		IsBuiltin:   true,
		Category:    "coding",
	},
	{
		Name:        "qwen_code",
		DisplayName: "Qwen Code",
		ProjectDir:  ".qwen_code/skills",
		GlobalDir:   ".qwen_code/skills",
		DetectPaths: []string{".qwen_code"},
		IsBuiltin:   true,
		Category:    "coding",
	},
	{
		Name:        "trae_cn",
		DisplayName: "Trae CN",
		ProjectDir:  ".trae_cn/skills",
		GlobalDir:   ".trae_cn/skills",
		DetectPaths: []string{".trae_cn"},
		IsBuiltin:   true,
		Category:    "coding",
	},
	{
		Name:        "zencoder",
		DisplayName: "Zencoder",
		ProjectDir:  ".zencoder/skills",
		GlobalDir:   ".zencoder/skills",
		DetectPaths: []string{".zencoder"},
		IsBuiltin:   true,
		Category:    "coding",
	},
	{
		Name:        "adal",
		DisplayName: "Adal",
		ProjectDir:  ".adal/skills",
		GlobalDir:   ".adal/skills",
		DetectPaths: []string{".adal"},
		IsBuiltin:   true,
		Category:    "coding",
	},

	// ── Lobster agents (6) ──────────────────────────────────────────────
	{
		Name:        "openclaw",
		DisplayName: "OpenClaw",
		ProjectDir:  ".openclaw/skills",
		GlobalDir:   ".openclaw/skills",
		DetectPaths: []string{".openclaw"},
		IsBuiltin:   true,
		Category:    "lobster",
	},
	{
		Name:        "hermes",
		DisplayName: "Hermes",
		ProjectDir:  ".hermes/skills",
		GlobalDir:   ".hermes/skills",
		DetectPaths: []string{".hermes"},
		IsBuiltin:   true,
		Category:    "lobster",
	},
	{
		Name:        "qclaw",
		DisplayName: "QClaw",
		ProjectDir:  ".qclaw/skills",
		GlobalDir:   ".qclaw/skills",
		DetectPaths: []string{".qclaw"},
		IsBuiltin:   true,
		Category:    "lobster",
	},
	{
		Name:        "easyclaw",
		DisplayName: "EasyClaw",
		ProjectDir:  ".easyclaw/skills",
		GlobalDir:   ".easyclaw/skills",
		DetectPaths: []string{".easyclaw"},
		IsBuiltin:   true,
		Category:    "lobster",
	},
	{
		Name:        "autoclaw",
		DisplayName: "AutoClaw",
		ProjectDir:  ".autoclaw/skills",
		GlobalDir:   ".autoclaw/skills",
		DetectPaths: []string{".autoclaw"},
		IsBuiltin:   true,
		Category:    "lobster",
	},
	{
		Name:        "workbuddy",
		DisplayName: "WorkBuddy",
		ProjectDir:  ".workbuddy/skills",
		GlobalDir:   ".workbuddy/skills",
		DetectPaths: []string{".workbuddy"},
		IsBuiltin:   true,
		Category:    "lobster",
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
