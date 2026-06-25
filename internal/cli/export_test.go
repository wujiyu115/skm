package cli

import (
	"encoding/json"
	"testing"
)

func TestExportSkills(t *testing.T) {
	cfg := testConfig(t)
	insertTestSkill(t, cfg.Store, "skill-a")
	insertTestSkill(t, cfg.Store, "skill-b")

	// Add a tag.
	if err := cfg.Store.AddTag("test-skill-a", "dev"); err != nil {
		t.Fatal(err)
	}

	// Build the export entries the same way exportSkills does.
	skills, err := cfg.Store.ListSkills()
	if err != nil {
		t.Fatal(err)
	}

	entries := make([]exportEntry, 0, len(skills))
	for _, sk := range skills {
		tags, _ := cfg.Store.ListSkillTags(sk.ID)
		if tags == nil {
			tags = []string{}
		}
		entries = append(entries, exportEntry{
			Name:        sk.Name,
			Description: sk.Description,
			SourceType:  sk.SourceType,
			SourceRef:   sk.SourceRef,
			Enabled:     sk.Enabled,
			Tags:        tags,
		})
	}

	data, err := json.MarshalIndent(entries, "", "  ")
	if err != nil {
		t.Fatalf("json.MarshalIndent error: %v", err)
	}

	// Verify it's valid JSON and has the right structure.
	var parsed []exportEntry
	if err := json.Unmarshal(data, &parsed); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}

	if len(parsed) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(parsed))
	}

	// Skills are ordered by name.
	if parsed[0].Name != "skill-a" {
		t.Errorf("first entry name = %q, want %q", parsed[0].Name, "skill-a")
	}
	if parsed[1].Name != "skill-b" {
		t.Errorf("second entry name = %q, want %q", parsed[1].Name, "skill-b")
	}

	// Check tags on skill-a.
	if len(parsed[0].Tags) != 1 || parsed[0].Tags[0] != "dev" {
		t.Errorf("skill-a tags = %v, want [dev]", parsed[0].Tags)
	}

	// Check all required fields are present.
	for _, e := range parsed {
		if e.Name == "" {
			t.Error("entry missing name")
		}
		if e.SourceType == "" {
			t.Error("entry missing source_type")
		}
	}
}

func TestExportEmpty(t *testing.T) {
	cfg := testConfig(t)

	skills, err := cfg.Store.ListSkills()
	if err != nil {
		t.Fatal(err)
	}

	entries := make([]exportEntry, 0)
	for _, sk := range skills {
		tags, _ := cfg.Store.ListSkillTags(sk.ID)
		if tags == nil {
			tags = []string{}
		}
		entries = append(entries, exportEntry{
			Name:        sk.Name,
			Description: sk.Description,
			SourceType:  sk.SourceType,
			SourceRef:   sk.SourceRef,
			Enabled:     sk.Enabled,
			Tags:        tags,
		})
	}

	data, err := json.MarshalIndent(entries, "", "  ")
	if err != nil {
		t.Fatal(err)
	}

	var parsed []exportEntry
	if err := json.Unmarshal(data, &parsed); err != nil {
		t.Fatalf("empty export is not valid JSON: %v", err)
	}
	if len(parsed) != 0 {
		t.Errorf("expected 0 entries, got %d", len(parsed))
	}
}
