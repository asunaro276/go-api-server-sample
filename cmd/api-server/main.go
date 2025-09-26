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

	"go-api-server-sample/cmd/api-server/internal/config"
	"go-api-server-sample/cmd/api-server/internal/container"
	"go-api-server-sample/cmd/api-server/internal/router"
	"go-api-server-sample/internal/infrastructure/database/migrations"
	"go-api-server-sample/internal/infrastructure/database/postgres"
)

func main() {
	// Initialize logger
	logger := log.New(os.Stdout, "[API-SERVER] ", log.LstdFlags|log.Lshortfile)
	logger.Println("Starting Go API Server...")

	// Parse command line arguments
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "migrate":
			runMigration(logger)
			return
		case "migrate-reset":
			runMigrationReset(logger)
			return
		case "version":
			fmt.Println("Go API Server v1.0.0")
			return
		case "help":
			printHelp()
			return
		}
	}

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		logger.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize database connection
	dbConfig := &postgres.Config{
		Host:     cfg.Database.Host,
		Port:     cfg.Database.Port,
		User:     cfg.Database.User,
		Password: cfg.Database.Password,
		DBName:   cfg.Database.DBName,
		SSLMode:  cfg.Database.SSLMode,
		TimeZone: "UTC",
	}

	connectionManager := postgres.NewConnectionManager(dbConfig)
	db, err := connectionManager.Connect()
	if err != nil {
		logger.Fatalf("Failed to connect to database: %v", err)
	}
	defer func() {
		if err := connectionManager.Close(); err != nil {
			logger.Printf("Error closing database connection: %v", err)
		}
	}()

	// Optimize connection pool for production
	if cfg.App.Environment == "production" {
		if err := connectionManager.OptimizeConnectionPool(20, 200, time.Hour*2); err != nil {
			logger.Printf("Warning: Failed to optimize connection pool: %v", err)
		}
	}

	// Run migrations if enabled
	if cfg.Database.AutoMigrate {
		migrationManager := migrations.NewMigrationManager(db)
		if err := migrationManager.Migrate(); err != nil {
			logger.Fatalf("Failed to run migrations: %v", err)
		}
	}

	// Initialize dependency injection container
	appContainer := container.NewContainer(db, logger)
	defer func() {
		if err := appContainer.Cleanup(); err != nil {
			logger.Printf("Error during container cleanup: %v", err)
		}
	}()

	// Setup HTTP router
	httpRouter := router.SetupRouter(appContainer)

	// Create HTTP server
	server := &http.Server{
		Addr:         cfg.ServerAddress(),
		Handler:      httpRouter,
		ReadTimeout:  time.Duration(cfg.Server.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(cfg.Server.WriteTimeout) * time.Second,
		IdleTimeout:  time.Duration(cfg.Server.IdleTimeout) * time.Second,
	}

	// Start server in a goroutine
	go func() {
		logger.Printf("Starting HTTP server on %s", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Println("Shutting down server...")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Fatalf("Server forced to shutdown: %v", err)
	}

	logger.Println("Server exited gracefully")
}

// runMigration runs database migrations
func runMigration(logger *log.Logger) {
	logger.Println("Running database migrations...")

	cfg, err := config.Load()
	if err != nil {
		logger.Fatalf("Failed to load configuration: %v", err)
	}

	dbConfig := &postgres.Config{
		Host:     cfg.Database.Host,
		Port:     cfg.Database.Port,
		User:     cfg.Database.User,
		Password: cfg.Database.Password,
		DBName:   cfg.Database.DBName,
		SSLMode:  cfg.Database.SSLMode,
		TimeZone: "UTC",
	}

	connectionManager := postgres.NewConnectionManager(dbConfig)
	db, err := connectionManager.Connect()
	if err != nil {
		logger.Fatalf("Failed to connect to database: %v", err)
	}
	defer connectionManager.Close()

	migrationManager := migrations.NewMigrationManager(db)
	if err := migrationManager.Migrate(); err != nil {
		logger.Fatalf("Failed to run migrations: %v", err)
	}

	logger.Println("Migrations completed successfully")
}

// runMigrationReset resets the database (drops and recreates all tables)
func runMigrationReset(logger *log.Logger) {
	logger.Println("WARNING: Resetting database (this will delete all data)...")

	cfg, err := config.Load()
	if err != nil {
		logger.Fatalf("Failed to load configuration: %v", err)
	}

	// Safety check - only allow in development
	if cfg.App.Environment == "production" {
		logger.Fatal("Migration reset is not allowed in production environment")
	}

	dbConfig := &postgres.Config{
		Host:     cfg.Database.Host,
		Port:     cfg.Database.Port,
		User:     cfg.Database.User,
		Password: cfg.Database.Password,
		DBName:   cfg.Database.DBName,
		SSLMode:  cfg.Database.SSLMode,
		TimeZone: "UTC",
	}

	connectionManager := postgres.NewConnectionManager(dbConfig)
	db, err := connectionManager.Connect()
	if err != nil {
		logger.Fatalf("Failed to connect to database: %v", err)
	}
	defer connectionManager.Close()

	migrationManager := migrations.NewMigrationManager(db)
	if err := migrationManager.Reset(); err != nil {
		logger.Fatalf("Failed to reset database: %v", err)
	}

	logger.Println("Database reset completed successfully")
}

// printHelp prints help information
func printHelp() {
	fmt.Println("Go API Server")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  go run cmd/api-server/main.go [command]")
	fmt.Println()
	fmt.Println("Available Commands:")
	fmt.Println("  migrate        Run database migrations")
	fmt.Println("  migrate-reset  Reset database (development only)")
	fmt.Println("  version        Show version information")
	fmt.Println("  help           Show this help message")
	fmt.Println()
	fmt.Println("Environment Variables:")
	fmt.Println("  DB_HOST        Database host (default: localhost)")
	fmt.Println("  DB_PORT        Database port (default: 5432)")
	fmt.Println("  DB_USER        Database user")
	fmt.Println("  DB_PASSWORD    Database password")
	fmt.Println("  DB_NAME        Database name")
	fmt.Println("  SERVER_HOST    Server host (default: localhost)")
	fmt.Println("  SERVER_PORT    Server port (default: 8080)")
	fmt.Println("  APP_ENV        Application environment (default: development)")
}