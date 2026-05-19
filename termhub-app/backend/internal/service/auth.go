// Package service contains business logic for authentication.
package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/termhub/backend/internal/config"
	"github.com/termhub/backend/internal/domain"
	"github.com/termhub/backend/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

// bcryptCost is the cost factor for bcrypt hashing.
const bcryptCost = 12

// Token type identifiers embedded in JWT claims.
const (
	tokenTypeAccess  = "access"
	tokenTypeRefresh = "refresh"
)

// Sentinel errors for the auth service.
var (
	ErrInvalidCredentials = errors.New("auth: invalid credentials")
	ErrInvalidToken       = errors.New("auth: invalid or expired token")
	ErrUserNotFound       = errors.New("auth: user not found")
	ErrUsernameExists     = errors.New("auth: username already exists")
	ErrInvalidRole        = errors.New("auth: invalid role")
	ErrInvalidInput       = errors.New("auth: invalid input")
)

// TokenClaims represents the claims embedded in a JWT token.
type TokenClaims struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	Type     string `json:"type"`
	jwt.RegisteredClaims
}

// AuthService handles user authentication and token management.
type AuthService struct {
	userRepo   *repository.UserRepository
	jwtSecret  []byte
	accessTTL  time.Duration
	refreshTTL time.Duration
}

// NewAuthService creates a new AuthService.
func NewAuthService(userRepo *repository.UserRepository, cfg *config.Config) *AuthService {
	return &AuthService{
		userRepo:   userRepo,
		jwtSecret:  []byte(cfg.JWTSecret),
		accessTTL:  cfg.JWTAccessTTL,
		refreshTTL: cfg.JWTRefreshTTL,
	}
}

// Register creates a new user account.
func (s *AuthService) Register(ctx context.Context, username, password, role string) (*domain.User, error) {
	if len(username) < 3 || len(username) > 32 {
		return nil, fmt.Errorf("auth: username must be between 3 and 32 characters: %w", ErrInvalidInput)
	}
	if len(password) < 8 {
		return nil, fmt.Errorf("auth: password must be at least 8 characters: %w", ErrInvalidInput)
	}

	userRole := domain.Role(role)
	if !userRole.IsValid() {
		return nil, fmt.Errorf("auth: role %q is not valid: %w", role, ErrInvalidRole)
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcryptCost)
	if err != nil {
		return nil, fmt.Errorf("auth: failed to hash password: %w", err)
	}

	now := time.Now().UTC()
	user := &domain.User{
		ID:           uuid.New(),
		Username:     username,
		PasswordHash: string(hash),
		Role:         userRole,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		if errors.Is(err, repository.ErrDuplicate) {
			return nil, fmt.Errorf("auth: username %q already taken: %w", username, ErrUsernameExists)
		}
		return nil, fmt.Errorf("auth: failed to create user: %w", err)
	}

	return user, nil
}

// Login authenticates a user by username and password.
func (s *AuthService) Login(ctx context.Context, username, password string) (accessToken, refreshToken string, user *domain.User, err error) {
	user, err = s.userRepo.GetByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return "", "", nil, ErrInvalidCredentials
		}
		return "", "", nil, fmt.Errorf("auth: failed to find user: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", "", nil, ErrInvalidCredentials
	}

	accessToken, err = s.generateToken(user, tokenTypeAccess, s.accessTTL)
	if err != nil {
		return "", "", nil, fmt.Errorf("auth: failed to generate access token: %w", err)
	}

	refreshToken, err = s.generateToken(user, tokenTypeRefresh, s.refreshTTL)
	if err != nil {
		return "", "", nil, fmt.Errorf("auth: failed to generate refresh token: %w", err)
	}

	return accessToken, refreshToken, user, nil
}

// RefreshToken validates a refresh token and issues a new token pair.
func (s *AuthService) RefreshToken(ctx context.Context, refreshToken string) (newAccessToken, newRefreshToken string, err error) {
	claims, err := s.ValidateToken(refreshToken)
	if err != nil {
		return "", "", fmt.Errorf("auth: invalid refresh token: %w", ErrInvalidToken)
	}

	if claims.Type != tokenTypeRefresh {
		return "", "", fmt.Errorf("auth: token is not a refresh token: %w", ErrInvalidToken)
	}

	// Fetch the current user to ensure they still exist and get latest role.
	userID, err := uuid.Parse(claims.UserID)
	if err != nil {
		return "", "", fmt.Errorf("auth: invalid user id in token: %w", ErrInvalidToken)
	}

	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return "", "", fmt.Errorf("auth: user no longer exists: %w", ErrInvalidToken)
		}
		return "", "", fmt.Errorf("auth: failed to fetch user for refresh: %w", err)
	}

	newAccessToken, err = s.generateToken(user, tokenTypeAccess, s.accessTTL)
	if err != nil {
		return "", "", fmt.Errorf("auth: failed to generate new access token: %w", err)
	}

	newRefreshToken, err = s.generateToken(user, tokenTypeRefresh, s.refreshTTL)
	if err != nil {
		return "", "", fmt.Errorf("auth: failed to generate new refresh token: %w", err)
	}

	return newAccessToken, newRefreshToken, nil
}

// ValidateToken parses and validates a JWT token string.
func (s *AuthService) ValidateToken(tokenString string) (*TokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("auth: unexpected signing method: %v", token.Header["alg"])
		}
		return s.jwtSecret, nil
	})
	if err != nil {
		return nil, fmt.Errorf("auth: failed to parse token: %w", err)
	}

	claims, ok := token.Claims.(*TokenClaims)
	if !ok || !token.Valid {
		return nil, ErrInvalidToken
	}

	return claims, nil
}

// GetUserByID retrieves a user by their ID string.
func (s *AuthService) GetUserByID(ctx context.Context, id string) (*domain.User, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("auth: invalid user id format: %w", ErrUserNotFound)
	}

	user, err := s.userRepo.GetByID(ctx, uid)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("auth: failed to get user: %w", err)
	}

	return user, nil
}

// GetUserRepo returns the user repository for direct access when needed.
func (s *AuthService) GetUserRepo() *repository.UserRepository {
	return s.userRepo
}

// IssueTokens generates a fresh access+refresh token pair for the given user.
// Used after registration so the user is immediately logged in.
func (s *AuthService) IssueTokens(user *domain.User) (accessToken, refreshToken string, err error) {
	accessToken, err = s.generateToken(user, tokenTypeAccess, s.accessTTL)
	if err != nil {
		return "", "", fmt.Errorf("auth: failed to generate access token: %w", err)
	}
	refreshToken, err = s.generateToken(user, tokenTypeRefresh, s.refreshTTL)
	if err != nil {
		return "", "", fmt.Errorf("auth: failed to generate refresh token: %w", err)
	}
	return accessToken, refreshToken, nil
}

// generateToken creates a signed JWT with the given user data, type, and TTL.
func (s *AuthService) generateToken(user *domain.User, tokenType string, ttl time.Duration) (string, error) {
	now := time.Now().UTC()
	claims := TokenClaims{
		UserID:   user.ID.String(),
		Username: user.Username,
		Role:     string(user.Role),
		Type:     tokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(ttl)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    "termhub",
			Subject:   user.ID.String(),
			ID:        uuid.New().String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(s.jwtSecret)
	if err != nil {
		return "", fmt.Errorf("auth: failed to sign token: %w", err)
	}

	return signed, nil
}
