package container

import (
	"log"

	"gorm.io/gorm"

	"go-api-server-sample/cmd/api-server/internal/application"
	"go-api-server-sample/cmd/api-server/internal/controller"
	"go-api-server-sample/cmd/api-server/internal/middleware"
	"go-api-server-sample/internal/domain/repositories"
	"go-api-server-sample/internal/domain/services"
	infraRepos "go-api-server-sample/internal/infrastructure/repositories"
)

// Container holds all application dependencies
type Container struct {
	// Infrastructure
	DB *gorm.DB

	// Repositories
	UserRepository repositories.UserRepository

	// Domain Services
	UserDomainService *services.UserDomainService

	// Application Services
	UserService application.UserServiceInterface

	// Controllers
	UserController *controller.UserController

	// Middleware
	ValidationMiddleware  *middleware.ValidationMiddleware
	ErrorHandlerMiddleware *middleware.ErrorHandlerMiddleware

	// Logger
	Logger *log.Logger
}

// NewContainer creates and initializes a new dependency injection container
func NewContainer(db *gorm.DB, logger *log.Logger) *Container {
	container := &Container{
		DB:     db,
		Logger: logger,
	}

	container.initializeRepositories()
	container.initializeDomainServices()
	container.initializeApplicationServices()
	container.initializeControllers()
	container.initializeMiddleware()

	return container
}

// initializeRepositories creates repository implementations
func (c *Container) initializeRepositories() {
	c.UserRepository = infraRepos.NewUserRepository(c.DB)
}

// initializeDomainServices creates domain service instances
func (c *Container) initializeDomainServices() {
	c.UserDomainService = services.NewUserDomainService(c.UserRepository)
}

// initializeApplicationServices creates application service instances
func (c *Container) initializeApplicationServices() {
	c.UserService = application.NewUserService(c.UserRepository, c.UserDomainService)
}

// initializeControllers creates controller instances
func (c *Container) initializeControllers() {
	c.UserController = controller.NewUserController(c.UserService)
}

// initializeMiddleware creates middleware instances
func (c *Container) initializeMiddleware() {
	c.ValidationMiddleware = middleware.NewValidationMiddleware()
	c.ErrorHandlerMiddleware = middleware.NewErrorHandlerMiddleware(c.Logger)
}

// GetUserController returns the user controller instance
func (c *Container) GetUserController() *controller.UserController {
	return c.UserController
}

// GetValidationMiddleware returns the validation middleware instance
func (c *Container) GetValidationMiddleware() *middleware.ValidationMiddleware {
	return c.ValidationMiddleware
}

// GetErrorHandlerMiddleware returns the error handler middleware instance
func (c *Container) GetErrorHandlerMiddleware() *middleware.ErrorHandlerMiddleware {
	return c.ErrorHandlerMiddleware
}

// GetLogger returns the logger instance
func (c *Container) GetLogger() *log.Logger {
	return c.Logger
}

// Cleanup performs cleanup operations for the container
func (c *Container) Cleanup() error {
	c.Logger.Println("Performing container cleanup...")

	// Close database connection if needed
	if c.DB != nil {
		sqlDB, err := c.DB.DB()
		if err == nil {
			if err := sqlDB.Close(); err != nil {
				c.Logger.Printf("Error closing database connection: %v", err)
				return err
			}
		}
	}

	c.Logger.Println("Container cleanup completed")
	return nil
}

// HealthCheck performs health checks on all container components
func (c *Container) HealthCheck() error {
	c.Logger.Println("Performing container health check...")

	// Check database connection
	if c.DB != nil {
		sqlDB, err := c.DB.DB()
		if err != nil {
			return err
		}
		if err := sqlDB.Ping(); err != nil {
			return err
		}
	}

	c.Logger.Println("Container health check passed")
	return nil
}