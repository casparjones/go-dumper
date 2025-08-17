package router

import (
	"database/sql"

	"github.com/casparjones/go-dumper/internal/backup"
	"github.com/casparjones/go-dumper/internal/config"
	"github.com/casparjones/go-dumper/internal/http/handlers"
	"github.com/casparjones/go-dumper/internal/http/middleware"
	"github.com/casparjones/go-dumper/internal/store"
	"github.com/gin-gonic/gin"
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
	jobsHandler := handlers.NewJobsHandler(repo, dumper)
	configHandler := handlers.NewConfigHandler(repo)
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
			targets.POST("/discover", targetsHandler.DiscoverDatabases)
		}

		backups := api.Group("/backups")
		{
			backups.GET("", backupsHandler.GetAllBackups)
			backups.GET("/:id/download", backupsHandler.DownloadBackup)
			backups.POST("/:id/restore", backupsHandler.RestoreBackup)
			backups.DELETE("/:id", backupsHandler.DeleteBackup)
		}

		jobs := api.Group("/jobs")
		{
			jobs.GET("", jobsHandler.GetJobs)
			jobs.POST("", jobsHandler.CreateJob)
			jobs.GET("/:id", jobsHandler.GetJob)
			jobs.PUT("/:id", jobsHandler.UpdateJob)
			jobs.DELETE("/:id", jobsHandler.DeleteJob)
			jobs.POST("/:id/run", jobsHandler.RunJobNow)
		}

		config := api.Group("/config")
		{
			config.GET("", configHandler.GetAllConfigs)
			config.POST("", configHandler.SetConfig)
			config.GET("/:key", configHandler.GetConfig)
			config.GET("/theme", configHandler.GetTheme)
			config.POST("/theme", configHandler.SetTheme)
		}
	}

	r.Static("/assets", "./web/public/assets")
	r.StaticFile("/", "./web/public/index.html")
	r.NoRoute(func(c *gin.Context) {
		c.File("./web/public/index.html")
	})

	return r
}
