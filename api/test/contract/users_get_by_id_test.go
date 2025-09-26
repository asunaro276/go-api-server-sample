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

// TestGETUserByIDContract tests GET /api/v1/users/{id} endpoint contract
func TestGETUserByIDContract(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("should return user by valid ID", func(t *testing.T) {
		// This test will fail until implementation is complete
		userID := 1
		w := httptest.NewRecorder()
		httpReq, err := http.NewRequest("GET", fmt.Sprintf("/api/v1/users/%d", userID), nil)
		require.NoError(t, err)

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
		assert.Contains(t, response, "name")
		assert.Contains(t, response, "email")
		assert.Contains(t, response, "created_at")
		assert.Contains(t, response, "updated_at")

		// Verify ID matches requested ID
		assert.Equal(t, float64(userID), response["id"])

		// Verify email format
		email, ok := response["email"].(string)
		assert.True(t, ok)
		assert.Contains(t, email, "@")

		// Verify name is not empty
		name, ok := response["name"].(string)
		assert.True(t, ok)
		assert.NotEmpty(t, name)
	})

	t.Run("should return 404 for non-existent user", func(t *testing.T) {
		userID := 99999
		w := httptest.NewRecorder()
		httpReq, err := http.NewRequest("GET", fmt.Sprintf("/api/v1/users/%d", userID), nil)
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
		httpReq, err := http.NewRequest("GET", "/api/v1/users/invalid", nil)
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
				httpReq, err := http.NewRequest("GET", fmt.Sprintf("/api/v1/users/%d", userID), nil)
				require.NoError(t, err)

				router := gin.New()
				router.ServeHTTP(w, httpReq)

				// Expect 400 Bad Request for invalid ID
				assert.Equal(t, http.StatusBadRequest, w.Code)
			})
		}
	})
}