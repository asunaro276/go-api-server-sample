package performance

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	"go-api-server-sample/cmd/api-server/internal/api/content"
	"go-api-server-sample/cmd/api-server/internal/infrastructure/repositories"
	"go-api-server-sample/cmd/api-server/internal/middleware"
	"go-api-server-sample/internal/domain/entities"

	"github.com/gin-gonic/gin"
)

func setupContentCreateBenchmark(b *testing.B) (*httptest.Server, *http.Client) {
	// ルーター設定
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.CORS())

	contentRepo := repositories.NewContentRepository(getDB())
	contentAPI := content.NewContentAPI(contentRepo)

	v1 := r.Group("/api/v1")
	v1.Use(middleware.ErrorHandler())

	contents := v1.Group("/contents")
	{
		contents.POST("", contentAPI.Create)
	}

	// テストサーバー起動
	server := httptest.NewServer(r)
	httpClient := &http.Client{Timeout: 10 * time.Second}

	return server, httpClient
}

// BenchmarkContentCreate はコンテンツ作成のベンチマーク
func BenchmarkContentCreate(b *testing.B) {
	cleanupDB(b)
	server, httpClient := setupContentCreateBenchmark(b)
	defer server.Close()

	b.Run("シングルリクエスト", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			reqBody := map[string]string{
				"title":        fmt.Sprintf("ベンチマークタイトル%d", i),
				"body":         "ベンチマーク本文",
				"content_type": "article",
				"author":       "ベンチマーク作成者",
			}
			jsonBytes, _ := json.Marshal(reqBody)

			resp, err := httpClient.Post(
				server.URL+"/api/v1/contents",
				"application/json",
				bytes.NewBuffer(jsonBytes),
			)
			if err != nil {
				b.Fatal(err)
			}

			if resp.StatusCode != http.StatusCreated {
				resp.Body.Close()
				b.Fatalf("unexpected status code: %d", resp.StatusCode)
			}

			var response entities.Content
			if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
				resp.Body.Close()
				b.Fatal(err)
			}
			resp.Body.Close()

			if response.ID == 0 {
				b.Fatal("content ID is 0")
			}
		}
	})

	var counter atomic.Int64
	b.Run("並行リクエスト", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()

		b.RunParallel(func(pb *testing.PB) {
			client := &http.Client{Timeout: 10 * time.Second}
			for pb.Next() {
				cnt := counter.Add(1)
				reqBody := map[string]string{
					"title":        fmt.Sprintf("ベンチマークタイトル%d", cnt),
					"body":         "ベンチマーク本文",
					"content_type": "article",
					"author":       "ベンチマーク作成者",
				}
				jsonBytes, _ := json.Marshal(reqBody)

				resp, err := client.Post(
					server.URL+"/api/v1/contents",
					"application/json",
					bytes.NewBuffer(jsonBytes),
				)
				if err != nil {
					b.Error(err)
					continue
				}

				if resp.StatusCode != http.StatusCreated {
					resp.Body.Close()
					b.Errorf("unexpected status code: %d", resp.StatusCode)
					continue
				}

				var response entities.Content
				if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
					resp.Body.Close()
					b.Error(err)
					continue
				}
				resp.Body.Close()

				if response.ID == 0 {
					b.Error("content ID is 0")
				}
			}
		})
	})
}
