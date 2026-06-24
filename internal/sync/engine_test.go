// internal/sync/engine_test.go
package sync

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSyncSkill_Symlink(t *testing.T) {
	src := t.TempDir()
	os.WriteFile(filepath.Join(src, "SKILL.md"), []byte("# test"), 0644)

	target := filepath.Join(t.TempDir(), "my-skill")

	err := SyncSkill(src, target, "symlink")
	if err != nil {
		t.Fatalf("SyncSkill() error: %v", err)
	}

	link, err := os.Readlink(target)
	if err != nil {
		t.Fatalf("Readlink error: %v", err)
	}
	if link != src {
		t.Errorf("symlink points to %q, want %q", link, src)
	}
}

func TestSyncSkill_Copy(t *testing.T) {
	src := t.TempDir()
	os.WriteFile(filepath.Join(src, "SKILL.md"), []byte("# test"), 0644)

	target := filepath.Join(t.TempDir(), "my-skill")

	err := SyncSkill(src, target, "copy")
	if err != nil {
		t.Fatalf("SyncSkill() error: %v", err)
	}

	if _, err := os.Lstat(target); err != nil {
		t.Fatal("target not created")
	}
	info, _ := os.Lstat(target)
	if info.Mode()&os.ModeSymlink != 0 {
		t.Error("copy mode should not create symlink")
	}
	if _, err := os.Stat(filepath.Join(target, "SKILL.md")); err != nil {
		t.Error("SKILL.md not copied")
	}
}

func TestUnsync(t *testing.T) {
	src := t.TempDir()
	target := filepath.Join(t.TempDir(), "link")
	os.Symlink(src, target)

	if err := Unsync(target); err != nil {
		t.Fatalf("Unsync() error: %v", err)
	}
	if _, err := os.Lstat(target); !os.IsNotExist(err) {
		t.Error("target should be removed")
	}
}

func TestIsCurrent_Symlink(t *testing.T) {
	src := t.TempDir()
	target := filepath.Join(t.TempDir(), "link")
	os.Symlink(src, target)

	if !IsCurrent(src, target, "symlink") {
		t.Error("should be current when symlink points to source")
	}

	wrongSrc := t.TempDir()
	if IsCurrent(wrongSrc, target, "symlink") {
		t.Error("should not be current when symlink points elsewhere")
	}
}

func TestIsCurrent_Missing(t *testing.T) {
	if IsCurrent("/src", "/nonexistent", "symlink") {
		t.Error("missing target should not be current")
	}
}

func TestSyncSkill_CreatesParentDirs(t *testing.T) {
	src := t.TempDir()
	os.WriteFile(filepath.Join(src, "SKILL.md"), []byte("# test"), 0644)

	target := filepath.Join(t.TempDir(), "deep", "nested", "my-skill")

	err := SyncSkill(src, target, "symlink")
	if err != nil {
		t.Fatalf("SyncSkill() error: %v", err)
	}
	if _, err := os.Lstat(target); err != nil {
		t.Error("target not created")
	}
}

func TestSyncSkill_TargetInsideSource(t *testing.T) {
	src := t.TempDir()
	target := filepath.Join(src, "sub", "skill")

	err := SyncSkill(src, target, "symlink")
	if err == nil {
		t.Fatal("should reject target inside source")
	}
}
