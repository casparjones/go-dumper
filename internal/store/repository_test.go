package store

import (
	"database/sql"
	"os"
	"testing"
	"time"

	_ "modernc.org/sqlite"
)

func setupTestDB(t *testing.T) (*sql.DB, *Repository) {
	// Create temporary database
	tempFile, err := os.CreateTemp("", "test*.db")
	if err != nil {
		t.Fatal(err)
	}
	tempFile.Close()
	
	t.Cleanup(func() {
		os.Remove(tempFile.Name())
	})

	db, err := InitDB(tempFile.Name())
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		db.Close()
	})

	return db, NewRepository(db)
}

func setupTestEncryption(t *testing.T) {
	testKey := "dGVzdGtleTEyMzQ1Njc4OTBhYmNkZWZnaGlqa2xtbm9wcXJzdHV2d3h5ej0="
	os.Setenv("APP_ENC_KEY", testKey)
	t.Cleanup(func() {
		os.Unsetenv("APP_ENC_KEY")
	})
}

func TestTargetCRUD(t *testing.T) {
	setupTestEncryption(t)
	_, repo := setupTestDB(t)

	// Create target
	encryptedPass, err := EncryptPassword("testpass")
	if err != nil {
		t.Fatal(err)
	}

	target := &Target{
		Name:          "Test Target",
		Host:          "localhost",
		Port:          3306,
		DBName:        "testdb",
		User:          "testuser",
		PasswordEnc:   encryptedPass,
		Comment:       "Test comment",
		ScheduleTime:  "02:30",
		RetentionDays: 7,
		AutoCompress:  true,
	}

	// Test Create
	err = repo.CreateTarget(target)
	if err != nil {
		t.Fatalf("CreateTarget failed: %v", err)
	}

	if target.ID == 0 {
		t.Fatal("Target ID not set after creation")
	}

	// Test GetTarget
	retrieved, err := repo.GetTarget(target.ID)
	if err != nil {
		t.Fatalf("GetTarget failed: %v", err)
	}

	if retrieved.Name != target.Name {
		t.Errorf("Name mismatch: expected %q, got %q", target.Name, retrieved.Name)
	}
	if retrieved.Host != target.Host {
		t.Errorf("Host mismatch: expected %q, got %q", target.Host, retrieved.Host)
	}
	if retrieved.Port != target.Port {
		t.Errorf("Port mismatch: expected %d, got %d", target.Port, retrieved.Port)
	}

	// Test GetTargets
	targets, err := repo.GetTargets()
	if err != nil {
		t.Fatalf("GetTargets failed: %v", err)
	}

	if len(targets) != 1 {
		t.Fatalf("Expected 1 target, got %d", len(targets))
	}

	// Test Update
	target.Name = "Updated Target"
	target.Comment = "Updated comment"
	err = repo.UpdateTarget(target)
	if err != nil {
		t.Fatalf("UpdateTarget failed: %v", err)
	}

	updated, err := repo.GetTarget(target.ID)
	if err != nil {
		t.Fatalf("GetTarget after update failed: %v", err)
	}

	if updated.Name != "Updated Target" {
		t.Errorf("Name not updated: expected %q, got %q", "Updated Target", updated.Name)
	}
	if updated.Comment != "Updated comment" {
		t.Errorf("Comment not updated: expected %q, got %q", "Updated comment", updated.Comment)
	}

	// Test Delete
	err = repo.DeleteTarget(target.ID)
	if err != nil {
		t.Fatalf("DeleteTarget failed: %v", err)
	}

	_, err = repo.GetTarget(target.ID)
	if err == nil {
		t.Fatal("Expected error when getting deleted target")
	}

	// Test GetTargets after delete
	targets, err = repo.GetTargets()
	if err != nil {
		t.Fatalf("GetTargets after delete failed: %v", err)
	}

	if len(targets) != 0 {
		t.Fatalf("Expected 0 targets after delete, got %d", len(targets))
	}
}

func TestBackupCRUD(t *testing.T) {
	setupTestEncryption(t)
	_, repo := setupTestDB(t)

	// First create a target
	encryptedPass, err := EncryptPassword("testpass")
	if err != nil {
		t.Fatal(err)
	}

	target := &Target{
		Name:        "Test Target",
		Host:        "localhost",
		Port:        3306,
		DBName:      "testdb",
		User:        "testuser",
		PasswordEnc: encryptedPass,
	}

	err = repo.CreateTarget(target)
	if err != nil {
		t.Fatal(err)
	}

	// Create backup
	startTime := time.Now()
	backup := &Backup{
		TargetID:  target.ID,
		StartedAt: startTime,
		Status:    BackupStatusRunning,
	}

	// Test Create
	err = repo.CreateBackup(backup)
	if err != nil {
		t.Fatalf("CreateBackup failed: %v", err)
	}

	if backup.ID == 0 {
		t.Fatal("Backup ID not set after creation")
	}

	// Test GetBackup
	retrieved, err := repo.GetBackup(backup.ID)
	if err != nil {
		t.Fatalf("GetBackup failed: %v", err)
	}

	if retrieved.TargetID != backup.TargetID {
		t.Errorf("TargetID mismatch: expected %d, got %d", backup.TargetID, retrieved.TargetID)
	}
	if retrieved.Status != backup.Status {
		t.Errorf("Status mismatch: expected %q, got %q", backup.Status, retrieved.Status)
	}

	// Test Update
	finishTime := time.Now()
	backup.FinishedAt = &finishTime
	backup.Status = BackupStatusSuccess
	backup.SizeBytes = 12345
	backup.FilePath = "/tmp/backup.sql.gz"

	err = repo.UpdateBackup(backup)
	if err != nil {
		t.Fatalf("UpdateBackup failed: %v", err)
	}

	updated, err := repo.GetBackup(backup.ID)
	if err != nil {
		t.Fatalf("GetBackup after update failed: %v", err)
	}

	if updated.Status != BackupStatusSuccess {
		t.Errorf("Status not updated: expected %q, got %q", BackupStatusSuccess, updated.Status)
	}
	if updated.SizeBytes != 12345 {
		t.Errorf("SizeBytes not updated: expected %d, got %d", 12345, updated.SizeBytes)
	}
	if updated.FilePath != "/tmp/backup.sql.gz" {
		t.Errorf("FilePath not updated: expected %q, got %q", "/tmp/backup.sql.gz", updated.FilePath)
	}

	// Test GetBackupsByTarget
	backups, err := repo.GetBackupsByTarget(target.ID)
	if err != nil {
		t.Fatalf("GetBackupsByTarget failed: %v", err)
	}

	if len(backups) != 1 {
		t.Fatalf("Expected 1 backup, got %d", len(backups))
	}

	// Create another backup to test ordering
	backup2 := &Backup{
		TargetID:  target.ID,
		StartedAt: time.Now().Add(1 * time.Hour),
		Status:    BackupStatusSuccess,
	}
	err = repo.CreateBackup(backup2)
	if err != nil {
		t.Fatal(err)
	}

	backups, err = repo.GetBackupsByTarget(target.ID)
	if err != nil {
		t.Fatal(err)
	}

	if len(backups) != 2 {
		t.Fatalf("Expected 2 backups, got %d", len(backups))
	}

	// Should be ordered by started_at DESC
	if backups[0].ID != backup2.ID {
		t.Error("Backups not ordered correctly by started_at DESC")
	}

	// Test Delete
	err = repo.DeleteBackup(backup.ID)
	if err != nil {
		t.Fatalf("DeleteBackup failed: %v", err)
	}

	_, err = repo.GetBackup(backup.ID)
	if err == nil {
		t.Fatal("Expected error when getting deleted backup")
	}

	// Test DeleteOldBackups
	cutoff := time.Now().Add(-1 * time.Minute)
	err = repo.DeleteOldBackups(target.ID, cutoff)
	if err != nil {
		t.Fatalf("DeleteOldBackups failed: %v", err)
	}

	// backup2 should still exist since it's newer than cutoff
	backups, err = repo.GetBackupsByTarget(target.ID)
	if err != nil {
		t.Fatal(err)
	}

	if len(backups) != 1 {
		t.Fatalf("Expected 1 backup after DeleteOldBackups, got %d", len(backups))
	}
}

func TestForeignKeyConstraints(t *testing.T) {
	setupTestEncryption(t)
	_, repo := setupTestDB(t)

	// Create a target
	encryptedPass, err := EncryptPassword("testpass")
	if err != nil {
		t.Fatal(err)
	}

	target := &Target{
		Name:        "Test Target",
		Host:        "localhost",
		Port:        3306,
		DBName:      "testdb",
		User:        "testuser",
		PasswordEnc: encryptedPass,
	}

	err = repo.CreateTarget(target)
	if err != nil {
		t.Fatal(err)
	}

	// Create a backup
	backup := &Backup{
		TargetID:  target.ID,
		StartedAt: time.Now(),
		Status:    BackupStatusSuccess,
	}

	err = repo.CreateBackup(backup)
	if err != nil {
		t.Fatal(err)
	}

	// Delete the target - should also delete the backup due to CASCADE
	err = repo.DeleteTarget(target.ID)
	if err != nil {
		t.Fatal(err)
	}

	// Backup should be gone too
	_, err = repo.GetBackup(backup.ID)
	if err == nil {
		t.Fatal("Expected backup to be deleted when target is deleted")
	}
}