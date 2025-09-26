package migrations

import (
	"fmt"
	"log"

	"gorm.io/gorm"

	"go-api-server-sample/internal/domain/entities"
)

// MigrationManager handles database migrations
type MigrationManager struct {
	db *gorm.DB
}

// NewMigrationManager creates a new migration manager
func NewMigrationManager(db *gorm.DB) *MigrationManager {
	return &MigrationManager{
		db: db,
	}
}

// AutoMigrate runs GORM auto-migration for all entities
func (m *MigrationManager) AutoMigrate() error {
	log.Println("Starting database auto-migration...")

	// List of all entities to migrate
	entities := []interface{}{
		&entities.User{},
	}

	// Run auto-migration for each entity
	for _, entity := range entities {
		if err := m.db.AutoMigrate(entity); err != nil {
			return fmt.Errorf("failed to auto-migrate entity %T: %w", entity, err)
		}
		log.Printf("Auto-migration completed for entity: %T", entity)
	}

	log.Println("Database auto-migration completed successfully")
	return nil
}

// CreateIndexes creates additional indexes that auto-migration might miss
func (m *MigrationManager) CreateIndexes() error {
	log.Println("Creating additional database indexes...")

	// Create additional indexes if needed
	indexes := []string{
		// Users table indexes
		"CREATE INDEX IF NOT EXISTS idx_users_email_active ON users(email) WHERE deleted_at IS NULL;",
		"CREATE INDEX IF NOT EXISTS idx_users_created_at ON users(created_at);",
	}

	for _, indexSQL := range indexes {
		if err := m.db.Exec(indexSQL).Error; err != nil {
			log.Printf("Warning: Failed to create index with SQL: %s, Error: %v", indexSQL, err)
			// Continue with other indexes even if one fails
		} else {
			log.Printf("Index created successfully: %s", indexSQL)
		}
	}

	log.Println("Additional indexes creation completed")
	return nil
}

// DropAllTables drops all tables (for development/testing purposes)
func (m *MigrationManager) DropAllTables() error {
	log.Println("WARNING: Dropping all database tables...")

	// List of tables to drop (in reverse dependency order)
	tables := []interface{}{
		&entities.User{},
	}

	for _, table := range tables {
		if err := m.db.Migrator().DropTable(table); err != nil {
			return fmt.Errorf("failed to drop table for entity %T: %w", table, err)
		}
		log.Printf("Table dropped for entity: %T", table)
	}

	log.Println("All tables dropped successfully")
	return nil
}

// Migrate runs a full migration (drop, create, and index)
func (m *MigrationManager) Migrate() error {
	if err := m.AutoMigrate(); err != nil {
		return err
	}

	if err := m.CreateIndexes(); err != nil {
		return err
	}

	return nil
}

// Reset drops all tables and recreates them (for development/testing)
func (m *MigrationManager) Reset() error {
	log.Println("WARNING: Resetting database (dropping and recreating all tables)...")

	if err := m.DropAllTables(); err != nil {
		return err
	}

	if err := m.Migrate(); err != nil {
		return err
	}

	log.Println("Database reset completed successfully")
	return nil
}

// CheckMigrationStatus checks if migrations are up to date
func (m *MigrationManager) CheckMigrationStatus() error {
	log.Println("Checking migration status...")

	// Check if all required tables exist
	entities := []interface{}{
		&entities.User{},
	}

	for _, entity := range entities {
		if !m.db.Migrator().HasTable(entity) {
			return fmt.Errorf("table for entity %T does not exist", entity)
		}
	}

	log.Println("Migration status check completed - all tables exist")
	return nil
}