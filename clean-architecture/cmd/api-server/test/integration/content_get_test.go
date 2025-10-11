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

type ContentGetIntegrationTestSuite struct {
	suite.Suite
	container  *postgres.PostgresContainer
	db         *gorm.DB
	server     *httptest.Server
	httpClient *http.Client
}

func (suite *ContentGetIntegrationTestSuite) SetupSuite() {
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

func (suite *ContentGetIntegrationTestSuite) TearDownSuite() {
	ctx := context.Background()
	if suite.server != nil {
		suite.server.Close()
	}
	if suite.container != nil {
		suite.container.Terminate(ctx)
	}
}

func (suite *ContentGetIntegrationTestSuite) SetupSubTest() {
	// テストデータクリーンアップ
	suite.db.Exec("DELETE FROM contents")
}

func (suite *ContentGetIntegrationTestSuite) setupRouter() *gin.Engine {
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
		contents.GET("/:id", contentAPI.GetByID)
	}

	return r
}

func (suite *ContentGetIntegrationTestSuite) TestGetByID() {
	suite.Run("存在するIDでコンテンツを取得できる", func() {
		// Given: DBに直接テストデータを挿入
		testContent := &entities.Content{
			Title:       "テストタイトル",
			Body:        "テスト本文",
			ContentType: "article",
			Author:      "テスト作成者",
		}
		err := suite.db.Create(testContent).Error
		suite.Require().NoError(err)

		// When: HTTPリクエストで取得
		resp, err := suite.httpClient.Get(
			fmt.Sprintf("%s/api/v1/contents/%d", suite.server.URL, testContent.ID),
		)
		suite.Require().NoError(err)
		defer resp.Body.Close()

		// Then
		assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

		var response entities.Content
		err = json.NewDecoder(resp.Body).Decode(&response)
		assert.NoError(suite.T(), err)
		assert.Equal(suite.T(), testContent.ID, response.ID)
		assert.Equal(suite.T(), "テストタイトル", response.Title)
		assert.Equal(suite.T(), "テスト本文", response.Body)
		assert.Equal(suite.T(), "article", response.ContentType)
		assert.Equal(suite.T(), "テスト作成者", response.Author)
	})

	suite.Run("存在しないIDでは404エラー", func() {
		// When
		resp, err := suite.httpClient.Get(
			suite.server.URL + "/api/v1/contents/99999",
		)
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
		resp, err := suite.httpClient.Get(
			suite.server.URL + "/api/v1/contents/invalid",
		)
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
		resp, err := suite.httpClient.Get(
			suite.server.URL + "/api/v1/contents/-1",
		)
		suite.Require().NoError(err)
		defer resp.Body.Close()

		// Then
		assert.Equal(suite.T(), http.StatusBadRequest, resp.StatusCode)
	})
}

func TestContentGetIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(ContentGetIntegrationTestSuite))
}
