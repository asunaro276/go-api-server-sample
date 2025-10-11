package performance

import (
	"bytes"
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

func setupContentUpdateBenchmark(b *testing.B) (*httptest.Server, *http.Client, uint) {
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
		contents.PUT("/:id", contentAPI.Update)
	}

	// テストサーバー起動
	server := httptest.NewServer(r)
	httpClient := &http.Client{Timeout: 10 * time.Second}

	return server, httpClient, testContent.ID
}

// BenchmarkContentUpdate はコンテンツ更新のベンチマーク
func BenchmarkContentUpdate(b *testing.B) {
	cleanupDB(b)
	server, httpClient, contentID := setupContentUpdateBenchmark(b)
	defer server.Close()

	url := fmt.Sprintf("%s/api/v1/contents/%d", server.URL, contentID)

	b.Run("シングルリクエスト", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			reqBody := map[string]string{
				"title":        fmt.Sprintf("更新タイトル%d", i),
				"body":         "更新本文",
				"content_type": "article",
				"author":       "更新作成者",
			}
			jsonBytes, _ := json.Marshal(reqBody)

			req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(jsonBytes))
			if err != nil {
				b.Fatal(err)
			}
			req.Header.Set("Content-Type", "application/json")

			resp, err := httpClient.Do(req)
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
			counter := 0
			for pb.Next() {
				counter++
				reqBody := map[string]string{
					"title":        fmt.Sprintf("更新タイトル%d", counter),
					"body":         "更新本文",
					"content_type": "article",
					"author":       "更新作成者",
				}
				jsonBytes, _ := json.Marshal(reqBody)

				req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(jsonBytes))
				if err != nil {
					b.Error(err)
					continue
				}
				req.Header.Set("Content-Type", "application/json")

				resp, err := client.Do(req)
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
