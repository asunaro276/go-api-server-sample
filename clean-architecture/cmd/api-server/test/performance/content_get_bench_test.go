package performance

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"go-api-server-sample/cmd/api-server/internal/api/content"
	"go-api-server-sample/cmd/api-server/internal/infrastructure/repositories"
	"go-api-server-sample/cmd/api-server/internal/middleware"
	"go-api-server-sample/internal/domain/entities"

	"github.com/gin-gonic/gin"
)

func setupContentGetBenchmark(b *testing.B) (*httptest.Server, *http.Client, uint) {
	// テスト用コンテンツを事前に作成
	testContent := createTestContent(b, "ベンチマーク用コンテンツ", "ベンチマーク用本文")

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
		contents.GET("/:id", contentAPI.GetByID)
	}

	// テストサーバー起動
	server := httptest.NewServer(r)
	httpClient := &http.Client{Timeout: 10 * time.Second}

	return server, httpClient, testContent.ID
}

// BenchmarkContentGet はコンテンツ取得のベンチマーク
func BenchmarkContentGet(b *testing.B) {
	cleanupDB(b)
	server, httpClient, contentID := setupContentGetBenchmark(b)
	defer server.Close()

	url := fmt.Sprintf("%s/api/v1/contents/%d", server.URL, contentID)

	b.Run("シングルリクエスト", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			resp, err := httpClient.Get(url)
			if err != nil {
				b.Fatal(err)
			}

			if resp.StatusCode != http.StatusOK {
				resp.Body.Close()
				b.Fatalf("unexpected status code: %d", resp.StatusCode)
			}

			var response entities.Content
			if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
				resp.Body.Close()
				b.Fatal(err)
			}
			resp.Body.Close()

			if response.ID != contentID {
				b.Fatalf("unexpected content ID: %d", response.ID)
			}
		}
	})

	b.Run("並行リクエスト", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()

		b.RunParallel(func(pb *testing.PB) {
			client := &http.Client{Timeout: 10 * time.Second}
			for pb.Next() {
				resp, err := client.Get(url)
				if err != nil {
					b.Error(err)
					continue
				}

				if resp.StatusCode != http.StatusOK {
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

				if response.ID != contentID {
					b.Errorf("unexpected content ID: %d", response.ID)
				}
			}
		})
	})
}
