package router

import (
	"database/sql"

	"github.com/gin-gonic/gin"
	"github.com/user/go-dumper/internal/backup"
	"github.com/user/go-dumper/internal/config"
	"github.com/user/go-dumper/internal/http/handlers"
	"github.com/user/go-dumper/internal/http/middleware"
	"github.com/user/go-dumper/internal/store"
)

func New(db *sql.DB) *gin.Engine {
	if gin.Mode() != gin.DebugMode {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(middleware.CORS())

	repo := store.NewRepository(db)
	backupDir := config.GetEnv("BACKUP_DIR", "/data/backups")
	
	dumper := backup.NewDumper(repo, backupDir)
	restorer := backup.NewRestorer(repo)

	targetsHandler := handlers.NewTargetsHandler(repo, dumper)
	backupsHandler := handlers.NewBackupsHandler(repo, restorer)
	healthHandler := handlers.NewHealthHandler(db)

	r.GET("/healthz", healthHandler.Healthz)
	r.GET("/readyz", healthHandler.Readyz)

	api := r.Group("/api")
	api.Use(middleware.BasicAuth())
	{
		targets := api.Group("/targets")
		{
			targets.GET("", targetsHandler.GetTargets)
			targets.POST("", targetsHandler.CreateTarget)
			targets.GET("/:id", targetsHandler.GetTarget)
			targets.PUT("/:id", targetsHandler.UpdateTarget)
			targets.DELETE("/:id", targetsHandler.DeleteTarget)
			targets.POST("/:id/backup", targetsHandler.CreateBackup)
			targets.GET("/:id/backups", targetsHandler.GetTargetBackups)
		}

		backups := api.Group("/backups")
		{
			backups.GET("/:id/download", backupsHandler.DownloadBackup)
			backups.POST("/:id/restore", backupsHandler.RestoreBackup)
			backups.DELETE("/:id", backupsHandler.DeleteBackup)
		}
	}

	r.Static("/assets", "./web/public/assets")
	r.StaticFile("/", "./web/public/index.html")
	r.NoRoute(func(c *gin.Context) {
		c.File("./web/public/index.html")
	})

	return r
}