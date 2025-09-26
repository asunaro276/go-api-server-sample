package services

import (
	"context"
	"fmt"

	"go-api-server-sample/internal/domain/entities"
	"go-api-server-sample/internal/domain/repositories"
)

// UserDomainService provides domain-level business logic for users
type UserDomainService struct {
	userRepo repositories.UserRepository
}

// NewUserDomainService creates a new UserDomainService
func NewUserDomainService(userRepo repositories.UserRepository) *UserDomainService {
	return &UserDomainService{
		userRepo: userRepo,
	}
}

// ValidateUserForCreation validates a user before creation
func (s *UserDomainService) ValidateUserForCreation(ctx context.Context, user *entities.User) error {
	// Validate entity fields
	if err := user.Validate(); err != nil {
		return err
	}

	// Check if email already exists
	exists, err := s.userRepo.ExistsByEmail(ctx, user.Email)
	if err != nil {
		return fmt.Errorf("failed to check email existence: %w", err)
	}
	if exists {
		return entities.ErrUserEmailExists
	}

	return nil
}

// ValidateUserForUpdate validates a user before update
func (s *UserDomainService) ValidateUserForUpdate(ctx context.Context, user *entities.User) error {
	// Validate entity fields
	if err := user.Validate(); err != nil {
		return err
	}

	// Check if email is taken by another user
	if user.Email != "" {
		exists, err := s.userRepo.ExistsByEmailExcludingID(ctx, user.Email, user.ID)
		if err != nil {
			return fmt.Errorf("failed to check email existence: %w", err)
		}
		if exists {
			return entities.ErrUserEmailExists
		}
	}

	return nil
}

// ValidateUserID validates a user ID
func (s *UserDomainService) ValidateUserID(id uint) error {
	if id == 0 {
		return entities.ErrInvalidUserID
	}
	return nil
}

// CanDeleteUser checks if a user can be deleted
func (s *UserDomainService) CanDeleteUser(ctx context.Context, userID uint) error {
	if err := s.ValidateUserID(userID); err != nil {
		return err
	}

	// Check if user exists
	_, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	// Add any business rules for deletion here
	// For example: check if user has active orders, etc.

	return nil
}