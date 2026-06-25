package skill

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

var priorityDirs = []string{
	"skills",
	".claude/skills",
	".cursor/skills",
	".codex/skills",
}

const maxRecurseDepth = 5

// Discover finds all skills in a directory tree. If subpath is non-empty,
// the search is scoped to baseDir/subpath. It checks priority directories
// first (skills/, .claude/skills/, etc.) and falls back to recursive search
// up to maxRecurseDepth. Skills are deduplicated by name.
func Discover(baseDir string, subpath string) ([]*Skill, error) {
	searchDir := baseDir
	if subpath != "" {
		found := false
		candidate := filepath.Join(baseDir, subpath)
		if info, err := os.Stat(candidate); err == nil && info.IsDir() {
			searchDir = candidate
			found = true
		} else {
			for _, pd := range priorityDirs {
				alt := filepath.Join(baseDir, pd, subpath)
				if info, err := os.Stat(alt); err == nil && info.IsDir() {
					searchDir = alt
					found = true
					break
				}
			}
		}
		if !found {
			return nil, fmt.Errorf("subpath %q not found in repository", subpath)
		}
	}

	// If the search dir itself is a skill, return it directly.
	if HasSkillMD(searchDir) {
		s, err := ParseSkillMD(searchDir)
		if err != nil {
			return nil, err
		}
		s.CentralPath = searchDir
		return []*Skill{s}, nil
	}

	seen := map[string]bool{}
	var skills []*Skill

	addSkillsFromDir := func(dir string) {
		entries, err := os.ReadDir(dir)
		if err != nil {
			return
		}
		for _, e := range entries {
			if !e.IsDir() {
				continue
			}
			childDir := filepath.Join(dir, e.Name())
			if HasSkillMD(childDir) {
				s, err := ParseSkillMD(childDir)
				if err != nil {
					continue
				}
				if seen[s.Name] {
					continue
				}
				seen[s.Name] = true
				s.CentralPath = childDir
				skills = append(skills, s)
			}
		}
	}

	// Check immediate children of searchDir.
	addSkillsFromDir(searchDir)

	// Check priority directories.
	for _, pd := range priorityDirs {
		d := filepath.Join(searchDir, pd)
		if info, err := os.Stat(d); err == nil && info.IsDir() {
			addSkillsFromDir(d)
		}
	}

	// If we found skills in priority locations, return early.
	if len(skills) > 0 {
		return skills, nil
	}

	// Fallback: recursive walk.
	filepath.WalkDir(searchDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if !d.IsDir() {
			return nil
		}

		rel, _ := filepath.Rel(searchDir, path)
		depth := len(strings.Split(rel, string(filepath.Separator)))
		if depth > maxRecurseDepth {
			return filepath.SkipDir
		}

		if d.Name() == ".git" || d.Name() == "node_modules" {
			return filepath.SkipDir
		}

		if HasSkillMD(path) {
			s, err := ParseSkillMD(path)
			if err != nil {
				return nil
			}
			if !seen[s.Name] {
				seen[s.Name] = true
				s.CentralPath = path
				skills = append(skills, s)
			}
			return filepath.SkipDir
		}

		return nil
	})

	return skills, nil
}
