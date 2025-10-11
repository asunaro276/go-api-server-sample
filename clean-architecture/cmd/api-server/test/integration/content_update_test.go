package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"go-api-server-sample/cmd/api-server/internal/api/health"
	"go-api-server-sample/cmd/api-server/internal/api/content"
	"go-api-server-sample/cmd/api-server/internal/infrastructure/repositories"
	"go-api-server-sample/cmd/api-server/internal/middleware"
	"go-api-server-sample/internal/domain/entities"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	postgresDriver "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type ContentUpdateIntegrationTestSuite struct {
	suite.Suite
	container  *postgres.PostgresContainer
	db         *gorm.DB
	server     *httptest.Server
	httpClient *http.Client
}

func (suite *ContentUpdateIntegrationTestSuite) SetupSuite() {
	ctx := context.Background()

	// PostgreSQLコンテナ起動
	container, err := postgres.Run(ctx,
		"postgres:15",
		postgres.WithDatabase("testdb"),
		postgres.WithUsername("testuser"),
		postgres.WithPassword("testpass"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(60*time.Second),
		),
	)
	suite.Require().NoError(err)
	suite.container = container

	// DB接続とマイグレーション実行
	connStr, err := container.ConnectionString(ctx, "sslmode=disable")
	suite.Require().NoError(err)

	suite.db, err = gorm.Open(postgresDriver.Open(connStr), &gorm.Config{})
	suite.Require().NoError(err)

	err = suite.db.AutoMigrate(&entities.Content{})
	suite.Require().NoError(err)

	// ルーター設定
	gin.SetMode(gin.TestMode)
	router := suite.setupRouter()

	// テストサーバー起動
	suite.server = httptest.NewServer(router)
	suite.httpClient = &http.Client{
		Timeout: 10 * time.Second,
	}
}

func (suite *ContentUpdateIntegrationTestSuite) TearDownSuite() {
	ctx := context.Background()
	if suite.server != nil {
		suite.server.Close()
	}
	if suite.container != nil {
		suite.container.Terminate(ctx)
	}
}

func (suite *ContentUpdateIntegrationTestSuite) SetupSubTest() {
	// テストデータクリーンアップ
	suite.db.Exec("DELETE FROM contents")
}

func (suite *ContentUpdateIntegrationTestSuite) setupRouter() *gin.Engine {
	r := gin.New()

	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(middleware.CORS())

	healthAPI := health.NewHealthAPI(suite.db)
	r.GET("/health", healthAPI.Check)

	// リポジトリとAPIを直接初期化
	contentRepo := repositories.NewContentRepository(suite.db)
	contentAPI := content.NewContentAPI(contentRepo)
	v1 := r.Group("/api/v1")
	v1.Use(middleware.ErrorHandler())

	contents := v1.Group("/contents")
	{
		contents.POST("", contentAPI.Create)
		contents.PUT("/:id", contentAPI.Update)
	}

	return r
}

func (suite *ContentUpdateIntegrationTestSuite) createContent(title, body, contentType, author string) uint {
	content := &entities.Content{
		Title:       title,
		Body:        body,
		ContentType: contentType,
		Author:      author,
	}
	err := suite.db.Create(content).Error
	suite.Require().NoError(err)
	return content.ID
}

func (suite *ContentUpdateIntegrationTestSuite) TestUpdate() {
	suite.Run("正常にコンテンツを更新できる", func() {
		// Given: テストデータ作成
		contentID := suite.createContent("元のタイトル", "元の本文", "article", "元の作成者")

		updateBody := map[string]string{
			"title":        "新しいタイトル",
			"body":         "新しい本文",
			"content_type": "blog",
			"author":       "新しい作成者",
		}
		jsonBytes, _ := json.Marshal(updateBody)

		// When: HTTPリクエストで更新
		req, _ := http.NewRequest(
			http.MethodPut,
			fmt.Sprintf("%s/api/v1/contents/%d", suite.server.URL, contentID),
			bytes.NewBuffer(jsonBytes),
		)
		req.Header.Set("Content-Type", "application/json")

		resp, err := suite.httpClient.Do(req)
		suite.Require().NoError(err)
		defer resp.Body.Close()

		// Then
		assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

		var response entities.Content
		err = json.NewDecoder(resp.Body).Decode(&response)
		assert.NoError(suite.T(), err)
		assert.Equal(suite.T(), contentID, response.ID)
		assert.Equal(suite.T(), "新しいタイトル", response.Title)
		assert.Equal(suite.T(), "新しい本文", response.Body)
		assert.Equal(suite.T(), "blog", response.ContentType)
		assert.Equal(suite.T(), "新しい作成者", response.Author)
	})

	suite.Run("存在しないIDでは404エラー", func() {
		// Given
		updateBody := map[string]string{
			"title":        "新しいタイトル",
			"body":         "新しい本文",
			"content_type": "blog",
			"author":       "新しい作成者",
		}
		jsonBytes, _ := json.Marshal(updateBody)

		// When
		req, _ := http.NewRequest(
			http.MethodPut,
			suite.server.URL+"/api/v1/contents/99999",
			bytes.NewBuffer(jsonBytes),
		)
		req.Header.Set("Content-Type", "application/json")

		resp, err := suite.httpClient.Do(req)
		suite.Require().NoError(err)
		defer resp.Body.Close()

		// Then
		assert.Equal(suite.T(), http.StatusNotFound, resp.StatusCode)

		var response map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&response)
		assert.NoError(suite.T(), err)
		assert.Equal(suite.T(), float64(http.StatusNotFound), response["code"])
		assert.Contains(suite.T(), response["message"], "指定されたコンテンツが見つかりません")
	})

	suite.Run("不正なID形式では400エラー", func() {
		// Given
		updateBody := map[string]string{
			"title":        "新しいタイトル",
			"body":         "新しい本文",
			"content_type": "blog",
			"author":       "新しい作成者",
		}
		jsonBytes, _ := json.Marshal(updateBody)

		// When
		req, _ := http.NewRequest(
			http.MethodPut,
			suite.server.URL+"/api/v1/contents/invalid",
			bytes.NewBuffer(jsonBytes),
		)
		req.Header.Set("Content-Type", "application/json")

		resp, err := suite.httpClient.Do(req)
		suite.Require().NoError(err)
		defer resp.Body.Close()

		// Then
		assert.Equal(suite.T(), http.StatusBadRequest, resp.StatusCode)

		var response map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&response)
		assert.NoError(suite.T(), err)
		assert.Equal(suite.T(), float64(http.StatusBadRequest), response["code"])
		assert.Contains(suite.T(), response["message"], "不正なIDです")
	})

	suite.Run("titleが空の場合はバリデーションエラー", func() {
		// Given: テストデータ作成
		contentID := suite.createContent("元のタイトル", "元の本文", "article", "元の作成者")

		updateBody := map[string]string{
			"title":        "",
			"body":         "新しい本文",
			"content_type": "blog",
			"author":       "新しい作成者",
		}
		jsonBytes, _ := json.Marshal(updateBody)

		// When
		req, _ := http.NewRequest(
			http.MethodPut,
			fmt.Sprintf("%s/api/v1/contents/%d", suite.server.URL, contentID),
			bytes.NewBuffer(jsonBytes),
		)
		req.Header.Set("Content-Type", "application/json")

		resp, err := suite.httpClient.Do(req)
		suite.Require().NoError(err)
		defer resp.Body.Close()

		// Then
		assert.Equal(suite.T(), http.StatusBadRequest, resp.StatusCode)

		var response map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&response)
		assert.NoError(suite.T(), err)
		assert.Equal(suite.T(), float64(http.StatusBadRequest), response["code"])
	})

	suite.Run("content_typeが不正な値の場合はバリデーションエラー", func() {
		// Given: テストデータ作成
		contentID := suite.createContent("元のタイトル", "元の本文", "article", "元の作成者")

		updateBody := map[string]string{
			"title":        "新しいタイトル",
			"body":         "新しい本文",
			"content_type": "invalid_type",
			"author":       "新しい作成者",
		}
		jsonBytes, _ := json.Marshal(updateBody)

		// When
		req, _ := http.NewRequest(
			http.MethodPut,
			fmt.Sprintf("%s/api/v1/contents/%d", suite.server.URL, contentID),
			bytes.NewBuffer(jsonBytes),
		)
		req.Header.Set("Content-Type", "application/json")

		resp, err := suite.httpClient.Do(req)
		suite.Require().NoError(err)
		defer resp.Body.Close()

		// Then
		assert.Equal(suite.T(), http.StatusBadRequest, resp.StatusCode)
	})

	suite.Run("JSONが不正な場合はバリデーションエラー", func() {
		// Given: テストデータ作成
		contentID := suite.createContent("元のタイトル", "元の本文", "article", "元の作成者")

		// When
		req, _ := http.NewRequest(
			http.MethodPut,
			fmt.Sprintf("%s/api/v1/contents/%d", suite.server.URL, contentID),
			bytes.NewBuffer([]byte("invalid json")),
		)
		req.Header.Set("Content-Type", "application/json")

		resp, err := suite.httpClient.Do(req)
		suite.Require().NoError(err)
		defer resp.Body.Close()

		// Then
		assert.Equal(suite.T(), http.StatusBadRequest, resp.StatusCode)
	})
}

func TestContentUpdateIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(ContentUpdateIntegrationTestSuite))
}
