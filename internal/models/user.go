package models

import (
	"time"
)

// UserRole represents user roles
type UserRole string

const (
	// SuperAdminRole represents super admin role
	SuperAdminRole UserRole = "super_admin"
	// AdminRole represents admin role
	AdminRole UserRole = "admin"
)

// User represents a user in the system
type User struct {
	ID        int        `json:"id" db:"id"`
	Username  string     `json:"username" db:"username" validate:"required,min=3,max=50"`
	Email     string     `json:"email" db:"email" validate:"required,email"`
	Password  string     `json:"-" db:"password"` // Never expose password in JSON
	Role      UserRole   `json:"role" db:"role" validate:"required,oneof=super_admin admin"`
	IsActive  bool       `json:"is_active" db:"is_active"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt time.Time  `json:"updated_at" db:"updated_at"`
	LastLogin *time.Time `json:"last_login,omitempty" db:"last_login"`
}

// LoginRequest represents login request
type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// RegisterRequest represents registration request (only for super admin)
type RegisterRequest struct {
	Username string   `json:"username" validate:"required,min=3,max=50"`
	Email    string   `json:"email" validate:"required,email"`
	Password string   `json:"password" validate:"required,min=8"`
	Role     UserRole `json:"role" validate:"required,oneof=super_admin admin"`
}

// ChangePasswordRequest represents change password request
type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" validate:"required"`
	NewPassword     string `json:"new_password" validate:"required,min=8"`
}

// UpdateUserRequest represents update user request
type UpdateUserRequest struct {
	Email    *string   `json:"email,omitempty" validate:"omitempty,email"`
	Role     *UserRole `json:"role,omitempty" validate:"omitempty,oneof=super_admin admin"`
	IsActive *bool     `json:"is_active,omitempty"`
}

// AuthResponse represents authentication response
type AuthResponse struct {
	User         *UserResponse `json:"user"`
	AccessToken  string        `json:"access_token"`
	RefreshToken string        `json:"refresh_token"`
	ExpiresIn    int64         `json:"expires_in"` // seconds
	TokenType    string        `json:"token_type"`
}

// UserResponse represents user data in responses (without sensitive info)
type UserResponse struct {
	ID        int        `json:"id"`
	Username  string     `json:"username"`
	Email     string     `json:"email"`
	Role      UserRole   `json:"role"`
	IsActive  bool       `json:"is_active"`
	CreatedAt time.Time  `json:"created_at"`
	LastLogin *time.Time `json:"last_login,omitempty"`
}

// RefreshTokenRequest represents refresh token request
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

// JWTClaims represents JWT claims
type JWTClaims struct {
	UserID   int      `json:"user_id"`
	Username string   `json:"username"`
	Role     UserRole `json:"role"`
	Type     string   `json:"type"` // "access" or "refresh"
}

// ToUserResponse converts User to UserResponse
func (u *User) ToUserResponse() *UserResponse {
	return &UserResponse{
		ID:        u.ID,
		Username:  u.Username,
		Email:     u.Email,
		Role:      u.Role,
		IsActive:  u.IsActive,
		CreatedAt: u.CreatedAt,
		LastLogin: u.LastLogin,
	}
}
