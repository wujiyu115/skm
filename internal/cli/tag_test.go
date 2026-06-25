package cli

import (
	"strings"
	"testing"
)

func TestTagAdd(t *testing.T) {
	cfg := testConfig(t)
	insertTestSkill(t, cfg.Store, "my-skill")

	if err := tagAdd(cfg, "my-skill", "productivity"); err != nil {
		t.Fatalf("tagAdd error: %v", err)
	}

	tags, err := cfg.Store.ListSkillTags("test-my-skill")
	if err != nil {
		t.Fatal(err)
	}
	if len(tags) != 1 || tags[0] != "productivity" {
		t.Errorf("tags = %v, want [productivity]", tags)
	}
}

func TestTagAdd_NotFound(t *testing.T) {
	cfg := testConfig(t)

	err := tagAdd(cfg, "nonexistent", "foo")
	if err == nil {
		t.Fatal("expected error for nonexistent skill")
	}
	if !strings.Contains(err.Error(), "nonexistent") {
		t.Errorf("error should mention skill name, got: %v", err)
	}
}

func TestTagAdd_Duplicate(t *testing.T) {
	cfg := testConfig(t)
	insertTestSkill(t, cfg.Store, "my-skill")

	if err := tagAdd(cfg, "my-skill", "dev"); err != nil {
		t.Fatal(err)
	}
	// Adding same tag again should not error (INSERT OR IGNORE).
	if err := tagAdd(cfg, "my-skill", "dev"); err != nil {
		t.Fatalf("duplicate tagAdd should not error, got: %v", err)
	}

	tags, err := cfg.Store.ListSkillTags("test-my-skill")
	if err != nil {
		t.Fatal(err)
	}
	if len(tags) != 1 {
		t.Errorf("expected 1 tag after duplicate add, got %d", len(tags))
	}
}

func TestTagRemove(t *testing.T) {
	cfg := testConfig(t)
	insertTestSkill(t, cfg.Store, "my-skill")

	if err := tagAdd(cfg, "my-skill", "test-tag"); err != nil {
		t.Fatal(err)
	}

	if err := tagRemove(cfg, "my-skill", "test-tag"); err != nil {
		t.Fatalf("tagRemove error: %v", err)
	}

	tags, err := cfg.Store.ListSkillTags("test-my-skill")
	if err != nil {
		t.Fatal(err)
	}
	if len(tags) != 0 {
		t.Errorf("expected 0 tags after remove, got %v", tags)
	}
}

func TestTagRemove_NotFound(t *testing.T) {
	cfg := testConfig(t)

	err := tagRemove(cfg, "nonexistent", "foo")
	if err == nil {
		t.Fatal("expected error for nonexistent skill")
	}
}

func TestTagRename(t *testing.T) {
	cfg := testConfig(t)
	insertTestSkill(t, cfg.Store, "s1")
	insertTestSkill(t, cfg.Store, "s2")

	if err := tagAdd(cfg, "s1", "old-name"); err != nil {
		t.Fatal(err)
	}
	if err := tagAdd(cfg, "s2", "old-name"); err != nil {
		t.Fatal(err)
	}

	if err := tagRename(cfg, "old-name", "new-name"); err != nil {
		t.Fatalf("tagRename error: %v", err)
	}

	// Both skills should now have "new-name" instead of "old-name".
	for _, id := range []string{"test-s1", "test-s2"} {
		tags, err := cfg.Store.ListSkillTags(id)
		if err != nil {
			t.Fatal(err)
		}
		if len(tags) != 1 || tags[0] != "new-name" {
			t.Errorf("skill %s tags = %v, want [new-name]", id, tags)
		}
	}

	// "old-name" should no longer appear in global tag list.
	allTags, _ := cfg.Store.ListTags()
	for _, tag := range allTags {
		if tag == "old-name" {
			t.Error("old-name should not appear in global tags after rename")
		}
	}
}

func TestTagList(t *testing.T) {
	cfg := testConfig(t)
	insertTestSkill(t, cfg.Store, "s1")
	insertTestSkill(t, cfg.Store, "s2")

	tagAdd(cfg, "s1", "alpha")
	tagAdd(cfg, "s1", "beta")
	tagAdd(cfg, "s2", "beta")
	tagAdd(cfg, "s2", "gamma")

	tags, err := cfg.Store.ListTags()
	if err != nil {
		t.Fatal(err)
	}

	want := []string{"alpha", "beta", "gamma"}
	if len(tags) != len(want) {
		t.Fatalf("tags = %v, want %v", tags, want)
	}
	for i, tag := range tags {
		if tag != want[i] {
			t.Errorf("tags[%d] = %q, want %q", i, tag, want[i])
		}
	}
}
