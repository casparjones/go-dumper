package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/casparjones/go-dumper/internal/backup"
	"github.com/casparjones/go-dumper/internal/store"
	"github.com/gin-gonic/gin"
)

type JobsHandler struct {
	repo    *store.Repository
	dumper  *backup.Dumper
}

func NewJobsHandler(repo *store.Repository, dumper *backup.Dumper) *JobsHandler {
	return &JobsHandler{
		repo:   repo,
		dumper: dumper,
	}
}

type CreateJobRequest struct {
	TargetID       int64                  `json:"target_id" binding:"required"`
	Name           string                 `json:"name" binding:"required"`
	Description    string                 `json:"description"`
	ScheduleConfig ScheduleConfig         `json:"schedule_config" binding:"required"`
	BackupOptions  BackupOptions          `json:"backup_options" binding:"required"`
	MetaConfig     map[string]interface{} `json:"meta_config"`
}

type ScheduleConfig struct {
	Frequency    string `json:"frequency" binding:"required"`
	Minutes      []int  `json:"minutes"`
	Hours        []int  `json:"hours"`
	Weekdays     []int  `json:"weekdays"`
	DaysOfMonth  []int  `json:"days_of_month"`
	Months       []int  `json:"months"`
}

type BackupOptions struct {
	Compress         bool     `json:"compress"`
	IncludeStructure bool     `json:"include_structure"`
	IncludeData      bool     `json:"include_data"`
	Databases        []string `json:"databases"`
}

type UpdateJobRequest struct {
	Name           string                 `json:"name" binding:"required"`
	Description    string                 `json:"description"`
	IsActive       bool                   `json:"is_active"`
	ScheduleConfig ScheduleConfig         `json:"schedule_config" binding:"required"`
	BackupOptions  BackupOptions          `json:"backup_options" binding:"required"`
	MetaConfig     map[string]interface{} `json:"meta_config"`
}

type JobResponse struct {
	*store.ScheduleJob
	Target *store.Target `json:"target,omitempty"`
}

// GetJobs returns all schedule jobs
func (h *JobsHandler) GetJobs(c *gin.Context) {
	jobs, err := h.repo.GetScheduleJobs()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get jobs"})
		return
	}

	// Enrich jobs with target information
	jobResponses := make([]*JobResponse, len(jobs))
	for i, job := range jobs {
		target, _ := h.repo.GetTarget(job.TargetID) // Ignore error for optional target info
		jobResponses[i] = &JobResponse{
			ScheduleJob: job,
			Target:      target,
		}
	}

	c.JSON(http.StatusOK, jobResponses)
}

// GetJob returns a specific schedule job
func (h *JobsHandler) GetJob(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid job ID"})
		return
	}

	job, err := h.repo.GetScheduleJob(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Job not found"})
		return
	}

	target, _ := h.repo.GetTarget(job.TargetID)
	response := &JobResponse{
		ScheduleJob: job,
		Target:      target,
	}

	c.JSON(http.StatusOK, response)
}

// CreateJob creates a new schedule job
func (h *JobsHandler) CreateJob(c *gin.Context) {
	var req CreateJobRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate target exists
	_, err := h.repo.GetTarget(req.TargetID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Target not found"})
		return
	}

	// Convert structs to JSON strings
	scheduleConfigJSON, err := json.Marshal(req.ScheduleConfig)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to serialize schedule config"})
		return
	}

	backupOptionsJSON, err := json.Marshal(req.BackupOptions)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to serialize backup options"})
		return
	}

	metaConfigJSON, err := json.Marshal(req.MetaConfig)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to serialize meta config"})
		return
	}

	// Calculate next run time
	nextRun := h.calculateNextRun(req.ScheduleConfig)

	job := &store.ScheduleJob{
		TargetID:       req.TargetID,
		Name:           req.Name,
		Description:    req.Description,
		IsActive:       true,
		ScheduleConfig: string(scheduleConfigJSON),
		BackupOptions:  string(backupOptionsJSON),
		MetaConfig:     string(metaConfigJSON),
		NextRunAt:      nextRun,
	}

	if err := h.repo.CreateScheduleJob(job); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create job"})
		return
	}

	target, _ := h.repo.GetTarget(job.TargetID)
	response := &JobResponse{
		ScheduleJob: job,
		Target:      target,
	}

	c.JSON(http.StatusCreated, response)
}

// UpdateJob updates an existing schedule job
func (h *JobsHandler) UpdateJob(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid job ID"})
		return
	}

	var req UpdateJobRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	job, err := h.repo.GetScheduleJob(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Job not found"})
		return
	}

	// Convert structs to JSON strings
	scheduleConfigJSON, err := json.Marshal(req.ScheduleConfig)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to serialize schedule config"})
		return
	}

	backupOptionsJSON, err := json.Marshal(req.BackupOptions)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to serialize backup options"})
		return
	}

	metaConfigJSON, err := json.Marshal(req.MetaConfig)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to serialize meta config"})
		return
	}

	// Update job fields
	job.Name = req.Name
	job.Description = req.Description
	job.IsActive = req.IsActive
	job.ScheduleConfig = string(scheduleConfigJSON)
	job.BackupOptions = string(backupOptionsJSON)
	job.MetaConfig = string(metaConfigJSON)

	// Recalculate next run time if schedule changed
	job.NextRunAt = h.calculateNextRun(req.ScheduleConfig)

	if err := h.repo.UpdateScheduleJob(job); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update job"})
		return
	}

	target, _ := h.repo.GetTarget(job.TargetID)
	response := &JobResponse{
		ScheduleJob: job,
		Target:      target,
	}

	c.JSON(http.StatusOK, response)
}

// DeleteJob deletes a schedule job
func (h *JobsHandler) DeleteJob(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid job ID"})
		return
	}

	_, err = h.repo.GetScheduleJob(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Job not found"})
		return
	}

	if err := h.repo.DeleteScheduleJob(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete job"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Job deleted successfully"})
}

// RunJobNow executes a job immediately
func (h *JobsHandler) RunJobNow(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid job ID"})
		return
	}

	job, err := h.repo.GetScheduleJob(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Job not found"})
		return
	}

	// Update job status to running
	now := time.Now()
	nextRun := h.calculateNextRun(parseScheduleConfig(job.ScheduleConfig))
	if err := h.repo.UpdateScheduleJobRunStatus(id, store.JobStatusRunning, "Manual execution started", &now, nextRun); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update job status"})
		return
	}

	// Start backup in background
	go h.executeJob(job)

	c.JSON(http.StatusOK, gin.H{
		"message": "Job execution started",
		"job_id":  id,
	})
}

// executeJob runs the actual backup process for a job
func (h *JobsHandler) executeJob(job *store.ScheduleJob) {
	startTime := time.Now()
	var status, notes string

	// Parse backup options
	var backupOptions BackupOptions
	if err := json.Unmarshal([]byte(job.BackupOptions), &backupOptions); err != nil {
		status = store.JobStatusFailed
		notes = fmt.Sprintf("Failed to parse backup options: %v", err)
	} else {
		// Execute backup
		_, err := h.dumper.CreateBackup(nil, job.TargetID)
		if err != nil {
			status = store.JobStatusFailed
			notes = fmt.Sprintf("Backup failed: %v", err)
		} else {
			status = store.JobStatusSuccess
			notes = "Backup completed successfully"
		}
	}

	// Calculate next run time
	scheduleConfig := parseScheduleConfig(job.ScheduleConfig)
	nextRun := h.calculateNextRun(scheduleConfig)

	// Update job status
	h.repo.UpdateScheduleJobRunStatus(job.ID, status, notes, &startTime, nextRun)
}

// calculateNextRun calculates the next execution time based on schedule config
func (h *JobsHandler) calculateNextRun(config ScheduleConfig) *time.Time {
	now := time.Now()
	
	switch config.Frequency {
	case "hourly":
		return h.calculateHourlyNext(now, config)
	case "daily":
		return h.calculateDailyNext(now, config)
	case "weekly":
		return h.calculateWeeklyNext(now, config)
	case "monthly":
		return h.calculateMonthlyNext(now, config)
	case "yearly":
		return h.calculateYearlyNext(now, config)
	default:
		// Default to daily at midnight
		next := now.AddDate(0, 0, 1)
		next = time.Date(next.Year(), next.Month(), next.Day(), 0, 0, 0, 0, next.Location())
		return &next
	}
}

func (h *JobsHandler) calculateHourlyNext(now time.Time, config ScheduleConfig) *time.Time {
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

func (h *JobsHandler) calculateDailyNext(now time.Time, config ScheduleConfig) *time.Time {
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

func (h *JobsHandler) calculateWeeklyNext(now time.Time, config ScheduleConfig) *time.Time {
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
				var targetTime time.Time
				if day == currentWeekday {
					// Today - check if any time is still available
					for _, hour := range config.Hours {
						for _, minute := range config.Minutes {
							targetTime = time.Date(now.Year(), now.Month(), now.Day(), hour, minute, 0, 0, now.Location())
							if targetTime.After(now) {
								return &targetTime
							}
						}
					}
				} else {
					// Future day this week
					daysAhead := day - currentWeekday
					targetDate := now.AddDate(0, 0, daysAhead)
					targetTime = time.Date(targetDate.Year(), targetDate.Month(), targetDate.Day(), config.Hours[0], config.Minutes[0], 0, 0, targetDate.Location())
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

func (h *JobsHandler) calculateMonthlyNext(now time.Time, config ScheduleConfig) *time.Time {
	// Similar logic for monthly - simplified for now
	next := now.AddDate(0, 1, 0)
	return &next
}

func (h *JobsHandler) calculateYearlyNext(now time.Time, config ScheduleConfig) *time.Time {
	// Similar logic for yearly - simplified for now
	next := now.AddDate(1, 0, 0)
	return &next
}

// parseScheduleConfig parses JSON schedule config
func parseScheduleConfig(configStr string) ScheduleConfig {
	var config ScheduleConfig
	json.Unmarshal([]byte(configStr), &config)
	return config
}