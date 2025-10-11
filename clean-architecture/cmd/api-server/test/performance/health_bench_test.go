package performance

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"go-api-server-sample/cmd/api-server/internal/api/health"
	"go-api-server-sample/cmd/api-server/internal/middleware"

	"github.com/gin-gonic/gin"
)

func setupHealthBenchmark(b *testing.B) (*httptest.Server, *http.Client) {
	// ルーター設定
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.CORS())

	healthAPI := health.NewHealthAPI(getDB())
	r.GET("/health", healthAPI.Check)

	// テストサーバー起動
	server := httptest.NewServer(r)
	httpClient := &http.Client{Timeout: 10 * time.Second}

	return server, httpClient
}

// BenchmarkHealthCheck はヘルスチェックエンドポイントのベンチマーク
func BenchmarkHealthCheck(b *testing.B) {
	server, httpClient := setupHealthBenchmark(b)
	defer server.Close()

	b.Run("シングルリクエスト", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			resp, err := httpClient.Get(server.URL + "/health")
			if err != nil {
				b.Fatal(err)
			}

			var response map[string]interface{}
			if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
				resp.Body.Close()
				b.Fatal(err)
			}
			resp.Body.Close()

			if response["status"] != "healthy" {
				b.Fatalf("unexpected status: %v", response["status"])
			}
		}
	})

	b.Run("並行リクエスト", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()

		b.RunParallel(func(pb *testing.PB) {
			client := &http.Client{Timeout: 10 * time.Second}
			for pb.Next() {
				resp, err := client.Get(server.URL + "/health")
				if err != nil {
					b.Error(err)
					continue
				}

				var response map[string]interface{}
				if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
					resp.Body.Close()
					b.Error(err)
					continue
				}
				resp.Body.Close()

				if response["status"] != "healthy" {
					b.Errorf("unexpected status: %v", response["status"])
				}
			}
		})
	})
}
