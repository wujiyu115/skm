package store

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNew_CreatesDatabase(t *testing.T) {
	dir := t.TempDir()
	dbPath := filepath.Join(dir, "test.db")

	s, err := New(dbPath)
	if err != nil {
		t.Fatalf("New() error: %v", err)
	}
	defer s.Close()

	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		t.Fatal("database file not created")
	}
}

func TestNew_RunsMigrations(t *testing.T) {
	dir := t.TempDir()
	s, err := New(filepath.Join(dir, "test.db"))
	if err != nil {
		t.Fatalf("New() error: %v", err)
	}
	defer s.Close()

	var version int
	err = s.db.QueryRow("PRAGMA user_version").Scan(&version)
	if err != nil {
		t.Fatalf("PRAGMA user_version error: %v", err)
	}
	if version != 3 {
		t.Fatalf("expected user_version=3, got %d", version)
	}
}

func TestNew_CreatesSkillsTable(t *testing.T) {
	dir := t.TempDir()
	s, err := New(filepath.Join(dir, "test.db"))
	if err != nil {
		t.Fatalf("New() error: %v", err)
	}
	defer s.Close()

	_, err = s.db.Exec("SELECT id, name, description, source_type, source_ref, central_path, content_hash, enabled, installed_at, updated_at FROM skills LIMIT 0")
	if err != nil {
		t.Fatalf("skills table missing or wrong schema: %v", err)
	}
}

func TestNew_MigrationV2_SkillColumns(t *testing.T) {
	dir := t.TempDir()
	s, err := New(filepath.Join(dir, "test.db"))
	if err != nil {
		t.Fatalf("New() error: %v", err)
	}
	defer s.Close()

	// Verify the new columns exist on the skills table
	_, err = s.db.Exec("SELECT source_branch, source_subpath, remote_revision, update_status, last_checked_at FROM skills LIMIT 0")
	if err != nil {
		t.Fatalf("v2 skill columns missing: %v", err)
	}
}

func TestNew_MigrationV2_SkillTagsTable(t *testing.T) {
	dir := t.TempDir()
	s, err := New(filepath.Join(dir, "test.db"))
	if err != nil {
		t.Fatalf("New() error: %v", err)
	}
	defer s.Close()

	_, err = s.db.Exec("SELECT skill_id, tag FROM skill_tags LIMIT 0")
	if err != nil {
		t.Fatalf("skill_tags table missing or wrong schema: %v", err)
	}
}

func TestNew_MigrationV2_AuditLogTable(t *testing.T) {
	dir := t.TempDir()
	s, err := New(filepath.Join(dir, "test.db"))
	if err != nil {
		t.Fatalf("New() error: %v", err)
	}
	defer s.Close()

	_, err = s.db.Exec("SELECT id, action, target, detail, created_at FROM audit_log LIMIT 0")
	if err != nil {
		t.Fatalf("audit_log table missing or wrong schema: %v", err)
	}
}

func TestNew_WALMode(t *testing.T) {
	dir := t.TempDir()
	s, err := New(filepath.Join(dir, "test.db"))
	if err != nil {
		t.Fatalf("New() error: %v", err)
	}
	defer s.Close()

	var mode string
	err = s.db.QueryRow("PRAGMA journal_mode").Scan(&mode)
	if err != nil {
		t.Fatalf("PRAGMA journal_mode error: %v", err)
	}
	if mode != "wal" {
		t.Fatalf("expected journal_mode=wal, got %s", mode)
	}
}
