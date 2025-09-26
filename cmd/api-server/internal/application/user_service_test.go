package application

import (
	"context"
	"errors"
	"testing"

	"go-api-server-sample/internal/domain/entities"
	"go-api-server-sample/internal/domain/services"
)

// MockUserRepository is a mock implementation of UserRepository for testing
type MockUserRepository struct {
	users           map[uint]*entities.User
	nextID          uint
	existsByEmail   map[string]bool
	shouldFailQuery bool
}

func NewMockUserRepository() *MockUserRepository {
	return &MockUserRepository{
		users:         make(map[uint]*entities.User),
		nextID:        1,
		existsByEmail: make(map[string]bool),
	}
}

func (m *MockUserRepository) Create(ctx context.Context, user *entities.User) error {
	if m.shouldFailQuery {
		return errors.New("database error")
	}
	if m.existsByEmail[user.Email] {
		return entities.ErrUserEmailExists
	}
	user.ID = m.nextID
	m.nextID++
	m.users[user.ID] = user
	m.existsByEmail[user.Email] = true
	return nil
}

func (m *MockUserRepository) GetByID(ctx context.Context, id uint) (*entities.User, error) {
	if m.shouldFailQuery {
		return nil, errors.New("database error")
	}
	user, exists := m.users[id]
	if !exists {
		return nil, entities.ErrUserNotFound
	}
	return user, nil
}

func (m *MockUserRepository) GetByEmail(ctx context.Context, email string) (*entities.User, error) {
	if m.shouldFailQuery {
		return nil, errors.New("database error")
	}
	for _, user := range m.users {
		if user.Email == email {
			return user, nil
		}
	}
	return nil, entities.ErrUserNotFound
}

func (m *MockUserRepository) Update(ctx context.Context, user *entities.User) error {
	if m.shouldFailQuery {
		return errors.New("database error")
	}
	if _, exists := m.users[user.ID]; !exists {
		return entities.ErrUserNotFound
	}
	
	// Check for email conflicts with other users
	for id, existingUser := range m.users {
		if id != user.ID && existingUser.Email == user.Email {
			return entities.ErrUserEmailExists
		}
	}
	
	m.users[user.ID] = user
	return nil
}

func (m *MockUserRepository) Delete(ctx context.Context, id uint) error {
	if m.shouldFailQuery {
		return errors.New("database error")
	}
	if _, exists := m.users[id]; !exists {
		return entities.ErrUserNotFound
	}
	delete(m.users, id)
	return nil
}

func (m *MockUserRepository) List(ctx context.Context, limit, offset int) ([]*entities.User, error) {
	if m.shouldFailQuery {
		return nil, errors.New("database error")
	}
	var users []*entities.User
	for _, user := range m.users {
		users = append(users, user)
	}
	
	// Simple pagination simulation
	start := offset
	if start >= len(users) {
		return []*entities.User{}, nil
	}
	
	end := start + limit
	if end > len(users) {
		end = len(users)
	}
	
	return users[start:end], nil
}

func (m *MockUserRepository) Count(ctx context.Context) (int64, error) {
	if m.shouldFailQuery {
		return 0, errors.New("database error")
	}
	return int64(len(m.users)), nil
}

func (m *MockUserRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	if m.shouldFailQuery {
		return false, errors.New("database error")
	}
	return m.existsByEmail[email], nil
}

func (m *MockUserRepository) ExistsByEmailExcludingID(ctx context.Context, email string, excludeID uint) (bool, error) {
	if m.shouldFailQuery {
		return false, errors.New("database error")
	}
	for id, user := range m.users {
		if id != excludeID && user.Email == email {
			return true, nil
		}
	}
	return false, nil
}

func (m *MockUserRepository) SetShouldFailQuery(shouldFail bool) {
	m.shouldFailQuery = shouldFail
}

func TestUserService_CreateUser(t *testing.T) {
	tests := []struct {
		name    string
		user    *entities.User
		wantErr bool
		errType error
	}{
		{
			name: "valid user creation",
			user: &entities.User{
				Name:  "John Doe",
				Email: "john@example.com",
			},
			wantErr: false,
		},
		{
			name: "invalid user - empty name",
			user: &entities.User{
				Name:  "",
				Email: "john@example.com",
			},
			wantErr: true,
			errType: entities.ErrUserNameRequired,
		},
		{
			name: "invalid user - invalid email",
			user: &entities.User{
				Name:  "John Doe",
				Email: "invalid-email",
			},
			wantErr: true,
			errType: entities.ErrInvalidEmail,
		},
		{
			name: "duplicate email",
			user: &entities.User{
				Name:  "Jane Doe",
				Email: "john@example.com", // Will be duplicate after first test
			},
			wantErr: true,
			errType: entities.ErrUserEmailExists,
		},
	}

	mockRepo := NewMockUserRepository()
	domainService := services.NewUserDomainService(mockRepo)
	userService := NewUserService(mockRepo, domainService)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := userService.CreateUser(context.Background(), tt.user)
			
			if (err != nil) != tt.wantErr {
				t.Errorf("UserService.CreateUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				if tt.errType != nil && !errors.Is(err, tt.errType) {
					t.Errorf("UserService.CreateUser() error = %v, want %v", err, tt.errType)
				}
				if result != nil {
					t.Errorf("UserService.CreateUser() result should be nil on error, got %v", result)
				}
			} else {
				if result == nil {
					t.Errorf("UserService.CreateUser() result should not be nil on success")
				}
				if result.ID == 0 {
					t.Errorf("UserService.CreateUser() result ID should be set")
				}
			}
		})
	}
}

func TestUserService_GetUserByID(t *testing.T) {
	mockRepo := NewMockUserRepository()
	domainService := services.NewUserDomainService(mockRepo)
	userService := NewUserService(mockRepo, domainService)

	// Create a user first
	user := &entities.User{
		Name:  "John Doe",
		Email: "john@example.com",
	}
	createdUser, _ := userService.CreateUser(context.Background(), user)

	tests := []struct {
		name    string
		id      uint
		wantErr bool
		errType error
	}{
		{
			name:    "existing user",
			id:      createdUser.ID,
			wantErr: false,
		},
		{
			name:    "non-existing user",
			id:      999,
			wantErr: true,
			errType: entities.ErrUserNotFound,
		},
		{
			name:    "invalid ID",
			id:      0,
			wantErr: true,
			errType: entities.ErrInvalidUserID,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := userService.GetUserByID(context.Background(), tt.id)
			
			if (err != nil) != tt.wantErr {
				t.Errorf("UserService.GetUserByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				if tt.errType != nil && !errors.Is(err, tt.errType) {
					t.Errorf("UserService.GetUserByID() error = %v, want %v", err, tt.errType)
				}
			} else {
				if result == nil {
					t.Errorf("UserService.GetUserByID() result should not be nil")
				}
				if result.ID != tt.id {
					t.Errorf("UserService.GetUserByID() result ID = %v, want %v", result.ID, tt.id)
				}
			}
		})
	}
}

func TestUserService_UpdateUser(t *testing.T) {
	mockRepo := NewMockUserRepository()
	domainService := services.NewUserDomainService(mockRepo)
	userService := NewUserService(mockRepo, domainService)

	// Create a user first
	user := &entities.User{
		Name:  "John Doe",
		Email: "john@example.com",
	}
	createdUser, _ := userService.CreateUser(context.Background(), user)

	// Create another user for conflict testing
	user2 := &entities.User{
		Name:  "Jane Doe",
		Email: "jane@example.com",
	}
	_, _ = userService.CreateUser(context.Background(), user2)

	tests := []struct {
		name       string
		id         uint
		updateData *entities.User
		wantErr    bool
		errType    error
	}{
		{
			name: "valid update",
			id:   createdUser.ID,
			updateData: &entities.User{
				Name:  "John Updated",
				Email: "john.updated@example.com",
			},
			wantErr: false,
		},
		{
			name: "update with existing email",
			id:   createdUser.ID,
			updateData: &entities.User{
				Name:  "John Doe",
				Email: "jane@example.com", // Conflicts with user2
			},
			wantErr: true,
			errType: entities.ErrUserEmailExists,
		},
		{
			name: "invalid user ID",
			id:   0,
			updateData: &entities.User{
				Name:  "John Doe",
				Email: "john@example.com",
			},
			wantErr: true,
			errType: entities.ErrInvalidUserID,
		},
		{
			name: "non-existing user",
			id:   999,
			updateData: &entities.User{
				Name:  "John Doe",
				Email: "john@example.com",
			},
			wantErr: true,
			errType: entities.ErrUserNotFound,
		},
		{
			name: "invalid email format",
			id:   createdUser.ID,
			updateData: &entities.User{
				Name:  "John Doe",
				Email: "invalid-email",
			},
			wantErr: true,
			errType: entities.ErrInvalidEmail,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := userService.UpdateUser(context.Background(), tt.id, tt.updateData)
			
			if (err != nil) != tt.wantErr {
				t.Errorf("UserService.UpdateUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				if tt.errType != nil && !errors.Is(err, tt.errType) {
					t.Errorf("UserService.UpdateUser() error = %v, want %v", err, tt.errType)
				}
			} else {
				if result == nil {
					t.Errorf("UserService.UpdateUser() result should not be nil")
				}
				if result.ID != tt.id {
					t.Errorf("UserService.UpdateUser() result ID = %v, want %v", result.ID, tt.id)
				}
			}
		})
	}
}

func TestUserService_DeleteUser(t *testing.T) {
	mockRepo := NewMockUserRepository()
	domainService := services.NewUserDomainService(mockRepo)
	userService := NewUserService(mockRepo, domainService)

	// Create a user first
	user := &entities.User{
		Name:  "John Doe",
		Email: "john@example.com",
	}
	createdUser, _ := userService.CreateUser(context.Background(), user)

	tests := []struct {
		name    string
		id      uint
		wantErr bool
		errType error
	}{
		{
			name:    "valid deletion",
			id:      createdUser.ID,
			wantErr: false,
		},
		{
			name:    "invalid user ID",
			id:      0,
			wantErr: true,
			errType: entities.ErrInvalidUserID,
		},
		{
			name:    "non-existing user",
			id:      999,
			wantErr: true,
			errType: entities.ErrUserNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := userService.DeleteUser(context.Background(), tt.id)
			
			if (err != nil) != tt.wantErr {
				t.Errorf("UserService.DeleteUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && tt.errType != nil && !errors.Is(err, tt.errType) {
				t.Errorf("UserService.DeleteUser() error = %v, want %v", err, tt.errType)
			}
		})
	}
}

func TestUserService_ListUsers(t *testing.T) {
	mockRepo := NewMockUserRepository()
	domainService := services.NewUserDomainService(mockRepo)
	userService := NewUserService(mockRepo, domainService)

	// Create some users
	users := []*entities.User{
		{Name: "User 1", Email: "user1@example.com"},
		{Name: "User 2", Email: "user2@example.com"},
		{Name: "User 3", Email: "user3@example.com"},
	}

	for _, user := range users {
		_, _ = userService.CreateUser(context.Background(), user)
	}

	tests := []struct {
		name    string
		req     *ListUsersRequest
		wantLen int
		wantErr bool
	}{
		{
			name: "default pagination",
			req: &ListUsersRequest{
				Limit:  10,
				Offset: 0,
			},
			wantLen: 3,
			wantErr: false,
		},
		{
			name: "limit of 2",
			req: &ListUsersRequest{
				Limit:  2,
				Offset: 0,
			},
			wantLen: 2,
			wantErr: false,
		},
		{
			name: "offset of 2",
			req: &ListUsersRequest{
				Limit:  10,
				Offset: 2,
			},
			wantLen: 1,
			wantErr: false,
		},
		{
			name: "invalid limit (should be set to default)",
			req: &ListUsersRequest{
				Limit:  0,
				Offset: 0,
			},
			wantLen: 3,
			wantErr: false,
		},
		{
			name: "limit too high (should be capped)",
			req: &ListUsersRequest{
				Limit:  200,
				Offset: 0,
			},
			wantLen: 3,
			wantErr: false,
		},
		{
			name: "negative offset (should be set to 0)",
			req: &ListUsersRequest{
				Limit:  10,
				Offset: -5,
			},
			wantLen: 3,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := userService.ListUsers(context.Background(), tt.req)
			
			if (err != nil) != tt.wantErr {
				t.Errorf("UserService.ListUsers() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if result == nil {
					t.Errorf("UserService.ListUsers() result should not be nil")
					return
				}
				if len(result.Users) != tt.wantLen {
					t.Errorf("UserService.ListUsers() result length = %v, want %v", len(result.Users), tt.wantLen)
				}
				if result.Total != 3 {
					t.Errorf("UserService.ListUsers() total = %v, want %v", result.Total, 3)
				}
			}
		})
	}
}

// Benchmark tests
func BenchmarkUserService_CreateUser(b *testing.B) {
	mockRepo := NewMockUserRepository()
	domainService := services.NewUserDomainService(mockRepo)
	userService := NewUserService(mockRepo, domainService)

	user := &entities.User{
		Name:  "John Doe",
		Email: "john@example.com",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		user.Email = "john" + string(rune(i)) + "@example.com" // Unique email for each iteration
		_, _ = userService.CreateUser(context.Background(), user)
	}
}

func BenchmarkUserService_GetUserByID(b *testing.B) {
	mockRepo := NewMockUserRepository()
	domainService := services.NewUserDomainService(mockRepo)
	userService := NewUserService(mockRepo, domainService)

	// Create a user first
	user := &entities.User{
		Name:  "John Doe",
		Email: "john@example.com",
	}
	createdUser, _ := userService.CreateUser(context.Background(), user)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = userService.GetUserByID(context.Background(), createdUser.ID)
	}
}