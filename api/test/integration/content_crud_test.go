//go:build integration

package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"gorm.io/gorm"
)

type ContentCRUDIntegrationTestSuite struct {
	suite.Suite
	container *postgres.PostgresContainer
	db        *gorm.DB
	router    *gin.Engine
}

func (suite *ContentCRUDIntegrationTestSuite) SetupSuite() {
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

	// TODO: データベース接続とマイグレーション（実装後に有効化）
	// connStr, err := container.ConnectionString(ctx, "sslmode=disable")
	// suite.Require().NoError(err)
	//
	// suite.db, err = gorm.Open(postgres.Open(connStr), &gorm.Config{})
	// suite.Require().NoError(err)
	//
	// err = suite.db.AutoMigrate(&entities.Content{})
	// suite.Require().NoError(err)

	gin.SetMode(gin.TestMode)
	// TODO: ルーター設定（実装後に有効化）
	suite.router = setupTestRouter()
}

func (suite *ContentCRUDIntegrationTestSuite) TearDownSuite() {
	ctx := context.Background()
	if suite.container != nil {
		suite.container.Terminate(ctx)
	}
}

func (suite *ContentCRUDIntegrationTestSuite) SetupSubTest() {
	// 各テスト前にデータをクリーンアップ
	// TODO: 実装後に有効化
	// suite.db.Exec("DELETE FROM contents")
}

func (suite *ContentCRUDIntegrationTestSuite) TestContentCRUDWorkflow() {
	suite.Run("コンテンツCRUDの完全なワークフローが動作する", func() {
		// 1. 初期状態では空のリストが返される
		suite.Run("初期状態でコンテンツ一覧が空", func() {
			req, _ := http.NewRequest("GET", "/api/v1/contents", nil)
			w := httptest.NewRecorder()
			suite.router.ServeHTTP(w, req)

			assert.Equal(suite.T(), http.StatusOK, w.Code)

			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(suite.T(), err)

			contents, ok := response["contents"].([]interface{})
			assert.True(suite.T(), ok)
			assert.Empty(suite.T(), contents)
			assert.Equal(suite.T(), float64(0), response["total"])
		})

		// 2. コンテンツを作成
		var contentID float64
		suite.Run("新しいコンテンツを作成", func() {
			requestBody := map[string]interface{}{
				"title":        "統合テスト記事",
				"body":         "これは統合テスト用の記事です。",
				"content_type": "article",
				"author":       "統合テストユーザー",
			}

			jsonBody, _ := json.Marshal(requestBody)
			req, _ := http.NewRequest("POST", "/api/v1/contents", bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			suite.router.ServeHTTP(w, req)

			assert.Equal(suite.T(), http.StatusCreated, w.Code)

			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(suite.T(), err)

			contentID = response["id"].(float64)
			assert.Greater(suite.T(), contentID, float64(0))
			assert.Equal(suite.T(), requestBody["title"], response["title"])
		})

		// 3. 作成したコンテンツを詳細取得
		suite.Run("作成したコンテンツの詳細を取得", func() {
			req, _ := http.NewRequest("GET", fmt.Sprintf("/api/v1/contents/%.0f", contentID), nil)
			w := httptest.NewRecorder()
			suite.router.ServeHTTP(w, req)

			assert.Equal(suite.T(), http.StatusOK, w.Code)

			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(suite.T(), err)

			assert.Equal(suite.T(), contentID, response["id"])
			assert.Equal(suite.T(), "統合テスト記事", response["title"])
		})

		// 4. コンテンツ一覧で作成したコンテンツが表示される
		suite.Run("コンテンツ一覧に作成したコンテンツが表示される", func() {
			req, _ := http.NewRequest("GET", "/api/v1/contents", nil)
			w := httptest.NewRecorder()
			suite.router.ServeHTTP(w, req)

			assert.Equal(suite.T(), http.StatusOK, w.Code)

			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(suite.T(), err)

			contents, ok := response["contents"].([]interface{})
			assert.True(suite.T(), ok)
			assert.Len(suite.T(), contents, 1)
			assert.Equal(suite.T(), float64(1), response["total"])

			firstContent := contents[0].(map[string]interface{})
			assert.Equal(suite.T(), contentID, firstContent["id"])
		})

		// 5. コンテンツを更新
		suite.Run("コンテンツを更新", func() {
			updateBody := map[string]interface{}{
				"title":        "更新された統合テスト記事",
				"body":         "これは更新された統合テスト用の記事です。",
				"content_type": "blog",
				"author":       "更新者",
			}

			jsonBody, _ := json.Marshal(updateBody)
			req, _ := http.NewRequest("PUT", fmt.Sprintf("/api/v1/contents/%.0f", contentID), bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			suite.router.ServeHTTP(w, req)

			assert.Equal(suite.T(), http.StatusOK, w.Code)

			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(suite.T(), err)

			assert.Equal(suite.T(), contentID, response["id"])
			assert.Equal(suite.T(), updateBody["title"], response["title"])
			assert.Equal(suite.T(), updateBody["content_type"], response["content_type"])
		})

		// 6. コンテンツを削除
		suite.Run("コンテンツを削除", func() {
			req, _ := http.NewRequest("DELETE", fmt.Sprintf("/api/v1/contents/%.0f", contentID), nil)
			w := httptest.NewRecorder()
			suite.router.ServeHTTP(w, req)

			assert.Equal(suite.T(), http.StatusNoContent, w.Code)
		})

		// 7. 削除されたコンテンツは取得できない
		suite.Run("削除されたコンテンツは404エラー", func() {
			req, _ := http.NewRequest("GET", fmt.Sprintf("/api/v1/contents/%.0f", contentID), nil)
			w := httptest.NewRecorder()
			suite.router.ServeHTTP(w, req)

			assert.Equal(suite.T(), http.StatusNotFound, w.Code)
		})

		// 8. 削除後は再び空のリストになる
		suite.Run("削除後にコンテンツ一覧が空になる", func() {
			req, _ := http.NewRequest("GET", "/api/v1/contents", nil)
			w := httptest.NewRecorder()
			suite.router.ServeHTTP(w, req)

			assert.Equal(suite.T(), http.StatusOK, w.Code)

			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(suite.T(), err)

			contents, ok := response["contents"].([]interface{})
			assert.True(suite.T(), ok)
			assert.Empty(suite.T(), contents)
			assert.Equal(suite.T(), float64(0), response["total"])
		})
	})
}

func (suite *ContentCRUDIntegrationTestSuite) TestContentFiltering() {
	suite.Run("コンテンツのフィルタリング機能が動作する", func() {
		// テストデータ作成（複数のコンテンツタイプと作成者）
		testContents := []map[string]interface{}{
			{
				"title":        "記事1",
				"body":         "記事の内容",
				"content_type": "article",
				"author":       "作成者A",
			},
			{
				"title":        "ブログ1",
				"body":         "ブログの内容",
				"content_type": "blog",
				"author":       "作成者B",
			},
			{
				"title":        "記事2",
				"body":         "記事の内容2",
				"content_type": "article",
				"author":       "作成者A",
			},
		}

		// テストデータを作成
		for _, content := range testContents {
			jsonBody, _ := json.Marshal(content)
			req, _ := http.NewRequest("POST", "/api/v1/contents", bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			suite.router.ServeHTTP(w, req)
			assert.Equal(suite.T(), http.StatusCreated, w.Code)
		}

		// コンテンツタイプでフィルタリング
		suite.Run("コンテンツタイプでフィルタリング", func() {
			req, _ := http.NewRequest("GET", "/api/v1/contents?content_type=article", nil)
			w := httptest.NewRecorder()
			suite.router.ServeHTTP(w, req)

			assert.Equal(suite.T(), http.StatusOK, w.Code)

			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(suite.T(), err)

			contents, ok := response["contents"].([]interface{})
			assert.True(suite.T(), ok)
			assert.Len(suite.T(), contents, 2) // article タイプが2件
		})

		// 作成者でフィルタリング
		suite.Run("作成者でフィルタリング", func() {
			req, _ := http.NewRequest("GET", "/api/v1/contents?author=作成者A", nil)
			w := httptest.NewRecorder()
			suite.router.ServeHTTP(w, req)

			assert.Equal(suite.T(), http.StatusOK, w.Code)

			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(suite.T(), err)

			contents, ok := response["contents"].([]interface{})
			assert.True(suite.T(), ok)
			assert.Len(suite.T(), contents, 2) // 作成者Aが2件
		})
	})
}

func TestContentCRUDIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(ContentCRUDIntegrationTestSuite))
}

// setupTestRouter は統合テスト用のルーター設定
func setupTestRouter() *gin.Engine {
	// TODO: 実装が完了するまでは空のルーターを返す
	r := gin.New()
	return r
}