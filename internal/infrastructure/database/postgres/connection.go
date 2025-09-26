package postgres

import (
	"fmt"
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Config holds PostgreSQL connection configuration
type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
	TimeZone string
}

// ConnectionManager manages PostgreSQL database connections
type ConnectionManager struct {
	config *Config
	db     *gorm.DB
}

// NewConnectionManager creates a new connection manager
func NewConnectionManager(config *Config) *ConnectionManager {
	return &ConnectionManager{
		config: config,
	}
}

// Connect establishes a connection to PostgreSQL database
func (cm *ConnectionManager) Connect() (*gorm.DB, error) {
	if cm.db != nil {
		return cm.db, nil
	}

	dsn := cm.buildDSN()
	log.Printf("Connecting to PostgreSQL database: %s@%s:%s/%s", 
		cm.config.User, cm.config.Host, cm.config.Port, cm.config.DBName)

	// Configure GORM logger
	gormLogger := logger.New(
		log.New(log.Writer(), "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second,   // Slow SQL threshold
			LogLevel:                  logger.Silent, // Log level: Silent, Error, Warn, Info
			IgnoreRecordNotFoundError: true,          // Ignore ErrRecordNotFound error for logger
			Colorful:                  false,         // Disable color
		},
	)

	// Open database connection
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormLogger,
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Configure connection pool
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	// Configure connection pool settings
	sqlDB.SetMaxIdleConns(10)                // Maximum number of idle connections
	sqlDB.SetMaxOpenConns(100)               // Maximum number of open connections
	sqlDB.SetConnMaxLifetime(time.Hour)      // Maximum connection lifetime

	// Test the connection
	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	cm.db = db
	log.Println("Successfully connected to PostgreSQL database")
	return db, nil
}

// Close closes the database connection
func (cm *ConnectionManager) Close() error {
	if cm.db == nil {
		return nil
	}

	sqlDB, err := cm.db.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	if err := sqlDB.Close(); err != nil {
		return fmt.Errorf("failed to close database connection: %w", err)
	}

	cm.db = nil
	log.Println("Database connection closed")
	return nil
}

// GetDB returns the GORM database instance
func (cm *ConnectionManager) GetDB() *gorm.DB {
	return cm.db
}

// IsConnected checks if the database connection is active
func (cm *ConnectionManager) IsConnected() bool {
	if cm.db == nil {
		return false
	}

	sqlDB, err := cm.db.DB()
	if err != nil {
		return false
	}

	return sqlDB.Ping() == nil
}

// buildDSN builds the PostgreSQL Data Source Name
func (cm *ConnectionManager) buildDSN() string {
	timeZone := cm.config.TimeZone
	if timeZone == "" {
		timeZone = "UTC"
	}

	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
		cm.config.Host,
		cm.config.User,
		cm.config.Password,
		cm.config.DBName,
		cm.config.Port,
		cm.config.SSLMode,
		timeZone,
	)
}

// HealthCheck performs a database health check
func (cm *ConnectionManager) HealthCheck() error {
	if !cm.IsConnected() {
		return fmt.Errorf("database connection is not active")
	}

	sqlDB, err := cm.db.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	// Check connection pool stats
	stats := sqlDB.Stats()
	log.Printf("Database connection pool stats - Open: %d, InUse: %d, Idle: %d",
		stats.OpenConnections, stats.InUse, stats.Idle)

	return nil
}

// OptimizeConnectionPool optimizes the connection pool settings based on environment
func (cm *ConnectionManager) OptimizeConnectionPool(maxIdle, maxOpen int, maxLifetime time.Duration) error {
	if cm.db == nil {
		return fmt.Errorf("database connection not established")
	}

	sqlDB, err := cm.db.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	// Apply optimized settings
	sqlDB.SetMaxIdleConns(maxIdle)
	sqlDB.SetMaxOpenConns(maxOpen)
	sqlDB.SetConnMaxLifetime(maxLifetime)

	log.Printf("Connection pool optimized - MaxIdle: %d, MaxOpen: %d, MaxLifetime: %v",
		maxIdle, maxOpen, maxLifetime)

	return nil
}