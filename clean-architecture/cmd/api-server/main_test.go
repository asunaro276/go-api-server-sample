package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"go-api-server-sample/cmd/api-server/internal/container"
	"go-api-server-sample/internal/infrastructure/database"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type ContractTestSuite struct {
	suite.Suite
	router    *gin.Engine
	db        *gorm.DB
	contentID uint
}

func (suite *ContractTestSuite) SetupSuite() {
	gin.SetMode(gin.TestMode)

	// テスト用データベース接続
	db, err := database.Connect()
	suite.Require().NoError(err)
	suite.db = db

	// マイグレーション実行
	err = database.Migrate(db)
	suite.Require().NoError(err)

	// 依存性注入コンテナ作成
	deps := container.NewContainer(db)

	// 実際のルーター設定
	suite.router = setupRouter(deps)
}

func (suite *ContractTestSuite) SetupSubTest() {
	// 各サブテスト前にデータベースをクリーンアップ
	suite.db.Exec("DELETE FROM contents")
}

func (suite *ContractTestSuite) TestHealthCheckContract() {
	suite.Run("ヘルスチェックAPIが正常にレスポンスを返す", func() {
		req, _ := http.NewRequest("GET", "/health", nil)
		w := httptest.NewRecorder()
		suite.router.ServeHTTP(w, req)

		assert.Equal(suite.T(), http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(suite.T(), err)

		// 必須フィールドの存在確認
		assert.Contains(suite.T(), response, "status")
		assert.Contains(suite.T(), response, "timestamp")

		// ステータス値の検証
		status, ok := response["status"].(string)
		assert.True(suite.T(), ok)
		assert.Contains(suite.T(), []string{"healthy", "unhealthy"}, status)
	})
}

func (suite *ContractTestSuite) TestCreateContentContract() {
	suite.Run("コンテンツ作成APIが正常にレスポンスを返す", func() {
		requestBody := map[string]interface{}{
			"title":        "テスト記事",
			"body":         "これはテスト記事の本文です。",
			"content_type": "article",
			"author":       "テストユーザー",
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

		// 必須フィールドの存在確認
		assert.Contains(suite.T(), response, "id")
		assert.Contains(suite.T(), response, "title")
		assert.Contains(suite.T(), response, "body")
		assert.Contains(suite.T(), response, "content_type")
		assert.Contains(suite.T(), response, "author")
		assert.Contains(suite.T(), response, "created_at")
		assert.Contains(suite.T(), response, "updated_at")

		// レスポンス値の検証
		assert.Equal(suite.T(), requestBody["title"], response["title"])
		assert.Equal(suite.T(), requestBody["body"], response["body"])
		assert.Equal(suite.T(), requestBody["content_type"], response["content_type"])
		assert.Equal(suite.T(), requestBody["author"], response["author"])
	})
}

func (suite *ContractTestSuite) TestListContentsContract() {
	suite.Run("コンテンツ一覧取得APIが正常にレスポンスを返す", func() {
		req, _ := http.NewRequest("GET", "/api/v1/contents", nil)
		w := httptest.NewRecorder()
		suite.router.ServeHTTP(w, req)

		assert.Equal(suite.T(), http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(suite.T(), err)

		// 必須フィールドの存在確認
		assert.Contains(suite.T(), response, "contents")
		assert.Contains(suite.T(), response, "total")
		assert.Contains(suite.T(), response, "limit")
		assert.Contains(suite.T(), response, "offset")

		// 配列型の確認
		contents, ok := response["contents"].([]interface{})
		assert.True(suite.T(), ok)
		assert.NotNil(suite.T(), contents)
	})
}

func (suite *ContractTestSuite) TestGetContentByIDContract() {
	suite.Run("コンテンツ詳細取得APIが正常にレスポンスを返す", func() {
		// まずコンテンツを作成
		requestBody := map[string]interface{}{
			"title":        "テスト記事",
			"body":         "これはテスト記事の本文です。",
			"content_type": "article",
			"author":       "テストユーザー",
		}

		jsonBody, _ := json.Marshal(requestBody)
		createReq, _ := http.NewRequest("POST", "/api/v1/contents", bytes.NewBuffer(jsonBody))
		createReq.Header.Set("Content-Type", "application/json")
		createW := httptest.NewRecorder()
		suite.router.ServeHTTP(createW, createReq)

		// 作成されたコンテンツのIDを取得
		var createResponse map[string]interface{}
		err := json.Unmarshal(createW.Body.Bytes(), &createResponse)
		suite.Require().NoError(err)
		contentID := uint(createResponse["id"].(float64))

		// 詳細取得テスト
		req, _ := http.NewRequest("GET", fmt.Sprintf("/api/v1/contents/%d", contentID), nil)
		w := httptest.NewRecorder()
		suite.router.ServeHTTP(w, req)

		assert.Equal(suite.T(), http.StatusOK, w.Code)

		var response map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(suite.T(), err)

		// 必須フィールドの存在確認
		assert.Contains(suite.T(), response, "id")
		assert.Contains(suite.T(), response, "title")
		assert.Contains(suite.T(), response, "body")
		assert.Contains(suite.T(), response, "content_type")
		assert.Contains(suite.T(), response, "author")
		assert.Contains(suite.T(), response, "created_at")
		assert.Contains(suite.T(), response, "updated_at")
	})
}

func (suite *ContractTestSuite) TestUpdateContentContract() {
	suite.Run("コンテンツ更新APIが正常にレスポンスを返す", func() {
		// まずコンテンツを作成
		createRequestBody := map[string]interface{}{
			"title":        "元のテスト記事",
			"body":         "これは元のテスト記事の本文です。",
			"content_type": "article",
			"author":       "元の作成者",
		}

		jsonBody, _ := json.Marshal(createRequestBody)
		createReq, _ := http.NewRequest("POST", "/api/v1/contents", bytes.NewBuffer(jsonBody))
		createReq.Header.Set("Content-Type", "application/json")
		createW := httptest.NewRecorder()
		suite.router.ServeHTTP(createW, createReq)

		// 作成されたコンテンツのIDを取得
		var createResponse map[string]interface{}
		err := json.Unmarshal(createW.Body.Bytes(), &createResponse)
		suite.Require().NoError(err)
		contentID := uint(createResponse["id"].(float64))

		// 更新テスト
		updateRequestBody := map[string]interface{}{
			"title":        "更新されたテスト記事",
			"body":         "これは更新されたテスト記事の本文です。",
			"content_type": "blog",
			"author":       "更新者",
		}

		jsonBody, _ = json.Marshal(updateRequestBody)
		req, _ := http.NewRequest("PUT", fmt.Sprintf("/api/v1/contents/%d", contentID), bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		suite.router.ServeHTTP(w, req)

		assert.Equal(suite.T(), http.StatusOK, w.Code)

		var response map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(suite.T(), err)

		// 必須フィールドの存在確認
		assert.Contains(suite.T(), response, "id")
		assert.Contains(suite.T(), response, "title")
		assert.Contains(suite.T(), response, "body")
		assert.Contains(suite.T(), response, "content_type")
		assert.Contains(suite.T(), response, "author")
		assert.Contains(suite.T(), response, "created_at")
		assert.Contains(suite.T(), response, "updated_at")

		// 更新値の検証
		assert.Equal(suite.T(), updateRequestBody["title"], response["title"])
		assert.Equal(suite.T(), updateRequestBody["body"], response["body"])
		assert.Equal(suite.T(), updateRequestBody["content_type"], response["content_type"])
		assert.Equal(suite.T(), updateRequestBody["author"], response["author"])
	})
}

func (suite *ContractTestSuite) TestDeleteContentContract() {
	suite.Run("コンテンツ削除APIが正常にレスポンスを返す", func() {
		// まずコンテンツを作成
		createRequestBody := map[string]interface{}{
			"title":        "削除対象記事",
			"body":         "これは削除対象の記事の本文です。",
			"content_type": "article",
			"author":       "作成者",
		}

		jsonBody, _ := json.Marshal(createRequestBody)
		createReq, _ := http.NewRequest("POST", "/api/v1/contents", bytes.NewBuffer(jsonBody))
		createReq.Header.Set("Content-Type", "application/json")
		createW := httptest.NewRecorder()
		suite.router.ServeHTTP(createW, createReq)

		// 作成されたコンテンツのIDを取得
		var createResponse map[string]interface{}
		err := json.Unmarshal(createW.Body.Bytes(), &createResponse)
		suite.Require().NoError(err)
		contentID := uint(createResponse["id"].(float64))

		// 削除テスト
		req, _ := http.NewRequest("DELETE", fmt.Sprintf("/api/v1/contents/%d", contentID), nil)
		w := httptest.NewRecorder()
		suite.router.ServeHTTP(w, req)

		assert.Equal(suite.T(), http.StatusNoContent, w.Code)
		assert.Empty(suite.T(), w.Body.String())
	})
}

func TestContractTestSuite(t *testing.T) {
	suite.Run(t, new(ContractTestSuite))
}
