package backup

import (
	"bufio"
	"compress/gzip"
	"context"
	"database/sql"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/casparjones/go-dumper/internal/store"
	"github.com/go-sql-driver/mysql"
)

type Restorer struct {
	repo *store.Repository
}

func NewRestorer(repo *store.Repository) *Restorer {
	return &Restorer{
		repo: repo,
	}
}

func (r *Restorer) RestoreBackup(ctx context.Context, backupID int64) error {
	backup, err := r.repo.GetBackup(backupID)
	if err != nil {
		return fmt.Errorf("failed to get backup: %w", err)
	}

	target, err := r.repo.GetTarget(backup.TargetID)
	if err != nil {
		return fmt.Errorf("failed to get target: %w", err)
	}

	password, err := store.DecryptPassword(target.PasswordEnc)
	if err != nil {
		return fmt.Errorf("failed to decrypt password: %w", err)
	}

	// Validate backup has a database name (for new multi-database system)
	if backup.DatabaseName == "" {
		return fmt.Errorf("backup does not specify a database name - this may be an old backup format")
	}

	if _, err := os.Stat(backup.FilePath); os.IsNotExist(err) {
		return fmt.Errorf("backup file not found: %s", backup.FilePath)
	}

	file, err := os.Open(backup.FilePath)
	if err != nil {
		return fmt.Errorf("failed to open backup file: %w", err)
	}
	defer file.Close()

	var reader io.Reader = file

	if strings.HasSuffix(backup.FilePath, ".gz") {
		gzReader, err := gzip.NewReader(file)
		if err != nil {
			return fmt.Errorf("failed to create gzip reader: %w", err)
		}
		defer gzReader.Close()
		reader = gzReader
	}

	cfg := mysql.Config{
		User:   target.User,
		Passwd: password,
		Net:    "tcp",
		Addr:   fmt.Sprintf("%s:%d", target.Host, target.Port),
		DBName: backup.DatabaseName,
		Params: map[string]string{
			"charset":         "utf8mb4",
			"multiStatements": "true",
		},
		ParseTime:            true,
		AllowNativePasswords: true,
	}

	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		return fmt.Errorf("failed to connect to MySQL: %w", err)
	}
	defer db.Close()

	// Set connection settings to avoid timeouts
	db.SetConnMaxLifetime(30 * time.Second)
	db.SetMaxOpenConns(1)

	if err := db.Ping(); err != nil {
		return fmt.Errorf("failed to ping MySQL: %w", err)
	}

	// Verify that the target database exists
	if err := r.verifyDatabaseExists(ctx, db, backup.DatabaseName); err != nil {
		return fmt.Errorf("database verification failed: %w", err)
	}

	return r.executeSQLFile(ctx, db, reader)
}

func (r *Restorer) verifyDatabaseExists(ctx context.Context, db *sql.DB, dbName string) error {
	var exists bool
	query := "SELECT EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.SCHEMATA WHERE SCHEMA_NAME = ?)"
	
	err := db.QueryRow(query, dbName).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check if database exists: %w", err)
	}
	
	if !exists {
		return fmt.Errorf("target database '%s' does not exist - please create it first or enable 'Create database if it doesn't exist' option", dbName)
	}
	
	return nil
}

// RestoreBackupWithOptions allows restoring with additional options like creating database
func (r *Restorer) RestoreBackupWithOptions(ctx context.Context, backupID int64, createDatabase bool) error {
	backup, err := r.repo.GetBackup(backupID)
	if err != nil {
		return fmt.Errorf("failed to get backup: %w", err)
	}

	target, err := r.repo.GetTarget(backup.TargetID)
	if err != nil {
		return fmt.Errorf("failed to get target: %w", err)
	}

	password, err := store.DecryptPassword(target.PasswordEnc)
	if err != nil {
		return fmt.Errorf("failed to decrypt password: %w", err)
	}

	// Validate backup has a database name
	if backup.DatabaseName == "" {
		return fmt.Errorf("backup does not specify a database name - this may be an old backup format")
	}

	if _, err := os.Stat(backup.FilePath); os.IsNotExist(err) {
		return fmt.Errorf("backup file not found: %s", backup.FilePath)
	}

	// First connect without specifying database to check/create it
	if createDatabase {
		if err := r.ensureDatabaseExists(ctx, target, password, backup.DatabaseName); err != nil {
			return fmt.Errorf("failed to ensure database exists: %w", err)
		}
	}

	// Now proceed with normal restore
	return r.RestoreBackup(ctx, backupID)
}

func (r *Restorer) ensureDatabaseExists(ctx context.Context, target *store.Target, password, dbName string) error {
	// Connect without specifying database
	cfg := mysql.Config{
		User:   target.User,
		Passwd: password,
		Net:    "tcp",
		Addr:   fmt.Sprintf("%s:%d", target.Host, target.Port),
		Params: map[string]string{
			"charset": "utf8mb4",
		},
		ParseTime:            true,
		AllowNativePasswords: true,
	}

	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		return fmt.Errorf("failed to connect to MySQL: %w", err)
	}
	defer db.Close()

	// Set a longer timeout for the connection
	db.SetConnMaxLifetime(30 * time.Second)
	db.SetMaxOpenConns(1)
	
	if err := db.Ping(); err != nil {
		return fmt.Errorf("failed to ping MySQL: %w", err)
	}

	// Check if database exists
	var exists bool
	query := "SELECT EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.SCHEMATA WHERE SCHEMA_NAME = ?)"
	err = db.QueryRow(query, dbName).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check if database exists: %w", err)
	}

	// Create database if it doesn't exist
	if !exists {
		createQuery := fmt.Sprintf("CREATE DATABASE `%s` CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci", dbName)
		_, err = db.Exec(createQuery)
		if err != nil {
			return fmt.Errorf("failed to create database '%s': %w", dbName, err)
		}
	}

	return nil
}

func (r *Restorer) executeSQLFile(ctx context.Context, db *sql.DB, reader io.Reader) error {
	scanner := bufio.NewScanner(reader)
	scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024) // Increase buffer size for large statements

	var currentStatement strings.Builder
	lineNumber := 0

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		lineNumber++

		if line == "" || strings.HasPrefix(line, "--") || strings.HasPrefix(line, "/*") {
			continue
		}

		currentStatement.WriteString(line)
		currentStatement.WriteString(" ")

		if strings.HasSuffix(line, ";") {
			stmt := strings.TrimSpace(currentStatement.String())
			if stmt != "" && stmt != ";" {
				if err := r.executeStatement(ctx, db, stmt); err != nil {
					return fmt.Errorf("error at line %d: %w\nStatement: %s", lineNumber, err, stmt[:min(len(stmt), 100)])
				}
			}
			currentStatement.Reset()
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading SQL file: %w", err)
	}

	if currentStatement.Len() > 0 {
		stmt := strings.TrimSpace(currentStatement.String())
		if stmt != "" && stmt != ";" {
			if err := r.executeStatement(ctx, db, stmt); err != nil {
				return fmt.Errorf("error in final statement: %w\nStatement: %s", err, stmt[:min(len(stmt), 100)])
			}
		}
	}

	return nil
}

func (r *Restorer) executeStatement(ctx context.Context, db *sql.DB, statement string) error {
	statement = strings.TrimSpace(statement)
	if statement == "" {
		return nil
	}

	if strings.HasPrefix(strings.ToUpper(statement), "LOCK TABLES") ||
		strings.HasPrefix(strings.ToUpper(statement), "UNLOCK TABLES") {
		return nil
	}

	if strings.HasPrefix(statement, "/*!") && strings.HasSuffix(statement, "*/;") {
		versionComment := statement[3 : len(statement)-3]

		if len(versionComment) > 5 {
			versionStr := versionComment[:5]
			if strings.HasPrefix(versionStr, "40") || strings.HasPrefix(versionStr, "50") {
				actualSQL := strings.TrimSpace(versionComment[5:])
				if actualSQL != "" {
					_, err := db.ExecContext(ctx, actualSQL)
					return err
				}
			}
		}
		return nil
	}

	_, err := db.ExecContext(ctx, statement)
	if err != nil {
		if strings.Contains(err.Error(), "Unknown database") {
			return fmt.Errorf("database does not exist - please create it first: %w", err)
		}
		if strings.Contains(err.Error(), "Access denied") {
			return fmt.Errorf("access denied - check user permissions: %w", err)
		}
	}

	return err
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
