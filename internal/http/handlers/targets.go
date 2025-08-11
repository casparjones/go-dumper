package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/casparjones/go-dumper/internal/backup"
	"github.com/casparjones/go-dumper/internal/store"
	"github.com/gin-gonic/gin"
)

type TargetsHandler struct {
	repo   *store.Repository
	dumper *backup.Dumper
}

type CreateTargetRequest struct {
	Name              string   `json:"name" binding:"required"`
	Host              string   `json:"host" binding:"required"`
	Port              int      `json:"port" binding:"required"`
	User              string   `json:"user" binding:"required"`
	Password          string   `json:"password" binding:"required"`
	Comment           string   `json:"comment"`
	ScheduleTime      string   `json:"schedule_time"`
	RetentionDays     int      `json:"retention_days"`
	AutoCompress      bool     `json:"auto_compress"`
	DatabaseMode      string   `json:"database_mode"`
	SelectedDatabases []string `json:"selected_databases,omitempty"`
}

type UpdateTargetRequest struct {
	Name              string   `json:"name" binding:"required"`
	Host              string   `json:"host" binding:"required"`
	Port              int      `json:"port" binding:"required"`
	User              string   `json:"user" binding:"required"`
	Password          string   `json:"password,omitempty"`
	Comment           string   `json:"comment"`
	ScheduleTime      string   `json:"schedule_time"`
	RetentionDays     int      `json:"retention_days"`
	AutoCompress      bool     `json:"auto_compress"`
	DatabaseMode      string   `json:"database_mode"`
	SelectedDatabases []string `json:"selected_databases,omitempty"`
}

type TargetResponse struct {
	ID                int64    `json:"id"`
	Name              string   `json:"name"`
	Host              string   `json:"host"`
	Port              int      `json:"port"`
	User              string   `json:"user"`
	Comment           string   `json:"comment"`
	ScheduleTime      string   `json:"schedule_time"`
	RetentionDays     int      `json:"retention_days"`
	AutoCompress      bool     `json:"auto_compress"`
	DatabaseMode      string   `json:"database_mode"`
	SelectedDatabases []string `json:"selected_databases,omitempty"`
	CreatedAt         string   `json:"created_at"`
	UpdatedAt         string   `json:"updated_at"`
}

func NewTargetsHandler(repo *store.Repository, dumper *backup.Dumper) *TargetsHandler {
	return &TargetsHandler{
		repo:   repo,
		dumper: dumper,
	}
}

func (h *TargetsHandler) GetTargets(c *gin.Context) {
	targets, err := h.repo.GetTargets()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := make([]*TargetResponse, len(targets))
	for i, target := range targets {
		response[i] = h.targetToResponse(target)
	}

	c.JSON(http.StatusOK, response)
}

func (h *TargetsHandler) GetTarget(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid target ID"})
		return
	}

	target, err := h.repo.GetTarget(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Target not found"})
		return
	}

	c.JSON(http.StatusOK, h.targetToResponse(target))
}

func (h *TargetsHandler) CreateTarget(c *gin.Context) {
	var req CreateTargetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	encryptedPassword, err := store.EncryptPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to encrypt password"})
		return
	}

	// Set default database mode if not provided
	if req.DatabaseMode == "" {
		req.DatabaseMode = store.DatabaseModeAll
	}

	var selectedDatabasesJson string
	if req.DatabaseMode == store.DatabaseModeSelected && len(req.SelectedDatabases) > 0 {
		selectedDbBytes, err := json.Marshal(req.SelectedDatabases)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to serialize selected databases"})
			return
		}
		selectedDatabasesJson = string(selectedDbBytes)
	}

	target := &store.Target{
		Name:              req.Name,
		Host:              req.Host,
		Port:              req.Port,
		User:              req.User,
		PasswordEnc:       encryptedPassword,
		Comment:           req.Comment,
		ScheduleTime:      req.ScheduleTime,
		RetentionDays:     req.RetentionDays,
		AutoCompress:      req.AutoCompress,
		DatabaseMode:      req.DatabaseMode,
		SelectedDatabases: selectedDatabasesJson,
	}

	if target.RetentionDays <= 0 {
		target.RetentionDays = 30
	}

	if err := h.repo.CreateTarget(target); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, h.targetToResponse(target))
}

func (h *TargetsHandler) UpdateTarget(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid target ID"})
		return
	}

	target, err := h.repo.GetTarget(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Target not found"})
		return
	}

	var req UpdateTargetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	target.Name = req.Name
	target.Host = req.Host
	target.Port = req.Port
	target.User = req.User
	target.Comment = req.Comment
	target.ScheduleTime = req.ScheduleTime
	target.RetentionDays = req.RetentionDays
	target.AutoCompress = req.AutoCompress

	// Set default database mode if not provided
	if req.DatabaseMode == "" {
		req.DatabaseMode = store.DatabaseModeAll
	}
	target.DatabaseMode = req.DatabaseMode

	// Handle selected databases
	var selectedDatabasesJson string
	if req.DatabaseMode == store.DatabaseModeSelected && len(req.SelectedDatabases) > 0 {
		selectedDbBytes, err := json.Marshal(req.SelectedDatabases)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to serialize selected databases"})
			return
		}
		selectedDatabasesJson = string(selectedDbBytes)
	}
	target.SelectedDatabases = selectedDatabasesJson

	if req.Password != "" {
		encryptedPassword, err := store.EncryptPassword(req.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to encrypt password"})
			return
		}
		target.PasswordEnc = encryptedPassword
	}

	if target.RetentionDays <= 0 {
		target.RetentionDays = 30
	}

	if err := h.repo.UpdateTarget(target); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, h.targetToResponse(target))
}

func (h *TargetsHandler) DeleteTarget(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid target ID"})
		return
	}

	if err := h.repo.DeleteTarget(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Target deleted successfully"})
}

func (h *TargetsHandler) CreateBackup(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid target ID"})
		return
	}

	backup, err := h.dumper.CreateBackup(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{
		"message":   "Backup started",
		"backup_id": backup.ID,
		"status":    backup.Status,
	})
}

func (h *TargetsHandler) GetTargetBackups(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid target ID"})
		return
	}

	backups, err := h.repo.GetBackupsByTarget(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, backups)
}

type DiscoverDatabasesRequest struct {
	Host     string `json:"host" binding:"required"`
	Port     int    `json:"port" binding:"required"`
	User     string `json:"user" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (h *TargetsHandler) DiscoverDatabases(c *gin.Context) {
	var req DiscoverDatabasesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	databases, err := h.getDatabases(req.Host, req.Port, req.User, req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to discover databases: %v", err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"databases": databases})
}

func (h *TargetsHandler) getDatabases(host string, port int, user, password string) ([]store.DatabaseInfo, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/", user, password, host, port)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %v", err)
	}

	rows, err := db.Query("SHOW DATABASES")
	if err != nil {
		return nil, fmt.Errorf("failed to list databases: %v", err)
	}
	defer rows.Close()

	var databases []store.DatabaseInfo
	systemDatabases := map[string]bool{
		"information_schema": true,
		"performance_schema": true,
		"mysql":              true,
		"sys":                true,
	}

	for rows.Next() {
		var dbName string
		if err := rows.Scan(&dbName); err != nil {
			return nil, fmt.Errorf("failed to scan database name: %v", err)
		}

		// Skip system databases
		if !systemDatabases[dbName] {
			databases = append(databases, store.DatabaseInfo{Name: dbName})
		}
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over databases: %v", err)
	}

	return databases, nil
}

func (h *TargetsHandler) targetToResponse(target *store.Target) *TargetResponse {
	var selectedDatabases []string
	if target.DatabaseMode == store.DatabaseModeSelected && target.SelectedDatabases != "" {
		if err := json.Unmarshal([]byte(target.SelectedDatabases), &selectedDatabases); err != nil {
			// Log error but continue with empty array
			selectedDatabases = []string{}
		}
	}

	return &TargetResponse{
		ID:                target.ID,
		Name:              target.Name,
		Host:              target.Host,
		Port:              target.Port,
		User:              target.User,
		Comment:           target.Comment,
		ScheduleTime:      target.ScheduleTime,
		RetentionDays:     target.RetentionDays,
		AutoCompress:      target.AutoCompress,
		DatabaseMode:      target.DatabaseMode,
		SelectedDatabases: selectedDatabases,
		CreatedAt:         target.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt:         target.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}
}
