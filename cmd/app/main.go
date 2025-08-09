package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/user/go-dumper/internal/config"
	"github.com/user/go-dumper/internal/http/router"
	"github.com/user/go-dumper/internal/scheduler"
	"github.com/user/go-dumper/internal/store"
)

func main() {
	// Load environment variables from .env files
	if err := config.LoadEnvFiles(); err != nil {
		log.Printf("Warning: Failed to load .env files: %v", err)
	}

	// Load configuration
	cfg := config.Load()

	if err := os.MkdirAll(cfg.BackupDir, 0755); err != nil {
		log.Fatal("Failed to create backup directory:", err)
	}

	db, err := store.InitDB(cfg.SQLitePath)
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer db.Close()

	sched := scheduler.New(db)
	go sched.Start()
	defer sched.Stop()

	r := router.New(db)
	
	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: r,
	}

	go func() {
		log.Printf("Server starting on port %s", cfg.Port)
		log.Printf("Backup directory: %s", cfg.BackupDir)
		log.Printf("Database path: %s", cfg.SQLitePath)
		if cfg.AdminUser != "" {
			log.Printf("Basic authentication enabled for user: %s", cfg.AdminUser)
		}
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Server failed to start:", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exited")
}