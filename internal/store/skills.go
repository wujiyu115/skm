package store

import (
	"database/sql"
	"fmt"

	"github.com/ejoy/skm/internal/skill"
)

// Target represents a deployed skill in an agent's directory.
type Target struct {
	SkillID    string `json:"skill_id"`
	Agent      string `json:"agent"`
	TargetPath string `json:"target_path"`
	Mode       string `json:"mode"`
	SourceHash string `json:"source_hash"`
}

// InsertSkill adds a new skill record to the database.
func (s *Store) InsertSkill(sk *skill.Skill) error {
	_, err := s.db.Exec(
		`INSERT INTO skills (id, name, description, source_type, source_ref, central_path, content_hash, enabled)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		sk.ID, sk.Name, sk.Description, sk.SourceType, sk.SourceRef, sk.CentralPath, sk.ContentHash, sk.Enabled,
	)
	if err != nil {
		return fmt.Errorf("insert skill: %w", err)
	}
	return nil
}

// GetSkill retrieves a skill by name.
func (s *Store) GetSkill(name string) (*skill.Skill, error) {
	sk := &skill.Skill{}
	err := s.db.QueryRow(
		`SELECT id, name, description, source_type, source_ref, central_path, content_hash, enabled
		 FROM skills WHERE name = ?`, name,
	).Scan(&sk.ID, &sk.Name, &sk.Description, &sk.SourceType, &sk.SourceRef, &sk.CentralPath, &sk.ContentHash, &sk.Enabled)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("skill %q not found", name)
	}
	if err != nil {
		return nil, fmt.Errorf("get skill: %w", err)
	}
	return sk, nil
}

// GetSkillByID retrieves a skill by its unique ID.
func (s *Store) GetSkillByID(id string) (*skill.Skill, error) {
	sk := &skill.Skill{}
	err := s.db.QueryRow(
		`SELECT id, name, description, source_type, source_ref, central_path, content_hash, enabled
		 FROM skills WHERE id = ?`, id,
	).Scan(&sk.ID, &sk.Name, &sk.Description, &sk.SourceType, &sk.SourceRef, &sk.CentralPath, &sk.ContentHash, &sk.Enabled)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("skill with id %q not found", id)
	}
	if err != nil {
		return nil, fmt.Errorf("get skill by id: %w", err)
	}
	return sk, nil
}

// ListSkills returns all skills ordered by name.
func (s *Store) ListSkills() ([]*skill.Skill, error) {
	rows, err := s.db.Query(
		`SELECT id, name, description, source_type, source_ref, central_path, content_hash, enabled
		 FROM skills ORDER BY name`,
	)
	if err != nil {
		return nil, fmt.Errorf("list skills: %w", err)
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

// DeleteSkill removes a skill by ID. Returns an error if not found.
func (s *Store) DeleteSkill(id string) error {
	res, err := s.db.Exec("DELETE FROM skills WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("delete skill: %w", err)
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return fmt.Errorf("skill with id %q not found", id)
	}
	return nil
}

// DeleteSkillByName removes a skill by name. Returns an error if not found.
func (s *Store) DeleteSkillByName(name string) error {
	res, err := s.db.Exec("DELETE FROM skills WHERE name = ?", name)
	if err != nil {
		return fmt.Errorf("delete skill by name: %w", err)
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return fmt.Errorf("skill %q not found", name)
	}
	return nil
}

// UpdateSkillHash updates the content hash and updated_at timestamp for a skill.
func (s *Store) UpdateSkillHash(id, hash string) error {
	_, err := s.db.Exec(
		"UPDATE skills SET content_hash = ?, updated_at = datetime('now') WHERE id = ?",
		hash, id,
	)
	return err
}

// UpsertTarget inserts or updates a skill deployment target.
func (s *Store) UpsertTarget(skillID, agent, targetPath, mode, hash string) error {
	_, err := s.db.Exec(
		`INSERT INTO skill_targets (skill_id, agent, target_path, mode, source_hash)
		 VALUES (?, ?, ?, ?, ?)
		 ON CONFLICT(skill_id, agent) DO UPDATE SET
		   target_path = excluded.target_path,
		   mode = excluded.mode,
		   source_hash = excluded.source_hash,
		   synced_at = datetime('now')`,
		skillID, agent, targetPath, mode, hash,
	)
	return err
}

// ListTargets returns all deployment targets for a skill.
func (s *Store) ListTargets(skillID string) ([]Target, error) {
	rows, err := s.db.Query(
		"SELECT skill_id, agent, target_path, mode, source_hash FROM skill_targets WHERE skill_id = ?",
		skillID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var targets []Target
	for rows.Next() {
		var t Target
		if err := rows.Scan(&t.SkillID, &t.Agent, &t.TargetPath, &t.Mode, &t.SourceHash); err != nil {
			return nil, err
		}
		targets = append(targets, t)
	}
	return targets, rows.Err()
}

// DeleteTarget removes a specific skill-agent deployment target.
func (s *Store) DeleteTarget(skillID, agent string) error {
	_, err := s.db.Exec("DELETE FROM skill_targets WHERE skill_id = ? AND agent = ?", skillID, agent)
	return err
}
