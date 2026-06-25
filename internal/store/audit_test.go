package store

import (
	"path/filepath"
	"testing"
)

func newTestStore(t *testing.T) *Store {
	t.Helper()
	dir := t.TempDir()
	s, err := New(filepath.Join(dir, "test.db"))
	if err != nil {
		t.Fatalf("New() error: %v", err)
	}
	return s
}

func TestInsertAuditLog(t *testing.T) {
	s := newTestStore(t)
	defer s.Close()

	if err := s.InsertAuditLog("install", "test-skill", "installed from /tmp"); err != nil {
		t.Fatalf("InsertAuditLog() error: %v", err)
	}

	entries, err := s.ListAuditLog(10)
	if err != nil {
		t.Fatalf("ListAuditLog() error: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	e := entries[0]
	if e.Action != "install" {
		t.Fatalf("action = %q, want %q", e.Action, "install")
	}
	if e.Target != "test-skill" {
		t.Fatalf("target = %q, want %q", e.Target, "test-skill")
	}
	if e.Detail != "installed from /tmp" {
		t.Fatalf("detail = %q, want %q", e.Detail, "installed from /tmp")
	}
	if e.CreatedAt == "" {
		t.Fatal("created_at should not be empty")
	}
}

func TestListAuditLog_Order(t *testing.T) {
	s := newTestStore(t)
	defer s.Close()

	s.InsertAuditLog("action-1", "t1", "first")
	s.InsertAuditLog("action-2", "t2", "second")
	s.InsertAuditLog("action-3", "t3", "third")

	entries, err := s.ListAuditLog(0)
	if err != nil {
		t.Fatalf("ListAuditLog() error: %v", err)
	}
	if len(entries) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(entries))
	}
	// Newest first
	if entries[0].Action != "action-3" {
		t.Fatalf("first entry should be newest, got %q", entries[0].Action)
	}
	if entries[2].Action != "action-1" {
		t.Fatalf("last entry should be oldest, got %q", entries[2].Action)
	}
}

func TestListAuditLog_Limit(t *testing.T) {
	s := newTestStore(t)
	defer s.Close()

	for i := 0; i < 10; i++ {
		s.InsertAuditLog("action", "target", "detail")
	}

	entries, err := s.ListAuditLog(3)
	if err != nil {
		t.Fatalf("ListAuditLog() error: %v", err)
	}
	if len(entries) != 3 {
		t.Fatalf("expected 3 entries with limit=3, got %d", len(entries))
	}
}

func TestListAuditLog_Empty(t *testing.T) {
	s := newTestStore(t)
	defer s.Close()

	entries, err := s.ListAuditLog(10)
	if err != nil {
		t.Fatalf("ListAuditLog() error: %v", err)
	}
	if len(entries) != 0 {
		t.Fatalf("expected 0 entries, got %d", len(entries))
	}
}

func TestPruneAuditLog(t *testing.T) {
	s := newTestStore(t)
	defer s.Close()

	for i := 0; i < 10; i++ {
		s.InsertAuditLog("action", "target", "detail")
	}

	if err := s.PruneAuditLog(3); err != nil {
		t.Fatalf("PruneAuditLog() error: %v", err)
	}

	entries, err := s.ListAuditLog(0)
	if err != nil {
		t.Fatalf("ListAuditLog() error: %v", err)
	}
	if len(entries) != 3 {
		t.Fatalf("expected 3 entries after prune, got %d", len(entries))
	}
}

func TestPruneAuditLog_KeepsNewest(t *testing.T) {
	s := newTestStore(t)
	defer s.Close()

	s.InsertAuditLog("old", "t", "d")
	s.InsertAuditLog("middle", "t", "d")
	s.InsertAuditLog("newest", "t", "d")

	if err := s.PruneAuditLog(1); err != nil {
		t.Fatalf("PruneAuditLog() error: %v", err)
	}

	entries, _ := s.ListAuditLog(0)
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	if entries[0].Action != "newest" {
		t.Fatalf("expected newest entry to survive prune, got %q", entries[0].Action)
	}
}
