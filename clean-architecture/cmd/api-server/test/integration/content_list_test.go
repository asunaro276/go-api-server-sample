package integration

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"go-api-server-sample/cmd/api-server/internal/api/health"
	"go-api-server-sample/cmd/api-server/internal/container"
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

type ContentListIntegrationTestSuite struct {
	suite.Suite
	container  *postgres.PostgresContainer
	db         *gorm.DB
	server     *httptest.Server
	httpClient *http.Client
}

func (suite *ContentListIntegrationTestSuite) SetupSuite() {
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

func (suite *ContentListIntegrationTestSuite) TearDownSuite() {
	ctx := context.Background()
	if suite.server != nil {
		suite.server.Close()
	}
	if suite.container != nil {
		suite.container.Terminate(ctx)
	}
}

func (suite *ContentListIntegrationTestSuite) SetupSubTest() {
	// テストデータクリーンアップ
	suite.db.Exec("DELETE FROM contents")
}

func (suite *ContentListIntegrationTestSuite) setupRouter() *gin.Engine {
	r := gin.New()

	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(middleware.CORS())

	healthAPI := health.NewHealthAPI(suite.db)
	r.GET("/health", healthAPI.Check)

	deps := container.NewContainer(suite.db)
	v1 := r.Group("/api/v1")
	v1.Use(middleware.ErrorHandler())

	contents := v1.Group("/contents")
	{
		contents.POST("", deps.ContentAPI.Create)
		contents.GET("", deps.ContentAPI.List)
	}

	return r
}

func (suite *ContentListIntegrationTestSuite) createContent(title, body, contentType, author string) {
	content := &entities.Content{
		Title:       title,
		Body:        body,
		ContentType: contentType,
		Author:      author,
	}
	err := suite.db.Create(content).Error
	suite.Require().NoError(err)
}

func (suite *ContentListIntegrationTestSuite) TestList() {
	suite.Run("全件取得できる", func() {
		// Given: テストデータ作成
		suite.createContent("記事1", "本文1", "article", "作成者A")
		suite.createContent("ブログ1", "本文2", "blog", "作成者B")
		suite.createContent("記事2", "本文3", "article", "作成者A")

		// When: HTTPリクエストを送信
		resp, err := suite.httpClient.Get(suite.server.URL + "/api/v1/contents")
		suite.Require().NoError(err)
		defer resp.Body.Close()

		// Then
		assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

		var response map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&response)
		assert.NoError(suite.T(), err)
		assert.Equal(suite.T(), float64(3), response["total"])

		contentsList := response["contents"].([]interface{})
		assert.Len(suite.T(), contentsList, 3)
	})

	suite.Run("content_typeでフィルタリングできる", func() {
		// Given
		suite.createContent("記事1", "本文1", "article", "作成者A")
		suite.createContent("ブログ1", "本文2", "blog", "作成者B")
		suite.createContent("記事2", "本文3", "article", "作成者A")

		// When
		resp, err := suite.httpClient.Get(
			suite.server.URL + "/api/v1/contents?content_type=article",
		)
		suite.Require().NoError(err)
		defer resp.Body.Close()

		// Then
		assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

		var response map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&response)
		assert.NoError(suite.T(), err)
		assert.Equal(suite.T(), float64(2), response["total"])

		contentsList := response["contents"].([]interface{})
		assert.Len(suite.T(), contentsList, 2)
	})

	suite.Run("authorでフィルタリングできる", func() {
		// Given
		suite.createContent("記事1", "本文1", "article", "作成者A")
		suite.createContent("ブログ1", "本文2", "blog", "作成者B")
		suite.createContent("記事2", "本文3", "article", "作成者A")

		// When
		resp, err := suite.httpClient.Get(
			suite.server.URL + "/api/v1/contents?author=作成者A",
		)
		suite.Require().NoError(err)
		defer resp.Body.Close()

		// Then
		assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

		var response map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&response)
		assert.NoError(suite.T(), err)
		assert.Equal(suite.T(), float64(2), response["total"])

		contentsList := response["contents"].([]interface{})
		assert.Len(suite.T(), contentsList, 2)
	})

	suite.Run("ページネーションが機能する", func() {
		// Given
		for i := 1; i <= 5; i++ {
			suite.createContent("記事", "本文", "article", "作成者")
		}

		// When
		resp, err := suite.httpClient.Get(
			suite.server.URL + "/api/v1/contents?limit=2&offset=0",
		)
		suite.Require().NoError(err)
		defer resp.Body.Close()

		// Then
		assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

		var response map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&response)
		assert.NoError(suite.T(), err)
		assert.Equal(suite.T(), float64(5), response["total"])
		assert.Equal(suite.T(), float64(2), response["limit"])
		assert.Equal(suite.T(), float64(0), response["offset"])

		contentsList := response["contents"].([]interface{})
		assert.Len(suite.T(), contentsList, 2)
	})

	suite.Run("不正なcontent_typeでは400エラー", func() {
		// When
		resp, err := suite.httpClient.Get(
			suite.server.URL + "/api/v1/contents?content_type=invalid",
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

	suite.Run("負のlimitでは400エラー", func() {
		// When
		resp, err := suite.httpClient.Get(
			suite.server.URL + "/api/v1/contents?limit=-1",
		)
		suite.Require().NoError(err)
		defer resp.Body.Close()

		// Then
		assert.Equal(suite.T(), http.StatusBadRequest, resp.StatusCode)
	})

	suite.Run("負のoffsetでは400エラー", func() {
		// When
		resp, err := suite.httpClient.Get(
			suite.server.URL + "/api/v1/contents?offset=-1",
		)
		suite.Require().NoError(err)
		defer resp.Body.Close()

		// Then
		assert.Equal(suite.T(), http.StatusBadRequest, resp.StatusCode)
	})
}

func TestContentListIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(ContentListIntegrationTestSuite))
}
