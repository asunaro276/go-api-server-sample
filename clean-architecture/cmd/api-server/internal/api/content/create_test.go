package content

import (
	"bytes"
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

type CreateContentTestSuite struct {
	suite.Suite
	container *postgres.PostgresContainer
	db        *gorm.DB
	api       *ContentAPI
}

func (suite *CreateContentTestSuite) SetupSuite() {
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

	// ContentAPI初期化（本物のリポジトリを使用）
	repo := infraRepos.NewContentRepository(suite.db)
	suite.api = NewContentAPI(repo)
}

func (suite *CreateContentTestSuite) TearDownSuite() {
	ctx := context.Background()
	if suite.container != nil {
		suite.container.Terminate(ctx)
	}
}

func (suite *CreateContentTestSuite) SetupSubTest() {
	// テストデータクリーンアップ
	suite.db.Exec("DELETE FROM contents")
}

func (suite *CreateContentTestSuite) TestCreate() {
	suite.Run("正常にコンテンツを作成できる", func() {
		// Given
		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		reqBody := map[string]string{
			"title":        "テストタイトル",
			"body":         "テスト本文",
			"content_type": "article",
			"author":       "テスト作成者",
		}
		jsonBytes, _ := json.Marshal(reqBody)
		c.Request = httptest.NewRequest(http.MethodPost, "/api/v1/contents", bytes.NewBuffer(jsonBytes))
		c.Request.Header.Set("Content-Type", "application/json")

		// When
		suite.api.Create(c)

		// Then
		assert.Equal(suite.T(), http.StatusCreated, w.Code)

		var response entities.Content
		err := json.Unmarshal(w.Body.Bytes(), &response)
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
		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		reqBody := map[string]string{
			"title":        "",
			"body":         "テスト本文",
			"content_type": "article",
			"author":       "テスト作成者",
		}
		jsonBytes, _ := json.Marshal(reqBody)
		c.Request = httptest.NewRequest(http.MethodPost, "/api/v1/contents", bytes.NewBuffer(jsonBytes))
		c.Request.Header.Set("Content-Type", "application/json")

		// When
		suite.api.Create(c)

		// Then
		assert.Equal(suite.T(), http.StatusBadRequest, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(suite.T(), err)
		assert.Equal(suite.T(), float64(http.StatusBadRequest), response["code"])
		assert.Contains(suite.T(), response["message"], "不正なリクエストです")
	})

	suite.Run("bodyが空の場合はバリデーションエラー", func() {
		// Given
		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		reqBody := map[string]string{
			"title":        "テストタイトル",
			"body":         "",
			"content_type": "article",
			"author":       "テスト作成者",
		}
		jsonBytes, _ := json.Marshal(reqBody)
		c.Request = httptest.NewRequest(http.MethodPost, "/api/v1/contents", bytes.NewBuffer(jsonBytes))
		c.Request.Header.Set("Content-Type", "application/json")

		// When
		suite.api.Create(c)

		// Then
		assert.Equal(suite.T(), http.StatusBadRequest, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(suite.T(), err)
		assert.Equal(suite.T(), float64(http.StatusBadRequest), response["code"])
	})

	suite.Run("content_typeが不正な値の場合はバリデーションエラー", func() {
		// Given
		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		reqBody := map[string]string{
			"title":        "テストタイトル",
			"body":         "テスト本文",
			"content_type": "invalid_type",
			"author":       "テスト作成者",
		}
		jsonBytes, _ := json.Marshal(reqBody)
		c.Request = httptest.NewRequest(http.MethodPost, "/api/v1/contents", bytes.NewBuffer(jsonBytes))
		c.Request.Header.Set("Content-Type", "application/json")

		// When
		suite.api.Create(c)

		// Then
		assert.Equal(suite.T(), http.StatusBadRequest, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(suite.T(), err)
		assert.Equal(suite.T(), float64(http.StatusBadRequest), response["code"])
	})

	suite.Run("authorが空の場合はバリデーションエラー", func() {
		// Given
		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		reqBody := map[string]string{
			"title":        "テストタイトル",
			"body":         "テスト本文",
			"content_type": "article",
			"author":       "",
		}
		jsonBytes, _ := json.Marshal(reqBody)
		c.Request = httptest.NewRequest(http.MethodPost, "/api/v1/contents", bytes.NewBuffer(jsonBytes))
		c.Request.Header.Set("Content-Type", "application/json")

		// When
		suite.api.Create(c)

		// Then
		assert.Equal(suite.T(), http.StatusBadRequest, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(suite.T(), err)
		assert.Equal(suite.T(), float64(http.StatusBadRequest), response["code"])
	})

	suite.Run("JSONが不正な場合はバリデーションエラー", func() {
		// Given
		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		c.Request = httptest.NewRequest(http.MethodPost, "/api/v1/contents", bytes.NewBuffer([]byte("invalid json")))
		c.Request.Header.Set("Content-Type", "application/json")

		// When
		suite.api.Create(c)

		// Then
		assert.Equal(suite.T(), http.StatusBadRequest, w.Code)
	})
}

func TestCreateContentTestSuite(t *testing.T) {
	suite.Run(t, new(CreateContentTestSuite))
}
