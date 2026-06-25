package store

import (
	"database/sql"
	"fmt"
)

// GetSetting returns the value for a single settings key.
// Returns ("", sql.ErrNoRows) if the key does not exist.
func (s *Store) GetSetting(key string) (string, error) {
	var value string
	err := s.db.QueryRow("SELECT value FROM settings WHERE key = ?", key).Scan(&value)
	if err != nil {
		return "", fmt.Errorf("get setting %q: %w", key, err)
	}
	return value, nil
}

// SetSetting upserts a single key-value pair in the settings table.
func (s *Store) SetSetting(key, value string) error {
	_, err := s.db.Exec(
		"INSERT OR REPLACE INTO settings (key, value) VALUES (?, ?)",
		key, value,
	)
	if err != nil {
		return fmt.Errorf("set setting %q: %w", key, err)
	}
	return nil
}

// ListSettings returns all settings as a key-value map.
func (s *Store) ListSettings() (map[string]string, error) {
	rows, err := s.db.Query("SELECT key, value FROM settings ORDER BY key")
	if err != nil {
		return nil, fmt.Errorf("list settings: %w", err)
	}
	defer rows.Close()

	result := make(map[string]string)
	for rows.Next() {
		var k, v string
		if err := rows.Scan(&k, &v); err != nil {
			return nil, fmt.Errorf("scan setting row: %w", err)
		}
		result[k] = v
	}
	return result, rows.Err()
}

// SetSettings upserts multiple key-value pairs in a single transaction.
func (s *Store) SetSettings(settings map[string]string) error {
	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare("INSERT OR REPLACE INTO settings (key, value) VALUES (?, ?)")
	if err != nil {
		return fmt.Errorf("prepare upsert: %w", err)
	}
	defer stmt.Close()

	for k, v := range settings {
		if _, err := stmt.Exec(k, v); err != nil {
			return fmt.Errorf("set setting %q: %w", k, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit settings: %w", err)
	}
	return nil
}

// DeleteSetting removes a single key from the settings table.
func (s *Store) DeleteSetting(key string) error {
	res, err := s.db.Exec("DELETE FROM settings WHERE key = ?", key)
	if err != nil {
		return fmt.Errorf("delete setting %q: %w", key, err)
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return sql.ErrNoRows
	}
	return nil
}
