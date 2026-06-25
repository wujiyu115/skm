package store

import "fmt"

// AddTag associates a tag with a skill. Duplicate tags are silently ignored.
func (s *Store) AddTag(skillID, tag string) error {
	_, err := s.db.Exec(
		"INSERT OR IGNORE INTO skill_tags (skill_id, tag) VALUES (?, ?)",
		skillID, tag,
	)
	if err != nil {
		return fmt.Errorf("add tag: %w", err)
	}
	return nil
}

// RemoveTag removes a specific tag from a skill.
func (s *Store) RemoveTag(skillID, tag string) error {
	_, err := s.db.Exec(
		"DELETE FROM skill_tags WHERE skill_id = ? AND tag = ?",
		skillID, tag,
	)
	if err != nil {
		return fmt.Errorf("remove tag: %w", err)
	}
	return nil
}

// SetSkillTags replaces all tags for a skill with the given set.
func (s *Store) SetSkillTags(skillID string, tags []string) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err := tx.Exec("DELETE FROM skill_tags WHERE skill_id = ?", skillID); err != nil {
		return fmt.Errorf("set skill tags (clear): %w", err)
	}

	for _, tag := range tags {
		if _, err := tx.Exec("INSERT INTO skill_tags (skill_id, tag) VALUES (?, ?)", skillID, tag); err != nil {
			return fmt.Errorf("set skill tags (insert %q): %w", tag, err)
		}
	}

	return tx.Commit()
}

// ListTags returns all distinct tags across all skills, ordered alphabetically.
func (s *Store) ListTags() ([]string, error) {
	rows, err := s.db.Query("SELECT DISTINCT tag FROM skill_tags ORDER BY tag")
	if err != nil {
		return nil, fmt.Errorf("list tags: %w", err)
	}
	defer rows.Close()

	var tags []string
	for rows.Next() {
		var tag string
		if err := rows.Scan(&tag); err != nil {
			return nil, err
		}
		tags = append(tags, tag)
	}
	return tags, rows.Err()
}

// ListSkillTags returns all tags for a specific skill, ordered alphabetically.
func (s *Store) ListSkillTags(skillID string) ([]string, error) {
	rows, err := s.db.Query(
		"SELECT tag FROM skill_tags WHERE skill_id = ? ORDER BY tag",
		skillID,
	)
	if err != nil {
		return nil, fmt.Errorf("list skill tags: %w", err)
	}
	defer rows.Close()

	var tags []string
	for rows.Next() {
		var tag string
		if err := rows.Scan(&tag); err != nil {
			return nil, err
		}
		tags = append(tags, tag)
	}
	return tags, rows.Err()
}

// RenameTag renames a tag across all skills.
// If a skill already has both oldTag and newTag, the conflicting row is
// silently ignored by UPDATE OR IGNORE, then the leftover old rows are
// cleaned up by a DELETE.
func (s *Store) RenameTag(oldTag, newTag string) error {
	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("rename tag: %w", err)
	}
	defer tx.Rollback()

	// Rename where there is no PK conflict; skip conflicts silently.
	if _, err := tx.Exec(
		"UPDATE OR IGNORE skill_tags SET tag = ? WHERE tag = ?",
		newTag, oldTag,
	); err != nil {
		return fmt.Errorf("rename tag (update): %w", err)
	}

	// Remove any remaining old-tag rows (these were conflicts).
	if _, err := tx.Exec(
		"DELETE FROM skill_tags WHERE tag = ?", oldTag,
	); err != nil {
		return fmt.Errorf("rename tag (cleanup): %w", err)
	}

	return tx.Commit()
}

// DeleteTag removes a tag from all skills.
func (s *Store) DeleteTag(tag string) error {
	_, err := s.db.Exec("DELETE FROM skill_tags WHERE tag = ?", tag)
	if err != nil {
		return fmt.Errorf("delete tag: %w", err)
	}
	return nil
}
