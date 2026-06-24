package fetcher

import (
	"fmt"
	"os"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

// CloneToTemp clones a git repository to a temporary directory.
// It performs a shallow clone (depth=1) for efficiency.
// The caller is responsible for cleaning up the temp directory via CleanupTemp.
func CloneToTemp(src ParsedSource) (string, error) {
	url := src.GitURL()
	if url == "" {
		return "", fmt.Errorf("not a git source: %s", src.OriginalInput)
	}

	tmpDir, err := os.MkdirTemp("", "skm-clone-")
	if err != nil {
		return "", fmt.Errorf("create temp dir: %w", err)
	}

	opts := &git.CloneOptions{
		URL:   url,
		Depth: 1,
	}
	if src.Branch != "" {
		opts.ReferenceName = plumbing.NewBranchReferenceName(src.Branch)
		opts.SingleBranch = true
	}

	if _, err := git.PlainClone(tmpDir, false, opts); err != nil {
		os.RemoveAll(tmpDir)
		return "", fmt.Errorf("git clone %s: %w", url, err)
	}

	return tmpDir, nil
}

// CleanupTemp removes a temporary directory created by CloneToTemp.
func CleanupTemp(dir string) {
	os.RemoveAll(dir)
}
