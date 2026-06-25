package cli

import (
	"testing"
)

func TestSearchByName(t *testing.T) {
	cfg := testConfig(t)
	insertTestSkill(t, cfg.Store, "git-helper")
	insertTestSkill(t, cfg.Store, "docker-tools")

	results, err := searchSkills(cfg, "git")
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Skill.Name != "git-helper" {
		t.Errorf("expected git-helper, got %s", results[0].Skill.Name)
	}
}

func TestSearchByDescription(t *testing.T) {
	cfg := testConfig(t)
	insertTestSkill(t, cfg.Store, "my-skill")

	// insertTestSkill sets description to "test skill my-skill"
	results, err := searchSkills(cfg, "test skill")
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
}

func TestSearchByTag(t *testing.T) {
	cfg := testConfig(t)
	insertTestSkill(t, cfg.Store, "alpha")
	insertTestSkill(t, cfg.Store, "beta")

	// Add a tag only to "alpha".
	if err := cfg.Store.AddTag("test-alpha", "productivity"); err != nil {
		t.Fatal(err)
	}

	results, err := searchSkills(cfg, "productivity")
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Skill.Name != "alpha" {
		t.Errorf("expected alpha, got %s", results[0].Skill.Name)
	}
}

func TestSearchCaseInsensitive(t *testing.T) {
	cfg := testConfig(t)
	insertTestSkill(t, cfg.Store, "MyTool")

	results, err := searchSkills(cfg, "mytool")
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result for case-insensitive search, got %d", len(results))
	}
}

func TestSearchNoMatch(t *testing.T) {
	cfg := testConfig(t)
	insertTestSkill(t, cfg.Store, "alpha")

	results, err := searchSkills(cfg, "zzz-nonexistent")
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 0 {
		t.Fatalf("expected 0 results, got %d", len(results))
	}
}
