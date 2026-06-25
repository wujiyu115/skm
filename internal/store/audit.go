package store

import "fmt"

// AuditEntry represents a single audit log record.
type AuditEntry struct {
	ID        int64  `json:"id"`
	Action    string `json:"action"`
	Target    string `json:"target"`
	Detail    string `json:"detail"`
	CreatedAt string `json:"created_at"`
}

// InsertAuditLog appends an entry to the audit log.
func (s *Store) InsertAuditLog(action, target, detail string) error {
	_, err := s.db.Exec(
		"INSERT INTO audit_log (action, target, detail) VALUES (?, ?, ?)",
		action, target, detail,
	)
	if err != nil {
		return fmt.Errorf("insert audit log: %w", err)
	}
	return nil
}

// ListAuditLog returns the most recent audit entries, newest first.
// A limit of 0 or less returns all entries.
func (s *Store) ListAuditLog(limit int) ([]AuditEntry, error) {
	var query string
	var args []interface{}

	if limit > 0 {
		query = "SELECT id, action, target, detail, created_at FROM audit_log ORDER BY id DESC LIMIT ?"
		args = append(args, limit)
	} else {
		query = "SELECT id, action, target, detail, created_at FROM audit_log ORDER BY id DESC"
	}

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("list audit log: %w", err)
	}
	defer rows.Close()

	var entries []AuditEntry
	for rows.Next() {
		var e AuditEntry
		if err := rows.Scan(&e.ID, &e.Action, &e.Target, &e.Detail, &e.CreatedAt); err != nil {
			return nil, err
		}
		entries = append(entries, e)
	}
	return entries, rows.Err()
}

// PruneAuditLog keeps only the most recent maxRows entries, deleting older ones.
func (s *Store) PruneAuditLog(maxRows int) error {
	_, err := s.db.Exec(
		`DELETE FROM audit_log WHERE id NOT IN (
			SELECT id FROM audit_log ORDER BY id DESC LIMIT ?
		)`, maxRows,
	)
	if err != nil {
		return fmt.Errorf("prune audit log: %w", err)
	}
	return nil
}
