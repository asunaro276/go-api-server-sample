//go:build integration

package integration

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"gorm.io/gorm"
)

type HealthCheckIntegrationTestSuite struct {
	suite.Suite
	container *postgres.PostgresContainer
	db        *gorm.DB
	router    *gin.Engine
}

func (suite *HealthCheckIntegrationTestSuite) SetupSuite() {
	ctx := context.Background()

	// PostgreSQLコンテナ起動
	container, err := postgres.RunContainer(ctx,
		postgres.WithImage("postgres:15"),
		postgres.WithDatabase("testdb"),
		postgres.WithUsername("testuser"),
		postgres.WithPassword("testpass"),
	)
	suite.Require().NoError(err)
	suite.container = container

	// TODO: データベース接続（実装後に有効化）
	// connStr, err := container.ConnectionString(ctx, "sslmode=disable")
	// suite.Require().NoError(err)
	//
	// suite.db, err = gorm.Open(postgres.Open(connStr), &gorm.Config{})
	// suite.Require().NoError(err)

	gin.SetMode(gin.TestMode)
	// TODO: ルーター設定（実装後に有効化）
	suite.router = setupHealthRouter()
}

func (suite *HealthCheckIntegrationTestSuite) TearDownSuite() {
	ctx := context.Background()
	if suite.container != nil {
		suite.container.Terminate(ctx)
	}
}

func (suite *HealthCheckIntegrationTestSuite) TestHealthCheckWithDatabase() {
	suite.Run("データベース接続ありのヘルスチェックが正常に動作する", func() {
		req, _ := http.NewRequest("GET", "/health", nil)
		w := httptest.NewRecorder()
		suite.router.ServeHTTP(w, req)

		assert.Equal(suite.T(), http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(suite.T(), err)

		// 基本的なフィールドの存在確認
		assert.Contains(suite.T(), response, "status")
		assert.Contains(suite.T(), response, "timestamp")

		// ステータスが healthy または unhealthy のいずれか
		status, ok := response["status"].(string)
		assert.True(suite.T(), ok)
		assert.Contains(suite.T(), []string{"healthy", "unhealthy"}, status)

		// タイムスタンプが適切な形式
		timestamp, ok := response["timestamp"].(string)
		assert.True(suite.T(), ok)
		assert.NotEmpty(suite.T(), timestamp)

		// タイムスタンプがパース可能
		_, err = time.Parse(time.RFC3339, timestamp)
		assert.NoError(suite.T(), err)

		// データベース接続状況の確認
		if database, exists := response["database"]; exists {
			dbStatus, ok := database.(string)
			assert.True(suite.T(), ok)
			assert.Contains(suite.T(), []string{"connected", "disconnected"}, dbStatus)
		}
	})
}

func (suite *HealthCheckIntegrationTestSuite) TestHealthCheckResponseTime() {
	suite.Run("ヘルスチェックのレスポンス時間が200ms未満", func() {
		start := time.Now()

		req, _ := http.NewRequest("GET", "/health", nil)
		w := httptest.NewRecorder()
		suite.router.ServeHTTP(w, req)

		duration := time.Since(start)

		assert.Equal(suite.T(), http.StatusOK, w.Code)
		assert.Less(suite.T(), duration, 200*time.Millisecond)
	})
}

func (suite *HealthCheckIntegrationTestSuite) TestHealthCheckConcurrency() {
	suite.Run("ヘルスチェックが並行リクエストに対応する", func() {
		const numRequests = 10
		responses := make(chan int, numRequests)

		// 並行してリクエストを送信
		for i := 0; i < numRequests; i++ {
			go func() {
				req, _ := http.NewRequest("GET", "/health", nil)
				w := httptest.NewRecorder()
				suite.router.ServeHTTP(w, req)
				responses <- w.Code
			}()
		}

		// 全てのレスポンスを収集
		successCount := 0
		for i := 0; i < numRequests; i++ {
			code := <-responses
			if code == http.StatusOK {
				successCount++
			}
		}

		// 全てのリクエストが成功することを確認
		assert.Equal(suite.T(), numRequests, successCount)
	})
}

func (suite *HealthCheckIntegrationTestSuite) TestHealthCheckHTTPMethods() {
	suite.Run("ヘルスチェックはGETメソッドのみ対応", func() {
		// GET以外のメソッドでは405エラーまたは404エラーが返される
		methods := []string{"POST", "PUT", "DELETE", "PATCH"}

		for _, method := range methods {
			suite.Run(fmt.Sprintf("%sメソッドは許可されない", method), func() {
				req, _ := http.NewRequest(method, "/health", nil)
				w := httptest.NewRecorder()
				suite.router.ServeHTTP(w, req)

				// 405 Method Not Allowed または 404 Not Found が期待される
				assert.Contains(suite.T(), []int{http.StatusMethodNotAllowed, http.StatusNotFound}, w.Code)
			})
		}
	})
}

func (suite *HealthCheckIntegrationTestSuite) TestHealthCheckContentType() {
	suite.Run("ヘルスチェックのレスポンスがJSON形式", func() {
		req, _ := http.NewRequest("GET", "/health", nil)
		w := httptest.NewRecorder()
		suite.router.ServeHTTP(w, req)

		assert.Equal(suite.T(), http.StatusOK, w.Code)

		// Content-Typeがapplication/jsonであることを確認
		contentType := w.Header().Get("Content-Type")
		assert.Contains(suite.T(), contentType, "application/json")

		// JSONとしてパース可能であることを確認
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(suite.T(), err)
	})
}

func (suite *HealthCheckIntegrationTestSuite) TestHealthCheckIdempotency() {
	suite.Run("ヘルスチェックが冪等性を持つ", func() {
		// 複数回同じリクエストを送信
		var firstResponse, secondResponse map[string]interface{}

		// 1回目のリクエスト
		req1, _ := http.NewRequest("GET", "/health", nil)
		w1 := httptest.NewRecorder()
		suite.router.ServeHTTP(w1, req1)

		assert.Equal(suite.T(), http.StatusOK, w1.Code)
		err := json.Unmarshal(w1.Body.Bytes(), &firstResponse)
		assert.NoError(suite.T(), err)

		// わずかな時間を空ける
		time.Sleep(10 * time.Millisecond)

		// 2回目のリクエスト
		req2, _ := http.NewRequest("GET", "/health", nil)
		w2 := httptest.NewRecorder()
		suite.router.ServeHTTP(w2, req2)

		assert.Equal(suite.T(), http.StatusOK, w2.Code)
		err = json.Unmarshal(w2.Body.Bytes(), &secondResponse)
		assert.NoError(suite.T(), err)

		// ステータスは同じであることを確認（タイムスタンプは異なって良い）
		assert.Equal(suite.T(), firstResponse["status"], secondResponse["status"])
		if firstResponse["database"] != nil && secondResponse["database"] != nil {
			assert.Equal(suite.T(), firstResponse["database"], secondResponse["database"])
		}
	})
}

func TestHealthCheckIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(HealthCheckIntegrationTestSuite))
}

// setupHealthRouter はヘルスチェック統合テスト用のルーター設定
func setupHealthRouter() *gin.Engine {
	// TODO: 実装が完了するまでは空のルーターを返す
	r := gin.New()
	return r
}