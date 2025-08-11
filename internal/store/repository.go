package store

import (
	"database/sql"
	"fmt"
	"time"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) CreateTarget(target *Target) error {
	query := `
		INSERT INTO targets (name, host, port, user, password_enc, comment, 
		                     schedule_time, retention_days, auto_compress, database_mode, 
		                     selected_databases, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	now := time.Now()
	target.CreatedAt = now
	target.UpdatedAt = now

	result, err := r.db.Exec(query, target.Name, target.Host, target.Port, target.User, 
		target.PasswordEnc, target.Comment, target.ScheduleTime, target.RetentionDays, 
		target.AutoCompress, target.DatabaseMode, target.SelectedDatabases, target.CreatedAt, target.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to create target: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert id: %w", err)
	}
	target.ID = id

	return nil
}

func (r *Repository) GetTargets() ([]*Target, error) {
	query := `
		SELECT id, name, host, port, user, password_enc, comment,
		       schedule_time, retention_days, auto_compress, database_mode, 
		       selected_databases, created_at, updated_at
		FROM targets ORDER BY name
	`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query targets: %w", err)
	}
	defer rows.Close()

	var targets []*Target
	for rows.Next() {
		target := &Target{}
		err := rows.Scan(&target.ID, &target.Name, &target.Host, &target.Port,
			&target.User, &target.PasswordEnc, &target.Comment,
			&target.ScheduleTime, &target.RetentionDays, &target.AutoCompress,
			&target.DatabaseMode, &target.SelectedDatabases, &target.CreatedAt, &target.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan target: %w", err)
		}
		targets = append(targets, target)
	}

	return targets, nil
}

func (r *Repository) GetTarget(id int64) (*Target, error) {
	query := `
		SELECT id, name, host, port, user, password_enc, comment,
		       schedule_time, retention_days, auto_compress, database_mode, 
		       selected_databases, created_at, updated_at
		FROM targets WHERE id = ?
	`
	target := &Target{}
	err := r.db.QueryRow(query, id).Scan(&target.ID, &target.Name, &target.Host,
		&target.Port, &target.User, &target.PasswordEnc, &target.Comment,
		&target.ScheduleTime, &target.RetentionDays, &target.AutoCompress,
		&target.DatabaseMode, &target.SelectedDatabases, &target.CreatedAt, &target.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("target not found")
		}
		return nil, fmt.Errorf("failed to get target: %w", err)
	}

	return target, nil
}

func (r *Repository) UpdateTarget(target *Target) error {
	query := `
		UPDATE targets SET name = ?, host = ?, port = ?, user = ?,
		                   password_enc = ?, comment = ?, schedule_time = ?,
		                   retention_days = ?, auto_compress = ?, database_mode = ?, 
		                   selected_databases = ?, updated_at = ?
		WHERE id = ?
	`
	target.UpdatedAt = time.Now()
	_, err := r.db.Exec(query, target.Name, target.Host, target.Port, target.User,
		target.PasswordEnc, target.Comment, target.ScheduleTime, target.RetentionDays, 
		target.AutoCompress, target.DatabaseMode, target.SelectedDatabases, target.UpdatedAt, target.ID)
	if err != nil {
		return fmt.Errorf("failed to update target: %w", err)
	}

	return nil
}

func (r *Repository) DeleteTarget(id int64) error {
	_, err := r.db.Exec("DELETE FROM targets WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("failed to delete target: %w", err)
	}
	return nil
}

func (r *Repository) CreateBackup(backup *Backup) error {
	query := `
		INSERT INTO backups (target_id, database_name, started_at, finished_at, size_bytes, status, file_path, notes)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`
	result, err := r.db.Exec(query, backup.TargetID, backup.DatabaseName, backup.StartedAt, backup.FinishedAt,
		backup.SizeBytes, backup.Status, backup.FilePath, backup.Notes)
	if err != nil {
		return fmt.Errorf("failed to create backup: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert id: %w", err)
	}
	backup.ID = id

	return nil
}

func (r *Repository) UpdateBackup(backup *Backup) error {
	query := `
		UPDATE backups SET database_name = ?, finished_at = ?, size_bytes = ?, status = ?, file_path = ?, notes = ?
		WHERE id = ?
	`
	_, err := r.db.Exec(query, backup.DatabaseName, backup.FinishedAt, backup.SizeBytes, backup.Status,
		backup.FilePath, backup.Notes, backup.ID)
	if err != nil {
		return fmt.Errorf("failed to update backup: %w", err)
	}
	return nil
}

func (r *Repository) GetBackupsByTarget(targetID int64) ([]*Backup, error) {
	query := `
		SELECT id, target_id, database_name, started_at, finished_at, size_bytes, status, file_path, notes
		FROM backups WHERE target_id = ? ORDER BY started_at DESC
	`
	rows, err := r.db.Query(query, targetID)
	if err != nil {
		return nil, fmt.Errorf("failed to query backups: %w", err)
	}
	defer rows.Close()

	var backups []*Backup
	for rows.Next() {
		backup := &Backup{}
		err := rows.Scan(&backup.ID, &backup.TargetID, &backup.DatabaseName, &backup.StartedAt, &backup.FinishedAt,
			&backup.SizeBytes, &backup.Status, &backup.FilePath, &backup.Notes)
		if err != nil {
			return nil, fmt.Errorf("failed to scan backup: %w", err)
		}
		backups = append(backups, backup)
	}

	return backups, nil
}

func (r *Repository) GetBackup(id int64) (*Backup, error) {
	query := `
		SELECT id, target_id, database_name, started_at, finished_at, size_bytes, status, file_path, notes
		FROM backups WHERE id = ?
	`
	backup := &Backup{}
	err := r.db.QueryRow(query, id).Scan(&backup.ID, &backup.TargetID, &backup.DatabaseName,
		&backup.StartedAt, &backup.FinishedAt, &backup.SizeBytes, &backup.Status, &backup.FilePath, &backup.Notes)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("backup not found")
		}
		return nil, fmt.Errorf("failed to get backup: %w", err)
	}

	return backup, nil
}

func (r *Repository) GetAllBackups() ([]*Backup, error) {
	query := `
		SELECT id, target_id, database_name, started_at, finished_at, size_bytes, status, file_path, notes
		FROM backups ORDER BY started_at DESC
	`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query backups: %w", err)
	}
	defer rows.Close()

	var backups []*Backup
	for rows.Next() {
		backup := &Backup{}
		err := rows.Scan(&backup.ID, &backup.TargetID, &backup.DatabaseName, &backup.StartedAt, &backup.FinishedAt,
			&backup.SizeBytes, &backup.Status, &backup.FilePath, &backup.Notes)
		if err != nil {
			return nil, fmt.Errorf("failed to scan backup: %w", err)
		}
		backups = append(backups, backup)
	}

	return backups, nil
}

func (r *Repository) DeleteBackup(id int64) error {
	_, err := r.db.Exec("DELETE FROM backups WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("failed to delete backup: %w", err)
	}
	return nil
}

func (r *Repository) DeleteOldBackups(targetID int64, cutoff time.Time) error {
	_, err := r.db.Exec("DELETE FROM backups WHERE target_id = ? AND started_at < ?", targetID, cutoff)
	if err != nil {
		return fmt.Errorf("failed to delete old backups: %w", err)
	}
	return nil
}