package scheduler

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/casparjones/go-dumper/internal/backup"
	"github.com/casparjones/go-dumper/internal/config"
	"github.com/casparjones/go-dumper/internal/store"
)

type Scheduler struct {
	repo     *store.Repository
	dumper   *backup.Dumper
	ticker   *time.Ticker
	stopChan chan bool
}

type ScheduleConfig struct {
	Frequency   string `json:"frequency"`
	Minutes     []int  `json:"minutes"`
	Hours       []int  `json:"hours"`
	Weekdays    []int  `json:"weekdays"`
	DaysOfMonth []int  `json:"days_of_month"`
	Months      []int  `json:"months"`
}

type BackupOptions struct {
	Compress         bool     `json:"compress"`
	IncludeStructure bool     `json:"include_structure"`
	IncludeData      bool     `json:"include_data"`
	Databases        []string `json:"databases"`
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
	s.ticker = time.NewTicker(10 * time.Second)

	log.Println("Job scheduler started")

	for {
		select {
		case <-s.ticker.C:
			s.checkAndRunScheduledJobs()
		case <-s.stopChan:
			s.ticker.Stop()
			log.Println("Job scheduler stopped")
			return
		}
	}
}

func (s *Scheduler) Stop() {
	if s.stopChan != nil {
		close(s.stopChan)
	}
}

func (s *Scheduler) checkAndRunScheduledJobs() {
	jobs, err := s.repo.GetActiveScheduleJobs()
	if err != nil {
		log.Printf("Failed to get active jobs for scheduling: %v", err)
		return
	}

	now := time.Now()

	for _, job := range jobs {
		// Check if job is due to run
		if s.isJobDue(job, now) {
			log.Printf("Starting scheduled job: %s (ID: %d)", job.Name, job.ID)

			// Update job status to running
			s.updateJobRunStatus(job.ID, store.JobStatusRunning, "Job execution started", &now, nil)

			// Execute job in background
			go s.executeJob(job)
		}
	}
}

func (s *Scheduler) isJobDue(job *store.ScheduleJob, now time.Time) bool {
	// If next_run_at is not set or job is already running, skip
	if job.NextRunAt == nil || job.LastRunStatus == store.JobStatusRunning {
		return false
	}

	// Check if it's time to run (within 1 minute window)
	return now.After(*job.NextRunAt)
}

func (s *Scheduler) executeJob(job *store.ScheduleJob) {
	startTime := time.Now()
	var status, notes string

	// Parse backup options
	var backupOptions BackupOptions
	if err := json.Unmarshal([]byte(job.BackupOptions), &backupOptions); err != nil {
		status = store.JobStatusFailed
		notes = fmt.Sprintf("Failed to parse backup options: %v", err)
		log.Printf("Job %d failed: %s", job.ID, notes)
	} else {
		// Execute backup
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
		defer cancel()

		_, err := s.dumper.CreateBackup(ctx, job.TargetID)
		if err != nil {
			status = store.JobStatusFailed
			notes = fmt.Sprintf("Backup failed: %v", err)
			log.Printf("Job %d failed: %s", job.ID, notes)
		} else {
			status = store.JobStatusSuccess
			notes = "Backup completed successfully"
			log.Printf("Job %d completed successfully", job.ID)
		}
	}

	// Calculate next run time
	nextRun := s.calculateNextRun(job)

	// Update job status
	if err := s.updateJobRunStatus(job.ID, status, notes, &startTime, nextRun); err != nil {
		log.Printf("Failed to update job %d status: %v", job.ID, err)
	}
}

func (s *Scheduler) updateJobRunStatus(jobID int64, status, notes string, lastRunAt, nextRunAt *time.Time) error {
	return s.repo.UpdateScheduleJobRunStatus(jobID, status, notes, lastRunAt, nextRunAt)
}

func (s *Scheduler) calculateNextRun(job *store.ScheduleJob) *time.Time {
	var config ScheduleConfig
	if err := json.Unmarshal([]byte(job.ScheduleConfig), &config); err != nil {
		log.Printf("Failed to parse schedule config for job %d: %v", job.ID, err)
		return nil
	}

	now := time.Now()

	switch config.Frequency {
	case "hourly":
		return s.calculateHourlyNext(now, config)
	case "daily":
		return s.calculateDailyNext(now, config)
	case "weekly":
		return s.calculateWeeklyNext(now, config)
	case "monthly":
		return s.calculateMonthlyNext(now, config)
	case "yearly":
		return s.calculateYearlyNext(now, config)
	default:
		log.Printf("Unknown frequency for job %d: %s", job.ID, config.Frequency)
		return nil
	}
}

func (s *Scheduler) calculateHourlyNext(now time.Time, config ScheduleConfig) *time.Time {
	if len(config.Minutes) == 0 {
		config.Minutes = []int{0}
	}

	// Find next minute
	currentMinute := now.Minute()
	for _, minute := range config.Minutes {
		if minute > currentMinute {
			next := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), minute, 0, 0, now.Location())
			return &next
		}
	}

	// No minute found in current hour, use first minute of next hour
	next := now.Add(time.Hour)
	next = time.Date(next.Year(), next.Month(), next.Day(), next.Hour(), config.Minutes[0], 0, 0, next.Location())
	return &next
}

func (s *Scheduler) calculateDailyNext(now time.Time, config ScheduleConfig) *time.Time {
	if len(config.Hours) == 0 {
		config.Hours = []int{0}
	}
	if len(config.Minutes) == 0 {
		config.Minutes = []int{0}
	}

	// Find next time today
	for _, hour := range config.Hours {
		for _, minute := range config.Minutes {
			next := time.Date(now.Year(), now.Month(), now.Day(), hour, minute, 0, 0, now.Location())
			if next.After(now) {
				return &next
			}
		}
	}

	// No time found today, use first time tomorrow
	tomorrow := now.AddDate(0, 0, 1)
	next := time.Date(tomorrow.Year(), tomorrow.Month(), tomorrow.Day(), config.Hours[0], config.Minutes[0], 0, 0, tomorrow.Location())
	return &next
}

func (s *Scheduler) calculateWeeklyNext(now time.Time, config ScheduleConfig) *time.Time {
	if len(config.Weekdays) == 0 {
		config.Weekdays = []int{1} // Default to Monday
	}
	if len(config.Hours) == 0 {
		config.Hours = []int{0}
	}
	if len(config.Minutes) == 0 {
		config.Minutes = []int{0}
	}

	// Convert Sunday=0 to Monday=1 format
	currentWeekday := int(now.Weekday())
	if currentWeekday == 0 {
		currentWeekday = 7
	}

	// Check remaining days this week
	for day := currentWeekday; day <= 7; day++ {
		for _, targetWeekday := range config.Weekdays {
			if targetWeekday == day {
				if day == currentWeekday {
					// Today - check if any time is still available
					for _, hour := range config.Hours {
						for _, minute := range config.Minutes {
							targetTime := time.Date(now.Year(), now.Month(), now.Day(), hour, minute, 0, 0, now.Location())
							if targetTime.After(now) {
								return &targetTime
							}
						}
					}
				} else {
					// Future day this week
					daysAhead := day - currentWeekday
					targetDate := now.AddDate(0, 0, daysAhead)
					targetTime := time.Date(targetDate.Year(), targetDate.Month(), targetDate.Day(), config.Hours[0], config.Minutes[0], 0, 0, targetDate.Location())
					return &targetTime
				}
			}
		}
	}

	// No day found this week, find first day next week
	for _, weekday := range config.Weekdays {
		daysAhead := (7 - currentWeekday) + weekday
		targetDate := now.AddDate(0, 0, daysAhead)
		targetTime := time.Date(targetDate.Year(), targetDate.Month(), targetDate.Day(), config.Hours[0], config.Minutes[0], 0, 0, targetDate.Location())
		return &targetTime
	}

	// Fallback
	next := now.AddDate(0, 0, 7)
	return &next
}

func (s *Scheduler) calculateMonthlyNext(now time.Time, config ScheduleConfig) *time.Time {
	if len(config.DaysOfMonth) == 0 {
		config.DaysOfMonth = []int{1} // Default to 1st of month
	}
	if len(config.Hours) == 0 {
		config.Hours = []int{0}
	}
	if len(config.Minutes) == 0 {
		config.Minutes = []int{0}
	}

	// Check remaining days this month
	currentDay := now.Day()
	for _, targetDay := range config.DaysOfMonth {
		if targetDay >= currentDay {
			if targetDay == currentDay {
				// Today - check if any time is still available
				for _, hour := range config.Hours {
					for _, minute := range config.Minutes {
						targetTime := time.Date(now.Year(), now.Month(), now.Day(), hour, minute, 0, 0, now.Location())
						if targetTime.After(now) {
							return &targetTime
						}
					}
				}
			} else {
				// Future day this month - check if day exists in current month
				targetTime := time.Date(now.Year(), now.Month(), targetDay, config.Hours[0], config.Minutes[0], 0, 0, now.Location())
				if targetTime.Month() == now.Month() { // Day exists in current month
					return &targetTime
				}
			}
		}
	}

	// No valid day found this month, use first day of next month
	nextMonth := now.AddDate(0, 1, 0)
	firstDay := config.DaysOfMonth[0]
	targetTime := time.Date(nextMonth.Year(), nextMonth.Month(), firstDay, config.Hours[0], config.Minutes[0], 0, 0, nextMonth.Location())

	// Ensure the day exists in the target month
	if targetTime.Month() != nextMonth.Month() {
		// Day doesn't exist (e.g., Feb 30), use last day of month
		lastDay := time.Date(nextMonth.Year(), nextMonth.Month()+1, 0, config.Hours[0], config.Minutes[0], 0, 0, nextMonth.Location())
		return &lastDay
	}

	return &targetTime
}

func (s *Scheduler) calculateYearlyNext(now time.Time, config ScheduleConfig) *time.Time {
	if len(config.Months) == 0 {
		config.Months = []int{1} // Default to January
	}
	if len(config.DaysOfMonth) == 0 {
		config.DaysOfMonth = []int{1}
	}
	if len(config.Hours) == 0 {
		config.Hours = []int{0}
	}
	if len(config.Minutes) == 0 {
		config.Minutes = []int{0}
	}

	// Check remaining months this year
	currentMonth := int(now.Month())
	for _, targetMonth := range config.Months {
		if targetMonth >= currentMonth {
			for _, targetDay := range config.DaysOfMonth {
				if targetMonth == currentMonth && targetDay >= now.Day() {
					if targetDay == now.Day() {
						// Today - check if any time is still available
						for _, hour := range config.Hours {
							for _, minute := range config.Minutes {
								targetTime := time.Date(now.Year(), time.Month(targetMonth), now.Day(), hour, minute, 0, 0, now.Location())
								if targetTime.After(now) {
									return &targetTime
								}
							}
						}
					} else {
						// Future day this month/year
						targetTime := time.Date(now.Year(), time.Month(targetMonth), targetDay, config.Hours[0], config.Minutes[0], 0, 0, now.Location())
						if targetTime.Month() == time.Month(targetMonth) { // Day exists
							return &targetTime
						}
					}
				} else if targetMonth > currentMonth {
					// Future month this year
					targetTime := time.Date(now.Year(), time.Month(targetMonth), targetDay, config.Hours[0], config.Minutes[0], 0, 0, now.Location())
					if targetTime.Month() == time.Month(targetMonth) { // Day exists
						return &targetTime
					}
				}
			}
		}
	}

	// No valid date found this year, use first date next year
	nextYear := now.AddDate(1, 0, 0)
	targetTime := time.Date(nextYear.Year(), time.Month(config.Months[0]), config.DaysOfMonth[0], config.Hours[0], config.Minutes[0], 0, 0, nextYear.Location())

	// Ensure the day exists in the target month
	if targetTime.Month() != time.Month(config.Months[0]) {
		// Day doesn't exist, use last day of month
		lastDay := time.Date(nextYear.Year(), time.Month(config.Months[0])+1, 0, config.Hours[0], config.Minutes[0], 0, 0, nextYear.Location())
		return &lastDay
	}

	return &targetTime
}
