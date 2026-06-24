package store

import (
	"database/sql"
	"fmt"

	"github.com/ejoy/skm/internal/skill"
)

// Group represents a named collection of skills.
type Group struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// InsertGroup adds a new group record to the database.
func (s *Store) InsertGroup(id, name, description string) error {
	_, err := s.db.Exec(
		"INSERT INTO groups (id, name, description) VALUES (?, ?, ?)",
		id, name, description,
	)
	return err
}

// GetGroup retrieves a group by name.
func (s *Store) GetGroup(name string) (*Group, error) {
	g := &Group{}
	err := s.db.QueryRow(
		"SELECT id, name, description FROM groups WHERE name = ?", name,
	).Scan(&g.ID, &g.Name, &g.Description)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("group %q not found", name)
	}
	return g, err
}

// GetGroupByID retrieves a group by its unique ID.
func (s *Store) GetGroupByID(id string) (*Group, error) {
	g := &Group{}
	err := s.db.QueryRow(
		"SELECT id, name, description FROM groups WHERE id = ?", id,
	).Scan(&g.ID, &g.Name, &g.Description)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("group not found")
	}
	return g, err
}

// ListGroups returns all groups ordered by name.
func (s *Store) ListGroups() ([]Group, error) {
	rows, err := s.db.Query("SELECT id, name, description FROM groups ORDER BY name")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var groups []Group
	for rows.Next() {
		var g Group
		if err := rows.Scan(&g.ID, &g.Name, &g.Description); err != nil {
			return nil, err
		}
		groups = append(groups, g)
	}
	return groups, rows.Err()
}

// DeleteGroup removes a group by ID. Cascade deletes group_skills via FK.
func (s *Store) DeleteGroup(id string) error {
	_, err := s.db.Exec("DELETE FROM groups WHERE id = ?", id)
	return err
}

// UpdateGroup updates a group's name and description.
func (s *Store) UpdateGroup(id, name, description string) error {
	_, err := s.db.Exec("UPDATE groups SET name = ?, description = ? WHERE id = ?", name, description, id)
	return err
}

// AddSkillToGroup associates a skill with a group at the given sort order.
func (s *Store) AddSkillToGroup(groupID, skillID string, sortOrder int) error {
	_, err := s.db.Exec(
		"INSERT OR IGNORE INTO group_skills (group_id, skill_id, sort_order) VALUES (?, ?, ?)",
		groupID, skillID, sortOrder,
	)
	return err
}

// RemoveSkillFromGroup removes the association between a skill and a group.
func (s *Store) RemoveSkillFromGroup(groupID, skillID string) error {
	_, err := s.db.Exec("DELETE FROM group_skills WHERE group_id = ? AND skill_id = ?", groupID, skillID)
	return err
}

// ListGroupSkills returns all skills in a group, ordered by sort_order then name.
func (s *Store) ListGroupSkills(groupID string) ([]*skill.Skill, error) {
	rows, err := s.db.Query(
		`SELECT s.id, s.name, s.description, s.source_type, s.source_ref, s.central_path, s.content_hash, s.enabled
		 FROM skills s
		 JOIN group_skills gs ON gs.skill_id = s.id
		 WHERE gs.group_id = ?
		 ORDER BY gs.sort_order, s.name`, groupID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var skills []*skill.Skill
	for rows.Next() {
		sk := &skill.Skill{}
		if err := rows.Scan(&sk.ID, &sk.Name, &sk.Description, &sk.SourceType, &sk.SourceRef, &sk.CentralPath, &sk.ContentHash, &sk.Enabled); err != nil {
			return nil, err
		}
		skills = append(skills, sk)
	}
	return skills, rows.Err()
}
