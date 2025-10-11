package performance

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"go-api-server-sample/cmd/api-server/internal/api/content"
	"go-api-server-sample/cmd/api-server/internal/infrastructure/repositories"
	"go-api-server-sample/cmd/api-server/internal/middleware"

	"github.com/gin-gonic/gin"
)

func setupContentListBenchmark(b *testing.B) (*httptest.Server, *http.Client) {
	// テスト用コンテンツを100件作成
	createTestContents(b, 100)

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
		contents.GET("", contentAPI.List)
	}

	// テストサーバー起動
	server := httptest.NewServer(r)
	httpClient := &http.Client{Timeout: 10 * time.Second}

	return server, httpClient
}

// BenchmarkContentList はコンテンツ一覧取得のベンチマーク
func BenchmarkContentList(b *testing.B) {
	cleanupDB(b)
	server, httpClient := setupContentListBenchmark(b)
	defer server.Close()

	b.Run("シングルリクエスト", func(b *testing.B) {
		url := server.URL + "/api/v1/contents"
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

			var response content.ListContentsResponse
			if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
				resp.Body.Close()
				b.Fatal(err)
			}
			resp.Body.Close()

			if len(response.Contents) == 0 {
				b.Fatal("no contents returned")
			}
		}
	})

	b.Run("ページネーション付きリクエスト", func(b *testing.B) {
		url := server.URL + "/api/v1/contents?limit=10&offset=0"
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

			var response content.ListContentsResponse
			if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
				resp.Body.Close()
				b.Fatal(err)
			}
			resp.Body.Close()

			if len(response.Contents) == 0 {
				b.Fatal("no contents returned")
			}
		}
	})

	b.Run("並行リクエスト", func(b *testing.B) {
		url := server.URL + "/api/v1/contents"
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

				var response content.ListContentsResponse
				if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
					resp.Body.Close()
					b.Error(err)
					continue
				}
				resp.Body.Close()

				if len(response.Contents) == 0 {
					b.Error("no contents returned")
				}
			}
		})
	})
}
