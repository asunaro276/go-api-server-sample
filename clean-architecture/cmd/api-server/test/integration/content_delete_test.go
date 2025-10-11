package integration

import (
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

type ContentDeleteIntegrationTestSuite struct {
	suite.Suite
	container  *postgres.PostgresContainer
	db         *gorm.DB
	server     *httptest.Server
	httpClient *http.Client
}

func (suite *ContentDeleteIntegrationTestSuite) SetupSuite() {
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

func (suite *ContentDeleteIntegrationTestSuite) TearDownSuite() {
	ctx := context.Background()
	if suite.server != nil {
		suite.server.Close()
	}
	if suite.container != nil {
		suite.container.Terminate(ctx)
	}
}

func (suite *ContentDeleteIntegrationTestSuite) SetupSubTest() {
	// テストデータクリーンアップ
	suite.db.Exec("DELETE FROM contents")
}

func (suite *ContentDeleteIntegrationTestSuite) setupRouter() *gin.Engine {
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
		contents.GET("", contentAPI.List)
		contents.GET("/:id", contentAPI.GetByID)
		contents.PUT("/:id", contentAPI.Update)
		contents.DELETE("/:id", contentAPI.Delete)
	}

	return r
}

func (suite *ContentDeleteIntegrationTestSuite) createContent(title, body, contentType, author string) uint {
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

func (suite *ContentDeleteIntegrationTestSuite) TestDelete() {
	suite.Run("正常にコンテンツを削除できる", func() {
		// Given: テストデータ作成
		contentID := suite.createContent("削除対象", "削除対象本文", "article", "作成者")

		// When: HTTPリクエストで削除
		req, _ := http.NewRequest(
			http.MethodDelete,
			fmt.Sprintf("%s/api/v1/contents/%d", suite.server.URL, contentID),
			nil,
		)

		resp, err := suite.httpClient.Do(req)
		suite.Require().NoError(err)
		defer resp.Body.Close()

		// Then
		assert.Equal(suite.T(), http.StatusNoContent, resp.StatusCode)

		// 削除確認：再度取得を試みる
		getResp, err := suite.httpClient.Get(
			fmt.Sprintf("%s/api/v1/contents/%d", suite.server.URL, contentID),
		)
		suite.Require().NoError(err)
		defer getResp.Body.Close()
		assert.Equal(suite.T(), http.StatusNotFound, getResp.StatusCode)
	})

	suite.Run("存在しないIDでは404エラー", func() {
		// When
		req, _ := http.NewRequest(
			http.MethodDelete,
			suite.server.URL+"/api/v1/contents/99999",
			nil,
		)

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
		// When
		req, _ := http.NewRequest(
			http.MethodDelete,
			suite.server.URL+"/api/v1/contents/invalid",
			nil,
		)

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

	suite.Run("負のIDでは400エラー", func() {
		// When
		req, _ := http.NewRequest(
			http.MethodDelete,
			suite.server.URL+"/api/v1/contents/-1",
			nil,
		)

		resp, err := suite.httpClient.Do(req)
		suite.Require().NoError(err)
		defer resp.Body.Close()

		// Then
		assert.Equal(suite.T(), http.StatusBadRequest, resp.StatusCode)
	})
}

func TestContentDeleteIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(ContentDeleteIntegrationTestSuite))
}
