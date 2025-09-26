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

// UserRegistrationTestSuite tests the complete user registration flow
type UserRegistrationTestSuite struct {
	suite.Suite
	router *gin.Engine
	// db     *gorm.DB // Will be initialized when database is implemented
}

// SetupSuite runs once before all tests in the suite
func (suite *UserRegistrationTestSuite) SetupSuite() {
	gin.SetMode(gin.TestMode)
	// Initialize test database connection
	// suite.db = setupTestDatabase() // Will be implemented later

	// Initialize router with all middleware and routes
	// suite.router = setupTestRouter(suite.db) // Will be implemented later
	suite.router = gin.New() // Placeholder until actual router is implemented
}

// TearDownSuite runs once after all tests in the suite
func (suite *UserRegistrationTestSuite) TearDownSuite() {
	// Clean up test database
	// cleanupTestDatabase(suite.db) // Will be implemented later
}

// SetupTest runs before each test method
func (suite *UserRegistrationTestSuite) SetupTest() {
	// Reset database state before each test
	// resetTestDatabase(suite.db) // Will be implemented later
}

// TestCompleteUserRegistrationFlow tests the entire user registration workflow
func (suite *UserRegistrationTestSuite) TestCompleteUserRegistrationFlow() {
	// This test will fail until implementation is complete
	t := suite.T()

	// Step 1: Register new user
	newUser := map[string]interface{}{
		"name":  "統合テストユーザー",
		"email": "integration@example.com",
	}

	body, err := json.Marshal(newUser)
	require.NoError(t, err)

	w := httptest.NewRecorder()
	req, err := http.NewRequest("POST", "/api/v1/users", bytes.NewBuffer(body))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	suite.router.ServeHTTP(w, req)

	// Should succeed with 201 Created
	assert.Equal(t, http.StatusCreated, w.Code)

	// Parse response to get user ID
	var createdUser map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &createdUser)
	require.NoError(t, err)

	userID := createdUser["id"]
	assert.NotNil(t, userID)
	assert.Equal(t, "統合テストユーザー", createdUser["name"])
	assert.Equal(t, "integration@example.com", createdUser["email"])

	// Step 2: Verify user appears in users list
	w2 := httptest.NewRecorder()
	req2, err := http.NewRequest("GET", "/api/v1/users", nil)
	require.NoError(t, err)

	suite.router.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusOK, w2.Code)

	var listResponse map[string]interface{}
	err = json.Unmarshal(w2.Body.Bytes(), &listResponse)
	require.NoError(t, err)

	users := listResponse["users"].([]interface{})
	assert.Greater(t, len(users), 0)

	// Find our created user in the list
	found := false
	for _, user := range users {
		userMap := user.(map[string]interface{})
		if userMap["email"] == "integration@example.com" {
			found = true
			break
		}
	}
	assert.True(t, found, "Created user should appear in users list")

	// Step 3: Verify user can be retrieved by ID
	w3 := httptest.NewRecorder()
	req3, err := http.NewRequest("GET", fmt.Sprintf("/api/v1/users/%v", userID), nil)
	require.NoError(t, err)

	suite.router.ServeHTTP(w3, req3)
	assert.Equal(t, http.StatusOK, w3.Code)

	var retrievedUser map[string]interface{}
	err = json.Unmarshal(w3.Body.Bytes(), &retrievedUser)
	require.NoError(t, err)

	assert.Equal(t, userID, retrievedUser["id"])
	assert.Equal(t, "統合テストユーザー", retrievedUser["name"])
	assert.Equal(t, "integration@example.com", retrievedUser["email"])
}

// TestUserRegistrationValidation tests validation during registration
func (suite *UserRegistrationTestSuite) TestUserRegistrationValidation() {
	t := suite.T()

	testCases := []struct {
		name           string
		requestBody    map[string]interface{}
		expectedStatus int
		expectedError  string
	}{
		{
			name: "missing name",
			requestBody: map[string]interface{}{
				"email": "test@example.com",
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "VALIDATION_ERROR",
		},
		{
			name: "missing email",
			requestBody: map[string]interface{}{
				"name": "Test User",
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "VALIDATION_ERROR",
		},
		{
			name: "invalid email format",
			requestBody: map[string]interface{}{
				"name":  "Test User",
				"email": "invalid-email",
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "VALIDATION_ERROR",
		},
		{
			name: "empty name",
			requestBody: map[string]interface{}{
				"name":  "",
				"email": "test@example.com",
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "VALIDATION_ERROR",
		},
		{
			name: "name too long",
			requestBody: map[string]interface{}{
				"name":  string(make([]byte, 101)), // 101 characters, exceeds 100 limit
				"email": "test@example.com",
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "VALIDATION_ERROR",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			body, err := json.Marshal(tc.requestBody)
			require.NoError(t, err)

			w := httptest.NewRecorder()
			req, err := http.NewRequest("POST", "/api/v1/users", bytes.NewBuffer(body))
			require.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")

			suite.router.ServeHTTP(w, req)

			assert.Equal(t, tc.expectedStatus, w.Code)

			if tc.expectedStatus != http.StatusCreated {
				var errorResponse map[string]interface{}
				err = json.Unmarshal(w.Body.Bytes(), &errorResponse)
				require.NoError(t, err)
				assert.Equal(t, tc.expectedError, errorResponse["code"])
			}
		})
	}
}

// TestDuplicateEmailRegistration tests that duplicate emails are rejected
func (suite *UserRegistrationTestSuite) TestDuplicateEmailRegistration() {
	t := suite.T()

	// First, create a user
	firstUser := map[string]interface{}{
		"name":  "First User",
		"email": "duplicate@example.com",
	}

	body, err := json.Marshal(firstUser)
	require.NoError(t, err)

	w := httptest.NewRecorder()
	req, err := http.NewRequest("POST", "/api/v1/users", bytes.NewBuffer(body))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	suite.router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)

	// Now try to create another user with the same email
	secondUser := map[string]interface{}{
		"name":  "Second User",
		"email": "duplicate@example.com", // Same email
	}

	body2, err := json.Marshal(secondUser)
	require.NoError(t, err)

	w2 := httptest.NewRecorder()
	req2, err := http.NewRequest("POST", "/api/v1/users", bytes.NewBuffer(body2))
	require.NoError(t, err)
	req2.Header.Set("Content-Type", "application/json")

	suite.router.ServeHTTP(w2, req2)

	// Should fail with 409 Conflict
	assert.Equal(t, http.StatusConflict, w2.Code)

	var errorResponse map[string]interface{}
	err = json.Unmarshal(w2.Body.Bytes(), &errorResponse)
	require.NoError(t, err)
	assert.Equal(t, "EMAIL_ALREADY_EXISTS", errorResponse["code"])
}

// Run the test suite
func TestUserRegistrationSuite(t *testing.T) {
	suite.Run(t, new(UserRegistrationTestSuite))
}