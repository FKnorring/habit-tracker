package auth

import (
	"context"
	"habit-tracker/server/db"
	"net/http"
	"strings"
)

// ContextKey type for context values
type ContextKey string

const (
	// UserContextKey is the key for user info in the request context
	UserContextKey ContextKey = "user"
	// UserIDContextKey is the key for user ID in the request context
	UserIDContextKey ContextKey = "userID"
)

// AuthMiddleware creates a middleware function for validating JWT tokens
func AuthMiddleware(authService *AuthService) func(http.Handler) http.Handler {
	return authService.AuthMiddleware
}

// OptionalAuthMiddleware creates a middleware function for optional JWT validation
func OptionalAuthMiddleware(authService *AuthService) func(http.Handler) http.Handler {
	return authService.OptionalAuthMiddleware
}

// AuthMiddleware creates middleware that validates JWT tokens
func (s *AuthService) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract token from Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header required", http.StatusUnauthorized)
			return
		}

		// Check Bearer token format
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			http.Error(w, "Invalid authorization format. Use: Bearer <token>", http.StatusUnauthorized)
			return
		}

		tokenString := parts[1]

		// Validate the token and get user
		user, err := s.GetUserFromToken(tokenString)
		if err != nil {
			switch err {
			case ErrExpiredToken:
				http.Error(w, "Token has expired", http.StatusUnauthorized)
			case ErrInvalidToken:
				http.Error(w, "Invalid token", http.StatusUnauthorized)
			default:
				http.Error(w, "Authentication failed", http.StatusUnauthorized)
			}
			return
		}

		// Add user to request context
		ctx := context.WithValue(r.Context(), UserContextKey, user)
		ctx = context.WithValue(ctx, UserIDContextKey, user.ID)

		// Call the next handler with the enhanced context
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// OptionalAuthMiddleware creates middleware that validates JWT tokens but doesn't require them
func (s *AuthService) OptionalAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract token from Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader != "" {
			// Check Bearer token format
			parts := strings.Split(authHeader, " ")
			if len(parts) == 2 && strings.ToLower(parts[0]) == "bearer" {
				tokenString := parts[1]

				// Try to validate the token and get user
				user, err := s.GetUserFromToken(tokenString)
				if err == nil {
					// Add user to request context
					ctx := context.WithValue(r.Context(), UserContextKey, user)
					ctx = context.WithValue(ctx, UserIDContextKey, user.ID)
					r = r.WithContext(ctx)
				}
			}
		}

		// Call the next handler (with or without user context)
		next.ServeHTTP(w, r)
	})
}

// GetUserFromContext extracts the user from the request context
func GetUserFromContext(ctx context.Context) *db.User {
	if ctx == nil {
		return nil
	}
	user, _ := ctx.Value(UserContextKey).(*db.User)
	return user
}

// GetUserIDFromContext extracts the user ID from the request context
func GetUserIDFromContext(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	userID, _ := ctx.Value(UserIDContextKey).(string)
	return userID
}

// SetUserInContext adds a user to the context
func SetUserInContext(ctx context.Context, user *db.User) context.Context {
	ctx = context.WithValue(ctx, UserContextKey, user)
	ctx = context.WithValue(ctx, UserIDContextKey, user.ID)
	return ctx
}

// Legacy functions for backward compatibility
func GetUserFromRequest(r *http.Request) (*db.User, bool) {
	user, ok := r.Context().Value(UserContextKey).(*db.User)
	return user, ok
}

func GetUserIDFromRequest(r *http.Request) (string, bool) {
	userID, ok := r.Context().Value(UserIDContextKey).(string)
	return userID, ok
}
