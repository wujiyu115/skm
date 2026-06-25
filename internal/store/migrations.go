package store

import "fmt"

const currentVersion = 4

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

	if version < 2 {
		if err := s.migrateV2(); err != nil {
			return err
		}
	}

	if version < 3 {
		if err := s.migrateV3(); err != nil {
			return err
		}
	}

	if version < 4 {
		if err := s.migrateV4(); err != nil {
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
	}

	for _, stmt := range stmts {
		if _, err := tx.Exec(stmt); err != nil {
			return fmt.Errorf("migration v1: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("migration v1 commit: %w", err)
	}

	// PRAGMA is not transactional in SQLite; execute after commit.
	if _, err := s.db.Exec("PRAGMA user_version = 1"); err != nil {
		return fmt.Errorf("set user_version: %w", err)
	}

	return nil
}

func (s *Store) migrateV2() error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmts := []string{
		// Add columns to skills
		`ALTER TABLE skills ADD COLUMN source_branch TEXT DEFAULT ''`,
		`ALTER TABLE skills ADD COLUMN source_subpath TEXT DEFAULT ''`,
		`ALTER TABLE skills ADD COLUMN remote_revision TEXT DEFAULT ''`,
		`ALTER TABLE skills ADD COLUMN update_status TEXT DEFAULT 'unknown'`,
		`ALTER TABLE skills ADD COLUMN last_checked_at TEXT`,

		// Skill tags
		`CREATE TABLE IF NOT EXISTS skill_tags (
			skill_id TEXT NOT NULL REFERENCES skills(id) ON DELETE CASCADE,
			tag      TEXT NOT NULL,
			PRIMARY KEY (skill_id, tag)
		)`,
		`CREATE INDEX IF NOT EXISTS idx_skill_tags_tag ON skill_tags(tag)`,

		// Audit log
		`CREATE TABLE IF NOT EXISTS audit_log (
			id         INTEGER PRIMARY KEY AUTOINCREMENT,
			action     TEXT NOT NULL,
			target     TEXT,
			detail     TEXT,
			created_at TEXT DEFAULT (datetime('now'))
		)`,
		`CREATE INDEX IF NOT EXISTS idx_audit_log_created ON audit_log(created_at)`,
	}

	for _, stmt := range stmts {
		if _, err := tx.Exec(stmt); err != nil {
			return fmt.Errorf("migration v2: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("migration v2 commit: %w", err)
	}

	if _, err := s.db.Exec("PRAGMA user_version = 2"); err != nil {
		return fmt.Errorf("set user_version: %w", err)
	}

	return nil
}

func (s *Store) migrateV3() error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmts := []string{
		`ALTER TABLE agents ADD COLUMN category TEXT DEFAULT ''`,
	}

	for _, stmt := range stmts {
		if _, err := tx.Exec(stmt); err != nil {
			return fmt.Errorf("migration v3: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("migration v3 commit: %w", err)
	}

	if _, err := s.db.Exec("PRAGMA user_version = 3"); err != nil {
		return fmt.Errorf("set user_version: %w", err)
	}

	return nil
}

func (s *Store) migrateV4() error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmts := []string{
		`CREATE TABLE IF NOT EXISTS projects (
			id         TEXT PRIMARY KEY,
			name       TEXT NOT NULL,
			path       TEXT UNIQUE NOT NULL,
			created_at TEXT DEFAULT (datetime('now'))
		)`,
	}

	for _, stmt := range stmts {
		if _, err := tx.Exec(stmt); err != nil {
			return fmt.Errorf("migration v4: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("migration v4 commit: %w", err)
	}

	if _, err := s.db.Exec("PRAGMA user_version = 4"); err != nil {
		return fmt.Errorf("set user_version: %w", err)
	}

	return nil
}
