package entities

import (
	"testing"
	"time"

	"gorm.io/gorm"
)

func TestUser_Validate(t *testing.T) {
	tests := []struct {
		name    string
		user    User
		wantErr bool
		errType error
	}{
		{
			name: "valid user",
			user: User{
				Name:  "John Doe",
				Email: "john@example.com",
			},
			wantErr: false,
		},
		{
			name: "empty name",
			user: User{
				Name:  "",
				Email: "john@example.com",
			},
			wantErr: true,
			errType: ErrUserNameRequired,
		},
		{
			name: "whitespace only name",
			user: User{
				Name:  "   ",
				Email: "john@example.com",
			},
			wantErr: true,
			errType: ErrUserNameRequired,
		},
		{
			name: "name too long",
			user: User{
				Name:  "a" + string(make([]byte, 100)), // 101 characters
				Email: "john@example.com",
			},
			wantErr: true,
			errType: ErrUserNameTooLong,
		},
		{
			name: "empty email",
			user: User{
				Name:  "John Doe",
				Email: "",
			},
			wantErr: true,
			errType: ErrInvalidEmail,
		},
		{
			name: "invalid email format - no @",
			user: User{
				Name:  "John Doe",
				Email: "johngmail.com",
			},
			wantErr: true,
			errType: ErrInvalidEmail,
		},
		{
			name: "invalid email format - no dot",
			user: User{
				Name:  "John Doe",
				Email: "john@gmail",
			},
			wantErr: true,
			errType: ErrInvalidEmail,
		},
		{
			name: "email too long",
			user: User{
				Name:  "John Doe",
				Email: string(make([]byte, 250)) + "@example.com", // > 255 characters
			},
			wantErr: true,
			errType: ErrInvalidEmail,
		},
		{
			name: "name with leading/trailing spaces (should be trimmed)",
			user: User{
				Name:  "  John Doe  ",
				Email: "john@example.com",
			},
			wantErr: false,
		},
		{
			name: "email with leading/trailing spaces and mixed case (should be normalized)",
			user: User{
				Name:  "John Doe",
				Email: "  JOHN@EXAMPLE.COM  ",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.user.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("User.Validate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && err != tt.errType {
				t.Errorf("User.Validate() error = %v, want %v", err, tt.errType)
			}

			// Check that normalization occurred for valid cases
			if !tt.wantErr {
				if tt.user.Name != "John Doe" && tt.name == "name with leading/trailing spaces (should be trimmed)" {
					t.Errorf("Expected name to be trimmed, got %q", tt.user.Name)
				}
				if tt.user.Email != "john@example.com" && tt.name == "email with leading/trailing spaces and mixed case (should be normalized)" {
					t.Errorf("Expected email to be normalized, got %q", tt.user.Email)
				}
			}
		})
	}
}

func TestUser_TableName(t *testing.T) {
	user := User{}
	if got := user.TableName(); got != "users" {
		t.Errorf("User.TableName() = %v, want %v", got, "users")
	}
}

func TestUser_IsDeleted(t *testing.T) {
	tests := []struct {
		name string
		user User
		want bool
	}{
		{
			name: "not deleted",
			user: User{
				DeletedAt: gorm.DeletedAt{},
			},
			want: false,
		},
		{
			name: "deleted",
			user: User{
				DeletedAt: gorm.DeletedAt{
					Time:  time.Now(),
					Valid: true,
				},
			},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.user.IsDeleted(); got != tt.want {
				t.Errorf("User.IsDeleted() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUser_UpdateFrom(t *testing.T) {
	tests := []struct {
		name     string
		original User
		update   User
		expected User
	}{
		{
			name: "update both name and email",
			original: User{
				ID:    1,
				Name:  "John Doe",
				Email: "john@example.com",
			},
			update: User{
				Name:  "Jane Doe",
				Email: "jane@example.com",
			},
			expected: User{
				ID:    1,
				Name:  "Jane Doe",
				Email: "jane@example.com",
			},
		},
		{
			name: "update only name",
			original: User{
				ID:    1,
				Name:  "John Doe",
				Email: "john@example.com",
			},
			update: User{
				Name:  "Jane Doe",
				Email: "",
			},
			expected: User{
				ID:    1,
				Name:  "Jane Doe",
				Email: "john@example.com",
			},
		},
		{
			name: "update only email",
			original: User{
				ID:    1,
				Name:  "John Doe",
				Email: "john@example.com",
			},
			update: User{
				Name:  "",
				Email: "jane@example.com",
			},
			expected: User{
				ID:    1,
				Name:  "John Doe",
				Email: "jane@example.com",
			},
		},
		{
			name: "update with empty values (no change)",
			original: User{
				ID:    1,
				Name:  "John Doe",
				Email: "john@example.com",
			},
			update: User{
				Name:  "",
				Email: "",
			},
			expected: User{
				ID:    1,
				Name:  "John Doe",
				Email: "john@example.com",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.original.UpdateFrom(&tt.update)
			
			if tt.original.ID != tt.expected.ID {
				t.Errorf("Expected ID %v, got %v", tt.expected.ID, tt.original.ID)
			}
			if tt.original.Name != tt.expected.Name {
				t.Errorf("Expected Name %v, got %v", tt.expected.Name, tt.original.Name)
			}
			if tt.original.Email != tt.expected.Email {
				t.Errorf("Expected Email %v, got %v", tt.expected.Email, tt.original.Email)
			}
		})
	}
}

func TestUser_BeforeCreate(t *testing.T) {
	tests := []struct {
		name    string
		user    User
		wantErr bool
	}{
		{
			name: "valid user",
			user: User{
				Name:  "John Doe",
				Email: "john@example.com",
			},
			wantErr: false,
		},
		{
			name: "invalid user",
			user: User{
				Name:  "",
				Email: "john@example.com",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.user.BeforeCreate(nil)
			if (err != nil) != tt.wantErr {
				t.Errorf("User.BeforeCreate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUser_BeforeUpdate(t *testing.T) {
	tests := []struct {
		name    string
		user    User
		wantErr bool
	}{
		{
			name: "valid user",
			user: User{
				Name:  "John Doe",
				Email: "john@example.com",
			},
			wantErr: false,
		},
		{
			name: "invalid user",
			user: User{
				Name:  "",
				Email: "john@example.com",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.user.BeforeUpdate(nil)
			if (err != nil) != tt.wantErr {
				t.Errorf("User.BeforeUpdate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// Benchmark tests
func BenchmarkUser_Validate(b *testing.B) {
	user := User{
		Name:  "John Doe",
		Email: "john@example.com",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = user.Validate()
	}
}

func BenchmarkUser_UpdateFrom(b *testing.B) {
	original := User{
		ID:    1,
		Name:  "John Doe",
		Email: "john@example.com",
	}
	update := User{
		Name:  "Jane Doe",
		Email: "jane@example.com",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		original.UpdateFrom(&update)
	}
}