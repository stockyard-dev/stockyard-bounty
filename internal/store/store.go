package store

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"time"

	_ "modernc.org/sqlite"
)

type DB struct { db *sql.DB }

type Issue struct {
	ID           string   `json:"id"`
	Title        string   `json:"title"`
	Description  string   `json:"description"`
	Priority     string   `json:"priority"`
	Status       string   `json:"status"`
	Assignee     string   `json:"assignee"`
	CreatedAt    string   `json:"created_at"`
}

func Open(dataDir string) (*DB, error) {
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return nil, err
	}
	dsn := filepath.Join(dataDir, "bounty.db") + "?_journal_mode=WAL&_busy_timeout=5000"
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, err
	}
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS issues (
			id TEXT PRIMARY KEY,\n\t\t\ttitle TEXT DEFAULT '',\n\t\t\tdescription TEXT DEFAULT '',\n\t\t\tpriority TEXT DEFAULT 'medium',\n\t\t\tstatus TEXT DEFAULT 'open',\n\t\t\tassignee TEXT DEFAULT '',
			created_at TEXT DEFAULT (datetime('now'))
		)`)
	if err != nil {
		return nil, fmt.Errorf("migrate: %w", err)
	}
	return &DB{db: db}, nil
}

func (d *DB) Close() error { return d.db.Close() }

func genID() string { return fmt.Sprintf("%d", time.Now().UnixNano()) }

func (d *DB) Create(e *Issue) error {
	e.ID = genID()
	e.CreatedAt = time.Now().UTC().Format(time.RFC3339)
	_, err := d.db.Exec(`INSERT INTO issues (id, title, description, priority, status, assignee, created_at) VALUES (?, ?, ?, ?, ?, ?, ?)`,
		e.ID, e.Title, e.Description, e.Priority, e.Status, e.Assignee, e.CreatedAt)
	return err
}

func (d *DB) Get(id string) *Issue {
	row := d.db.QueryRow(`SELECT id, title, description, priority, status, assignee, created_at FROM issues WHERE id=?`, id)
	var e Issue
	if err := row.Scan(&e.ID, &e.Title, &e.Description, &e.Priority, &e.Status, &e.Assignee, &e.CreatedAt); err != nil {
		return nil
	}
	return &e
}

func (d *DB) List() []Issue {
	rows, err := d.db.Query(`SELECT id, title, description, priority, status, assignee, created_at FROM issues ORDER BY created_at DESC`)
	if err != nil {
		return nil
	}
	defer rows.Close()
	var result []Issue
	for rows.Next() {
		var e Issue
		if err := rows.Scan(&e.ID, &e.Title, &e.Description, &e.Priority, &e.Status, &e.Assignee, &e.CreatedAt); err != nil {
			continue
		}
		result = append(result, e)
	}
	return result
}

func (d *DB) Delete(id string) error {
	_, err := d.db.Exec(`DELETE FROM issues WHERE id=?`, id)
	return err
}

func (d *DB) Count() int {
	var n int
	d.db.QueryRow(`SELECT COUNT(*) FROM issues`).Scan(&n)
	return n
}
