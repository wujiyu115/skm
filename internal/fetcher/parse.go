package fetcher

import (
	"path"
	"regexp"
	"strings"
)

// ParsedSource represents a parsed skill source input.
type ParsedSource struct {
	Type          string // "github-tree", "github", "github-shorthand", "local"
	Owner         string
	Repo          string
	Branch        string
	Subpath       string
	SkillName     string
	OriginalInput string
}

var (
	githubTreeRe = regexp.MustCompile(`^https?://github\.com/([^/]+)/([^/]+)/tree/([^/]+)/(.+)$`)
	githubRepoRe = regexp.MustCompile(`^https?://github\.com/([^/]+)/([^/]+?)(?:\.git)?$`)
	shorthandRe  = regexp.MustCompile(`^([a-zA-Z0-9_.-]+)/([a-zA-Z0-9_.-]+)(/.*)?$`)
)

// Parse takes a skill source input string and returns a structured ParsedSource.
// Supported formats:
//   - GitHub tree URL: https://github.com/owner/repo/tree/branch/path
//   - GitHub repo URL: https://github.com/owner/repo[.git]
//   - GitHub shorthand: owner/repo[/subpath]
//   - Local path: ./path, ../path, /absolute/path
func Parse(input string) ParsedSource {
	src := ParsedSource{OriginalInput: input}

	if m := githubTreeRe.FindStringSubmatch(input); m != nil {
		src.Type = "github-tree"
		src.Owner = m[1]
		src.Repo = m[2]
		src.Branch = m[3]
		src.Subpath = m[4]
		src.SkillName = path.Base(src.Subpath)
		return src
	}

	if m := githubRepoRe.FindStringSubmatch(input); m != nil {
		src.Type = "github"
		src.Owner = m[1]
		src.Repo = m[2]
		src.SkillName = m[2]
		return src
	}

	if isLocalPath(input) {
		src.Type = "local"
		src.SkillName = path.Base(input)
		src.Subpath = input
		return src
	}

	if m := shorthandRe.FindStringSubmatch(input); m != nil {
		src.Type = "github-shorthand"
		src.Owner = m[1]
		src.Repo = m[2]
		if m[3] != "" {
			src.Subpath = strings.TrimPrefix(m[3], "/")
			src.SkillName = path.Base(src.Subpath)
		} else {
			src.SkillName = m[2]
		}
		return src
	}

	src.Type = "local"
	src.SkillName = input
	return src
}

func isLocalPath(input string) bool {
	return strings.HasPrefix(input, "./") ||
		strings.HasPrefix(input, "../") ||
		strings.HasPrefix(input, "/")
}

// GitURL returns the HTTPS clone URL for GitHub sources, or empty string for local sources.
func (s ParsedSource) GitURL() string {
	switch s.Type {
	case "github-tree", "github", "github-shorthand":
		return "https://github.com/" + s.Owner + "/" + s.Repo + ".git"
	default:
		return ""
	}
}
