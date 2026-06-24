package fetcher

import "testing"

func TestParse_GitHubTreeURL(t *testing.T) {
	src := Parse("https://github.com/anthropics/skills/tree/main/skills/mcp-builder")
	if src.Type != "github-tree" {
		t.Errorf("Type = %q, want github-tree", src.Type)
	}
	if src.Owner != "anthropics" || src.Repo != "skills" {
		t.Errorf("Owner/Repo = %s/%s", src.Owner, src.Repo)
	}
	if src.Branch != "main" {
		t.Errorf("Branch = %q, want main", src.Branch)
	}
	if src.Subpath != "skills/mcp-builder" {
		t.Errorf("Subpath = %q", src.Subpath)
	}
	if src.SkillName != "mcp-builder" {
		t.Errorf("SkillName = %q", src.SkillName)
	}
}

func TestParse_GitHubRepoURL(t *testing.T) {
	src := Parse("https://github.com/owner/repo")
	if src.Type != "github" {
		t.Errorf("Type = %q, want github", src.Type)
	}
	if src.Owner != "owner" || src.Repo != "repo" {
		t.Errorf("Owner/Repo = %s/%s", src.Owner, src.Repo)
	}
	if src.SkillName != "repo" {
		t.Errorf("SkillName = %q, want repo", src.SkillName)
	}
}

func TestParse_GitHubRepoURL_DotGit(t *testing.T) {
	src := Parse("https://github.com/owner/repo.git")
	if src.Repo != "repo" {
		t.Errorf("Repo = %q, want repo (strip .git)", src.Repo)
	}
}

func TestParse_GitHubShorthand(t *testing.T) {
	src := Parse("anthropics/skills/skills/mcp-builder")
	if src.Type != "github-shorthand" {
		t.Errorf("Type = %q, want github-shorthand", src.Type)
	}
	if src.Owner != "anthropics" || src.Repo != "skills" {
		t.Errorf("Owner/Repo = %s/%s", src.Owner, src.Repo)
	}
	if src.Subpath != "skills/mcp-builder" {
		t.Errorf("Subpath = %q", src.Subpath)
	}
}

func TestParse_GitHubShorthandSimple(t *testing.T) {
	src := Parse("owner/repo")
	if src.Type != "github-shorthand" {
		t.Errorf("Type = %q, want github-shorthand", src.Type)
	}
	if src.SkillName != "repo" {
		t.Errorf("SkillName = %q, want repo", src.SkillName)
	}
}

func TestParse_LocalPath(t *testing.T) {
	for _, input := range []string{"./my-skill", "/abs/path/skill", "../relative"} {
		src := Parse(input)
		if src.Type != "local" {
			t.Errorf("Parse(%q).Type = %q, want local", input, src.Type)
		}
	}
}
