package dtos

import (
	"time"

	"go-api-server-sample/internal/domain/entities"
)

// CreateUserRequest represents the request for creating a user
type CreateUserRequest struct {
	Name  string `json:"name" binding:"required,max=100" example:"田中太郎"`
	Email string `json:"email" binding:"required,email,max=255" example:"tanaka@example.com"`
}

// UpdateUserRequest represents the request for updating a user
type UpdateUserRequest struct {
	Name  string `json:"name,omitempty" binding:"omitempty,max=100" example:"田中次郎"`
	Email string `json:"email,omitempty" binding:"omitempty,email,max=255" example:"tanaka.jiro@example.com"`
}

// UserResponse represents the response for a user
type UserResponse struct {
	ID        uint      `json:"id" example:"1"`
	Name      string    `json:"name" example:"田中太郎"`
	Email     string    `json:"email" example:"tanaka@example.com"`
	CreatedAt time.Time `json:"created_at" example:"2023-12-01T12:00:00Z"`
	UpdatedAt time.Time `json:"updated_at" example:"2023-12-01T12:00:00Z"`
}

// UserListResponse represents the response for listing users
type UserListResponse struct {
	Users  []UserResponse `json:"users"`
	Total  int64          `json:"total" example:"100"`
	Limit  int            `json:"limit" example:"10"`
	Offset int            `json:"offset" example:"0"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Code    string        `json:"code" example:"VALIDATION_ERROR"`
	Message string        `json:"message" example:"バリデーションエラーが発生しました"`
	Details []ErrorDetail `json:"details,omitempty"`
}

// ErrorDetail represents detailed error information
type ErrorDetail struct {
	Field   string `json:"field" example:"email"`
	Message string `json:"message" example:"有効なメールアドレスを入力してください"`
}

// ToEntity converts CreateUserRequest to User entity
func (r *CreateUserRequest) ToEntity() *entities.User {
	return &entities.User{
		Name:  r.Name,
		Email: r.Email,
	}
}

// ToEntity converts UpdateUserRequest to User entity
func (r *UpdateUserRequest) ToEntity() *entities.User {
	return &entities.User{
		Name:  r.Name,
		Email: r.Email,
	}
}

// FromEntity converts User entity to UserResponse
func FromEntity(user *entities.User) *UserResponse {
	return &UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

// FromEntities converts slice of User entities to slice of UserResponse
func FromEntities(users []*entities.User) []UserResponse {
	responses := make([]UserResponse, len(users))
	for i, user := range users {
		responses[i] = *FromEntity(user)
	}
	return responses
}

// NewErrorResponse creates a new error response
func NewErrorResponse(code, message string) *ErrorResponse {
	return &ErrorResponse{
		Code:    code,
		Message: message,
	}
}

// NewValidationErrorResponse creates a validation error response
func NewValidationErrorResponse(message string, details []ErrorDetail) *ErrorResponse {
	return &ErrorResponse{
		Code:    "VALIDATION_ERROR",
		Message: message,
		Details: details,
	}
}

// MapDomainErrorToResponse maps domain errors to appropriate error responses
func MapDomainErrorToResponse(err error) *ErrorResponse {
	switch err {
	case entities.ErrUserNotFound:
		return NewErrorResponse("USER_NOT_FOUND", "指定されたユーザーが見つかりません")
	case entities.ErrUserEmailExists:
		return NewErrorResponse("EMAIL_ALREADY_EXISTS", "このメールアドレスは既に使用されています")
	case entities.ErrInvalidEmail:
		return NewErrorResponse("VALIDATION_ERROR", "有効なメールアドレスを入力してください")
	case entities.ErrUserNameRequired:
		return NewErrorResponse("VALIDATION_ERROR", "ユーザー名は必須です")
	case entities.ErrUserNameTooLong:
		return NewErrorResponse("VALIDATION_ERROR", "ユーザー名は100文字以内で入力してください")
	case entities.ErrInvalidUserID:
		return NewErrorResponse("VALIDATION_ERROR", "無効なユーザーIDです")
	default:
		return NewErrorResponse("INTERNAL_ERROR", "内部サーバーエラーが発生しました")
	}
}