package contract

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestGETUsersListContract tests GET /api/v1/users endpoint contract
func TestGETUsersListContract(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("should return users list with pagination", func(t *testing.T) {
		// This test will fail until implementation is complete
		w := httptest.NewRecorder()
		httpReq, err := http.NewRequest("GET", "/api/v1/users?limit=10&offset=0", nil)
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
		assert.Contains(t, response, "users")
		assert.Contains(t, response, "total")
		assert.Contains(t, response, "limit")
		assert.Contains(t, response, "offset")

		// Verify users array
		users, ok := response["users"].([]interface{})
		assert.True(t, ok)
		assert.LessOrEqual(t, len(users), 10) // Should respect limit

		// Verify each user structure if users exist
		if len(users) > 0 {
			user := users[0].(map[string]interface{})
			assert.Contains(t, user, "id")
			assert.Contains(t, user, "name")
			assert.Contains(t, user, "email")
			assert.Contains(t, user, "created_at")
			assert.Contains(t, user, "updated_at")
		}
	})

	t.Run("should handle pagination parameters", func(t *testing.T) {
		w := httptest.NewRecorder()
		httpReq, err := http.NewRequest("GET", "/api/v1/users?limit=5&offset=10", nil)
		require.NoError(t, err)

		router := gin.New()
		router.ServeHTTP(w, httpReq)

		// Expect 200 OK
		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		// Verify pagination parameters are reflected
		assert.Equal(t, float64(5), response["limit"])
		assert.Equal(t, float64(10), response["offset"])
	})

	t.Run("should use default pagination when not specified", func(t *testing.T) {
		w := httptest.NewRecorder()
		httpReq, err := http.NewRequest("GET", "/api/v1/users", nil)
		require.NoError(t, err)

		router := gin.New()
		router.ServeHTTP(w, httpReq)

		// Expect 200 OK
		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		// Verify default pagination
		assert.Equal(t, float64(10), response["limit"]) // Default limit
		assert.Equal(t, float64(0), response["offset"])  // Default offset
	})

	t.Run("should validate pagination limits", func(t *testing.T) {
		w := httptest.NewRecorder()
		httpReq, err := http.NewRequest("GET", "/api/v1/users?limit=150&offset=0", nil)
		require.NoError(t, err)

		router := gin.New()
		router.ServeHTTP(w, httpReq)

		// Should still return 200 but with capped limit
		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		// Limit should be capped at 100
		assert.LessOrEqual(t, response["limit"], float64(100))
	})
}