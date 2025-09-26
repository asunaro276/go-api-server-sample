package controller

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/gin-gonic/gin"

	"go-api-server-sample/cmd/api-server/internal/application"
	"go-api-server-sample/internal/domain/entities"
)

// MockUserService is a mock implementation of UserServiceInterface for testing
type MockUserService struct {
	users          map[uint]*entities.User
	nextID         uint
	shouldFailCall bool
	errorToReturn  error
}

func NewMockUserService() *MockUserService {
	return &MockUserService{
		users:  make(map[uint]*entities.User),
		nextID: 1,
	}
}

func (m *MockUserService) CreateUser(ctx context.Context, user *entities.User) (*entities.User, error) {
	if m.shouldFailCall {
		return nil, m.errorToReturn
	}
	user.ID = m.nextID
	m.nextID++
	m.users[user.ID] = user
	return user, nil
}

func (m *MockUserService) GetUserByID(ctx context.Context, id uint) (*entities.User, error) {
	if m.shouldFailCall {
		return nil, m.errorToReturn
	}
	user, exists := m.users[id]
	if !exists {
		return nil, entities.ErrUserNotFound
	}
	return user, nil
}

func (m *MockUserService) UpdateUser(ctx context.Context, id uint, updateData *entities.User) (*entities.User, error) {
	if m.shouldFailCall {
		return nil, m.errorToReturn
	}
	user, exists := m.users[id]
	if !exists {
		return nil, entities.ErrUserNotFound
	}
	user.UpdateFrom(updateData)
	return user, nil
}

func (m *MockUserService) DeleteUser(ctx context.Context, id uint) error {
	if m.shouldFailCall {
		return m.errorToReturn
	}
	if _, exists := m.users[id]; !exists {
		return entities.ErrUserNotFound
	}
	delete(m.users, id)
	return nil
}

func (m *MockUserService) ListUsers(ctx context.Context, req *application.ListUsersRequest) (*application.ListUsersResponse, error) {
	if m.shouldFailCall {
		return nil, m.errorToReturn
	}

	var users []*entities.User
	for _, user := range m.users {
		users = append(users, user)
	}

	// Simple pagination
	start := req.Offset
	if start >= len(users) {
		users = []*entities.User{}
	} else {
		end := start + req.Limit
		if end > len(users) {
			end = len(users)
		}
		users = users[start:end]
	}

	return &application.ListUsersResponse{
		Users:  users,
		Total:  int64(len(m.users)),
		Limit:  req.Limit,
		Offset: req.Offset,
	}, nil
}

func (m *MockUserService) SetShouldFailCall(shouldFail bool, err error) {
	m.shouldFailCall = shouldFail
	m.errorToReturn = err
}

func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	return router
}

func TestUserController_CreateUser(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    string
		expectedStatus int
		shouldFail     bool
		serviceError   error
	}{
		{
			name:           "valid user creation",
			requestBody:    `{"name":"John Doe","email":"john@example.com"}`,
			expectedStatus: http.StatusCreated,
			shouldFail:     false,
		},
		{
			name:           "invalid JSON",
			requestBody:    `{"name":"John Doe","email":}`,
			expectedStatus: http.StatusBadRequest,
			shouldFail:     false,
		},
		{
			name:           "missing required fields",
			requestBody:    `{"name":""}`,
			expectedStatus: http.StatusBadRequest,
			shouldFail:     false,
		},
		{
			name:           "service returns user not found error",
			requestBody:    `{"name":"John Doe","email":"john@example.com"}`,
			expectedStatus: http.StatusNotFound,
			shouldFail:     true,
			serviceError:   entities.ErrUserNotFound,
		},
		{
			name:           "service returns email exists error",
			requestBody:    `{"name":"John Doe","email":"john@example.com"}`,
			expectedStatus: http.StatusConflict,
			shouldFail:     true,
			serviceError:   entities.ErrUserEmailExists,
		},
		{
			name:           "service returns validation error",
			requestBody:    `{"name":"John Doe","email":"john@example.com"}`,
			expectedStatus: http.StatusBadRequest,
			shouldFail:     true,
			serviceError:   entities.ErrInvalidEmail,
		},
		{
			name:           "service returns internal error",
			requestBody:    `{"name":"John Doe","email":"john@example.com"}`,
			expectedStatus: http.StatusInternalServerError,
			shouldFail:     true,
			serviceError:   errors.New("internal error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := NewMockUserService()
			if tt.shouldFail {
				mockService.SetShouldFailCall(true, tt.serviceError)
			}

			controller := NewUserController(mockService)
			router := setupTestRouter()
			router.POST("/users", controller.CreateUser)

			req, _ := http.NewRequest("POST", "/users", bytes.NewBufferString(tt.requestBody))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			// Check response is valid JSON
			var response map[string]interface{}
			if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
				t.Errorf("Response is not valid JSON: %v", err)
			}
		})
	}
}

func TestUserController_GetUsers(t *testing.T) {
	tests := []struct {
		name           string
		queryParams    string
		expectedStatus int
		shouldFail     bool
		serviceError   error
	}{
		{
			name:           "get users without params",
			queryParams:    "",
			expectedStatus: http.StatusOK,
			shouldFail:     false,
		},
		{
			name:           "get users with limit and offset",
			queryParams:    "?limit=5&offset=0",
			expectedStatus: http.StatusOK,
			shouldFail:     false,
		},
		{
			name:           "get users with invalid params (should use defaults)",
			queryParams:    "?limit=abc&offset=xyz",
			expectedStatus: http.StatusOK,
			shouldFail:     false,
		},
		{
			name:           "service returns error",
			queryParams:    "",
			expectedStatus: http.StatusInternalServerError,
			shouldFail:     true,
			serviceError:   errors.New("internal error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := NewMockUserService()
			if tt.shouldFail {
				mockService.SetShouldFailCall(true, tt.serviceError)
			}

			controller := NewUserController(mockService)
			router := setupTestRouter()
			router.GET("/users", controller.GetUsers)

			req, _ := http.NewRequest("GET", "/users"+tt.queryParams, nil)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			// Check response is valid JSON
			var response map[string]interface{}
			if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
				t.Errorf("Response is not valid JSON: %v", err)
			}
		})
	}
}

func TestUserController_GetUserByID(t *testing.T) {
	mockService := NewMockUserService()
	
	// Create a user first
	user := &entities.User{Name: "John Doe", Email: "john@example.com"}
	createdUser, _ := mockService.CreateUser(context.Background(), user)

	tests := []struct {
		name           string
		userID         string
		expectedStatus int
		shouldFail     bool
		serviceError   error
	}{
		{
			name:           "get existing user",
			userID:         strconv.Itoa(int(createdUser.ID)),
			expectedStatus: http.StatusOK,
			shouldFail:     false,
		},
		{
			name:           "get non-existing user",
			userID:         "999",
			expectedStatus: http.StatusNotFound,
			shouldFail:     false, // Service will return ErrUserNotFound
		},
		{
			name:           "invalid user ID",
			userID:         "invalid",
			expectedStatus: http.StatusBadRequest,
			shouldFail:     false,
		},
		{
			name:           "zero user ID",
			userID:         "0",
			expectedStatus: http.StatusBadRequest,
			shouldFail:     false,
		},
		{
			name:           "service returns error",
			userID:         strconv.Itoa(int(createdUser.ID)),
			expectedStatus: http.StatusInternalServerError,
			shouldFail:     true,
			serviceError:   errors.New("internal error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.shouldFail {
				mockService.SetShouldFailCall(true, tt.serviceError)
			} else {
				mockService.SetShouldFailCall(false, nil)
			}

			controller := NewUserController(mockService)
			router := setupTestRouter()
			router.GET("/users/:id", controller.GetUserByID)

			req, _ := http.NewRequest("GET", "/users/"+tt.userID, nil)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			// Check response is valid JSON
			var response map[string]interface{}
			if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
				t.Errorf("Response is not valid JSON: %v", err)
			}
		})
	}
}

func TestUserController_UpdateUser(t *testing.T) {
	mockService := NewMockUserService()
	
	// Create a user first
	user := &entities.User{Name: "John Doe", Email: "john@example.com"}
	createdUser, _ := mockService.CreateUser(context.Background(), user)

	tests := []struct {
		name           string
		userID         string
		requestBody    string
		expectedStatus int
		shouldFail     bool
		serviceError   error
	}{
		{
			name:           "valid user update",
			userID:         strconv.Itoa(int(createdUser.ID)),
			requestBody:    `{"name":"John Updated","email":"john.updated@example.com"}`,
			expectedStatus: http.StatusOK,
			shouldFail:     false,
		},
		{
			name:           "update non-existing user",
			userID:         "999",
			requestBody:    `{"name":"John Updated","email":"john.updated@example.com"}`,
			expectedStatus: http.StatusNotFound,
			shouldFail:     false, // Service will return ErrUserNotFound
		},
		{
			name:           "invalid user ID",
			userID:         "invalid",
			requestBody:    `{"name":"John Updated","email":"john.updated@example.com"}`,
			expectedStatus: http.StatusBadRequest,
			shouldFail:     false,
		},
		{
			name:           "invalid JSON",
			userID:         strconv.Itoa(int(createdUser.ID)),
			requestBody:    `{"name":"John Updated","email":}`,
			expectedStatus: http.StatusBadRequest,
			shouldFail:     false,
		},
		{
			name:           "service returns error",
			userID:         strconv.Itoa(int(createdUser.ID)),
			requestBody:    `{"name":"John Updated","email":"john.updated@example.com"}`,
			expectedStatus: http.StatusInternalServerError,
			shouldFail:     true,
			serviceError:   errors.New("internal error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.shouldFail {
				mockService.SetShouldFailCall(true, tt.serviceError)
			} else {
				mockService.SetShouldFailCall(false, nil)
			}

			controller := NewUserController(mockService)
			router := setupTestRouter()
			router.PUT("/users/:id", controller.UpdateUser)

			req, _ := http.NewRequest("PUT", "/users/"+tt.userID, bytes.NewBufferString(tt.requestBody))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			// Check response is valid JSON
			var response map[string]interface{}
			if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
				t.Errorf("Response is not valid JSON: %v", err)
			}
		})
	}
}

func TestUserController_DeleteUser(t *testing.T) {
	mockService := NewMockUserService()
	
	// Create a user first
	user := &entities.User{Name: "John Doe", Email: "john@example.com"}
	createdUser, _ := mockService.CreateUser(context.Background(), user)

	tests := []struct {
		name           string
		userID         string
		expectedStatus int
		shouldFail     bool
		serviceError   error
	}{
		{
			name:           "valid user deletion",
			userID:         strconv.Itoa(int(createdUser.ID)),
			expectedStatus: http.StatusNoContent,
			shouldFail:     false,
		},
		{
			name:           "delete non-existing user",
			userID:         "999",
			expectedStatus: http.StatusNotFound,
			shouldFail:     false, // Service will return ErrUserNotFound
		},
		{
			name:           "invalid user ID",
			userID:         "invalid",
			expectedStatus: http.StatusBadRequest,
			shouldFail:     false,
		},
		{
			name:           "service returns error",
			userID:         strconv.Itoa(int(createdUser.ID)),
			expectedStatus: http.StatusInternalServerError,
			shouldFail:     true,
			serviceError:   errors.New("internal error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.shouldFail {
				mockService.SetShouldFailCall(true, tt.serviceError)
			} else {
				mockService.SetShouldFailCall(false, nil)
			}

			controller := NewUserController(mockService)
			router := setupTestRouter()
			router.DELETE("/users/:id", controller.DeleteUser)

			req, _ := http.NewRequest("DELETE", "/users/"+tt.userID, nil)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			// For successful deletion, check that body is empty
			if tt.expectedStatus == http.StatusNoContent && w.Body.Len() != 0 {
				t.Errorf("Expected empty body for successful deletion, got %s", w.Body.String())
			}

			// For error cases, check response is valid JSON
			if tt.expectedStatus != http.StatusNoContent {
				var response map[string]interface{}
				if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
					t.Errorf("Response is not valid JSON: %v", err)
				}
			}
		})
	}
}

func TestUserController_parseUserID(t *testing.T) {
	tests := []struct {
		name      string
		paramID   string
		expectErr bool
		expectedID uint
	}{
		{
			name:       "valid ID",
			paramID:    "123",
			expectErr:  false,
			expectedID: 123,
		},
		{
			name:      "zero ID",
			paramID:   "0",
			expectErr: true,
		},
		{
			name:      "negative ID",
			paramID:   "-1",
			expectErr: true,
		},
		{
			name:      "invalid string",
			paramID:   "abc",
			expectErr: true,
		},
		{
			name:      "empty string",
			paramID:   "",
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := NewMockUserService()
			controller := NewUserController(mockService)
			
			// Create a test context
			gin.SetMode(gin.TestMode)
			c, _ := gin.CreateTestContext(httptest.NewRecorder())
			c.Params = []gin.Param{{Key: "id", Value: tt.paramID}}

			id, err := controller.parseUserID(c)

			if (err != nil) != tt.expectErr {
				t.Errorf("parseUserID() error = %v, expectErr %v", err, tt.expectErr)
				return
			}

			if !tt.expectErr && id != tt.expectedID {
				t.Errorf("parseUserID() id = %v, want %v", id, tt.expectedID)
			}
		})
	}
}

// Benchmark tests
func BenchmarkUserController_CreateUser(b *testing.B) {
	mockService := NewMockUserService()
	controller := NewUserController(mockService)
	router := setupTestRouter()
	router.POST("/users", controller.CreateUser)

	requestBody := `{"name":"John Doe","email":"john@example.com"}`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req, _ := http.NewRequest("POST", "/users", bytes.NewBufferString(requestBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
	}
}

func BenchmarkUserController_GetUserByID(b *testing.B) {
	mockService := NewMockUserService()
	controller := NewUserController(mockService)
	router := setupTestRouter()
	router.GET("/users/:id", controller.GetUserByID)

	// Create a user first
	user := &entities.User{Name: "John Doe", Email: "john@example.com"}
	createdUser, _ := mockService.CreateUser(context.Background(), user)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req, _ := http.NewRequest("GET", "/users/"+strconv.Itoa(int(createdUser.ID)), nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
	}
}