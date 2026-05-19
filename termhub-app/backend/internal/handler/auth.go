// Package handler provides HTTP handlers for authentication endpoints.
package handler

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
	"github.com/termhub/backend/internal/domain"
	"github.com/termhub/backend/internal/middleware"
	"github.com/termhub/backend/internal/service"
)

// NOTE: auditService is used in Login and Register to log authentication events.

// LoginRequest represents the login request body.
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// RegisterRequest represents the registration request body.
type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

// RefreshRequest represents the token refresh request body.
type RefreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

// AuthResponse represents the authentication response body.
type AuthResponse struct {
	AccessToken  string       `json:"access_token"`
	RefreshToken string       `json:"refresh_token"`
	User         *domain.User `json:"user"`
}

// AuthHandler handles authentication-related HTTP requests.
type AuthHandler struct {
	authService  *service.AuthService
	auditService *service.AuditService
	jwtSecret    string
}

// NewAuthHandler creates a new AuthHandler.
func NewAuthHandler(authService *service.AuthService, auditService *service.AuditService, jwtSecret string) *AuthHandler {
	return &AuthHandler{authService: authService, auditService: auditService, jwtSecret: jwtSecret}
}

// Routes returns a chi.Router with all auth routes mounted.
func (h *AuthHandler) Routes() chi.Router {
	r := chi.NewRouter()
	r.Post("/login", h.Login)
	r.Post("/register", h.Register)
	r.Post("/refresh", h.Refresh)
	r.With(middleware.Auth(h.jwtSecret)).Get("/me", h.Me)
	return r
}

// Login authenticates a user and returns access and refresh tokens.
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := readJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Username == "" || req.Password == "" {
		writeError(w, http.StatusBadRequest, "username and password are required")
		return
	}

	accessToken, refreshToken, user, err := h.authService.Login(r.Context(), req.Username, req.Password)
	if err != nil {
		if errors.Is(err, service.ErrInvalidCredentials) {
			writeError(w, http.StatusUnauthorized, "invalid username or password")
			return
		}
		log.Error().Err(err).Str("username", req.Username).Msg("login failed")
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	// Audit log the successful login.
	userID := &user.ID
	if auditErr := h.auditService.Log(r.Context(), userID, nil, domain.AuditActionUserLogin, map[string]interface{}{
		"username":   user.Username,
		"ip_address": r.RemoteAddr,
	}); auditErr != nil {
		log.Error().Err(auditErr).Str("user_id", user.ID.String()).Msg("failed to write login audit log")
	}

	writeJSON(w, http.StatusOK, AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         user,
	})
}

// Register creates a new user account.
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := readJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Username == "" || req.Password == "" {
		writeError(w, http.StatusBadRequest, "username and password are required")
		return
	}

	if len(req.Username) < 3 || len(req.Username) > 32 {
		writeError(w, http.StatusBadRequest, "username must be between 3 and 32 characters")
		return
	}

	if len(req.Password) < 8 {
		writeError(w, http.StatusBadRequest, "password must be at least 8 characters")
		return
	}

	// Default role to viewer if not specified.
	role := req.Role
	if role == "" {
		role = string(domain.RoleViewer)
	}

	if !domain.Role(role).IsValid() {
		writeError(w, http.StatusBadRequest, "invalid role: must be one of admin, operator, viewer, restricted")
		return
	}

	user, err := h.authService.Register(r.Context(), req.Username, req.Password, role)
	if err != nil {
		if errors.Is(err, service.ErrUsernameExists) {
			writeError(w, http.StatusConflict, "username already exists")
			return
		}
		if errors.Is(err, service.ErrInvalidInput) {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		if errors.Is(err, service.ErrInvalidRole) {
			writeError(w, http.StatusBadRequest, "invalid role")
			return
		}
		log.Error().Err(err).Str("username", req.Username).Msg("registration failed")
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	// Audit log the successful registration.
	userID := &user.ID
	if auditErr := h.auditService.Log(r.Context(), userID, nil, domain.AuditActionUserRegister, map[string]interface{}{
		"username":   user.Username,
		"role":       string(user.Role),
		"ip_address": r.RemoteAddr,
	}); auditErr != nil {
		log.Error().Err(auditErr).Str("user_id", user.ID.String()).Msg("failed to write register audit log")
	}

	writeJSON(w, http.StatusCreated, map[string]interface{}{
		"user": user,
	})
}

// Refresh issues a new access token from a valid refresh token.
func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	var req RefreshRequest
	if err := readJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.RefreshToken == "" {
		writeError(w, http.StatusBadRequest, "refresh_token is required")
		return
	}

	accessToken, refreshToken, err := h.authService.RefreshToken(r.Context(), req.RefreshToken)
	if err != nil {
		if errors.Is(err, service.ErrInvalidToken) {
			writeError(w, http.StatusUnauthorized, "invalid or expired refresh token")
			return
		}
		log.Error().Err(err).Msg("token refresh failed")
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

// Me returns the currently authenticated user's profile.
func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetUserFromContext(r.Context())
	if claims == nil {
		writeError(w, http.StatusUnauthorized, "authentication required")
		return
	}

	user, err := h.authService.GetUserByID(r.Context(), claims.UserID)
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			writeError(w, http.StatusNotFound, "user not found")
			return
		}
		log.Error().Err(err).Str("user_id", claims.UserID).Msg("failed to get user profile")
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"user": user,
	})
}
