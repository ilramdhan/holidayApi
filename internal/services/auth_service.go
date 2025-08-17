package services

import (
	"fmt"
	"strings"

	"holidayapi/internal/models"
	"holidayapi/internal/repository"
)

// UserRepository interface (forward declaration)
type UserRepository interface {
	Create(user *models.User) error
	GetByID(id int) (*models.User, error)
	GetByUsername(username string) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
	Update(id int, user *models.User) error
	Delete(id int) error
	UpdateLastLogin(id int) error
	GetAll() ([]models.User, error)
	HashPassword(password string) (string, error)
	CheckPassword(hashedPassword, password string) error
	ChangePassword(userID int, newPassword string) error
}

// AuthService handles authentication operations
type AuthService interface {
	Login(req models.LoginRequest, ipAddress, userAgent string) (*models.AuthResponse, error)
	Register(req models.RegisterRequest, createdBy *models.User) (*models.User, error)
	RefreshToken(req models.RefreshTokenRequest) (*models.AuthResponse, error)
	ChangePassword(userID int, req models.ChangePasswordRequest) error
	GetUserProfile(userID int) (*models.UserResponse, error)
	UpdateUserProfile(userID int, req models.UpdateUserRequest) (*models.UserResponse, error)
	GetAllUsers() ([]models.UserResponse, error)
	DeleteUser(userID int, deletedBy *models.User) error
}

// authService implements AuthService
type authService struct {
	userRepo   UserRepository
	auditRepo  repository.AuditRepository
	jwtService JWTService
}

// NewAuthService creates a new auth service
func NewAuthService(userRepo UserRepository, auditRepo repository.AuditRepository, jwtService JWTService) AuthService {
	return &authService{
		userRepo:   userRepo,
		auditRepo:  auditRepo,
		jwtService: jwtService,
	}
}

// Login authenticates a user and returns tokens
func (s *authService) Login(req models.LoginRequest, ipAddress, userAgent string) (*models.AuthResponse, error) {
	// Get user by username
	user, err := s.userRepo.GetByUsername(req.Username)
	if err != nil {
		// Log failed login attempt
		s.logAudit(nil, req.Username, models.ActionLoginFailed, models.ResourceAuth,
			fmt.Sprintf("Login failed: user not found"), ipAddress, userAgent, false)
		return nil, fmt.Errorf("invalid credentials")
	}

	// Check password
	if err := s.userRepo.CheckPassword(user.Password, req.Password); err != nil {
		// Log failed login attempt
		s.logAudit(&user.ID, user.Username, models.ActionLoginFailed, models.ResourceAuth,
			fmt.Sprintf("Login failed: invalid password"), ipAddress, userAgent, false)
		return nil, fmt.Errorf("invalid credentials")
	}

	// Check if user is active
	if !user.IsActive {
		s.logAudit(&user.ID, user.Username, models.ActionLoginFailed, models.ResourceAuth,
			fmt.Sprintf("Login failed: account deactivated"), ipAddress, userAgent, false)
		return nil, fmt.Errorf("account is deactivated")
	}

	// Generate tokens
	authResponse, err := s.jwtService.GenerateTokens(user)
	if err != nil {
		s.logAudit(&user.ID, user.Username, models.ActionLoginFailed, models.ResourceAuth,
			fmt.Sprintf("Login failed: token generation error"), ipAddress, userAgent, false)
		return nil, fmt.Errorf("failed to generate tokens: %w", err)
	}

	// Update last login
	if err := s.userRepo.UpdateLastLogin(user.ID); err != nil {
		// Log but don't fail the login
		fmt.Printf("Failed to update last login for user %d: %v\n", user.ID, err)
	}

	// Log successful login
	s.logAudit(&user.ID, user.Username, models.ActionLogin, models.ResourceAuth,
		fmt.Sprintf("User logged in successfully"), ipAddress, userAgent, true)

	return authResponse, nil
}

// Register creates a new user (only super admin can register new users)
func (s *authService) Register(req models.RegisterRequest, createdBy *models.User) (*models.User, error) {
	// Validate password strength
	if err := s.validatePassword(req.Password); err != nil {
		return nil, err
	}

	// Check if username already exists
	if _, err := s.userRepo.GetByUsername(req.Username); err == nil {
		return nil, fmt.Errorf("username already exists")
	}

	// Check if email already exists
	if _, err := s.userRepo.GetByEmail(req.Email); err == nil {
		return nil, fmt.Errorf("email already exists")
	}

	// Create user
	user := &models.User{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password, // Will be hashed in repository
		Role:     req.Role,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Log user creation
	s.logAudit(&createdBy.ID, createdBy.Username, models.ActionUserCreate, models.ResourceUser,
		fmt.Sprintf("Created new user: %s (role: %s)", user.Username, user.Role), "", "", true)

	// Clear password before returning
	user.Password = ""
	return user, nil
}

// RefreshToken generates new tokens using refresh token
func (s *authService) RefreshToken(req models.RefreshTokenRequest) (*models.AuthResponse, error) {
	authResponse, err := s.jwtService.RefreshTokens(req.RefreshToken, s.userRepo)
	if err != nil {
		return nil, fmt.Errorf("failed to refresh token: %w", err)
	}

	// Log token refresh
	s.logAudit(nil, authResponse.User.Username, models.ActionTokenRefresh, models.ResourceAuth,
		"Token refreshed successfully", "", "", true)

	return authResponse, nil
}

// ChangePassword changes user's password
func (s *authService) ChangePassword(userID int, req models.ChangePasswordRequest) error {
	// Get user
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	// Verify current password
	if err := s.userRepo.CheckPassword(user.Password, req.CurrentPassword); err != nil {
		s.logAudit(&userID, user.Username, models.ActionPasswordChange, models.ResourceUser,
			"Password change failed: invalid current password", "", "", false)
		return fmt.Errorf("current password is incorrect")
	}

	// Validate new password
	if err := s.validatePassword(req.NewPassword); err != nil {
		return err
	}

	// Change password
	if err := s.userRepo.ChangePassword(userID, req.NewPassword); err != nil {
		s.logAudit(&userID, user.Username, models.ActionPasswordChange, models.ResourceUser,
			"Password change failed: database error", "", "", false)
		return fmt.Errorf("failed to change password: %w", err)
	}

	// Log password change
	s.logAudit(&userID, user.Username, models.ActionPasswordChange, models.ResourceUser,
		"Password changed successfully", "", "", true)

	return nil
}

// GetUserProfile gets user profile
func (s *authService) GetUserProfile(userID int) (*models.UserResponse, error) {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	return user.ToUserResponse(), nil
}

// UpdateUserProfile updates user profile
func (s *authService) UpdateUserProfile(userID int, req models.UpdateUserRequest) (*models.UserResponse, error) {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	// Update fields if provided
	if req.Email != nil {
		user.Email = *req.Email
	}
	if req.Role != nil {
		user.Role = *req.Role
	}
	if req.IsActive != nil {
		user.IsActive = *req.IsActive
	}

	if err := s.userRepo.Update(userID, user); err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	// Log user update
	s.logAudit(&userID, user.Username, models.ActionUserUpdate, models.ResourceUser,
		"User profile updated", "", "", true)

	return user.ToUserResponse(), nil
}

// GetAllUsers gets all users (admin only)
func (s *authService) GetAllUsers() ([]models.UserResponse, error) {
	users, err := s.userRepo.GetAll()
	if err != nil {
		return nil, fmt.Errorf("failed to get users: %w", err)
	}

	var userResponses []models.UserResponse
	for _, user := range users {
		userResponses = append(userResponses, *user.ToUserResponse())
	}

	return userResponses, nil
}

// DeleteUser deletes a user (soft delete)
func (s *authService) DeleteUser(userID int, deletedBy *models.User) error {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	if err := s.userRepo.Delete(userID); err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	// Log user deletion
	s.logAudit(&deletedBy.ID, deletedBy.Username, models.ActionUserDelete, models.ResourceUser,
		fmt.Sprintf("Deleted user: %s", user.Username), "", "", true)

	return nil
}

// validatePassword validates password strength
func (s *authService) validatePassword(password string) error {
	if len(password) < 8 {
		return fmt.Errorf("password must be at least 8 characters long")
	}

	hasUpper := false
	hasLower := false
	hasDigit := false
	hasSpecial := false

	for _, char := range password {
		switch {
		case 'A' <= char && char <= 'Z':
			hasUpper = true
		case 'a' <= char && char <= 'z':
			hasLower = true
		case '0' <= char && char <= '9':
			hasDigit = true
		case strings.ContainsRune("!@#$%^&*()_+-=[]{}|;:,.<>?", char):
			hasSpecial = true
		}
	}

	if !hasUpper {
		return fmt.Errorf("password must contain at least one uppercase letter")
	}
	if !hasLower {
		return fmt.Errorf("password must contain at least one lowercase letter")
	}
	if !hasDigit {
		return fmt.Errorf("password must contain at least one digit")
	}
	if !hasSpecial {
		return fmt.Errorf("password must contain at least one special character")
	}

	return nil
}

// logAudit logs an audit entry
func (s *authService) logAudit(userID *int, username string, action models.AuditAction, resource models.AuditResource, details, ipAddress, userAgent string, success bool) {
	auditLog := &models.AuditLog{
		UserID:    userID,
		Username:  username,
		Action:    action,
		Resource:  resource,
		Details:   details,
		IPAddress: ipAddress,
		UserAgent: userAgent,
		Success:   success,
	}

	if err := s.auditRepo.Create(auditLog); err != nil {
		fmt.Printf("Failed to create audit log: %v\n", err)
	}
}
