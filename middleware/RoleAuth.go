package middleware

import (
	"net/http"
)

// RoleAuth returns a middleware that allows access only to users with the specified role
func RoleAuth(requiredRole string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get the role from context (set by JWTAuth)
			roleCtx := r.Context().Value("role")
			if roleCtx == nil {
				http.Error(w, "Role not found in token", http.StatusUnauthorized)
				return
			}

			role, ok := roleCtx.(string)
			if !ok {
				http.Error(w, "Invalid role type", http.StatusUnauthorized)
				return
			}

			// Check if user has required role
			if role != requiredRole {
				http.Error(w, "Forbidden: insufficient permissions", http.StatusForbidden)
				return
			}

			// User has required role, continue
			next.ServeHTTP(w, r)
		})
	}
}
