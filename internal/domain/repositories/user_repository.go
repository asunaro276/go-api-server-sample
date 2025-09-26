package repositories

import (
	"context"

	"go-api-server-sample/internal/domain/entities"
)

// UserRepository defines the interface for user data operations
type UserRepository interface {
	// Create creates a new user
	Create(ctx context.Context, user *entities.User) error

	// GetByID retrieves a user by ID
	GetByID(ctx context.Context, id uint) (*entities.User, error)

	// GetByEmail retrieves a user by email
	GetByEmail(ctx context.Context, email string) (*entities.User, error)

	// Update updates an existing user
	Update(ctx context.Context, user *entities.User) error

	// Delete soft deletes a user by ID
	Delete(ctx context.Context, id uint) error

	// List retrieves a list of users with pagination
	List(ctx context.Context, limit, offset int) ([]*entities.User, error)

	// Count returns the total number of users (excluding soft deleted)
	Count(ctx context.Context) (int64, error)

	// ExistsByEmail checks if a user with the given email exists
	ExistsByEmail(ctx context.Context, email string) (bool, error)

	// ExistsByEmailExcludingID checks if a user with the given email exists, excluding the specified ID
	ExistsByEmailExcludingID(ctx context.Context, email string, excludeID uint) (bool, error)
}