package performance

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type APIPerformanceTestSuite struct {
	suite.Suite
	router *gin.Engine
}

func (suite *APIPerformanceTestSuite) SetupSuite() {
	gin.SetMode(gin.TestMode)
	// TODO: 実装完了後に実際のルーターを使用
	suite.router = setupPerformanceRouter()
}

func (suite *APIPerformanceTestSuite) TestHealthCheckPerformance() {
	suite.Run("ヘルスチェックAPIのレスポンス時間が200ms未満", func() {
		const targetTime = 200 * time.Millisecond
		const iterations = 100

		var totalDuration time.Duration
		successCount := 0

		for i := 0; i < iterations; i++ {
			start := time.Now()

			req, _ := http.NewRequest("GET", "/health", nil)
			w := httptest.NewRecorder()
			suite.router.ServeHTTP(w, req)

			duration := time.Since(start)
			totalDuration += duration

			if w.Code == http.StatusOK {
				successCount++
			}

			// 各リクエストが目標時間以内
			assert.Less(suite.T(), duration, targetTime,
				fmt.Sprintf("リクエスト %d が目標時間 %v を超過: %v", i+1, targetTime, duration))
		}

		// 平均レスポンス時間
		avgDuration := totalDuration / time.Duration(iterations)
		suite.T().Logf("平均レスポンス時間: %v", avgDuration)
		assert.Less(suite.T(), avgDuration, targetTime)

		// 成功率
		successRate := float64(successCount) / float64(iterations) * 100
		suite.T().Logf("成功率: %.2f%%", successRate)
		assert.GreaterOrEqual(suite.T(), successRate, 95.0)
	})
}

func (suite *APIPerformanceTestSuite) TestContentListPerformance() {
	suite.Run("コンテンツ一覧取得APIのレスポンス時間が200ms未満", func() {
		const targetTime = 200 * time.Millisecond
		const iterations = 50

		var totalDuration time.Duration
		successCount := 0

		for i := 0; i < iterations; i++ {
			start := time.Now()

			req, _ := http.NewRequest("GET", "/api/v1/contents", nil)
			w := httptest.NewRecorder()
			suite.router.ServeHTTP(w, req)

			duration := time.Since(start)
			totalDuration += duration

			if w.Code == http.StatusOK {
				successCount++
			}

			assert.Less(suite.T(), duration, targetTime,
				fmt.Sprintf("リクエスト %d が目標時間 %v を超過: %v", i+1, targetTime, duration))
		}

		avgDuration := totalDuration / time.Duration(iterations)
		suite.T().Logf("平均レスポンス時間: %v", avgDuration)
		assert.Less(suite.T(), avgDuration, targetTime)

		successRate := float64(successCount) / float64(iterations) * 100
		suite.T().Logf("成功率: %.2f%%", successRate)
		assert.GreaterOrEqual(suite.T(), successRate, 95.0)
	})
}

func (suite *APIPerformanceTestSuite) TestContentCreatePerformance() {
	suite.Run("コンテンツ作成APIのレスポンス時間が200ms未満", func() {
		const targetTime = 200 * time.Millisecond
		const iterations = 30

		var totalDuration time.Duration
		successCount := 0

		for i := 0; i < iterations; i++ {
			requestBody := map[string]interface{}{
				"title":        fmt.Sprintf("パフォーマンステスト記事 %d", i+1),
				"body":         "パフォーマンステスト用の記事本文です。",
				"content_type": "article",
				"author":       "パフォーマンステスター",
			}

			jsonBody, _ := json.Marshal(requestBody)
			start := time.Now()

			req, _ := http.NewRequest("POST", "/api/v1/contents", bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			suite.router.ServeHTTP(w, req)

			duration := time.Since(start)
			totalDuration += duration

			if w.Code == http.StatusCreated {
				successCount++
			}

			assert.Less(suite.T(), duration, targetTime,
				fmt.Sprintf("リクエスト %d が目標時間 %v を超過: %v", i+1, targetTime, duration))
		}

		avgDuration := totalDuration / time.Duration(iterations)
		suite.T().Logf("平均レスポンス時間: %v", avgDuration)
		assert.Less(suite.T(), avgDuration, targetTime)

		successRate := float64(successCount) / float64(iterations) * 100
		suite.T().Logf("成功率: %.2f%%", successRate)
		assert.GreaterOrEqual(suite.T(), successRate, 95.0)
	})
}

func (suite *APIPerformanceTestSuite) TestConcurrentRequests() {
	suite.Run("並行リクエストに対する性能", func() {
		const concurrency = 10
		const requestsPerWorker = 20
		const targetTime = 200 * time.Millisecond

		results := make(chan time.Duration, concurrency*requestsPerWorker)
		start := time.Now()

		// 並行ワーカー起動
		for i := 0; i < concurrency; i++ {
			go func(workerID int) {
				for j := 0; j < requestsPerWorker; j++ {
					reqStart := time.Now()

					req, _ := http.NewRequest("GET", "/health", nil)
					w := httptest.NewRecorder()
					suite.router.ServeHTTP(w, req)

					results <- time.Since(reqStart)
				}
			}(i)
		}

		// 結果収集
		var totalDuration time.Duration
		successCount := 0
		totalRequests := concurrency * requestsPerWorker

		for i := 0; i < totalRequests; i++ {
			duration := <-results
			totalDuration += duration
			if duration < targetTime {
				successCount++
			}
		}

		overallDuration := time.Since(start)
		avgDuration := totalDuration / time.Duration(totalRequests)
		successRate := float64(successCount) / float64(totalRequests) * 100

		suite.T().Logf("全体実行時間: %v", overallDuration)
		suite.T().Logf("平均レスポンス時間: %v", avgDuration)
		suite.T().Logf("目標時間内成功率: %.2f%%", successRate)
		suite.T().Logf("スループット: %.2f req/sec", float64(totalRequests)/overallDuration.Seconds())

		assert.GreaterOrEqual(suite.T(), successRate, 90.0)
		assert.Less(suite.T(), avgDuration, targetTime)
	})
}

func (suite *APIPerformanceTestSuite) BenchmarkHealthCheck() {
	suite.T().Helper()
	b := testing.B{}
	b.Run("ヘルスチェックAPIベンチマーク", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			req, _ := http.NewRequest("GET", "/health", nil)
			w := httptest.NewRecorder()
			suite.router.ServeHTTP(w, req)

			if w.Code != http.StatusOK {
				b.Errorf("Expected status 200, got %d", w.Code)
			}
		}
	})
}

func (suite *APIPerformanceTestSuite) BenchmarkContentList() {
	suite.T().Helper()
	b := testing.B{}
	b.Run("コンテンツ一覧取得APIベンチマーク", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			req, _ := http.NewRequest("GET", "/api/v1/contents", nil)
			w := httptest.NewRecorder()
			suite.router.ServeHTTP(w, req)

			if w.Code != http.StatusOK {
				b.Errorf("Expected status 200, got %d", w.Code)
			}
		}
	})
}

func TestAPIPerformanceTestSuite(t *testing.T) {
	suite.Run(t, new(APIPerformanceTestSuite))
}

// setupPerformanceRouter はパフォーマンステスト用のルーター設定
func setupPerformanceRouter() *gin.Engine {
	// TODO: 実装が完了するまでは空のルーターを返す
	r := gin.New()
	return r
}
