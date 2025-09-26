package contract

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestPOSTUsersContract tests POST /api/v1/users endpoint contract
func TestPOSTUsersContract(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("should create user with valid request", func(t *testing.T) {
		// This test will fail until implementation is complete
		req := map[string]interface{}{
			"name":  "田中太郎",
			"email": "tanaka@example.com",
		}

		body, err := json.Marshal(req)
		require.NoError(t, err)

		w := httptest.NewRecorder()
		httpReq, err := http.NewRequest("POST", "/api/v1/users", bytes.NewBuffer(body))
		require.NoError(t, err)
		httpReq.Header.Set("Content-Type", "application/json")

		// This will fail until router is implemented
		router := gin.New()
		router.ServeHTTP(w, httpReq)

		// Expect 201 Created
		assert.Equal(t, http.StatusCreated, w.Code)

		// Verify response structure
		var response map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		// Verify required fields in response
		assert.Contains(t, response, "id")
		assert.Equal(t, "田中太郎", response["name"])
		assert.Equal(t, "tanaka@example.com", response["email"])
		assert.Contains(t, response, "created_at")
		assert.Contains(t, response, "updated_at")
	})

	t.Run("should return 400 for invalid request", func(t *testing.T) {
		// Test missing name
		req := map[string]interface{}{
			"email": "tanaka@example.com",
		}

		body, err := json.Marshal(req)
		require.NoError(t, err)

		w := httptest.NewRecorder()
		httpReq, err := http.NewRequest("POST", "/api/v1/users", bytes.NewBuffer(body))
		require.NoError(t, err)
		httpReq.Header.Set("Content-Type", "application/json")

		router := gin.New()
		router.ServeHTTP(w, httpReq)

		// Expect 400 Bad Request
		assert.Equal(t, http.StatusBadRequest, w.Code)

		// Verify error response structure
		var errorResponse map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &errorResponse)
		require.NoError(t, err)

		assert.Contains(t, errorResponse, "code")
		assert.Contains(t, errorResponse, "message")
	})

	t.Run("should return 409 for duplicate email", func(t *testing.T) {
		req := map[string]interface{}{
			"name":  "田中次郎",
			"email": "existing@example.com",
		}

		body, err := json.Marshal(req)
		require.NoError(t, err)

		w := httptest.NewRecorder()
		httpReq, err := http.NewRequest("POST", "/api/v1/users", bytes.NewBuffer(body))
		require.NoError(t, err)
		httpReq.Header.Set("Content-Type", "application/json")

		router := gin.New()
		router.ServeHTTP(w, httpReq)

		// Expect 409 Conflict
		assert.Equal(t, http.StatusConflict, w.Code)

		var errorResponse map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &errorResponse)
		require.NoError(t, err)

		assert.Contains(t, errorResponse, "code")
		assert.Equal(t, "EMAIL_ALREADY_EXISTS", errorResponse["code"])
	})
}