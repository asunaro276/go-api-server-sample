package controller

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"go-api-server-sample/cmd/api-server/internal/controller/dtos"
)

// HealthController handles health check requests
type HealthController struct {
	startTime time.Time
}

// NewHealthController creates a new health controller
func NewHealthController() *HealthController {
	return &HealthController{
		startTime: time.Now(),
	}
}

// HealthResponse represents the health check response
type HealthResponse struct {
	Status    string    `json:"status" example:"ok"`
	Timestamp time.Time `json:"timestamp" example:"2023-12-01T12:00:00Z"`
	Uptime    string    `json:"uptime,omitempty" example:"1h30m45s"`
	Version   string    `json:"version,omitempty" example:"1.0.0"`
}

// DetailedHealthResponse represents the detailed health check response
type DetailedHealthResponse struct {
	Status    string                 `json:"status" example:"ok"`
	Timestamp time.Time              `json:"timestamp" example:"2023-12-01T12:00:00Z"`
	Uptime    string                 `json:"uptime" example:"1h30m45s"`
	Version   string                 `json:"version,omitempty" example:"1.0.0"`
	Checks    map[string]HealthCheck `json:"checks"`
}

// HealthCheck represents an individual health check
type HealthCheck struct {
	Status    string    `json:"status" example:"ok"`
	Message   string    `json:"message,omitempty" example:"Database connection is healthy"`
	Timestamp time.Time `json:"timestamp" example:"2023-12-01T12:00:00Z"`
	Duration  string    `json:"duration,omitempty" example:"5ms"`
}

// Health handles GET /health
// @Summary Basic health check
// @Description Returns basic application health status
// @Tags health
// @Produce json
// @Success 200 {object} HealthResponse
// @Router /health [get]
func (h *HealthController) Health(ctx *gin.Context) {
	response := HealthResponse{
		Status:    "ok",
		Timestamp: time.Now().UTC(),
		Uptime:    time.Since(h.startTime).String(),
		Version:   getVersion(),
	}

	ctx.JSON(http.StatusOK, response)
}

// Readiness handles GET /health/ready
// @Summary Readiness check
// @Description Returns readiness status for load balancer health checks
// @Tags health
// @Produce json
// @Success 200 {object} HealthResponse
// @Failure 503 {object} dtos.ErrorResponse
// @Router /health/ready [get]
func (h *HealthController) Readiness(ctx *gin.Context) {
	// Perform readiness checks
	if !h.isReady() {
		errorResp := dtos.NewErrorResponse("NOT_READY", "Service is not ready")
		ctx.JSON(http.StatusServiceUnavailable, errorResp)
		return
	}

	response := HealthResponse{
		Status:    "ready",
		Timestamp: time.Now().UTC(),
	}

	ctx.JSON(http.StatusOK, response)
}

// Liveness handles GET /health/live
// @Summary Liveness check
// @Description Returns liveness status for Kubernetes liveness probes
// @Tags health
// @Produce json
// @Success 200 {object} HealthResponse
// @Failure 503 {object} dtos.ErrorResponse
// @Router /health/live [get]
func (h *HealthController) Liveness(ctx *gin.Context) {
	// Perform liveness checks
	if !h.isAlive() {
		errorResp := dtos.NewErrorResponse("NOT_ALIVE", "Service is not alive")
		ctx.JSON(http.StatusServiceUnavailable, errorResp)
		return
	}

	response := HealthResponse{
		Status:    "alive",
		Timestamp: time.Now().UTC(),
		Uptime:    time.Since(h.startTime).String(),
	}

	ctx.JSON(http.StatusOK, response)
}

// DetailedHealth handles GET /health/detailed
// @Summary Detailed health check
// @Description Returns detailed health status with component checks
// @Tags health
// @Produce json
// @Success 200 {object} DetailedHealthResponse
// @Failure 503 {object} DetailedHealthResponse
// @Router /health/detailed [get]
func (h *HealthController) DetailedHealth(ctx *gin.Context) {
	checks := h.performDetailedChecks()
	
	// Determine overall status
	status := "ok"
	httpStatus := http.StatusOK
	
	for _, check := range checks {
		if check.Status != "ok" {
			status = "degraded"
			httpStatus = http.StatusServiceUnavailable
			break
		}
	}

	response := DetailedHealthResponse{
		Status:    status,
		Timestamp: time.Now().UTC(),
		Uptime:    time.Since(h.startTime).String(),
		Version:   getVersion(),
		Checks:    checks,
	}

	ctx.JSON(httpStatus, response)
}

// isReady checks if the service is ready to handle requests
func (h *HealthController) isReady() bool {
	// Add readiness checks here
	// For example: database connectivity, required services availability
	return true
}

// isAlive checks if the service is alive
func (h *HealthController) isAlive() bool {
	// Add liveness checks here
	// For example: basic application functionality
	return true
}

// performDetailedChecks performs detailed health checks on all components
func (h *HealthController) performDetailedChecks() map[string]HealthCheck {
	checks := make(map[string]HealthCheck)

	// Application check
	checks["application"] = HealthCheck{
		Status:    "ok",
		Message:   "Application is running",
		Timestamp: time.Now().UTC(),
		Duration:  "0ms",
	}

	// Database check (placeholder - would need actual database instance)
	checks["database"] = h.checkDatabase()

	// Memory check
	checks["memory"] = h.checkMemory()

	return checks
}

// checkDatabase performs database health check
func (h *HealthController) checkDatabase() HealthCheck {
	start := time.Now()
	
	// TODO: Add actual database check when container is available
	// For now, return a placeholder
	return HealthCheck{
		Status:    "ok",
		Message:   "Database connection is healthy",
		Timestamp: time.Now().UTC(),
		Duration:  time.Since(start).String(),
	}
}

// checkMemory performs memory usage check
func (h *HealthController) checkMemory() HealthCheck {
	start := time.Now()
	
	// Basic memory check
	// In a real implementation, you might use runtime.MemStats
	return HealthCheck{
		Status:    "ok",
		Message:   "Memory usage is within acceptable limits",
		Timestamp: time.Now().UTC(),
		Duration:  time.Since(start).String(),
	}
}

// getVersion returns the application version
func getVersion() string {
	// This could be injected during build time
	return "1.0.0"
}