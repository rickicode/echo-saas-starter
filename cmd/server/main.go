package main

import (
	"context"
	"echo-saas-starter/internal/config"
	"echo-saas-starter/internal/core"
	"echo-saas-starter/internal/database"
	"echo-saas-starter/internal/plugins"
	"echo-saas-starter/internal/server"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize database
	db, err := database.NewDB(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Create plugin manager and register plugins
	pm := core.NewPluginManager()
	plugins.RegisterAll(pm)

	// Initialize all plugins
	if err := pm.InitAll(cfg, db); err != nil {
		log.Fatalf("Failed to initialize plugins: %v", err)
	}

	// Run migrations
	if err := pm.MigrateAll(db); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Create and start server
	srv := server.NewServer(cfg, db, pm)

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := srv.Start(); err != nil {
			log.Printf("Server stopped: %v", err)
		}
	}()

	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	if err := pm.Shutdown(); err != nil {
		log.Printf("Plugin shutdown error: %v", err)
	}

	log.Println("Server exited")
}
