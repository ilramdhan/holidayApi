package services

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"holidayapi/internal/models"
)

// JWTService handles JWT token operations
type JWTService interface {
	GenerateTokens(user *models.User) (*models.AuthResponse, error)
	ValidateAccessToken(tokenString string) (*models.JWTClaims, error)
	ValidateRefreshToken(tokenString string) (*models.JWTClaims, error)
	RefreshTokens(refreshToken string, userRepo UserRepository) (*models.AuthResponse, error)
}

// jwtService implements JWTService
type jwtService struct {
	secretKey       []byte
	accessTokenTTL  time.Duration
	refreshTokenTTL time.Duration
}

// NewJWTService creates a new JWT service
func NewJWTService(secretKey string, accessTokenTTL, refreshTokenTTL time.Duration) JWTService {
	return &jwtService{
		secretKey:       []byte(secretKey),
		accessTokenTTL:  accessTokenTTL,
		refreshTokenTTL: refreshTokenTTL,
	}
}

// GenerateTokens generates access and refresh tokens for a user
func (s *jwtService) GenerateTokens(user *models.User) (*models.AuthResponse, error) {
	// Generate access token
	accessToken, err := s.generateToken(user, "access", s.accessTokenTTL)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	// Generate refresh token
	refreshToken, err := s.generateToken(user, "refresh", s.refreshTokenTTL)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	return &models.AuthResponse{
		User:         user.ToUserResponse(),
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int64(s.accessTokenTTL.Seconds()),
		TokenType:    "Bearer",
	}, nil
}

// generateToken generates a JWT token
func (s *jwtService) generateToken(user *models.User, tokenType string, ttl time.Duration) (string, error) {
	now := time.Now()
	claims := jwt.MapClaims{
		"user_id":  user.ID,
		"username": user.Username,
		"role":     user.Role,
		"type":     tokenType,
		"iat":      now.Unix(),
		"exp":      now.Add(ttl).Unix(),
		"iss":      "holidayapi",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.secretKey)
}

// ValidateAccessToken validates an access token
func (s *jwtService) ValidateAccessToken(tokenString string) (*models.JWTClaims, error) {
	return s.validateToken(tokenString, "access")
}

// ValidateRefreshToken validates a refresh token
func (s *jwtService) ValidateRefreshToken(tokenString string) (*models.JWTClaims, error) {
	return s.validateToken(tokenString, "refresh")
}

// validateToken validates a JWT token
func (s *jwtService) validateToken(tokenString, expectedType string) (*models.JWTClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.secretKey, nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}

	// Check token type
	tokenType, ok := claims["type"].(string)
	if !ok || tokenType != expectedType {
		return nil, fmt.Errorf("invalid token type")
	}

	// Check issuer
	issuer, ok := claims["iss"].(string)
	if !ok || issuer != "holidayapi" {
		return nil, fmt.Errorf("invalid token issuer")
	}

	// Extract claims
	userID, ok := claims["user_id"].(float64)
	if !ok {
		return nil, fmt.Errorf("invalid user_id claim")
	}

	username, ok := claims["username"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid username claim")
	}

	role, ok := claims["role"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid role claim")
	}

	return &models.JWTClaims{
		UserID:   int(userID),
		Username: username,
		Role:     models.UserRole(role),
		Type:     tokenType,
	}, nil
}

// RefreshTokens generates new tokens using a refresh token
func (s *jwtService) RefreshTokens(refreshToken string, userRepo UserRepository) (*models.AuthResponse, error) {
	claims, err := s.ValidateRefreshToken(refreshToken)
	if err != nil {
		return nil, fmt.Errorf("invalid refresh token: %w", err)
	}

	// Fetch fresh user data from database
	user, err := userRepo.GetByID(claims.UserID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	// Check if user is still active
	if !user.IsActive {
		return nil, fmt.Errorf("user account is deactivated")
	}

	return s.GenerateTokens(user)
}
