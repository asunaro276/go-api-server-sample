package application

import (
	"context"
	"fmt"

	"go-api-server-sample/internal/domain/entities"
	"go-api-server-sample/internal/domain/repositories"
	"go-api-server-sample/internal/domain/services"
)

// UserService provides application-level user operations
type UserService struct {
	userRepo        repositories.UserRepository
	userDomainSvc   *services.UserDomainService
}

// NewUserService creates a new UserService
func NewUserService(userRepo repositories.UserRepository, userDomainSvc *services.UserDomainService) *UserService {
	return &UserService{
		userRepo:      userRepo,
		userDomainSvc: userDomainSvc,
	}
}

// CreateUser creates a new user
func (s *UserService) CreateUser(ctx context.Context, user *entities.User) (*entities.User, error) {
	// Validate user for creation
	if err := s.userDomainSvc.ValidateUserForCreation(ctx, user); err != nil {
		return nil, err
	}

	// Create user
	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}

// GetUserByID retrieves a user by ID
func (s *UserService) GetUserByID(ctx context.Context, id uint) (*entities.User, error) {
	if err := s.userDomainSvc.ValidateUserID(id); err != nil {
		return nil, err
	}

	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

// UpdateUser updates an existing user
func (s *UserService) UpdateUser(ctx context.Context, id uint, updateData *entities.User) (*entities.User, error) {
	if err := s.userDomainSvc.ValidateUserID(id); err != nil {
		return nil, err
	}

	// Get existing user
	existingUser, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get existing user: %w", err)
	}

	// Update only provided fields
	existingUser.UpdateFrom(updateData)

	// Validate updated user
	if err := s.userDomainSvc.ValidateUserForUpdate(ctx, existingUser); err != nil {
		return nil, err
	}

	// Update user
	if err := s.userRepo.Update(ctx, existingUser); err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	return existingUser, nil
}

// DeleteUser soft deletes a user
func (s *UserService) DeleteUser(ctx context.Context, id uint) error {
	// Validate deletion
	if err := s.userDomainSvc.CanDeleteUser(ctx, id); err != nil {
		return err
	}

	// Delete user
	if err := s.userRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}

// ListUsersRequest represents the request for listing users
type ListUsersRequest struct {
	Limit  int
	Offset int
}

// ListUsersResponse represents the response for listing users
type ListUsersResponse struct {
	Users  []*entities.User `json:"users"`
	Total  int64            `json:"total"`
	Limit  int              `json:"limit"`
	Offset int              `json:"offset"`
}

// ListUsers retrieves a paginated list of users
func (s *UserService) ListUsers(ctx context.Context, req *ListUsersRequest) (*ListUsersResponse, error) {
	// Validate and set defaults
	if req.Limit <= 0 {
		req.Limit = 10
	}
	if req.Limit > 100 {
		req.Limit = 100 // Cap at 100
	}
	if req.Offset < 0 {
		req.Offset = 0
	}

	// Get users
	users, err := s.userRepo.List(ctx, req.Limit, req.Offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}

	// Get total count
	total, err := s.userRepo.Count(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to count users: %w", err)
	}

	return &ListUsersResponse{
		Users:  users,
		Total:  total,
		Limit:  req.Limit,
		Offset: req.Offset,
	}, nil
}