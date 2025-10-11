package performance

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"go-api-server-sample/cmd/api-server/internal/api/content"
	"go-api-server-sample/cmd/api-server/internal/infrastructure/repositories"
	"go-api-server-sample/cmd/api-server/internal/middleware"

	"github.com/gin-gonic/gin"
)

func setupContentDeleteBenchmark(b *testing.B, count int) (*httptest.Server, *http.Client, []uint) {
	// テスト用コンテンツを事前に作成
	contents := createTestContents(b, count)
	contentIDs := make([]uint, len(contents))
	for i, c := range contents {
		contentIDs[i] = c.ID
	}

	// ルーター設定
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.CORS())

	contentRepo := repositories.NewContentRepository(getDB())
	contentAPI := content.NewContentAPI(contentRepo)

	v1 := r.Group("/api/v1")
	v1.Use(middleware.ErrorHandler())

	contents_group := v1.Group("/contents")
	{
		contents_group.DELETE("/:id", contentAPI.Delete)
	}

	// テストサーバー起動
	server := httptest.NewServer(r)
	httpClient := &http.Client{Timeout: 10 * time.Second}

	return server, httpClient, contentIDs
}

// BenchmarkContentDelete はコンテンツ削除のベンチマーク
func BenchmarkContentDelete(b *testing.B) {
	b.Run("シングルリクエスト", func(b *testing.B) {
		cleanupDB(b)
		server, httpClient, contentIDs := setupContentDeleteBenchmark(b, b.N)
		defer server.Close()

		b.ReportAllocs()
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			url := fmt.Sprintf("%s/api/v1/contents/%d", server.URL, contentIDs[i])
			req, err := http.NewRequest(http.MethodDelete, url, nil)
			if err != nil {
				b.Fatal(err)
			}

			resp, err := httpClient.Do(req)
			if err != nil {
				b.Fatal(err)
			}

			if resp.StatusCode != http.StatusNoContent {
				resp.Body.Close()
				b.Fatalf("unexpected status code: %d", resp.StatusCode)
			}
			resp.Body.Close()
		}
	})

	b.Run("並行リクエスト", func(b *testing.B) {
		cleanupDB(b)
		server, _, contentIDs := setupContentDeleteBenchmark(b, b.N*10)
		defer server.Close()

		var idIndex int
		var idIndexLock sync.Mutex

		b.ReportAllocs()
		b.ResetTimer()

		b.RunParallel(func(pb *testing.PB) {
			client := &http.Client{Timeout: 10 * time.Second}
			for pb.Next() {
				idIndexLock.Lock()
				if idIndex >= len(contentIDs) {
					idIndexLock.Unlock()
					b.Error("no more content IDs available")
					break
				}
				contentID := contentIDs[idIndex]
				idIndex++
				idIndexLock.Unlock()

				url := fmt.Sprintf("%s/api/v1/contents/%d", server.URL, contentID)
				req, err := http.NewRequest(http.MethodDelete, url, nil)
				if err != nil {
					b.Error(err)
					continue
				}

				resp, err := client.Do(req)
				if err != nil {
					b.Error(err)
					continue
				}

				if resp.StatusCode != http.StatusNoContent {
					resp.Body.Close()
					b.Errorf("unexpected status code: %d", resp.StatusCode)
					continue
				}
				resp.Body.Close()
			}
		})
	})
}
