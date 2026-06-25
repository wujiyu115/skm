package store

import (
	"database/sql"
	"fmt"
)

// Project represents a workspace directory tracked by SKM.
type Project struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Path      string `json:"path"`
	CreatedAt string `json:"created_at"`
}

// InsertProject adds a new project record to the database.
func (s *Store) InsertProject(id, name, path string) error {
	_, err := s.db.Exec(
		"INSERT INTO projects (id, name, path) VALUES (?, ?, ?)",
		id, name, path,
	)
	return err
}

// GetProjectByID retrieves a project by its unique ID.
func (s *Store) GetProjectByID(id string) (*Project, error) {
	p := &Project{}
	err := s.db.QueryRow(
		"SELECT id, name, path, created_at FROM projects WHERE id = ?", id,
	).Scan(&p.ID, &p.Name, &p.Path, &p.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("project not found")
	}
	return p, err
}

// ListProjects returns all projects ordered by name.
func (s *Store) ListProjects() ([]Project, error) {
	rows, err := s.db.Query("SELECT id, name, path, created_at FROM projects ORDER BY name")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projects []Project
	for rows.Next() {
		var p Project
		if err := rows.Scan(&p.ID, &p.Name, &p.Path, &p.CreatedAt); err != nil {
			return nil, err
		}
		projects = append(projects, p)
	}
	return projects, rows.Err()
}

// DeleteProject removes a project by ID.
func (s *Store) DeleteProject(id string) error {
	_, err := s.db.Exec("DELETE FROM projects WHERE id = ?", id)
	return err
}
