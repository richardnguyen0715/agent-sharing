package middleware

import (
	"net/http"

	"github.com/termhub/backend/internal/domain"
)

// RequireRole returns middleware that restricts access to users with one of the specified roles.
// It must be used after Auth middleware, which sets the user claims in the context.
func RequireRole(roles ...domain.Role) func(http.Handler) http.Handler {
	// Build a set for O(1) role lookup.
	allowed := make(map[string]struct{}, len(roles))
	for _, r := range roles {
		allowed[string(r)] = struct{}{}
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims := GetUserFromContext(r.Context())
			if claims == nil {
				writeAuthError(w, http.StatusUnauthorized, "authentication required")
				return
			}

			if _, ok := allowed[claims.Role]; !ok {
				writeAuthError(w, http.StatusForbidden, "insufficient permissions")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
