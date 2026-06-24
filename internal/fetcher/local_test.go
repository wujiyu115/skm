package fetcher

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCopyLocal(t *testing.T) {
	src := t.TempDir()
	os.WriteFile(filepath.Join(src, "SKILL.md"), []byte("# test"), 0644)
	sub := filepath.Join(src, "sub")
	os.MkdirAll(sub, 0755)
	os.WriteFile(filepath.Join(sub, "helper.go"), []byte("package sub"), 0644)

	dst := filepath.Join(t.TempDir(), "copied")
	err := CopyLocal(src, dst)
	if err != nil {
		t.Fatalf("CopyLocal() error: %v", err)
	}

	if _, err := os.Stat(filepath.Join(dst, "SKILL.md")); err != nil {
		t.Error("SKILL.md not copied")
	}
	if _, err := os.Stat(filepath.Join(dst, "sub", "helper.go")); err != nil {
		t.Error("sub/helper.go not copied")
	}
}

func TestCopyLocal_SkipsGitDir(t *testing.T) {
	src := t.TempDir()
	os.WriteFile(filepath.Join(src, "SKILL.md"), []byte("# test"), 0644)
	os.MkdirAll(filepath.Join(src, ".git"), 0755)
	os.WriteFile(filepath.Join(src, ".git", "HEAD"), []byte("ref"), 0644)

	dst := filepath.Join(t.TempDir(), "copied")
	CopyLocal(src, dst)

	if _, err := os.Stat(filepath.Join(dst, ".git")); !os.IsNotExist(err) {
		t.Error(".git should not be copied")
	}
}
