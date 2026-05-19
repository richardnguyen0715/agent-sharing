// Package middleware provides HTTP middleware for authentication.
package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

// contextKey is a custom type for context keys to avoid collisions.
type contextKey string

// UserContextKey is the key used to store user claims in the request context.
const UserContextKey contextKey = "user"

// UserClaims represents the authenticated user's identity extracted from a JWT.
type UserClaims struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	Type     string `json:"type"`
	jwt.RegisteredClaims
}

// GetUserFromContext retrieves the authenticated user's claims from the context.
// Returns nil if no user is present.
func GetUserFromContext(ctx context.Context) *UserClaims {
	claims, ok := ctx.Value(UserContextKey).(*UserClaims)
	if !ok {
		return nil
	}
	return claims
}

// Auth returns middleware that validates JWT tokens from the Authorization header.
// It extracts the user identity and injects it into the request context.
func Auth(jwtSecret string) func(http.Handler) http.Handler {
	secret := []byte(jwtSecret)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				writeAuthError(w, http.StatusUnauthorized, "missing authorization header")
				return
			}

			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || !strings.EqualFold(parts[0], "bearer") {
				writeAuthError(w, http.StatusUnauthorized, "invalid authorization header format")
				return
			}

			tokenString := parts[1]
			if tokenString == "" {
				writeAuthError(w, http.StatusUnauthorized, "empty token")
				return
			}

			token, err := jwt.ParseWithClaims(tokenString, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
				}
				return secret, nil
			})
			if err != nil {
				writeAuthError(w, http.StatusUnauthorized, "invalid or expired token")
				return
			}

			claims, ok := token.Claims.(*UserClaims)
			if !ok || !token.Valid {
				writeAuthError(w, http.StatusUnauthorized, "invalid token claims")
				return
			}

			// Reject refresh tokens used as access tokens
			if claims.Type == "refresh" {
				writeAuthError(w, http.StatusUnauthorized, "refresh token cannot be used for authentication")
				return
			}

			ctx := context.WithValue(r.Context(), UserContextKey, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// writeAuthError writes a JSON error response for authentication failures.
func writeAuthError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}
