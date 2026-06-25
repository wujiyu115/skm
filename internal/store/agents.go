package store

import (
	"encoding/json"

	"github.com/ejoy/skm/internal/agent"
)

// SeedBuiltinAgents inserts built-in agent adapters, ignoring duplicates.
func (s *Store) SeedBuiltinAgents(agents []agent.Adapter) error {
	for _, a := range agents {
		detectJSON, _ := json.Marshal(a.DetectPaths)
		_, err := s.db.Exec(
			`INSERT OR IGNORE INTO agents (name, display_name, project_dir, global_dir, detect_paths, is_builtin, enabled, category)
			 VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
			a.Name, a.DisplayName, a.ProjectDir, a.GlobalDir, string(detectJSON), 1, 1, a.Category,
		)
		if err != nil {
			return err
		}
	}
	// Update category for existing agents that predate v3
	for _, a := range agents {
		s.db.Exec("UPDATE agents SET category = ? WHERE name = ? AND category = ''", a.Category, a.Name)
	}
	return nil
}

// ListAgents returns all enabled agents ordered by name.
func (s *Store) ListAgents() ([]agent.Adapter, error) {
	rows, err := s.db.Query(
		"SELECT name, display_name, project_dir, global_dir, detect_paths, is_builtin, COALESCE(category,'') FROM agents WHERE enabled = 1 ORDER BY name",
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var agents []agent.Adapter
	for rows.Next() {
		var a agent.Adapter
		var detectJSON string
		if err := rows.Scan(&a.Name, &a.DisplayName, &a.ProjectDir, &a.GlobalDir, &detectJSON, &a.IsBuiltin, &a.Category); err != nil {
			return nil, err
		}
		json.Unmarshal([]byte(detectJSON), &a.DetectPaths)
		agents = append(agents, a)
	}
	return agents, rows.Err()
}

// InsertAgent adds a custom (non-builtin) agent.
func (s *Store) InsertAgent(a agent.Adapter) error {
	detectJSON, _ := json.Marshal(a.DetectPaths)
	_, err := s.db.Exec(
		`INSERT INTO agents (name, display_name, project_dir, global_dir, detect_paths, is_builtin, enabled, category)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		a.Name, a.DisplayName, a.ProjectDir, a.GlobalDir, string(detectJSON), 0, 1, a.Category,
	)
	return err
}

// DeleteAgent removes a custom agent. Built-in agents are protected from deletion.
func (s *Store) DeleteAgent(name string) error {
	_, err := s.db.Exec("DELETE FROM agents WHERE name = ? AND is_builtin = 0", name)
	return err
}
