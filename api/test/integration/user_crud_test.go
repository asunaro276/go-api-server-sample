// +build integration

package integration

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
	"github.com/stretchr/testify/suite"
)

// UserCRUDTestSuite tests complete CRUD operations with database
type UserCRUDTestSuite struct {
	suite.Suite
	router *gin.Engine
	// db     *gorm.DB // Will be initialized when database is implemented
}

// SetupSuite runs once before all tests in the suite
func (suite *UserCRUDTestSuite) SetupSuite() {
	gin.SetMode(gin.TestMode)
	// Initialize test database connection
	// suite.db = setupTestDatabase() // Will be implemented later

	// Initialize router with all middleware and routes
	// suite.router = setupTestRouter(suite.db) // Will be implemented later
	suite.router = gin.New() // Placeholder until actual router is implemented
}

// TearDownSuite runs once after all tests in the suite
func (suite *UserCRUDTestSuite) TearDownSuite() {
	// Clean up test database
	// cleanupTestDatabase(suite.db) // Will be implemented later
}

// SetupTest runs before each test method
func (suite *UserCRUDTestSuite) SetupTest() {
	// Reset database state before each test
	// resetTestDatabase(suite.db) // Will be implemented later
}

// TestCompleteCRUDFlow tests the complete Create, Read, Update, Delete flow
func (suite *UserCRUDTestSuite) TestCompleteCRUDFlow() {
	// This test will fail until implementation is complete
	t := suite.T()

	// Step 1: CREATE - Create a new user
	newUser := map[string]interface{}{
		"name":  "CRUD Test User",
		"email": "crud@example.com",
	}

	body, err := json.Marshal(newUser)
	require.NoError(t, err)

	w := httptest.NewRecorder()
	req, err := http.NewRequest("POST", "/api/v1/users", bytes.NewBuffer(body))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	suite.router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)

	var createdUser map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &createdUser)
	require.NoError(t, err)

	userID := createdUser["id"]
	assert.NotNil(t, userID)
	originalCreatedAt := createdUser["created_at"]
	originalUpdatedAt := createdUser["updated_at"]

	// Step 2: READ - Get the created user by ID
	w2 := httptest.NewRecorder()
	req2, err := http.NewRequest("GET", fmt.Sprintf("/api/v1/users/%v", userID), nil)
	require.NoError(t, err)

	suite.router.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusOK, w2.Code)

	var retrievedUser map[string]interface{}
	err = json.Unmarshal(w2.Body.Bytes(), &retrievedUser)
	require.NoError(t, err)

	assert.Equal(t, userID, retrievedUser["id"])
	assert.Equal(t, "CRUD Test User", retrievedUser["name"])
	assert.Equal(t, "crud@example.com", retrievedUser["email"])
	assert.Equal(t, originalCreatedAt, retrievedUser["created_at"])

	// Step 3: UPDATE - Update the user's information
	updateData := map[string]interface{}{
		"name":  "Updated CRUD User",
		"email": "updated-crud@example.com",
	}

	updateBody, err := json.Marshal(updateData)
	require.NoError(t, err)

	w3 := httptest.NewRecorder()
	req3, err := http.NewRequest("PUT", fmt.Sprintf("/api/v1/users/%v", userID), bytes.NewBuffer(updateBody))
	require.NoError(t, err)
	req3.Header.Set("Content-Type", "application/json")

	suite.router.ServeHTTP(w3, req3)
	assert.Equal(t, http.StatusOK, w3.Code)

	var updatedUser map[string]interface{}
	err = json.Unmarshal(w3.Body.Bytes(), &updatedUser)
	require.NoError(t, err)

	assert.Equal(t, userID, updatedUser["id"])
	assert.Equal(t, "Updated CRUD User", updatedUser["name"])
	assert.Equal(t, "updated-crud@example.com", updatedUser["email"])
	assert.Equal(t, originalCreatedAt, updatedUser["created_at"]) // Should not change
	assert.NotEqual(t, originalUpdatedAt, updatedUser["updated_at"]) // Should be updated

	// Step 4: READ after UPDATE - Verify changes persisted
	w4 := httptest.NewRecorder()
	req4, err := http.NewRequest("GET", fmt.Sprintf("/api/v1/users/%v", userID), nil)
	require.NoError(t, err)

	suite.router.ServeHTTP(w4, req4)
	assert.Equal(t, http.StatusOK, w4.Code)

	var verifyUser map[string]interface{}
	err = json.Unmarshal(w4.Body.Bytes(), &verifyUser)
	require.NoError(t, err)

	assert.Equal(t, "Updated CRUD User", verifyUser["name"])
	assert.Equal(t, "updated-crud@example.com", verifyUser["email"])

	// Step 5: DELETE - Soft delete the user
	w5 := httptest.NewRecorder()
	req5, err := http.NewRequest("DELETE", fmt.Sprintf("/api/v1/users/%v", userID), nil)
	require.NoError(t, err)

	suite.router.ServeHTTP(w5, req5)
	assert.Equal(t, http.StatusNoContent, w5.Code)
	assert.Empty(t, w5.Body.String())

	// Step 6: READ after DELETE - Should return 404
	w6 := httptest.NewRecorder()
	req6, err := http.NewRequest("GET", fmt.Sprintf("/api/v1/users/%v", userID), nil)
	require.NoError(t, err)

	suite.router.ServeHTTP(w6, req6)
	assert.Equal(t, http.StatusNotFound, w6.Code)

	var notFoundError map[string]interface{}
	err = json.Unmarshal(w6.Body.Bytes(), &notFoundError)
	require.NoError(t, err)
	assert.Equal(t, "USER_NOT_FOUND", notFoundError["code"])
}

// TestPartialUpdateOperations tests partial update scenarios
func (suite *UserCRUDTestSuite) TestPartialUpdateOperations() {
	t := suite.T()

	// Create initial user
	newUser := map[string]interface{}{
		"name":  "Partial Update User",
		"email": "partial@example.com",
	}

	body, err := json.Marshal(newUser)
	require.NoError(t, err)

	w := httptest.NewRecorder()
	req, err := http.NewRequest("POST", "/api/v1/users", bytes.NewBuffer(body))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	suite.router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)

	var createdUser map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &createdUser)
	require.NoError(t, err)
	userID := createdUser["id"]

	// Test 1: Update only name
	updateNameOnly := map[string]interface{}{
		"name": "Updated Name Only",
	}

	updateBody, err := json.Marshal(updateNameOnly)
	require.NoError(t, err)

	w2 := httptest.NewRecorder()
	req2, err := http.NewRequest("PUT", fmt.Sprintf("/api/v1/users/%v", userID), bytes.NewBuffer(updateBody))
	require.NoError(t, err)
	req2.Header.Set("Content-Type", "application/json")

	suite.router.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusOK, w2.Code)

	var updatedUser map[string]interface{}
	err = json.Unmarshal(w2.Body.Bytes(), &updatedUser)
	require.NoError(t, err)

	assert.Equal(t, "Updated Name Only", updatedUser["name"])
	assert.Equal(t, "partial@example.com", updatedUser["email"]) // Should remain unchanged

	// Test 2: Update only email
	updateEmailOnly := map[string]interface{}{
		"email": "updated-partial@example.com",
	}

	updateBody2, err := json.Marshal(updateEmailOnly)
	require.NoError(t, err)

	w3 := httptest.NewRecorder()
	req3, err := http.NewRequest("PUT", fmt.Sprintf("/api/v1/users/%v", userID), bytes.NewBuffer(updateBody2))
	require.NoError(t, err)
	req3.Header.Set("Content-Type", "application/json")

	suite.router.ServeHTTP(w3, req3)
	assert.Equal(t, http.StatusOK, w3.Code)

	var updatedUser2 map[string]interface{}
	err = json.Unmarshal(w3.Body.Bytes(), &updatedUser2)
	require.NoError(t, err)

	assert.Equal(t, "Updated Name Only", updatedUser2["name"]) // Should remain from previous update
	assert.Equal(t, "updated-partial@example.com", updatedUser2["email"])
}

// TestPaginationInList tests pagination functionality in list operations
func (suite *UserCRUDTestSuite) TestPaginationInList() {
	t := suite.T()

	// Create multiple users for pagination testing
	userCount := 25
	createdUsers := make([]interface{}, 0, userCount)

	for i := 0; i < userCount; i++ {
		user := map[string]interface{}{
			"name":  fmt.Sprintf("User %d", i+1),
			"email": fmt.Sprintf("user%d@example.com", i+1),
		}

		body, err := json.Marshal(user)
		require.NoError(t, err)

		w := httptest.NewRecorder()
		req, err := http.NewRequest("POST", "/api/v1/users", bytes.NewBuffer(body))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		suite.router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusCreated, w.Code)

		var createdUser map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &createdUser)
		require.NoError(t, err)
		createdUsers = append(createdUsers, createdUser["id"])
	}

	// Test pagination with different limits and offsets
	testCases := []struct {
		name          string
		limit         int
		offset        int
		expectedCount int
	}{
		{"first page", 10, 0, 10},
		{"second page", 10, 10, 10},
		{"third page", 10, 20, 5}, // Only 5 remaining
		{"large limit", 100, 0, 25}, // All users
		{"small limit", 5, 0, 5},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req, err := http.NewRequest("GET", fmt.Sprintf("/api/v1/users?limit=%d&offset=%d", tc.limit, tc.offset), nil)
			require.NoError(t, err)

			suite.router.ServeHTTP(w, req)
			assert.Equal(t, http.StatusOK, w.Code)

			var response map[string]interface{}
			err = json.Unmarshal(w.Body.Bytes(), &response)
			require.NoError(t, err)

			users := response["users"].([]interface{})
			assert.Len(t, users, tc.expectedCount)
			assert.Equal(t, float64(tc.limit), response["limit"])
			assert.Equal(t, float64(tc.offset), response["offset"])
			assert.GreaterOrEqual(t, response["total"].(float64), float64(userCount))
		})
	}
}

// TestDatabaseConstraints tests database-level constraints
func (suite *UserCRUDTestSuite) TestDatabaseConstraints() {
	t := suite.T()

	// Test email uniqueness at database level
	user1 := map[string]interface{}{
		"name":  "User One",
		"email": "unique@example.com",
	}

	body1, err := json.Marshal(user1)
	require.NoError(t, err)

	w1 := httptest.NewRecorder()
	req1, err := http.NewRequest("POST", "/api/v1/users", bytes.NewBuffer(body1))
	require.NoError(t, err)
	req1.Header.Set("Content-Type", "application/json")

	suite.router.ServeHTTP(w1, req1)
	assert.Equal(t, http.StatusCreated, w1.Code)

	// Try to create another user with same email
	user2 := map[string]interface{}{
		"name":  "User Two",
		"email": "unique@example.com", // Same email
	}

	body2, err := json.Marshal(user2)
	require.NoError(t, err)

	w2 := httptest.NewRecorder()
	req2, err := http.NewRequest("POST", "/api/v1/users", bytes.NewBuffer(body2))
	require.NoError(t, err)
	req2.Header.Set("Content-Type", "application/json")

	suite.router.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusConflict, w2.Code)

	var errorResponse map[string]interface{}
	err = json.Unmarshal(w2.Body.Bytes(), &errorResponse)
	require.NoError(t, err)
	assert.Equal(t, "EMAIL_ALREADY_EXISTS", errorResponse["code"])
}

// TestSoftDeleteBehavior tests soft delete functionality
func (suite *UserCRUDTestSuite) TestSoftDeleteBehavior() {
	t := suite.T()

	// Create a user
	newUser := map[string]interface{}{
		"name":  "Soft Delete Test",
		"email": "softdelete@example.com",
	}

	body, err := json.Marshal(newUser)
	require.NoError(t, err)

	w := httptest.NewRecorder()
	req, err := http.NewRequest("POST", "/api/v1/users", bytes.NewBuffer(body))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	suite.router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)

	var createdUser map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &createdUser)
	require.NoError(t, err)
	userID := createdUser["id"]

	// Delete the user (soft delete)
	w2 := httptest.NewRecorder()
	req2, err := http.NewRequest("DELETE", fmt.Sprintf("/api/v1/users/%v", userID), nil)
	require.NoError(t, err)

	suite.router.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusNoContent, w2.Code)

	// Verify user no longer appears in list
	w3 := httptest.NewRecorder()
	req3, err := http.NewRequest("GET", "/api/v1/users", nil)
	require.NoError(t, err)

	suite.router.ServeHTTP(w3, req3)
	assert.Equal(t, http.StatusOK, w3.Code)

	var listResponse map[string]interface{}
	err = json.Unmarshal(w3.Body.Bytes(), &listResponse)
	require.NoError(t, err)

	users := listResponse["users"].([]interface{})
	for _, user := range users {
		userMap := user.(map[string]interface{})
		assert.NotEqual(t, userID, userMap["id"], "Deleted user should not appear in list")
	}

	// Verify user cannot be retrieved by ID
	w4 := httptest.NewRecorder()
	req4, err := http.NewRequest("GET", fmt.Sprintf("/api/v1/users/%v", userID), nil)
	require.NoError(t, err)

	suite.router.ServeHTTP(w4, req4)
	assert.Equal(t, http.StatusNotFound, w4.Code)

	// Verify email can be reused after soft delete
	newUserSameEmail := map[string]interface{}{
		"name":  "New User Same Email",
		"email": "softdelete@example.com", // Same email as deleted user
	}

	body5, err := json.Marshal(newUserSameEmail)
	require.NoError(t, err)

	w5 := httptest.NewRecorder()
	req5, err := http.NewRequest("POST", "/api/v1/users", bytes.NewBuffer(body5))
	require.NoError(t, err)
	req5.Header.Set("Content-Type", "application/json")

	suite.router.ServeHTTP(w5, req5)
	assert.Equal(t, http.StatusCreated, w5.Code) // Should succeed
}

// Run the test suite
func TestUserCRUDSuite(t *testing.T) {
	suite.Run(t, new(UserCRUDTestSuite))
}