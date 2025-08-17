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
	user TEXT NOT NULL,
	password_enc TEXT NOT NULL,
	comment TEXT DEFAULT '',
	schedule_time TEXT DEFAULT '',
	retention_days INTEGER DEFAULT 30,
	auto_compress BOOLEAN DEFAULT 1,
	database_mode TEXT NOT NULL DEFAULT 'all',
	selected_databases TEXT DEFAULT '',
	created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
	updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS backups (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	target_id INTEGER NOT NULL,
	database_name TEXT NOT NULL DEFAULT '',
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

CREATE TABLE IF NOT EXISTS schedule_jobs (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	target_id INTEGER NOT NULL,
	name TEXT NOT NULL,
	description TEXT DEFAULT '',
	is_active BOOLEAN DEFAULT 1,
	schedule_config TEXT NOT NULL,
	backup_options TEXT NOT NULL,
	meta_config TEXT DEFAULT '{}',
	last_run_at DATETIME,
	last_run_status TEXT DEFAULT '',
	last_run_notes TEXT DEFAULT '',
	next_run_at DATETIME,
	created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
	updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
	FOREIGN KEY (target_id) REFERENCES targets(id) ON DELETE CASCADE
);

CREATE TRIGGER IF NOT EXISTS update_schedule_jobs_timestamp 
AFTER UPDATE ON schedule_jobs
FOR EACH ROW
BEGIN
	UPDATE schedule_jobs SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
END;

CREATE TABLE IF NOT EXISTS app_config (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	key TEXT NOT NULL UNIQUE,
	value TEXT NOT NULL,
	created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
	updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TRIGGER IF NOT EXISTS update_app_config_timestamp 
AFTER UPDATE ON app_config
FOR EACH ROW
BEGIN
	UPDATE app_config SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
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

	// Enable foreign key constraints
	if _, err := db.Exec("PRAGMA foreign_keys = ON"); err != nil {
		return nil, fmt.Errorf("failed to enable foreign keys: %w", err)
	}

	if _, err := db.Exec(schema); err != nil {
		return nil, fmt.Errorf("failed to create schema: %w", err)
	}

	// Run migrations for existing databases
	if err := runMigrations(db); err != nil {
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	return db, nil
}

func runMigrations(db *sql.DB) error {
	// Check for dbname column in targets table
	rows, err := db.Query("PRAGMA table_info(targets)")
	if err != nil {
		return fmt.Errorf("failed to check table info: %w", err)
	}
	defer rows.Close()

	var hasDbName, hasDatabaseMode, hasSelectedDatabases bool
	for rows.Next() {
		var cid int
		var name, dataType string
		var notNull, pk int
		var defaultValue interface{}

		err := rows.Scan(&cid, &name, &dataType, &notNull, &defaultValue, &pk)
		if err != nil {
			return fmt.Errorf("failed to scan column info: %w", err)
		}

		if name == "dbname" {
			hasDbName = true
		}
		if name == "database_mode" {
			hasDatabaseMode = true
		}
		if name == "selected_databases" {
			hasSelectedDatabases = true
		}
	}

	if hasDbName && (!hasDatabaseMode || !hasSelectedDatabases) {
		// Need to migrate from old schema
		fmt.Println("Migrating database schema...")

		// Add new columns if they don't exist
		if !hasDatabaseMode {
			if _, err := db.Exec("ALTER TABLE targets ADD COLUMN database_mode TEXT NOT NULL DEFAULT 'all'"); err != nil {
				return fmt.Errorf("failed to add database_mode column: %w", err)
			}
		}

		if !hasSelectedDatabases {
			if _, err := db.Exec("ALTER TABLE targets ADD COLUMN selected_databases TEXT DEFAULT ''"); err != nil {
				return fmt.Errorf("failed to add selected_databases column: %w", err)
			}
		}

		// For existing targets, set database_mode to 'all' and leave selected_databases empty
		if _, err := db.Exec("UPDATE targets SET database_mode = 'all' WHERE database_mode = ''"); err != nil {
			return fmt.Errorf("failed to update database_mode: %w", err)
		}

		fmt.Println("Database schema migration completed successfully.")
	}

	// Check if backups table needs database_name column
	rows, err = db.Query("PRAGMA table_info(backups)")
	if err != nil {
		return fmt.Errorf("failed to check backups table info: %w", err)
	}
	defer rows.Close()

	var hasDatabaseName bool
	for rows.Next() {
		var cid int
		var name, dataType string
		var notNull, pk int
		var defaultValue interface{}

		err := rows.Scan(&cid, &name, &dataType, &notNull, &defaultValue, &pk)
		if err != nil {
			return fmt.Errorf("failed to scan backup column info: %w", err)
		}

		if name == "database_name" {
			hasDatabaseName = true
		}
	}

	if !hasDatabaseName {
		if _, err := db.Exec("ALTER TABLE backups ADD COLUMN database_name TEXT NOT NULL DEFAULT ''"); err != nil {
			return fmt.Errorf("failed to add database_name column to backups: %w", err)
		}
		fmt.Println("Added database_name column to backups table.")
	}

	return nil
}
