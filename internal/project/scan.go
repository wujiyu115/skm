package project

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"

	"github.com/ejoy/skm/internal/agent"
	"github.com/ejoy/skm/internal/skill"
)

// ProjectSkill represents a skill installed in a project workspace.
type ProjectSkill struct {
	Agent        string `json:"agent"`
	AgentDisplay string `json:"agent_display"`
	SkillName    string `json:"skill_name"`
	Description  string `json:"description"`
	SkillPath    string `json:"skill_path"`
	Enabled      bool   `json:"enabled"`
}

// disabledDir returns the disabled-skills directory path corresponding to a
// given project directory. For example ".claude/skills" becomes
// ".claude/skills-disabled".
func disabledDir(projectDir string) string {
	return DisabledDirFor(projectDir)
}

// DisabledDirFor returns the disabled-skills directory path for a given
// agent project directory. Exported for use by other packages.
func DisabledDirFor(projectDir string) string {
	return filepath.Join(filepath.Dir(projectDir), filepath.Base(projectDir)+"-disabled")
}

// ScanSkills discovers all skills (enabled and disabled) installed in the
// project at projectPath for the given set of agent adapters. Results are
// sorted by agent name then skill name.
func ScanSkills(projectPath string, agents []agent.Adapter) ([]ProjectSkill, error) {
	var result []ProjectSkill

	for _, a := range agents {
		if a.ProjectDir == "" {
			continue
		}

		enabledPath := filepath.Join(projectPath, a.ProjectDir)
		disabledPath := filepath.Join(projectPath, disabledDir(a.ProjectDir))

		for _, dir := range []struct {
			path    string
			enabled bool
		}{
			{enabledPath, true},
			{disabledPath, false},
		} {
			entries, err := os.ReadDir(dir.path)
			if err != nil {
				// Directory doesn't exist — that's fine, skip it.
				continue
			}
			for _, e := range entries {
				childDir := filepath.Join(dir.path, e.Name())
				info, err := os.Stat(childDir)
				if err != nil || !info.IsDir() {
					continue
				}
				if !skill.HasSkillMD(childDir) {
					continue
				}
				s, err := skill.ParseSkillMD(childDir)
				if err != nil {
					continue
				}
				result = append(result, ProjectSkill{
					Agent:        a.Name,
					AgentDisplay: a.DisplayName,
					SkillName:    s.Name,
					Description:  s.Description,
					SkillPath:    childDir,
					Enabled:      dir.enabled,
				})
			}
		}
	}

	sort.Slice(result, func(i, j int) bool {
		if result[i].Agent != result[j].Agent {
			return result[i].Agent < result[j].Agent
		}
		return result[i].SkillName < result[j].SkillName
	})

	return result, nil
}

// ToggleSkill moves a skill directory between the enabled and disabled
// directories for the given agent. When enable is true the skill is moved from
// disabled to enabled; when false, from enabled to disabled.
func ToggleSkill(projectPath, agentProjectDir, skillName string, enable bool) error {
	var srcBase, dstBase string
	if enable {
		srcBase = disabledDir(agentProjectDir)
		dstBase = agentProjectDir
	} else {
		srcBase = agentProjectDir
		dstBase = disabledDir(agentProjectDir)
	}

	src := filepath.Join(projectPath, srcBase, skillName)
	dst := filepath.Join(projectPath, dstBase, skillName)

	if _, err := os.Stat(src); os.IsNotExist(err) {
		return fmt.Errorf("skill source not found: %s", src)
	}

	if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
		return fmt.Errorf("create destination parent: %w", err)
	}

	if err := os.Rename(src, dst); err != nil {
		return fmt.Errorf("move skill: %w", err)
	}

	return nil
}

// RemoveSkill removes a skill directory and all its contents.
func RemoveSkill(skillPath string) error {
	return os.RemoveAll(skillPath)
}

// AddFromCentral copies a skill from the central repository into the project
// workspace. It performs a recursive copy of all files and directories.
func AddFromCentral(centralPath, projectPath, agentProjectDir, skillName string) error {
	target := filepath.Join(projectPath, agentProjectDir, skillName)

	if _, err := os.Stat(target); err == nil {
		return fmt.Errorf("skill already exists at %s", target)
	}

	if err := os.MkdirAll(target, 0o755); err != nil {
		return fmt.Errorf("create target directory: %w", err)
	}

	err := filepath.WalkDir(centralPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		rel, err := filepath.Rel(centralPath, path)
		if err != nil {
			return err
		}

		dest := filepath.Join(target, rel)

		if d.IsDir() {
			return os.MkdirAll(dest, 0o755)
		}

		data, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("read %s: %w", path, err)
		}

		info, err := d.Info()
		if err != nil {
			return fmt.Errorf("stat %s: %w", path, err)
		}

		return os.WriteFile(dest, data, info.Mode())
	})

	if err != nil {
		// Clean up partial copy on failure.
		os.RemoveAll(target)
		return fmt.Errorf("copy skill: %w", err)
	}

	return nil
}
