package health

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

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

type HealthAPITestSuite struct {
	suite.Suite
	container *postgres.PostgresContainer
	db        *gorm.DB
	api       *HealthAPI
}

func (suite *HealthAPITestSuite) SetupSuite() {
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

	// HealthAPI初期化
	suite.api = NewHealthAPI(suite.db)
}

func (suite *HealthAPITestSuite) TearDownSuite() {
	ctx := context.Background()
	if suite.container != nil {
		suite.container.Terminate(ctx)
	}
}

func (suite *HealthAPITestSuite) TestCheck() {
	suite.Run("DB接続が正常な場合はhealthyを返す", func() {
		// Given
		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/health", nil)

		// When
		suite.api.Check(c)

		// Then
		assert.Equal(suite.T(), http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(suite.T(), err)
		assert.Equal(suite.T(), "healthy", response["status"])
		assert.Equal(suite.T(), "connected", response["database"])
		assert.NotNil(suite.T(), response["timestamp"])
	})

	suite.Run("DBがnilの場合でもステータスを返す", func() {
		// Given
		apiWithoutDB := NewHealthAPI(nil)

		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/health", nil)

		// When
		apiWithoutDB.Check(c)

		// Then
		assert.Equal(suite.T(), http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(suite.T(), err)
		assert.Equal(suite.T(), "healthy", response["status"])
		assert.NotNil(suite.T(), response["timestamp"])
	})
}

func TestHealthAPITestSuite(t *testing.T) {
	suite.Run(t, new(HealthAPITestSuite))
}
