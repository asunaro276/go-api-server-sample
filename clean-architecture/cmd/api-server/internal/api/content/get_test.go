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

type GetContentTestSuite struct {
	suite.Suite
	container *postgres.PostgresContainer
	db        *gorm.DB
	api       *ContentAPI
}

func (suite *GetContentTestSuite) SetupSuite() {
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

func (suite *GetContentTestSuite) TearDownSuite() {
	ctx := context.Background()
	if suite.container != nil {
		suite.container.Terminate(ctx)
	}
}

func (suite *GetContentTestSuite) SetupSubTest() {
	// テストデータクリーンアップ
	suite.db.Exec("DELETE FROM contents")
}

func (suite *GetContentTestSuite) TestGetByID() {
	suite.Run("存在するIDでコンテンツを取得できる", func() {
		// Given: テストデータ作成
		ctx := context.Background()
		repo := infraRepos.NewContentRepository(suite.db)
		content, _ := entities.NewContent("テストタイトル", "テスト本文", "article", "テスト作成者")
		err := repo.Create(ctx, content)
		suite.Require().NoError(err)

		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/api/v1/contents/1", nil)
		c.Params = gin.Params{
			{Key: "id", Value: "1"},
		}

		// When
		suite.api.GetByID(c)

		// Then
		assert.Equal(suite.T(), http.StatusOK, w.Code)

		var response entities.Content
		err = json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(suite.T(), err)
		assert.Equal(suite.T(), content.ID, response.ID)
		assert.Equal(suite.T(), "テストタイトル", response.Title)
		assert.Equal(suite.T(), "テスト本文", response.Body)
		assert.Equal(suite.T(), "article", response.ContentType)
		assert.Equal(suite.T(), "テスト作成者", response.Author)
	})

	suite.Run("存在しないIDでは404エラー", func() {
		// Given
		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/api/v1/contents/99999", nil)
		c.Params = gin.Params{
			{Key: "id", Value: "99999"},
		}

		// When
		suite.api.GetByID(c)

		// Then
		assert.Equal(suite.T(), http.StatusNotFound, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(suite.T(), err)
		assert.Equal(suite.T(), float64(http.StatusNotFound), response["code"])
		assert.Contains(suite.T(), response["message"], "指定されたコンテンツが見つかりません")
	})

	suite.Run("不正なID形式では400エラー", func() {
		// Given
		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/api/v1/contents/invalid", nil)
		c.Params = gin.Params{
			{Key: "id", Value: "invalid"},
		}

		// When
		suite.api.GetByID(c)

		// Then
		assert.Equal(suite.T(), http.StatusBadRequest, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(suite.T(), err)
		assert.Equal(suite.T(), float64(http.StatusBadRequest), response["code"])
		assert.Contains(suite.T(), response["message"], "不正なIDです")
	})

	suite.Run("負のIDでは400エラー", func() {
		// Given
		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/api/v1/contents/-1", nil)
		c.Params = gin.Params{
			{Key: "id", Value: "-1"},
		}

		// When
		suite.api.GetByID(c)

		// Then
		assert.Equal(suite.T(), http.StatusBadRequest, w.Code)
	})
}

func TestGetContentTestSuite(t *testing.T) {
	suite.Run(t, new(GetContentTestSuite))
}
