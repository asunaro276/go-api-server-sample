package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go-api-server-sample/cmd/api-server/internal/application"
)

type HealthController struct {
	healthCheckUseCase *application.HealthCheckUseCase
}

func NewHealthController(healthCheckUseCase *application.HealthCheckUseCase) *HealthController {
	return &HealthController{
		healthCheckUseCase: healthCheckUseCase,
	}
}

func (ctrl *HealthController) Check(c *gin.Context) {
	response := ctrl.healthCheckUseCase.Execute(c.Request.Context())

	statusCode := http.StatusOK
	if response.Status == "unhealthy" {
		statusCode = http.StatusServiceUnavailable
	}

	c.JSON(statusCode, response)
}
