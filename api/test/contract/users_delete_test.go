package contract

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestDELETEUserContract tests DELETE /api/v1/users/{id} endpoint contract
func TestDELETEUserContract(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("should delete user with valid ID", func(t *testing.T) {
		// This test will fail until implementation is complete
		userID := 1
		w := httptest.NewRecorder()
		httpReq, err := http.NewRequest("DELETE", fmt.Sprintf("/api/v1/users/%d", userID), nil)
		require.NoError(t, err)

		// This will fail until router is implemented
		router := gin.New()
		router.ServeHTTP(w, httpReq)

		// Expect 204 No Content
		assert.Equal(t, http.StatusNoContent, w.Code)

		// Response body should be empty for 204
		assert.Empty(t, w.Body.String())
	})

	t.Run("should return 404 for non-existent user", func(t *testing.T) {
		userID := 99999
		w := httptest.NewRecorder()
		httpReq, err := http.NewRequest("DELETE", fmt.Sprintf("/api/v1/users/%d", userID), nil)
		require.NoError(t, err)

		router := gin.New()
		router.ServeHTTP(w, httpReq)

		// Expect 404 Not Found
		assert.Equal(t, http.StatusNotFound, w.Code)

		// Verify error response structure
		var errorResponse map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &errorResponse)
		require.NoError(t, err)

		assert.Contains(t, errorResponse, "code")
		assert.Contains(t, errorResponse, "message")
		assert.Equal(t, "USER_NOT_FOUND", errorResponse["code"])
	})

	t.Run("should return 400 for invalid ID format", func(t *testing.T) {
		w := httptest.NewRecorder()
		httpReq, err := http.NewRequest("DELETE", "/api/v1/users/invalid", nil)
		require.NoError(t, err)

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

	t.Run("should return 400 for zero or negative ID", func(t *testing.T) {
		testCases := []int{0, -1, -999}

		for _, userID := range testCases {
			t.Run(fmt.Sprintf("ID_%d", userID), func(t *testing.T) {
				w := httptest.NewRecorder()
				httpReq, err := http.NewRequest("DELETE", fmt.Sprintf("/api/v1/users/%d", userID), nil)
				require.NoError(t, err)

				router := gin.New()
				router.ServeHTTP(w, httpReq)

				// Expect 400 Bad Request for invalid ID
				assert.Equal(t, http.StatusBadRequest, w.Code)
			})
		}
	})

	t.Run("should handle soft delete properly", func(t *testing.T) {
		// This test verifies that deletion is logical (soft delete)
		// After deletion, the user should not appear in GET requests
		userID := 1

		// First, delete the user
		w := httptest.NewRecorder()
		httpReq, err := http.NewRequest("DELETE", fmt.Sprintf("/api/v1/users/%d", userID), nil)
		require.NoError(t, err)

		router := gin.New()
		router.ServeHTTP(w, httpReq)

		// Expect 204 No Content
		assert.Equal(t, http.StatusNoContent, w.Code)

		// Then, try to GET the same user - should return 404
		w2 := httptest.NewRecorder()
		httpReq2, err := http.NewRequest("GET", fmt.Sprintf("/api/v1/users/%d", userID), nil)
		require.NoError(t, err)

		router.ServeHTTP(w2, httpReq2)

		// Should return 404 for soft-deleted user
		assert.Equal(t, http.StatusNotFound, w2.Code)
	})
}