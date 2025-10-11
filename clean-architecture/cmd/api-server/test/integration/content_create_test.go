package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"go-api-server-sample/cmd/api-server/internal/api/content"
	"go-api-server-sample/cmd/api-server/internal/api/health"
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

type ContentCreateIntegrationTestSuite struct {
	suite.Suite
	container  *postgres.PostgresContainer
	db         *gorm.DB
	server     *httptest.Server
	httpClient *http.Client
}

func (suite *ContentCreateIntegrationTestSuite) SetupSuite() {
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

func (suite *ContentCreateIntegrationTestSuite) TearDownSuite() {
	ctx := context.Background()
	if suite.server != nil {
		suite.server.Close()
	}
	if suite.container != nil {
		suite.container.Terminate(ctx)
	}
}

func (suite *ContentCreateIntegrationTestSuite) SetupSubTest() {
	// テストデータクリーンアップ
	suite.db.Exec("DELETE FROM contents")
}

func (suite *ContentCreateIntegrationTestSuite) setupRouter() *gin.Engine {
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
	}

	return r
}

func (suite *ContentCreateIntegrationTestSuite) TestCreate() {
	suite.Run("正常にコンテンツを作成できる", func() {
		// Given
		reqBody := map[string]string{
			"title":        "テストタイトル",
			"body":         "テスト本文",
			"content_type": "article",
			"author":       "テスト作成者",
		}
		jsonBytes, _ := json.Marshal(reqBody)

		// When: HTTPリクエストを送信
		resp, err := suite.httpClient.Post(
			suite.server.URL+"/api/v1/contents",
			"application/json",
			bytes.NewBuffer(jsonBytes),
		)
		suite.Require().NoError(err)
		defer resp.Body.Close()

		// Then
		assert.Equal(suite.T(), http.StatusCreated, resp.StatusCode)

		var response entities.Content
		err = json.NewDecoder(resp.Body).Decode(&response)
		assert.NoError(suite.T(), err)
		assert.Greater(suite.T(), response.ID, uint(0))
		assert.Equal(suite.T(), "テストタイトル", response.Title)
		assert.Equal(suite.T(), "テスト本文", response.Body)
		assert.Equal(suite.T(), "article", response.ContentType)
		assert.Equal(suite.T(), "テスト作成者", response.Author)
		assert.NotZero(suite.T(), response.CreatedAt)
		assert.NotZero(suite.T(), response.UpdatedAt)
	})

	suite.Run("titleが空の場合はバリデーションエラー", func() {
		// Given
		reqBody := map[string]string{
			"title":        "",
			"body":         "テスト本文",
			"content_type": "article",
			"author":       "テスト作成者",
		}
		jsonBytes, _ := json.Marshal(reqBody)

		// When
		resp, err := suite.httpClient.Post(
			suite.server.URL+"/api/v1/contents",
			"application/json",
			bytes.NewBuffer(jsonBytes),
		)
		suite.Require().NoError(err)
		defer resp.Body.Close()

		// Then
		assert.Equal(suite.T(), http.StatusBadRequest, resp.StatusCode)

		var response map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&response)
		assert.NoError(suite.T(), err)
		assert.Equal(suite.T(), float64(http.StatusBadRequest), response["code"])
		assert.Contains(suite.T(), response["message"], "不正なリクエストです")
	})

	suite.Run("bodyが空の場合はバリデーションエラー", func() {
		// Given
		reqBody := map[string]string{
			"title":        "テストタイトル",
			"body":         "",
			"content_type": "article",
			"author":       "テスト作成者",
		}
		jsonBytes, _ := json.Marshal(reqBody)

		// When
		resp, err := suite.httpClient.Post(
			suite.server.URL+"/api/v1/contents",
			"application/json",
			bytes.NewBuffer(jsonBytes),
		)
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
		// Given
		reqBody := map[string]string{
			"title":        "テストタイトル",
			"body":         "テスト本文",
			"content_type": "invalid_type",
			"author":       "テスト作成者",
		}
		jsonBytes, _ := json.Marshal(reqBody)

		// When
		resp, err := suite.httpClient.Post(
			suite.server.URL+"/api/v1/contents",
			"application/json",
			bytes.NewBuffer(jsonBytes),
		)
		suite.Require().NoError(err)
		defer resp.Body.Close()

		// Then
		assert.Equal(suite.T(), http.StatusBadRequest, resp.StatusCode)

		var response map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&response)
		assert.NoError(suite.T(), err)
		assert.Equal(suite.T(), float64(http.StatusBadRequest), response["code"])
	})

	suite.Run("authorが空の場合はバリデーションエラー", func() {
		// Given
		reqBody := map[string]string{
			"title":        "テストタイトル",
			"body":         "テスト本文",
			"content_type": "article",
			"author":       "",
		}
		jsonBytes, _ := json.Marshal(reqBody)

		// When
		resp, err := suite.httpClient.Post(
			suite.server.URL+"/api/v1/contents",
			"application/json",
			bytes.NewBuffer(jsonBytes),
		)
		suite.Require().NoError(err)
		defer resp.Body.Close()

		// Then
		assert.Equal(suite.T(), http.StatusBadRequest, resp.StatusCode)

		var response map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&response)
		assert.NoError(suite.T(), err)
		assert.Equal(suite.T(), float64(http.StatusBadRequest), response["code"])
	})

	suite.Run("JSONが不正な場合はバリデーションエラー", func() {
		// When
		resp, err := suite.httpClient.Post(
			suite.server.URL+"/api/v1/contents",
			"application/json",
			bytes.NewBuffer([]byte("invalid json")),
		)
		suite.Require().NoError(err)
		defer resp.Body.Close()

		// Then
		assert.Equal(suite.T(), http.StatusBadRequest, resp.StatusCode)
	})
}

func TestContentCreateIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(ContentCreateIntegrationTestSuite))
}
