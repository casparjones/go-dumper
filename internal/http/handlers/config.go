package handlers

import (
	"net/http"

	"github.com/casparjones/go-dumper/internal/store"
	"github.com/gin-gonic/gin"
)

type ConfigHandler struct {
	repo *store.Repository
}

func NewConfigHandler(repo *store.Repository) *ConfigHandler {
	return &ConfigHandler{repo: repo}
}

type SetConfigRequest struct {
	Key   string `json:"key" binding:"required"`
	Value string `json:"value" binding:"required"`
}

func (h *ConfigHandler) GetConfig(c *gin.Context) {
	key := c.Param("key")
	if key == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Key parameter is required"})
		return
	}

	config, err := h.repo.GetConfig(key)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if config == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Configuration not found"})
		return
	}

	c.JSON(http.StatusOK, config)
}

func (h *ConfigHandler) SetConfig(c *gin.Context) {
	var req SetConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.repo.SetConfig(req.Key, req.Value)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Configuration saved successfully"})
}

func (h *ConfigHandler) GetAllConfigs(c *gin.Context) {
	configs, err := h.repo.GetAllConfigs()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, configs)
}

func (h *ConfigHandler) GetTheme(c *gin.Context) {
	config, err := h.repo.GetConfig(store.ConfigKeyTheme)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Default theme if not set
	theme := "dim"
	if config != nil {
		theme = config.Value
	}

	c.JSON(http.StatusOK, gin.H{"theme": theme})
}

func (h *ConfigHandler) SetTheme(c *gin.Context) {
	var req struct {
		Theme string `json:"theme" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.repo.SetConfig(store.ConfigKeyTheme, req.Theme)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Theme saved successfully", "theme": req.Theme})
}