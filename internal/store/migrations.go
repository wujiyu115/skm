package store

import "fmt"

const currentVersion = 1

func (s *Store) migrate() error {
	var version int
	if err := s.db.QueryRow("PRAGMA user_version").Scan(&version); err != nil {
		return fmt.Errorf("read user_version: %w", err)
	}

	if version < 1 {
		if err := s.migrateV1(); err != nil {
			return err
		}
	}

	return nil
}

func (s *Store) migrateV1() error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmts := []string{
		`CREATE TABLE IF NOT EXISTS skills (
			id           TEXT PRIMARY KEY,
			name         TEXT UNIQUE NOT NULL,
			description  TEXT,
			source_type  TEXT NOT NULL,
			source_ref   TEXT,
			central_path TEXT UNIQUE NOT NULL,
			content_hash TEXT,
			enabled      INTEGER DEFAULT 1,
			installed_at TEXT DEFAULT (datetime('now')),
			updated_at   TEXT DEFAULT (datetime('now'))
		)`,
		`CREATE TABLE IF NOT EXISTS groups (
			id          TEXT PRIMARY KEY,
			name        TEXT UNIQUE NOT NULL,
			description TEXT,
			created_at  TEXT DEFAULT (datetime('now'))
		)`,
		`CREATE TABLE IF NOT EXISTS group_skills (
			group_id   TEXT NOT NULL REFERENCES groups(id) ON DELETE CASCADE,
			skill_id   TEXT NOT NULL REFERENCES skills(id) ON DELETE CASCADE,
			sort_order INTEGER DEFAULT 0,
			PRIMARY KEY (group_id, skill_id)
		)`,
		`CREATE TABLE IF NOT EXISTS skill_targets (
			skill_id    TEXT NOT NULL REFERENCES skills(id) ON DELETE CASCADE,
			agent       TEXT NOT NULL,
			target_path TEXT NOT NULL,
			mode        TEXT DEFAULT 'symlink',
			source_hash TEXT,
			synced_at   TEXT DEFAULT (datetime('now')),
			UNIQUE (skill_id, agent)
		)`,
		`CREATE TABLE IF NOT EXISTS agents (
			name         TEXT PRIMARY KEY,
			display_name TEXT NOT NULL,
			project_dir  TEXT NOT NULL,
			global_dir   TEXT NOT NULL,
			detect_paths TEXT,
			is_builtin   INTEGER DEFAULT 0,
			enabled      INTEGER DEFAULT 1
		)`,
		`CREATE TABLE IF NOT EXISTS settings (
			key   TEXT PRIMARY KEY,
			value TEXT
		)`,
		`PRAGMA user_version = 1`,
	}

	for _, stmt := range stmts {
		if _, err := tx.Exec(stmt); err != nil {
			return fmt.Errorf("migration v1: %w", err)
		}
	}

	return tx.Commit()
}
