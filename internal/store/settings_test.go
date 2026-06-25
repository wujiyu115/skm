package store

import (
	"database/sql"
	"errors"
	"testing"
)

func TestGetSetting_NotFound(t *testing.T) {
	s := newTestStore(t)
	defer s.Close()

	_, err := s.GetSetting("nonexistent")
	if !errors.Is(err, sql.ErrNoRows) {
		t.Fatalf("expected sql.ErrNoRows, got: %v", err)
	}
}

func TestSetAndGetSetting(t *testing.T) {
	s := newTestStore(t)
	defer s.Close()

	if err := s.SetSetting("theme", "dark"); err != nil {
		t.Fatalf("SetSetting() error: %v", err)
	}

	val, err := s.GetSetting("theme")
	if err != nil {
		t.Fatalf("GetSetting() error: %v", err)
	}
	if val != "dark" {
		t.Fatalf("expected 'dark', got %q", val)
	}
}

func TestSetSetting_Upsert(t *testing.T) {
	s := newTestStore(t)
	defer s.Close()

	if err := s.SetSetting("theme", "dark"); err != nil {
		t.Fatalf("SetSetting() error: %v", err)
	}
	if err := s.SetSetting("theme", "light"); err != nil {
		t.Fatalf("SetSetting() upsert error: %v", err)
	}

	val, err := s.GetSetting("theme")
	if err != nil {
		t.Fatalf("GetSetting() error: %v", err)
	}
	if val != "light" {
		t.Fatalf("expected 'light' after upsert, got %q", val)
	}
}

func TestListSettings_Empty(t *testing.T) {
	s := newTestStore(t)
	defer s.Close()

	m, err := s.ListSettings()
	if err != nil {
		t.Fatalf("ListSettings() error: %v", err)
	}
	if len(m) != 0 {
		t.Fatalf("expected empty map, got %d entries", len(m))
	}
}

func TestListSettings_MultipleKeys(t *testing.T) {
	s := newTestStore(t)
	defer s.Close()

	if err := s.SetSetting("theme", "dark"); err != nil {
		t.Fatal(err)
	}
	if err := s.SetSetting("language", "en"); err != nil {
		t.Fatal(err)
	}

	m, err := s.ListSettings()
	if err != nil {
		t.Fatalf("ListSettings() error: %v", err)
	}
	if len(m) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(m))
	}
	if m["theme"] != "dark" {
		t.Fatalf("expected theme=dark, got %q", m["theme"])
	}
	if m["language"] != "en" {
		t.Fatalf("expected language=en, got %q", m["language"])
	}
}

func TestSetSettings_Bulk(t *testing.T) {
	s := newTestStore(t)
	defer s.Close()

	batch := map[string]string{
		"theme":    "dark",
		"language": "zh",
		"text_size": "large",
	}
	if err := s.SetSettings(batch); err != nil {
		t.Fatalf("SetSettings() error: %v", err)
	}

	m, err := s.ListSettings()
	if err != nil {
		t.Fatalf("ListSettings() error: %v", err)
	}
	if len(m) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(m))
	}
	for k, v := range batch {
		if m[k] != v {
			t.Fatalf("expected %s=%q, got %q", k, v, m[k])
		}
	}
}

func TestSetSettings_UpsertExisting(t *testing.T) {
	s := newTestStore(t)
	defer s.Close()

	if err := s.SetSetting("theme", "dark"); err != nil {
		t.Fatal(err)
	}

	batch := map[string]string{
		"theme":    "light",
		"language": "en",
	}
	if err := s.SetSettings(batch); err != nil {
		t.Fatalf("SetSettings() upsert error: %v", err)
	}

	val, err := s.GetSetting("theme")
	if err != nil {
		t.Fatal(err)
	}
	if val != "light" {
		t.Fatalf("expected theme=light after bulk upsert, got %q", val)
	}
}

func TestDeleteSetting(t *testing.T) {
	s := newTestStore(t)
	defer s.Close()

	if err := s.SetSetting("theme", "dark"); err != nil {
		t.Fatal(err)
	}
	if err := s.DeleteSetting("theme"); err != nil {
		t.Fatalf("DeleteSetting() error: %v", err)
	}

	_, err := s.GetSetting("theme")
	if !errors.Is(err, sql.ErrNoRows) {
		t.Fatalf("expected sql.ErrNoRows after delete, got: %v", err)
	}
}

func TestDeleteSetting_NotFound(t *testing.T) {
	s := newTestStore(t)
	defer s.Close()

	err := s.DeleteSetting("nonexistent")
	if !errors.Is(err, sql.ErrNoRows) {
		t.Fatalf("expected sql.ErrNoRows, got: %v", err)
	}
}
