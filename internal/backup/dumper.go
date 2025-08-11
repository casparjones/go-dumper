package backup

import (
	"bufio"
	"compress/gzip"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/casparjones/go-dumper/internal/store"
	"github.com/go-sql-driver/mysql"
)

type Dumper struct {
	repo      *store.Repository
	backupDir string
}

type DumpOptions struct {
	Target       *store.Target
	DatabaseName string
	BackupID     int64
	Compress     bool
	BatchSize    int
}

func NewDumper(repo *store.Repository, backupDir string) *Dumper {
	return &Dumper{
		repo:      repo,
		backupDir: backupDir,
	}
}

func (d *Dumper) CreateBackup(ctx context.Context, targetID int64) (*store.Backup, error) {
	target, err := d.repo.GetTarget(targetID)
	if err != nil {
		return nil, fmt.Errorf("failed to get target: %w", err)
	}

	// Get databases to backup based on target configuration
	databases, err := d.getDatabasesForTarget(target)
	if err != nil {
		return nil, fmt.Errorf("failed to get databases for target: %w", err)
	}

	if len(databases) == 0 {
		return nil, fmt.Errorf("no databases found to backup")
	}

	// Create backup records for each database
	backups := make([]*store.Backup, 0, len(databases))
	for _, dbName := range databases {
		backup := &store.Backup{
			TargetID:     targetID,
			DatabaseName: dbName,
			StartedAt:    time.Now(),
			Status:       store.BackupStatusRunning,
		}

		if err := d.repo.CreateBackup(backup); err != nil {
			return nil, fmt.Errorf("failed to create backup record for %s: %w", dbName, err)
		}
		backups = append(backups, backup)
	}

	// Start backup process in background
	go func() {
		d.performMultipleDatabaseBackup(context.Background(), backups, target)
	}()

	// Return the first backup as reference (for API compatibility)
	return backups[0], nil
}

func (d *Dumper) getDatabasesForTarget(target *store.Target) ([]string, error) {
	password, err := store.DecryptPassword(target.PasswordEnc)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt password: %w", err)
	}

	// Connect without specifying a database
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
		return nil, fmt.Errorf("failed to connect to MySQL: %w", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping MySQL: %w", err)
	}

	if target.DatabaseMode == store.DatabaseModeAll {
		return d.getAllDatabases(db)
	} else if target.DatabaseMode == store.DatabaseModeSelected {
		return d.getSelectedDatabases(target)
	}

	return nil, fmt.Errorf("invalid database mode: %s", target.DatabaseMode)
}

func (d *Dumper) getAllDatabases(db *sql.DB) ([]string, error) {
	rows, err := db.Query("SHOW DATABASES")
	if err != nil {
		return nil, fmt.Errorf("failed to list databases: %w", err)
	}
	defer rows.Close()

	var databases []string
	systemDatabases := map[string]bool{
		"information_schema": true,
		"performance_schema": true,
		"mysql":              true,
		"sys":                true,
	}

	for rows.Next() {
		var dbName string
		if err := rows.Scan(&dbName); err != nil {
			return nil, fmt.Errorf("failed to scan database name: %w", err)
		}

		// Skip system databases
		if !systemDatabases[dbName] {
			databases = append(databases, dbName)
		}
	}

	return databases, nil
}

func (d *Dumper) getSelectedDatabases(target *store.Target) ([]string, error) {
	if target.SelectedDatabases == "" {
		return nil, fmt.Errorf("no databases selected")
	}

	var databases []string
	if err := json.Unmarshal([]byte(target.SelectedDatabases), &databases); err != nil {
		return nil, fmt.Errorf("failed to parse selected databases: %w", err)
	}

	return databases, nil
}

func (d *Dumper) performMultipleDatabaseBackup(ctx context.Context, backups []*store.Backup, target *store.Target) {
	password, err := store.DecryptPassword(target.PasswordEnc)
	if err != nil {
		for _, backup := range backups {
			d.updateBackupStatus(backup, store.BackupStatusFailed, fmt.Sprintf("Failed to decrypt password: %v", err))
		}
		return
	}

	// Process each database backup
	for _, backup := range backups {
		d.performSingleDatabaseBackup(ctx, backup, target, password)
	}

	// Cleanup old backups after all databases are processed
	d.cleanupOldBackups(target.ID, target.RetentionDays)
}

func (d *Dumper) performSingleDatabaseBackup(ctx context.Context, backup *store.Backup, target *store.Target, password string) {
	timestamp := backup.StartedAt.Format("2006-01-02_15-04-05")
	filename := fmt.Sprintf("%s_%s_%s.sql.gz", target.Name, backup.DatabaseName, timestamp)
	filepath := filepath.Join(d.backupDir, filename)

	options := &DumpOptions{
		Target:       target,
		DatabaseName: backup.DatabaseName,
		BackupID:     backup.ID,
		Compress:     target.AutoCompress,
		BatchSize:    1000,
	}

	size, err := d.dumpDatabase(ctx, options, filepath, password)
	if err != nil {
		d.updateBackupStatus(backup, store.BackupStatusFailed, err.Error())
		return
	}

	finishedAt := time.Now()
	backup.FinishedAt = &finishedAt
	backup.SizeBytes = size
	backup.Status = store.BackupStatusSuccess
	backup.FilePath = filepath

	if err := d.repo.UpdateBackup(backup); err != nil {
		d.updateBackupStatus(backup, store.BackupStatusFailed, fmt.Sprintf("Failed to update backup: %v", err))
		return
	}
}

func (d *Dumper) dumpDatabase(ctx context.Context, options *DumpOptions, outputPath, password string) (int64, error) {
	cfg := mysql.Config{
		User:                 options.Target.User,
		Passwd:               password,
		Net:                  "tcp",
		Addr:                 fmt.Sprintf("%s:%d", options.Target.Host, options.Target.Port),
		DBName:               options.DatabaseName,
		Params:               map[string]string{"charset": "utf8mb4"},
		ParseTime:            true,
		AllowNativePasswords: true,
	}

	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		return 0, fmt.Errorf("failed to connect to MySQL: %w", err)
	}
	defer db.Close()

	if err := db.PingContext(ctx); err != nil {
		return 0, fmt.Errorf("failed to ping MySQL: %w", err)
	}

	file, err := os.Create(outputPath)
	if err != nil {
		return 0, fmt.Errorf("failed to create output file: %w", err)
	}
	// WICHTIG: file.Close() NICHT sofort defer’n – erst ganz am Ende schließen
	// damit wir die Reihenfolge kontrollieren
	// defer file.Close()

	var (
		writer   io.Writer = file
		gzWriter *gzip.Writer
	)
	if options.Compress {
		gzWriter = gzip.NewWriter(file)
		writer = gzWriter
	}

	bufWriter := bufio.NewWriter(writer)

	// … deine Schreib-Logik bleibt gleich …

	if err := d.enableForeignKeyChecks(bufWriter); err != nil {
		return 0, fmt.Errorf("failed to enable foreign key checks: %w", err)
	}

	if err := bufWriter.Flush(); err != nil {
		return 0, fmt.Errorf("failed to flush buffer: %w", err)
	}
	if gzWriter != nil {
		if err := gzWriter.Close(); err != nil { // explizit schließen, damit alles auf die Datei geschrieben ist
			return 0, fmt.Errorf("failed to close gzip writer: %w", err)
		}
	}
	if err := file.Sync(); err != nil { // optional, sorgt für persistente Größe
		return 0, fmt.Errorf("failed to sync file: %w", err)
	}
	if err := file.Close(); err != nil {
		return 0, fmt.Errorf("failed to close file: %w", err)
	}

	stat, err := os.Stat(outputPath)
	if err != nil {
		return 0, fmt.Errorf("failed to get file stats: %w", err)
	}
	return stat.Size(), nil
}

func (d *Dumper) writeHeader(w io.Writer, target *store.Target, databaseName string) error {
	header := fmt.Sprintf(`-- MySQL dump created by go-dumper
-- Host: %s    Database: %s
-- ------------------------------------------------------
-- Server version	

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!50503 SET NAMES utf8mb4 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

`, target.Host, databaseName)

	_, err := w.Write([]byte(header))
	return err
}

func (d *Dumper) disableForeignKeyChecks(w io.Writer) error {
	_, err := w.Write([]byte("SET FOREIGN_KEY_CHECKS=0;\n\n"))
	return err
}

func (d *Dumper) enableForeignKeyChecks(w io.Writer) error {
	_, err := w.Write([]byte("\nSET FOREIGN_KEY_CHECKS=1;\n"))
	return err
}

func (d *Dumper) getTables(ctx context.Context, tx *sql.Tx) ([]string, error) {
	query := "SELECT table_name FROM information_schema.tables WHERE table_schema = DATABASE() AND table_type = 'BASE TABLE'"
	rows, err := tx.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tables []string
	for rows.Next() {
		var table string
		if err := rows.Scan(&table); err != nil {
			return nil, err
		}
		tables = append(tables, table)
	}

	return tables, nil
}

func (d *Dumper) getViews(ctx context.Context, tx *sql.Tx) ([]string, error) {
	query := "SELECT table_name FROM information_schema.views WHERE table_schema = DATABASE()"
	rows, err := tx.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var views []string
	for rows.Next() {
		var view string
		if err := rows.Scan(&view); err != nil {
			return nil, err
		}
		views = append(views, view)
	}

	return views, nil
}

func (d *Dumper) dumpTable(ctx context.Context, tx *sql.Tx, w io.Writer, table string, batchSize int) error {
	createTableSQL, err := d.getCreateTableSQL(ctx, tx, table)
	if err != nil {
		return fmt.Errorf("failed to get CREATE TABLE for %s: %w", table, err)
	}

	if _, err := w.Write([]byte(fmt.Sprintf("--\n-- Table structure for table `%s`\n--\n\n", table))); err != nil {
		return err
	}

	if _, err := w.Write([]byte(fmt.Sprintf("DROP TABLE IF EXISTS `%s`;\n", table))); err != nil {
		return err
	}

	if _, err := w.Write([]byte(createTableSQL + ";\n\n")); err != nil {
		return err
	}

	return d.dumpTableData(ctx, tx, w, table, batchSize)
}

func (d *Dumper) getCreateTableSQL(ctx context.Context, tx *sql.Tx, table string) (string, error) {
	var tableName, createSQL string
	err := tx.QueryRowContext(ctx, "SHOW CREATE TABLE `"+table+"`").Scan(&tableName, &createSQL)
	if err != nil {
		return "", err
	}
	return createSQL, nil
}

func (d *Dumper) dumpTableData(ctx context.Context, tx *sql.Tx, w io.Writer, table string, batchSize int) error {
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM `%s`", table)
	var count int64
	if err := tx.QueryRowContext(ctx, countQuery).Scan(&count); err != nil {
		return err
	}

	if count == 0 {
		return nil
	}

	if _, err := w.Write([]byte(fmt.Sprintf("--\n-- Dumping data for table `%s`\n--\n\n", table))); err != nil {
		return err
	}

	if _, err := w.Write([]byte(fmt.Sprintf("LOCK TABLES `%s` WRITE;\n", table))); err != nil {
		return err
	}

	columns, err := d.getTableColumns(ctx, tx, table)
	if err != nil {
		return err
	}

	query := fmt.Sprintf("SELECT %s FROM `%s`", strings.Join(columns, ", "), table)
	rows, err := tx.QueryContext(ctx, query)
	if err != nil {
		return err
	}
	defer rows.Close()

	values := make([]interface{}, len(columns))
	valuePtrs := make([]interface{}, len(columns))
	for i := range values {
		valuePtrs[i] = &values[i]
	}

	var insertValues []string
	rowCount := 0

	for rows.Next() {
		if err := rows.Scan(valuePtrs...); err != nil {
			return err
		}

		rowData := make([]string, len(columns))
		for i, val := range values {
			rowData[i] = d.formatValue(val)
		}

		insertValues = append(insertValues, fmt.Sprintf("(%s)", strings.Join(rowData, ", ")))
		rowCount++

		if rowCount%batchSize == 0 {
			if err := d.writeInsert(w, table, columns, insertValues); err != nil {
				return err
			}
			insertValues = insertValues[:0]
		}
	}

	if len(insertValues) > 0 {
		if err := d.writeInsert(w, table, columns, insertValues); err != nil {
			return err
		}
	}

	if _, err := w.Write([]byte("UNLOCK TABLES;\n\n")); err != nil {
		return err
	}

	return nil
}

func (d *Dumper) getTableColumns(ctx context.Context, tx *sql.Tx, table string) ([]string, error) {
	query := fmt.Sprintf("SHOW COLUMNS FROM `%s`", table)
	rows, err := tx.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var columns []string
	for rows.Next() {
		var field, typ, null, key, defaultVal, extra sql.NullString
		if err := rows.Scan(&field, &typ, &null, &key, &defaultVal, &extra); err != nil {
			return nil, err
		}
		columns = append(columns, fmt.Sprintf("`%s`", field.String))
	}

	return columns, nil
}

func (d *Dumper) writeInsert(w io.Writer, table string, columns []string, values []string) error {
	insert := fmt.Sprintf("INSERT INTO `%s` (%s) VALUES\n%s;\n",
		table,
		strings.Join(columns, ", "),
		strings.Join(values, ",\n"))

	_, err := w.Write([]byte(insert))
	return err
}

func (d *Dumper) formatValue(val interface{}) string {
	if val == nil {
		return "NULL"
	}

	switch v := val.(type) {
	case []byte:
		return "'" + d.escapeString(string(v)) + "'"
	case string:
		return "'" + d.escapeString(v) + "'"
	case time.Time:
		return "'" + v.Format("2006-01-02 15:04:05") + "'"
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		return fmt.Sprintf("%d", v)
	case float32, float64:
		return fmt.Sprintf("%g", v)
	case bool:
		if v {
			return "1"
		}
		return "0"
	default:
		return "'" + d.escapeString(fmt.Sprintf("%v", v)) + "'"
	}
}

func (d *Dumper) escapeString(s string) string {
	s = strings.ReplaceAll(s, "\\", "\\\\")
	s = strings.ReplaceAll(s, "'", "\\'")
	s = strings.ReplaceAll(s, "\"", "\\\"")
	s = strings.ReplaceAll(s, "\n", "\\n")
	s = strings.ReplaceAll(s, "\r", "\\r")
	s = strings.ReplaceAll(s, "\t", "\\t")
	return s
}

func (d *Dumper) dumpView(ctx context.Context, tx *sql.Tx, w io.Writer, view string) error {
	createViewSQL, err := d.getCreateViewSQL(ctx, tx, view)
	if err != nil {
		return fmt.Errorf("failed to get CREATE VIEW for %s: %w", view, err)
	}

	if _, err := w.Write([]byte(fmt.Sprintf("--\n-- View structure for view `%s`\n--\n\n", view))); err != nil {
		return err
	}

	if _, err := w.Write([]byte(fmt.Sprintf("DROP VIEW IF EXISTS `%s`;\n", view))); err != nil {
		return err
	}

	if _, err := w.Write([]byte(createViewSQL + ";\n\n")); err != nil {
		return err
	}

	return nil
}

func (d *Dumper) getCreateViewSQL(ctx context.Context, tx *sql.Tx, view string) (string, error) {
	var viewName, createSQL, charset, collation string
	err := tx.QueryRowContext(ctx, "SHOW CREATE VIEW `"+view+"`").Scan(&viewName, &createSQL, &charset, &collation)
	if err != nil {
		return "", err
	}
	return createSQL, nil
}

func (d *Dumper) updateBackupStatus(backup *store.Backup, status, notes string) {
	finishedAt := time.Now()
	backup.FinishedAt = &finishedAt
	backup.Status = status
	backup.Notes = notes
	d.repo.UpdateBackup(backup)
}

func (d *Dumper) cleanupOldBackups(targetID int64, retentionDays int) {
	if retentionDays <= 0 {
		return
	}

	cutoff := time.Now().AddDate(0, 0, -retentionDays)
	backups, err := d.repo.GetBackupsByTarget(targetID)
	if err != nil {
		return
	}

	for _, backup := range backups {
		if backup.StartedAt.Before(cutoff) && backup.FilePath != "" {
			os.Remove(backup.FilePath)
			d.repo.DeleteBackup(backup.ID)
		}
	}
}
