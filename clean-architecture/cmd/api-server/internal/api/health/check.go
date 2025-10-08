package health

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// HealthAPI はヘルスチェック関連のHTTPハンドラーを提供する構造体
type HealthAPI struct {
	db *gorm.DB
}

// NewHealthAPI はHealthAPIの新しいインスタンスを作成する
func NewHealthAPI(db *gorm.DB) *HealthAPI {
	return &HealthAPI{
		db: db,
	}
}

// HealthCheckResponse はヘルスチェックレスポンスの構造体
type HealthCheckResponse struct {
	Status    string    `json:"status"`
	Database  string    `json:"database,omitempty"`
	Timestamp time.Time `json:"timestamp"`
	Message   string    `json:"message,omitempty"`
}

// Check はヘルスチェックを実行するHTTPハンドラー
func (api *HealthAPI) Check(c *gin.Context) {
	response := &HealthCheckResponse{
		Timestamp: time.Now(),
	}

	var dbStatus string
	var healthyStatus = true

	if api.db != nil {
		sqlDB, err := api.db.DB()
		if err != nil {
			dbStatus = "disconnected"
			healthyStatus = false
		} else {
			if err := sqlDB.PingContext(c.Request.Context()); err != nil {
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

	c.JSON(http.StatusOK, response)
}
