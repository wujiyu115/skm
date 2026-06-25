package store

import "testing"

func TestInsertAndGetProject(t *testing.T) {
	s := testStore(t)

	if err := s.InsertProject("proj-1", "My Project", "/home/user/my-project"); err != nil {
		t.Fatalf("InsertProject() error: %v", err)
	}

	got, err := s.GetProjectByID("proj-1")
	if err != nil {
		t.Fatalf("GetProjectByID() error: %v", err)
	}
	if got.ID != "proj-1" {
		t.Errorf("ID = %q, want %q", got.ID, "proj-1")
	}
	if got.Name != "My Project" {
		t.Errorf("Name = %q, want %q", got.Name, "My Project")
	}
	if got.Path != "/home/user/my-project" {
		t.Errorf("Path = %q, want %q", got.Path, "/home/user/my-project")
	}
	if got.CreatedAt == "" {
		t.Error("CreatedAt is empty, want non-empty")
	}
}

func TestInsertProject_DuplicatePath(t *testing.T) {
	s := testStore(t)

	if err := s.InsertProject("proj-1", "First", "/same/path"); err != nil {
		t.Fatalf("InsertProject() error: %v", err)
	}
	if err := s.InsertProject("proj-2", "Second", "/same/path"); err == nil {
		t.Fatal("expected error on duplicate path insert")
	}
}

func TestListProjects(t *testing.T) {
	s := testStore(t)

	s.InsertProject("p1", "Beta", "/beta")
	s.InsertProject("p2", "Alpha", "/alpha")

	projects, err := s.ListProjects()
	if err != nil {
		t.Fatalf("ListProjects() error: %v", err)
	}
	if len(projects) != 2 {
		t.Fatalf("expected 2 projects, got %d", len(projects))
	}
	// Should be ordered by name
	if projects[0].Name != "Alpha" || projects[1].Name != "Beta" {
		t.Errorf("unexpected order: %q, %q", projects[0].Name, projects[1].Name)
	}
}

func TestListProjects_Empty(t *testing.T) {
	s := testStore(t)

	projects, err := s.ListProjects()
	if err != nil {
		t.Fatalf("ListProjects() error: %v", err)
	}
	if len(projects) != 0 {
		t.Fatalf("expected 0 projects, got %d", len(projects))
	}
}

func TestDeleteProject(t *testing.T) {
	s := testStore(t)

	s.InsertProject("del-1", "Delete Me", "/delete/me")

	if err := s.DeleteProject("del-1"); err != nil {
		t.Fatalf("DeleteProject() error: %v", err)
	}

	_, err := s.GetProjectByID("del-1")
	if err == nil {
		t.Fatal("expected error after delete")
	}
}
