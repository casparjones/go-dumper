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

	go func() {
		if err := h.restorer.RestoreBackup(c.Request.Context(), id); err != nil {
			// In a production system, you might want to log this error
			// or update a restore status in the database
		}
	}()

	c.JSON(http.StatusAccepted, gin.H{
		"message":   "Restore started",
		"backup_id": id,
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
