package scheduler

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/user/go-dumper/internal/backup"
	"github.com/user/go-dumper/internal/config"
	"github.com/user/go-dumper/internal/store"
)

type Scheduler struct {
	repo     *store.Repository
	dumper   *backup.Dumper
	ticker   *time.Ticker
	stopChan chan bool
}

func New(db *sql.DB) *Scheduler {
	repo := store.NewRepository(db)
	backupDir := config.GetEnv("BACKUP_DIR", "/data/backups")
	dumper := backup.NewDumper(repo, backupDir)

	return &Scheduler{
		repo:     repo,
		dumper:   dumper,
		stopChan: make(chan bool),
	}
}

func (s *Scheduler) Start() {
	s.ticker = time.NewTicker(1 * time.Minute)
	
	log.Println("Scheduler started")
	
	for {
		select {
		case <-s.ticker.C:
			s.checkAndRunScheduledBackups()
		case <-s.stopChan:
			s.ticker.Stop()
			log.Println("Scheduler stopped")
			return
		}
	}
}

func (s *Scheduler) Stop() {
	if s.stopChan != nil {
		close(s.stopChan)
	}
}

func (s *Scheduler) checkAndRunScheduledBackups() {
	targets, err := s.repo.GetTargets()
	if err != nil {
		log.Printf("Failed to get targets for scheduling: %v", err)
		return
	}

	now := time.Now().UTC()
	currentTime := fmt.Sprintf("%02d:%02d", now.Hour(), now.Minute())

	for _, target := range targets {
		if target.ScheduleTime == "" {
			continue
		}

		if target.ScheduleTime == currentTime {
			if s.shouldRunBackup(target.ID) {
				log.Printf("Starting scheduled backup for target: %s", target.Name)
				
				ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
				defer cancel()

				_, err := s.dumper.CreateBackup(ctx, target.ID)
				if err != nil {
					log.Printf("Failed to start scheduled backup for target %s: %v", target.Name, err)
				}
			}
		}
	}
}

func (s *Scheduler) shouldRunBackup(targetID int64) bool {
	backups, err := s.repo.GetBackupsByTarget(targetID)
	if err != nil {
		log.Printf("Failed to get backups for target %d: %v", targetID, err)
		return true
	}

	now := time.Now().UTC()
	today := now.Truncate(24 * time.Hour)

	for _, backup := range backups {
		backupDay := backup.StartedAt.Truncate(24 * time.Hour)
		
		if backupDay.Equal(today) && 
		   (backup.Status == store.BackupStatusSuccess || backup.Status == store.BackupStatusRunning) {
			return false
		}
	}

	return true
}

func parseScheduleTime(scheduleTime string) (hour, minute int, err error) {
	parts := strings.Split(scheduleTime, ":")
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("invalid schedule time format, expected HH:MM")
	}

	hour, err = strconv.Atoi(parts[0])
	if err != nil {
		return 0, 0, fmt.Errorf("invalid hour: %v", err)
	}

	minute, err = strconv.Atoi(parts[1])
	if err != nil {
		return 0, 0, fmt.Errorf("invalid minute: %v", err)
	}

	if hour < 0 || hour > 23 {
		return 0, 0, fmt.Errorf("hour must be between 0 and 23")
	}

	if minute < 0 || minute > 59 {
		return 0, 0, fmt.Errorf("minute must be between 0 and 59")
	}

	return hour, minute, nil
}