package store

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

const schema = `
CREATE TABLE IF NOT EXISTS targets (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	name TEXT NOT NULL UNIQUE,
	host TEXT NOT NULL,
	port INTEGER NOT NULL DEFAULT 3306,
	dbname TEXT NOT NULL,
	user TEXT NOT NULL,
	password_enc TEXT NOT NULL,
	comment TEXT DEFAULT '',
	schedule_time TEXT DEFAULT '',
	retention_days INTEGER DEFAULT 30,
	auto_compress BOOLEAN DEFAULT 1,
	created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
	updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS backups (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	target_id INTEGER NOT NULL,
	started_at DATETIME NOT NULL,
	finished_at DATETIME,
	size_bytes INTEGER DEFAULT 0,
	status TEXT NOT NULL DEFAULT 'running',
	file_path TEXT DEFAULT '',
	notes TEXT DEFAULT '',
	FOREIGN KEY (target_id) REFERENCES targets(id) ON DELETE CASCADE
);

CREATE TRIGGER IF NOT EXISTS update_targets_timestamp 
AFTER UPDATE ON targets
FOR EACH ROW
BEGIN
	UPDATE targets SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
END;
`

func InitDB(dbPath string) (*sql.DB, error) {
	if err := os.MkdirAll(filepath.Dir(dbPath), 0755); err != nil {
		return nil, fmt.Errorf("failed to create database directory: %w", err)
	}

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	if _, err := db.Exec("PRAGMA foreign_keys = ON"); err != nil {
		return nil, fmt.Errorf("failed to enable foreign keys: %w", err)
	}

	if _, err := db.Exec(schema); err != nil {
		return nil, fmt.Errorf("failed to create schema: %w", err)
	}

	return db, nil
}
