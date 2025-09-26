package router

import (
	"github.com/gin-gonic/gin"

	"go-api-server-sample/cmd/api-server/internal/container"
	"go-api-server-sample/cmd/api-server/internal/middleware"
)

// SetupRouter configures and returns the HTTP router
func SetupRouter(container *container.Container) *gin.Engine {
	// Create Gin router
	router := gin.New()

	// Add global middleware
	router.Use(gin.Logger())
	router.Use(container.GetErrorHandlerMiddleware().HandleErrors())
	router.Use(middleware.RecoveryWithLogger(container.GetLogger()))
	router.Use(middleware.CORS())
	router.Use(middleware.SecurityHeaders())
	router.Use(middleware.ContentTypeValidation())
	router.Use(middleware.RequestSizeLimit(1024 * 1024)) // 1MB limit

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		if err := container.HealthCheck(); err != nil {
			c.JSON(500, gin.H{
				"status": "error",
				"error":  err.Error(),
			})
			return
		}

		c.JSON(200, gin.H{
			"status":    "ok",
			"timestamp": gin.H{},
		})
	})

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		setupUserRoutes(v1, container)
	}

	return router
}

// setupUserRoutes configures user-related API routes
func setupUserRoutes(rg *gin.RouterGroup, container *container.Container) {
	userController := container.GetUserController()

	users := rg.Group("/users")
	{
		// POST /api/v1/users - Create user
		users.POST("", userController.CreateUser)

		// GET /api/v1/users - List users
		users.GET("", userController.GetUsers)

		// GET /api/v1/users/:id - Get user by ID
		users.GET("/:id", userController.GetUserByID)

		// PUT /api/v1/users/:id - Update user
		users.PUT("/:id", userController.UpdateUser)

		// DELETE /api/v1/users/:id - Delete user
		users.DELETE("/:id", userController.DeleteUser)
	}
}

// SetupTestRouter creates a router for testing purposes
func SetupTestRouter(container *container.Container) *gin.Engine {
	gin.SetMode(gin.TestMode)
	return SetupRouter(container)
}

// SetupProductionRouter creates a router optimized for production
func SetupProductionRouter(container *container.Container) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := SetupRouter(container)

	// Add production-specific middleware
	// e.g., rate limiting, authentication, etc.

	return router
}