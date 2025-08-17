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

// Schedule Jobs repository methods

func (r *Repository) CreateScheduleJob(job *ScheduleJob) error {
	query := `
		INSERT INTO schedule_jobs (target_id, name, description, is_active, schedule_config, 
		                          backup_options, meta_config, next_run_at, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	now := time.Now()
	job.CreatedAt = now
	job.UpdatedAt = now

	result, err := r.db.Exec(query, job.TargetID, job.Name, job.Description, job.IsActive,
		job.ScheduleConfig, job.BackupOptions, job.MetaConfig, job.NextRunAt, job.CreatedAt, job.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to create schedule job: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert id: %w", err)
	}
	job.ID = id

	return nil
}

func (r *Repository) GetScheduleJobs() ([]*ScheduleJob, error) {
	query := `
		SELECT id, target_id, name, description, is_active, schedule_config, backup_options,
		       meta_config, last_run_at, last_run_status, last_run_notes, next_run_at,
		       created_at, updated_at
		FROM schedule_jobs ORDER BY name
	`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query schedule jobs: %w", err)
	}
	defer rows.Close()

	var jobs []*ScheduleJob
	for rows.Next() {
		job := &ScheduleJob{}
		err := rows.Scan(&job.ID, &job.TargetID, &job.Name, &job.Description, &job.IsActive,
			&job.ScheduleConfig, &job.BackupOptions, &job.MetaConfig, &job.LastRunAt,
			&job.LastRunStatus, &job.LastRunNotes, &job.NextRunAt, &job.CreatedAt, &job.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan schedule job: %w", err)
		}
		jobs = append(jobs, job)
	}

	return jobs, nil
}

func (r *Repository) GetScheduleJob(id int64) (*ScheduleJob, error) {
	query := `
		SELECT id, target_id, name, description, is_active, schedule_config, backup_options,
		       meta_config, last_run_at, last_run_status, last_run_notes, next_run_at,
		       created_at, updated_at
		FROM schedule_jobs WHERE id = ?
	`
	job := &ScheduleJob{}
	err := r.db.QueryRow(query, id).Scan(&job.ID, &job.TargetID, &job.Name, &job.Description,
		&job.IsActive, &job.ScheduleConfig, &job.BackupOptions, &job.MetaConfig, &job.LastRunAt,
		&job.LastRunStatus, &job.LastRunNotes, &job.NextRunAt, &job.CreatedAt, &job.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("schedule job not found")
		}
		return nil, fmt.Errorf("failed to get schedule job: %w", err)
	}

	return job, nil
}

func (r *Repository) UpdateScheduleJob(job *ScheduleJob) error {
	query := `
		UPDATE schedule_jobs 
		SET target_id = ?, name = ?, description = ?, is_active = ?, schedule_config = ?,
		    backup_options = ?, meta_config = ?, next_run_at = ?, updated_at = ?
		WHERE id = ?
	`
	job.UpdatedAt = time.Now()
	_, err := r.db.Exec(query, job.TargetID, job.Name, job.Description, job.IsActive,
		job.ScheduleConfig, job.BackupOptions, job.MetaConfig, job.NextRunAt, job.UpdatedAt, job.ID)
	if err != nil {
		return fmt.Errorf("failed to update schedule job: %w", err)
	}
	return nil
}

func (r *Repository) UpdateScheduleJobRunStatus(id int64, status, notes string, lastRunAt, nextRunAt *time.Time) error {
	query := `
		UPDATE schedule_jobs 
		SET last_run_at = ?, last_run_status = ?, last_run_notes = ?, next_run_at = ?, updated_at = ?
		WHERE id = ?
	`
	now := time.Now()
	_, err := r.db.Exec(query, lastRunAt, status, notes, nextRunAt, now, id)
	if err != nil {
		return fmt.Errorf("failed to update schedule job run status: %w", err)
	}
	return nil
}

func (r *Repository) DeleteScheduleJob(id int64) error {
	_, err := r.db.Exec("DELETE FROM schedule_jobs WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("failed to delete schedule job: %w", err)
	}
	return nil
}

func (r *Repository) GetActiveScheduleJobs() ([]*ScheduleJob, error) {
	query := `
		SELECT id, target_id, name, description, is_active, schedule_config, backup_options,
		       meta_config, last_run_at, last_run_status, last_run_notes, next_run_at,
		       created_at, updated_at
		FROM schedule_jobs WHERE is_active = 1 ORDER BY next_run_at ASC
	`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query active schedule jobs: %w", err)
	}
	defer rows.Close()

	var jobs []*ScheduleJob
	for rows.Next() {
		job := &ScheduleJob{}
		err := rows.Scan(&job.ID, &job.TargetID, &job.Name, &job.Description, &job.IsActive,
			&job.ScheduleConfig, &job.BackupOptions, &job.MetaConfig, &job.LastRunAt,
			&job.LastRunStatus, &job.LastRunNotes, &job.NextRunAt, &job.CreatedAt, &job.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan active schedule job: %w", err)
		}
		jobs = append(jobs, job)
	}

	return jobs, nil
}

// Config methods
func (r *Repository) GetConfig(key string) (*AppConfig, error) {
	query := `SELECT id, key, value, created_at, updated_at FROM app_config WHERE key = ?`
	config := &AppConfig{}
	err := r.db.QueryRow(query, key).Scan(&config.ID, &config.Key, &config.Value, &config.CreatedAt, &config.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Config not found
		}
		return nil, fmt.Errorf("failed to get config: %w", err)
	}
	return config, nil
}

func (r *Repository) SetConfig(key, value string) error {
	// Check if config exists
	existing, err := r.GetConfig(key)
	if err != nil {
		return err
	}
	
	now := time.Now()
	if existing != nil {
		// Update existing config
		query := `UPDATE app_config SET value = ?, updated_at = ? WHERE key = ?`
		_, err = r.db.Exec(query, value, now, key)
		if err != nil {
			return fmt.Errorf("failed to update config: %w", err)
		}
	} else {
		// Create new config
		query := `INSERT INTO app_config (key, value, created_at, updated_at) VALUES (?, ?, ?, ?)`
		_, err = r.db.Exec(query, key, value, now, now)
		if err != nil {
			return fmt.Errorf("failed to create config: %w", err)
		}
	}
	
	return nil
}

func (r *Repository) GetAllConfigs() ([]*AppConfig, error) {
	query := `SELECT id, key, value, created_at, updated_at FROM app_config ORDER BY key`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query configs: %w", err)
	}
	defer rows.Close()

	var configs []*AppConfig
	for rows.Next() {
		config := &AppConfig{}
		err := rows.Scan(&config.ID, &config.Key, &config.Value, &config.CreatedAt, &config.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan config: %w", err)
		}
		configs = append(configs, config)
	}

	return configs, nil
}