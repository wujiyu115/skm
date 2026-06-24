// Package sync handles deploying skills to agent directories via symlink
// (default) or copy (fallback). Independent of store — purely filesystem
// operations.
package sync

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// SyncResult captures the outcome of a single skill sync operation.
type SyncResult struct {
	SkillName  string
	Agent      string
	TargetPath string
	Mode       string
	Status     string
	Error      error
}

// SyncSkill deploys a skill from sourcePath to targetPath using the given mode
// ("symlink" or "copy"). Parent directories are created as needed. Returns an
// error if target is inside source (would cause infinite recursion on copy).
func SyncSkill(sourcePath, targetPath, mode string) error {
	absSource, err := filepath.Abs(sourcePath)
	if err != nil {
		return err
	}
	absTarget, err := filepath.Abs(targetPath)
	if err != nil {
		return err
	}

	// Guard against target-inside-source (infinite copy, circular symlink).
	if strings.HasPrefix(absTarget, absSource+string(filepath.Separator)) {
		return fmt.Errorf("target %s is inside source %s", absTarget, absSource)
	}

	// Ensure parent directory exists.
	if err := os.MkdirAll(filepath.Dir(absTarget), 0755); err != nil {
		return fmt.Errorf("create parent dirs: %w", err)
	}

	// Remove any pre-existing target (symlink, dir, or file).
	if info, err := os.Lstat(absTarget); err == nil {
		if info.Mode()&os.ModeSymlink != 0 {
			os.Remove(absTarget)
		} else if info.IsDir() {
			os.RemoveAll(absTarget)
		} else {
			os.Remove(absTarget)
		}
	}

	if mode == "symlink" {
		return os.Symlink(absSource, absTarget)
	}

	return copyDir(absSource, absTarget)
}

// Unsync removes a previously synced skill at targetPath. Works for both
// symlinks and copied directories. No-op if the target does not exist.
func Unsync(targetPath string) error {
	info, err := os.Lstat(targetPath)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return err
	}

	if info.Mode()&os.ModeSymlink != 0 {
		return os.Remove(targetPath)
	}
	return os.RemoveAll(targetPath)
}

// IsCurrent reports whether targetPath is up to date with sourcePath for the
// given mode. For symlink mode it checks that the symlink destination matches
// the source. For copy mode it only checks that the target directory exists
// (full content hashing is left to higher-level callers).
func IsCurrent(sourcePath, targetPath, mode string) bool {
	info, err := os.Lstat(targetPath)
	if err != nil {
		return false
	}

	if mode == "symlink" {
		if info.Mode()&os.ModeSymlink == 0 {
			return false
		}
		link, err := os.Readlink(targetPath)
		if err != nil {
			return false
		}
		absSource, _ := filepath.Abs(sourcePath)
		absLink, _ := filepath.Abs(link)
		return absSource == absLink
	}

	// Copy mode: directory exists → considered current.
	return info.IsDir()
}

// copyDir recursively copies src to dst, skipping .git directories.
func copyDir(src, dst string) error {
	return filepath.WalkDir(src, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.Name() == ".git" && d.IsDir() {
			return filepath.SkipDir
		}

		rel, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		target := filepath.Join(dst, rel)

		if d.IsDir() {
			return os.MkdirAll(target, 0755)
		}
		return copyFile(path, target)
	})
}

// copyFile copies a single file from src to dst, creating parent directories
// as needed.
func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
		return err
	}

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	if _, err = io.Copy(out, in); err != nil {
		return err
	}
	return out.Close()
}
