package skill

import (
	"os"
	"path/filepath"
	"testing"
)

func TestHashDir_Deterministic(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "SKILL.md"), []byte("# test"), 0644)
	os.WriteFile(filepath.Join(dir, "helper.txt"), []byte("data"), 0644)

	h1, err := HashDir(dir)
	if err != nil {
		t.Fatalf("HashDir() error: %v", err)
	}
	h2, err := HashDir(dir)
	if err != nil {
		t.Fatalf("HashDir() second call error: %v", err)
	}
	if h1 != h2 {
		t.Fatalf("hashes differ: %s != %s", h1, h2)
	}
	if len(h1) != 64 {
		t.Fatalf("expected 64-char hex hash, got len=%d", len(h1))
	}
}

func TestHashDir_DifferentContent(t *testing.T) {
	d1 := t.TempDir()
	os.WriteFile(filepath.Join(d1, "SKILL.md"), []byte("# v1"), 0644)

	d2 := t.TempDir()
	os.WriteFile(filepath.Join(d2, "SKILL.md"), []byte("# v2"), 0644)

	h1, _ := HashDir(d1)
	h2, _ := HashDir(d2)
	if h1 == h2 {
		t.Fatal("different content should produce different hashes")
	}
}

func TestHashDir_SkipsGitDir(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "SKILL.md"), []byte("# test"), 0644)
	gitDir := filepath.Join(dir, ".git")
	os.MkdirAll(gitDir, 0755)
	os.WriteFile(filepath.Join(gitDir, "HEAD"), []byte("ref: refs/heads/main"), 0644)

	h1, _ := HashDir(dir)

	dir2 := t.TempDir()
	os.WriteFile(filepath.Join(dir2, "SKILL.md"), []byte("# test"), 0644)

	h2, _ := HashDir(dir2)
	if h1 != h2 {
		t.Fatal(".git should be ignored in hash")
	}
}
