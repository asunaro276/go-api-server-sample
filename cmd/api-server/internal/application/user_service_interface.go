package application

import (
	"context"

	"go-api-server-sample/internal/domain/entities"
)

// UserServiceInterface defines the interface for user service operations
type UserServiceInterface interface {
	CreateUser(ctx context.Context, user *entities.User) (*entities.User, error)
	GetUserByID(ctx context.Context, id uint) (*entities.User, error)
	UpdateUser(ctx context.Context, id uint, updateData *entities.User) (*entities.User, error)
	DeleteUser(ctx context.Context, id uint) error
	ListUsers(ctx context.Context, req *ListUsersRequest) (*ListUsersResponse, error)
}

// Ensure UserService implements the interface
var _ UserServiceInterface = (*UserService)(nil)