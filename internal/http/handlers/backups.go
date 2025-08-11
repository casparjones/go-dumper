package handlers

import (
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/casparjones/go-dumper/internal/backup"
	"github.com/casparjones/go-dumper/internal/store"
	"github.com/gin-gonic/gin"
)

type BackupsHandler struct {
	repo     *store.Repository
	restorer *backup.Restorer
}

func NewBackupsHandler(repo *store.Repository, restorer *backup.Restorer) *BackupsHandler {
	return &BackupsHandler{
		repo:     repo,
		restorer: restorer,
	}
}

type BackupWithTargetInfo struct {
	*store.Backup
	TargetName string `json:"target_name"`
	Host       string `json:"host"`
	Port       int    `json:"port"`
}

func (h *BackupsHandler) GetAllBackups(c *gin.Context) {
	backups, err := h.repo.GetAllBackups()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Get target information for each backup
	var backupsWithTargetInfo []BackupWithTargetInfo
	for _, backup := range backups {
		target, err := h.repo.GetTarget(backup.TargetID)
		if err != nil {
			// Skip this backup if we can't find the target
			continue
		}

		backupWithInfo := BackupWithTargetInfo{
			Backup:     backup,
			TargetName: target.Name,
			Host:       target.Host,
			Port:       target.Port,
		}
		backupsWithTargetInfo = append(backupsWithTargetInfo, backupWithInfo)
	}

	c.JSON(http.StatusOK, backupsWithTargetInfo)
}

func (h *BackupsHandler) DownloadBackup(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid backup ID"})
		return
	}

	backup, err := h.repo.GetBackup(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Backup not found"})
		return
	}

	if backup.Status != store.BackupStatusSuccess {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Backup is not complete or failed"})
		return
	}

	if _, err := os.Stat(backup.FilePath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Backup file not found"})
		return
	}

	filename := filepath.Base(backup.FilePath)
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Header("Content-Type", "application/gzip")

	c.File(backup.FilePath)
}

type RestoreBackupRequest struct {
	CreateDatabase bool `json:"create_database"`
}

func (h *BackupsHandler) RestoreBackup(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid backup ID"})
		return
	}

	backup, err := h.repo.GetBackup(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Backup not found"})
		return
	}

	if backup.Status != store.BackupStatusSuccess {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot restore incomplete or failed backup"})
		return
	}

	// Parse restore options from request body (optional)
	var req RestoreBackupRequest
	if c.Request.ContentLength > 0 {
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}
	}

	go func() {
		var restoreErr error
		if req.CreateDatabase {
			restoreErr = h.restorer.RestoreBackupWithOptions(c.Request.Context(), id, true)
		} else {
			restoreErr = h.restorer.RestoreBackup(c.Request.Context(), id)
		}
		
		if restoreErr != nil {
			// Log the error for debugging
			// In a production system, you might want to update a restore status in the database
			// or notify the user through websockets/polling
		}
	}()

	c.JSON(http.StatusAccepted, gin.H{
		"message":       "Restore started",
		"backup_id":     id,
		"database_name": backup.DatabaseName,
	})
}

func (h *BackupsHandler) DeleteBackup(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid backup ID"})
		return
	}

	backup, err := h.repo.GetBackup(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Backup not found"})
		return
	}

	if backup.FilePath != "" {
		os.Remove(backup.FilePath)
	}

	if err := h.repo.DeleteBackup(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Backup deleted successfully"})
}
