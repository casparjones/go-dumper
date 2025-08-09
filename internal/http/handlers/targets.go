package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/user/go-dumper/internal/backup"
	"github.com/user/go-dumper/internal/store"
)

type TargetsHandler struct {
	repo    *store.Repository
	dumper  *backup.Dumper
}

type CreateTargetRequest struct {
	Name          string `json:"name" binding:"required"`
	Host          string `json:"host" binding:"required"`
	Port          int    `json:"port" binding:"required"`
	DBName        string `json:"db_name" binding:"required"`
	User          string `json:"user" binding:"required"`
	Password      string `json:"password" binding:"required"`
	Comment       string `json:"comment"`
	ScheduleTime  string `json:"schedule_time"`
	RetentionDays int    `json:"retention_days"`
	AutoCompress  bool   `json:"auto_compress"`
}

type UpdateTargetRequest struct {
	Name          string `json:"name" binding:"required"`
	Host          string `json:"host" binding:"required"`
	Port          int    `json:"port" binding:"required"`
	DBName        string `json:"db_name" binding:"required"`
	User          string `json:"user" binding:"required"`
	Password      string `json:"password,omitempty"`
	Comment       string `json:"comment"`
	ScheduleTime  string `json:"schedule_time"`
	RetentionDays int    `json:"retention_days"`
	AutoCompress  bool   `json:"auto_compress"`
}

type TargetResponse struct {
	ID            int64  `json:"id"`
	Name          string `json:"name"`
	Host          string `json:"host"`
	Port          int    `json:"port"`
	DBName        string `json:"db_name"`
	User          string `json:"user"`
	Comment       string `json:"comment"`
	ScheduleTime  string `json:"schedule_time"`
	RetentionDays int    `json:"retention_days"`
	AutoCompress  bool   `json:"auto_compress"`
	CreatedAt     string `json:"created_at"`
	UpdatedAt     string `json:"updated_at"`
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

	target := &store.Target{
		Name:          req.Name,
		Host:          req.Host,
		Port:          req.Port,
		DBName:        req.DBName,
		User:          req.User,
		PasswordEnc:   encryptedPassword,
		Comment:       req.Comment,
		ScheduleTime:  req.ScheduleTime,
		RetentionDays: req.RetentionDays,
		AutoCompress:  req.AutoCompress,
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
	target.DBName = req.DBName
	target.User = req.User
	target.Comment = req.Comment
	target.ScheduleTime = req.ScheduleTime
	target.RetentionDays = req.RetentionDays
	target.AutoCompress = req.AutoCompress

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

func (h *TargetsHandler) targetToResponse(target *store.Target) *TargetResponse {
	return &TargetResponse{
		ID:            target.ID,
		Name:          target.Name,
		Host:          target.Host,
		Port:          target.Port,
		DBName:        target.DBName,
		User:          target.User,
		Comment:       target.Comment,
		ScheduleTime:  target.ScheduleTime,
		RetentionDays: target.RetentionDays,
		AutoCompress:  target.AutoCompress,
		CreatedAt:     target.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt:     target.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}
}