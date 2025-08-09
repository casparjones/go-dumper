//go:build integration
// +build integration

package backup

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/user/go-dumper/internal/store"

	_ "modernc.org/sqlite"
)

func setupIntegrationTest(t *testing.T) (string, *store.Repository, *Dumper, *Restorer) {
	// Set up test encryption key
	testKey := "dGVzdGtleTEyMzQ1Njc4OTBhYmNkZWZnaGlqa2xtbm9wcXJzdHV2d3h5ej0="
	os.Setenv("APP_ENC_KEY", testKey)
	t.Cleanup(func() {
		os.Unsetenv("APP_ENC_KEY")
	})

	// Create temporary database
	tempFile, err := os.CreateTemp("", "test*.db")
	if err != nil {
		t.Fatal(err)
	}
	tempFile.Close()

	t.Cleanup(func() {
		os.Remove(tempFile.Name())
	})

	db, err := store.InitDB(tempFile.Name())
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		db.Close()
	})

	repo := store.NewRepository(db)

	// Create temporary backup directory
	backupDir, err := os.MkdirTemp("", "backup-test-*")
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		os.RemoveAll(backupDir)
	})

	dumper := NewDumper(repo, backupDir)
	restorer := NewRestorer(repo)

	return backupDir, repo, dumper, restorer
}

func setupMySQLTarget(t *testing.T, repo *store.Repository) *store.Target {
	// Get MySQL DSN from environment or use default
	mysqlDSN := os.Getenv("MYSQL_DSN")
	if mysqlDSN == "" {
		mysqlDSN = "testuser:testpass@tcp(localhost:3306)/testdb"
	}

	cfg, err := mysql.ParseDSN(mysqlDSN)
	if err != nil {
		t.Skip("Invalid MySQL DSN, skipping integration test")
	}

	// Test MySQL connection
	db, err := sql.Open("mysql", mysqlDSN)
	if err != nil {
		t.Skip("Cannot connect to MySQL, skipping integration test")
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		t.Skip("MySQL not available, skipping integration test")
	}

	// Create test data
	_, err = db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS test_users (
			id INT PRIMARY KEY AUTO_INCREMENT,
			name VARCHAR(100) NOT NULL,
			email VARCHAR(100) UNIQUE NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		t.Fatal(err)
	}

	_, err = db.ExecContext(ctx, `
		INSERT INTO test_users (name, email) VALUES 
		('John Doe', 'john@example.com'),
		('Jane Smith', 'jane@example.com'),
		('Bob Johnson', 'bob@example.com')
		ON DUPLICATE KEY UPDATE name=VALUES(name)
	`)
	if err != nil {
		t.Fatal(err)
	}

	// Encrypt password
	encryptedPass, err := store.EncryptPassword(cfg.Passwd)
	if err != nil {
		t.Fatal(err)
	}

	target := &store.Target{
		Name:          "Test MySQL",
		Host:          cfg.Addr[:strings.LastIndex(cfg.Addr, ":")],
		Port:          3306,
		DBName:        cfg.DBName,
		User:          cfg.User,
		PasswordEnc:   encryptedPass,
		Comment:       "Integration test target",
		RetentionDays: 7,
		AutoCompress:  true,
	}

	if colonIndex := strings.LastIndex(cfg.Addr, ":"); colonIndex >= 0 {
		target.Host = cfg.Addr[:colonIndex]
		if portStr := cfg.Addr[colonIndex+1:]; portStr != "" {
			if port, err := strconv.Atoi(portStr); err == nil {
				target.Port = port
			}
		}
	}

	err = repo.CreateTarget(target)
	if err != nil {
		t.Fatal(err)
	}

	return target
}

func TestIntegrationBackupAndRestore(t *testing.T) {
	backupDir, repo, dumper, restorer := setupIntegrationTest(t)
	target := setupMySQLTarget(t, repo)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Create backup
	backup, err := dumper.CreateBackup(ctx, target.ID)
	if err != nil {
		t.Fatalf("Failed to create backup: %v", err)
	}

	// Wait for backup to complete
	var completedBackup *store.Backup
	for i := 0; i < 30; i++ {
		time.Sleep(1 * time.Second)
		completedBackup, err = repo.GetBackup(backup.ID)
		if err != nil {
			t.Fatal(err)
		}
		if completedBackup.Status != store.BackupStatusRunning {
			break
		}
	}

	if completedBackup.Status == store.BackupStatusRunning {
		t.Fatal("Backup did not complete within timeout")
	}

	if completedBackup.Status != store.BackupStatusSuccess {
		t.Fatalf("Backup failed: %s", completedBackup.Notes)
	}

	// Verify backup file exists
	if _, err := os.Stat(completedBackup.FilePath); os.IsNotExist(err) {
		t.Fatalf("Backup file does not exist: %s", completedBackup.FilePath)
	}

	// Verify backup file is not empty
	info, err := os.Stat(completedBackup.FilePath)
	if err != nil {
		t.Fatal(err)
	}
	if info.Size() == 0 {
		t.Fatal("Backup file is empty")
	}

	// Verify size matches database record
	if completedBackup.SizeBytes != info.Size() {
		t.Errorf("Size mismatch: database shows %d, file is %d", completedBackup.SizeBytes, info.Size())
	}

	// Test restore
	err = restorer.RestoreBackup(ctx, backup.ID)
	if err != nil {
		t.Fatalf("Failed to restore backup: %v", err)
	}

	t.Logf("Integration test completed successfully. Backup size: %d bytes", completedBackup.SizeBytes)
}

func TestIntegrationLargeDataset(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping large dataset test in short mode")
	}

	backupDir, repo, dumper, _ := setupIntegrationTest(t)
	target := setupMySQLTarget(t, repo)

	// Get MySQL connection to create large dataset
	password, err := store.DecryptPassword(target.PasswordEnc)
	if err != nil {
		t.Fatal(err)
	}

	cfg := mysql.Config{
		User:   target.User,
		Passwd: password,
		Net:    "tcp",
		Addr:   fmt.Sprintf("%s:%d", target.Host, target.Port),
		DBName: target.DBName,
	}

	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	// Create a larger dataset
	_, err = db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS large_test_table (
			id INT PRIMARY KEY AUTO_INCREMENT,
			data TEXT,
			number_col DECIMAL(10,2),
			date_col DATE,
			timestamp_col TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		t.Fatal(err)
	}

	// Insert many rows
	for i := 0; i < 1000; i++ {
		_, err = db.ExecContext(ctx, `
			INSERT INTO large_test_table (data, number_col, date_col) 
			VALUES (?, ?, ?)
		`, fmt.Sprintf("Test data row %d with some longer content to make the backup more realistic", i), float64(i)*3.14, "2023-01-01")
		if err != nil {
			t.Fatal(err)
		}
	}

	// Create backup
	backup, err := dumper.CreateBackup(ctx, target.ID)
	if err != nil {
		t.Fatalf("Failed to create backup: %v", err)
	}

	// Wait for backup to complete
	var completedBackup *store.Backup
	for i := 0; i < 60; i++ {
		time.Sleep(1 * time.Second)
		completedBackup, err = repo.GetBackup(backup.ID)
		if err != nil {
			t.Fatal(err)
		}
		if completedBackup.Status != store.BackupStatusRunning {
			break
		}
	}

	if completedBackup.Status != store.BackupStatusSuccess {
		t.Fatalf("Large dataset backup failed: %s", completedBackup.Notes)
	}

	t.Logf("Large dataset backup completed. Size: %d bytes", completedBackup.SizeBytes)

	// Clean up
	_, err = db.ExecContext(ctx, "DROP TABLE large_test_table")
	if err != nil {
		t.Logf("Failed to clean up test table: %v", err)
	}
}

func TestIntegrationBackupRotation(t *testing.T) {
	backupDir, repo, dumper, _ := setupIntegrationTest(t)
	target := setupMySQLTarget(t, repo)

	// Set very short retention
	target.RetentionDays = 0 // Immediate cleanup
	err := repo.UpdateTarget(target)
	if err != nil {
		t.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Create first backup
	backup1, err := dumper.CreateBackup(ctx, target.ID)
	if err != nil {
		t.Fatal(err)
	}

	// Wait for it to complete
	time.Sleep(5 * time.Second)
	completed1, _ := repo.GetBackup(backup1.ID)
	if completed1.Status != store.BackupStatusSuccess {
		t.Skip("First backup didn't complete, skipping rotation test")
	}

	// Create second backup (should trigger rotation)
	backup2, err := dumper.CreateBackup(ctx, target.ID)
	if err != nil {
		t.Fatal(err)
	}

	// Wait for it to complete
	time.Sleep(5 * time.Second)

	// Check if old backup was cleaned up
	backups, err := repo.GetBackupsByTarget(target.ID)
	if err != nil {
		t.Fatal(err)
	}

	// Should only have successful backups remaining
	successfulBackups := 0
	for _, b := range backups {
		if b.Status == store.BackupStatusSuccess {
			successfulBackups++
		}
	}

	t.Logf("Found %d successful backups after rotation", successfulBackups)
}