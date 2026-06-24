package skill

import (
	"crypto/sha256"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

var skipDirs = map[string]bool{
	".git": true, "__pycache__": true, "node_modules": true,
	".DS_Store": true, "Thumbs.db": true,
}

func HashDir(dir string) (string, error) {
	type entry struct {
		relPath string
		content []byte
	}

	var entries []entry

	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		name := d.Name()
		if d.IsDir() {
			if skipDirs[name] {
				return filepath.SkipDir
			}
			return nil
		}

		if skipDirs[name] {
			return nil
		}

		rel, err := filepath.Rel(dir, path)
		if err != nil {
			return err
		}
		rel = filepath.ToSlash(rel)

		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		entries = append(entries, entry{relPath: rel, content: content})
		return nil
	})
	if err != nil {
		return "", fmt.Errorf("walk directory: %w", err)
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].relPath < entries[j].relPath
	})

	h := sha256.New()
	for _, e := range entries {
		fmt.Fprintf(h, "%s\n%d\n", e.relPath, len(e.content))
		h.Write(e.content)
	}

	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

func IsSkillDir(dir string) bool {
	for _, name := range []string{"SKILL.md", "skill.md"} {
		p := filepath.Join(dir, name)
		if info, err := os.Stat(p); err == nil && !info.IsDir() {
			return true
		}
	}
	return false
}

func SanitizeName(name string) string {
	replacer := strings.NewReplacer(
		"<", "_", ">", "_", ":", "_", "\"", "_",
		"/", "_", "\\", "_", "|", "_", "?", "_", "*", "_",
	)
	name = replacer.Replace(name)
	name = strings.TrimSpace(name)
	name = strings.Trim(name, ".")
	if name == "" {
		name = "_"
	}
	return name
}
