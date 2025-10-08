package content

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"go-api-server-sample/internal/domain/entities"
	infraRepos "go-api-server-sample/internal/infrastructure/repositories"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	postgresDriver "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type ListContentsTestSuite struct {
	suite.Suite
	container *postgres.PostgresContainer
	db        *gorm.DB
	api       *ContentAPI
}

func (suite *ListContentsTestSuite) SetupSuite() {
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

	// ContentAPI初期化
	repo := infraRepos.NewContentRepository(suite.db)
	suite.api = NewContentAPI(repo)
}

func (suite *ListContentsTestSuite) TearDownSuite() {
	ctx := context.Background()
	if suite.container != nil {
		suite.container.Terminate(ctx)
	}
}

func (suite *ListContentsTestSuite) SetupSubTest() {
	// テストデータクリーンアップ
	suite.db.Exec("DELETE FROM contents")
}

func (suite *ListContentsTestSuite) TestList() {
	suite.Run("全件取得できる", func() {
		// Given: テストデータ作成
		ctx := context.Background()
		repo := infraRepos.NewContentRepository(suite.db)

		contents := []*entities.Content{
			func() *entities.Content { c, _ := entities.NewContent("記事1", "本文1", "article", "作成者A"); return c }(),
			func() *entities.Content { c, _ := entities.NewContent("ブログ1", "本文2", "blog", "作成者B"); return c }(),
			func() *entities.Content { c, _ := entities.NewContent("記事2", "本文3", "article", "作成者A"); return c }(),
		}

		for _, content := range contents {
			err := repo.Create(ctx, content)
			suite.Require().NoError(err)
		}

		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/api/v1/contents", nil)

		// When
		suite.api.List(c)

		// Then
		assert.Equal(suite.T(), http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(suite.T(), err)
		assert.Equal(suite.T(), float64(3), response["total"])

		contentsList := response["contents"].([]interface{})
		assert.Len(suite.T(), contentsList, 3)
	})

	suite.Run("content_typeでフィルタリングできる", func() {
		// Given
		ctx := context.Background()
		repo := infraRepos.NewContentRepository(suite.db)

		contents := []*entities.Content{
			func() *entities.Content { c, _ := entities.NewContent("記事1", "本文1", "article", "作成者A"); return c }(),
			func() *entities.Content { c, _ := entities.NewContent("ブログ1", "本文2", "blog", "作成者B"); return c }(),
			func() *entities.Content { c, _ := entities.NewContent("記事2", "本文3", "article", "作成者A"); return c }(),
		}

		for _, content := range contents {
			err := repo.Create(ctx, content)
			suite.Require().NoError(err)
		}

		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/api/v1/contents?content_type=article", nil)

		// When
		suite.api.List(c)

		// Then
		assert.Equal(suite.T(), http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(suite.T(), err)
		assert.Equal(suite.T(), float64(2), response["total"])

		contentsList := response["contents"].([]interface{})
		assert.Len(suite.T(), contentsList, 2)
	})

	suite.Run("authorでフィルタリングできる", func() {
		// Given
		ctx := context.Background()
		repo := infraRepos.NewContentRepository(suite.db)

		contents := []*entities.Content{
			func() *entities.Content { c, _ := entities.NewContent("記事1", "本文1", "article", "作成者A"); return c }(),
			func() *entities.Content { c, _ := entities.NewContent("ブログ1", "本文2", "blog", "作成者B"); return c }(),
			func() *entities.Content { c, _ := entities.NewContent("記事2", "本文3", "article", "作成者A"); return c }(),
		}

		for _, content := range contents {
			err := repo.Create(ctx, content)
			suite.Require().NoError(err)
		}

		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/api/v1/contents?author=作成者A", nil)

		// When
		suite.api.List(c)

		// Then
		assert.Equal(suite.T(), http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(suite.T(), err)
		assert.Equal(suite.T(), float64(2), response["total"])

		contentsList := response["contents"].([]interface{})
		assert.Len(suite.T(), contentsList, 2)
	})

	suite.Run("ページネーションが機能する", func() {
		// Given
		ctx := context.Background()
		repo := infraRepos.NewContentRepository(suite.db)

		for i := 1; i <= 5; i++ {
			content, _ := entities.NewContent("記事", "本文", "article", "作成者")
			err := repo.Create(ctx, content)
			suite.Require().NoError(err)
		}

		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/api/v1/contents?limit=2&offset=0", nil)

		// When
		suite.api.List(c)

		// Then
		assert.Equal(suite.T(), http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(suite.T(), err)
		assert.Equal(suite.T(), float64(5), response["total"])
		assert.Equal(suite.T(), float64(2), response["limit"])
		assert.Equal(suite.T(), float64(0), response["offset"])

		contentsList := response["contents"].([]interface{})
		assert.Len(suite.T(), contentsList, 2)
	})

	suite.Run("不正なcontent_typeでは400エラー", func() {
		// Given
		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/api/v1/contents?content_type=invalid", nil)

		// When
		suite.api.List(c)

		// Then
		assert.Equal(suite.T(), http.StatusBadRequest, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(suite.T(), err)
		assert.Equal(suite.T(), float64(http.StatusBadRequest), response["code"])
	})

	suite.Run("負のlimitでは400エラー", func() {
		// Given
		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/api/v1/contents?limit=-1", nil)

		// When
		suite.api.List(c)

		// Then
		assert.Equal(suite.T(), http.StatusBadRequest, w.Code)
	})

	suite.Run("負のoffsetでは400エラー", func() {
		// Given
		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/api/v1/contents?offset=-1", nil)

		// When
		suite.api.List(c)

		// Then
		assert.Equal(suite.T(), http.StatusBadRequest, w.Code)
	})
}

func TestListContentsTestSuite(t *testing.T) {
	suite.Run(t, new(ListContentsTestSuite))
}
