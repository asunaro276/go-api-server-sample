package entities

import (
	"errors"
	"strings"
	"time"

	"gorm.io/gorm"
)

// User represents a user entity in the domain
type User struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	Name      string         `gorm:"size:100;not null" json:"name" validate:"required,max=100"`
	Email     string         `gorm:"size:255;uniqueIndex;not null" json:"email" validate:"required,email"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// Domain errors
var (
	ErrUserNotFound      = errors.New("user not found")
	ErrUserEmailExists   = errors.New("user email already exists")
	ErrInvalidEmail      = errors.New("invalid email format")
	ErrUserNameRequired  = errors.New("user name is required")
	ErrUserNameTooLong   = errors.New("user name too long")
	ErrInvalidUserID     = errors.New("invalid user ID")
)

// TableName returns the table name for the User entity
func (User) TableName() string {
	return "users"
}

// Validate validates the user entity
func (u *User) Validate() error {
	if err := u.validateName(); err != nil {
		return err
	}
	if err := u.validateEmail(); err != nil {
		return err
	}
	return nil
}

// validateName validates the user name
func (u *User) validateName() error {
	u.Name = strings.TrimSpace(u.Name)
	if u.Name == "" {
		return ErrUserNameRequired
	}
	if len(u.Name) > 100 {
		return ErrUserNameTooLong
	}
	return nil
}

// validateEmail validates the user email
func (u *User) validateEmail() error {
	u.Email = strings.TrimSpace(strings.ToLower(u.Email))
	if u.Email == "" {
		return ErrInvalidEmail
	}
	// Basic email validation
	if !strings.Contains(u.Email, "@") || !strings.Contains(u.Email, ".") {
		return ErrInvalidEmail
	}
	if len(u.Email) > 255 {
		return ErrInvalidEmail
	}
	return nil
}

// BeforeCreate hook is called before creating a user
func (u *User) BeforeCreate(tx *gorm.DB) error {
	return u.Validate()
}

// BeforeUpdate hook is called before updating a user
func (u *User) BeforeUpdate(tx *gorm.DB) error {
	return u.Validate()
}

// IsDeleted returns true if the user is soft deleted
func (u *User) IsDeleted() bool {
	return u.DeletedAt.Valid
}

// UpdateFrom updates the user with non-zero values from another user
func (u *User) UpdateFrom(other *User) {
	if other.Name != "" {
		u.Name = other.Name
	}
	if other.Email != "" {
		u.Email = other.Email
	}
}