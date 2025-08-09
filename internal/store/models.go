package store

import (
	"time"
)

type Target struct {
	ID           int64     `json:"id" db:"id"`
	Name         string    `json:"name" db:"name"`
	Host         string    `json:"host" db:"host"`
	Port         int       `json:"port" db:"port"`
	DBName       string    `json:"db_name" db:"dbname"`
	User         string    `json:"user" db:"user"`
	PasswordEnc  string    `json:"-" db:"password_enc"`
	Comment      string    `json:"comment" db:"comment"`
	ScheduleTime string    `json:"schedule_time" db:"schedule_time"`
	RetentionDays int      `json:"retention_days" db:"retention_days"`
	AutoCompress bool      `json:"auto_compress" db:"auto_compress"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

type Backup struct {
	ID         int64     `json:"id" db:"id"`
	TargetID   int64     `json:"target_id" db:"target_id"`
	StartedAt  time.Time `json:"started_at" db:"started_at"`
	FinishedAt *time.Time `json:"finished_at" db:"finished_at"`
	SizeBytes  int64     `json:"size_bytes" db:"size_bytes"`
	Status     string    `json:"status" db:"status"`
	FilePath   string    `json:"file_path" db:"file_path"`
	Notes      string    `json:"notes" db:"notes"`
}

const (
	BackupStatusRunning = "running"
	BackupStatusSuccess = "success"
	BackupStatusFailed  = "failed"
)