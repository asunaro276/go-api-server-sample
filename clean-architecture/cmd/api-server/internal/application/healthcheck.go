package application

import (
	"context"
	"time"

	"gorm.io/gorm"
)

type HealthCheckUseCase struct {
	db *gorm.DB
}

type HealthCheckResponse struct {
	Status    string    `json:"status"`
	Database  string    `json:"database,omitempty"`
	Timestamp time.Time `json:"timestamp"`
	Message   string    `json:"message,omitempty"`
}

func NewHealthCheckUseCase(db *gorm.DB) *HealthCheckUseCase {
	return &HealthCheckUseCase{
		db: db,
	}
}

func (uc *HealthCheckUseCase) Execute(ctx context.Context) *HealthCheckResponse {
	response := &HealthCheckResponse{
		Timestamp: time.Now(),
	}

	var dbStatus string
	var healthyStatus = true

	if uc.db != nil {
		sqlDB, err := uc.db.DB()
		if err != nil {
			dbStatus = "disconnected"
			healthyStatus = false
		} else {
			if err := sqlDB.PingContext(ctx); err != nil {
				dbStatus = "disconnected"
				healthyStatus = false
			} else {
				dbStatus = "connected"
			}
		}
		response.Database = dbStatus
	}

	if healthyStatus {
		response.Status = "healthy"
	} else {
		response.Status = "unhealthy"
		response.Message = "データベース接続に問題があります"
	}

	return response
}
