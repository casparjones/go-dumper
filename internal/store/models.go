package store

import (
	"time"
)

type Target struct {
	ID                int64     `json:"id" db:"id"`
	Name              string    `json:"name" db:"name"`
	Host              string    `json:"host" db:"host"`
	Port              int       `json:"port" db:"port"`
	User              string    `json:"user" db:"user"`
	PasswordEnc       string    `json:"-" db:"password_enc"`
	Comment           string    `json:"comment" db:"comment"`
	ScheduleTime      string    `json:"schedule_time" db:"schedule_time"`
	RetentionDays     int       `json:"retention_days" db:"retention_days"`
	AutoCompress      bool      `json:"auto_compress" db:"auto_compress"`
	DatabaseMode      string    `json:"database_mode" db:"database_mode"` // "all" or "selected"
	SelectedDatabases string    `json:"selected_databases" db:"selected_databases"` // JSON array when mode="selected"
	CreatedAt         time.Time `json:"created_at" db:"created_at"`
	UpdatedAt         time.Time `json:"updated_at" db:"updated_at"`
}

type Backup struct {
	ID           int64      `json:"id" db:"id"`
	TargetID     int64      `json:"target_id" db:"target_id"`
	DatabaseName string     `json:"database_name" db:"database_name"`
	StartedAt    time.Time  `json:"started_at" db:"started_at"`
	FinishedAt   *time.Time `json:"finished_at" db:"finished_at"`
	SizeBytes    int64      `json:"size_bytes" db:"size_bytes"`
	Status       string     `json:"status" db:"status"`
	FilePath     string     `json:"file_path" db:"file_path"`
	Notes        string     `json:"notes" db:"notes"`
}

const (
	BackupStatusRunning = "running"
	BackupStatusSuccess = "success"
	BackupStatusFailed  = "failed"
)

const (
	DatabaseModeAll      = "all"
	DatabaseModeSelected = "selected"
)

type DatabaseInfo struct {
	Name string `json:"name"`
}

type ScheduleJob struct {
	ID              int64      `json:"id" db:"id"`
	TargetID        int64      `json:"target_id" db:"target_id"`
	Name            string     `json:"name" db:"name"`
	Description     string     `json:"description" db:"description"`
	IsActive        bool       `json:"is_active" db:"is_active"`
	ScheduleConfig  string     `json:"schedule_config" db:"schedule_config"`   // JSON with frequency, minutes, hours, etc.
	BackupOptions   string     `json:"backup_options" db:"backup_options"`     // JSON with compress, databases, etc.
	MetaConfig      string     `json:"meta_config" db:"meta_config"`           // JSON for future extensions (minio, nextcloud, etc.)
	LastRunAt       *time.Time `json:"last_run_at" db:"last_run_at"`
	LastRunStatus   string     `json:"last_run_status" db:"last_run_status"`
	LastRunNotes    string     `json:"last_run_notes" db:"last_run_notes"`
	NextRunAt       *time.Time `json:"next_run_at" db:"next_run_at"`
	CreatedAt       time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at" db:"updated_at"`
}

const (
	JobStatusPending = "pending"
	JobStatusRunning = "running"
	JobStatusSuccess = "success"
	JobStatusFailed  = "failed"
)

type AppConfig struct {
	ID        int64     `json:"id" db:"id"`
	Key       string    `json:"key" db:"key"`
	Value     string    `json:"value" db:"value"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

const (
	ConfigKeyTheme = "theme"
)