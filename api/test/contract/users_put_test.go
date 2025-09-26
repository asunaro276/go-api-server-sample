package contract

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestPUTUserContract tests PUT /api/v1/users/{id} endpoint contract
func TestPUTUserContract(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("should update user with valid request", func(t *testing.T) {
		// This test will fail until implementation is complete
		userID := 1
		req := map[string]interface{}{
			"name":  "田中次郎",
			"email": "tanaka.jiro@example.com",
		}

		body, err := json.Marshal(req)
		require.NoError(t, err)

		w := httptest.NewRecorder()
		httpReq, err := http.NewRequest("PUT", fmt.Sprintf("/api/v1/users/%d", userID), bytes.NewBuffer(body))
		require.NoError(t, err)
		httpReq.Header.Set("Content-Type", "application/json")

		// This will fail until router is implemented
		router := gin.New()
		router.ServeHTTP(w, httpReq)

		// Expect 200 OK
		assert.Equal(t, http.StatusOK, w.Code)

		// Verify response structure
		var response map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		// Verify required fields in response
		assert.Contains(t, response, "id")
		assert.Equal(t, "田中次郎", response["name"])
		assert.Equal(t, "tanaka.jiro@example.com", response["email"])
		assert.Contains(t, response, "created_at")
		assert.Contains(t, response, "updated_at")

		// Verify ID matches
		assert.Equal(t, float64(userID), response["id"])
	})

	t.Run("should allow partial updates", func(t *testing.T) {
		userID := 1
		// Update only name
		req := map[string]interface{}{
			"name": "田中三郎",
		}

		body, err := json.Marshal(req)
		require.NoError(t, err)

		w := httptest.NewRecorder()
		httpReq, err := http.NewRequest("PUT", fmt.Sprintf("/api/v1/users/%d", userID), bytes.NewBuffer(body))
		require.NoError(t, err)
		httpReq.Header.Set("Content-Type", "application/json")

		router := gin.New()
		router.ServeHTTP(w, httpReq)

		// Expect 200 OK
		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		// Verify name was updated
		assert.Equal(t, "田中三郎", response["name"])
		// Email should remain unchanged (not nil or empty)
		assert.Contains(t, response, "email")
		assert.NotEmpty(t, response["email"])
	})

	t.Run("should return 404 for non-existent user", func(t *testing.T) {
		userID := 99999
		req := map[string]interface{}{
			"name": "存在しないユーザー",
		}

		body, err := json.Marshal(req)
		require.NoError(t, err)

		w := httptest.NewRecorder()
		httpReq, err := http.NewRequest("PUT", fmt.Sprintf("/api/v1/users/%d", userID), bytes.NewBuffer(body))
		require.NoError(t, err)
		httpReq.Header.Set("Content-Type", "application/json")

		router := gin.New()
		router.ServeHTTP(w, httpReq)

		// Expect 404 Not Found
		assert.Equal(t, http.StatusNotFound, w.Code)

		var errorResponse map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &errorResponse)
		require.NoError(t, err)

		assert.Contains(t, errorResponse, "code")
		assert.Equal(t, "USER_NOT_FOUND", errorResponse["code"])
	})

	t.Run("should return 400 for invalid data", func(t *testing.T) {
		userID := 1
		req := map[string]interface{}{
			"email": "invalid-email", // Invalid email format
		}

		body, err := json.Marshal(req)
		require.NoError(t, err)

		w := httptest.NewRecorder()
		httpReq, err := http.NewRequest("PUT", fmt.Sprintf("/api/v1/users/%d", userID), bytes.NewBuffer(body))
		require.NoError(t, err)
		httpReq.Header.Set("Content-Type", "application/json")

		router := gin.New()
		router.ServeHTTP(w, httpReq)

		// Expect 400 Bad Request
		assert.Equal(t, http.StatusBadRequest, w.Code)

		var errorResponse map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &errorResponse)
		require.NoError(t, err)

		assert.Contains(t, errorResponse, "code")
		assert.Contains(t, errorResponse, "message")
	})

	t.Run("should return 409 for duplicate email", func(t *testing.T) {
		userID := 1
		req := map[string]interface{}{
			"email": "existing@example.com", // Email already in use by another user
		}

		body, err := json.Marshal(req)
		require.NoError(t, err)

		w := httptest.NewRecorder()
		httpReq, err := http.NewRequest("PUT", fmt.Sprintf("/api/v1/users/%d", userID), bytes.NewBuffer(body))
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