package auth

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"habit-tracker/server/db"
)

// RegisterRequest represents the registration payload
type RegisterRequest struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

// RegisterResponse contains the user data after successful registration
type RegisterResponse struct {
	User    UserResponse `json:"user"`
	Message string       `json:"message"`
}

// LoginRequest represents the login payload
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginResponse contains the JWT token after successful login
type LoginResponse struct {
	Token   string  `json:"token"`
	User    db.User `json:"user"`
	Message string  `json:"message"`
}

// UserResponse represents user data for API responses
type UserResponse struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	Username string `json:"username"`
}

// ErrorResponse represents error responses
type ErrorResponse struct {
	Error string `json:"error"`
}

// ProfileResponse contains user profile data
type ProfileResponse struct {
	User UserResponse `json:"user"`
}

// ValidateResponse contains token validation result
type ValidateResponse struct {
	Valid bool          `json:"valid"`
	User  *UserResponse `json:"user,omitempty"`
}

// RegisterHandler handles user registration
func (s *AuthService) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse the request body
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Invalid JSON"})
		return
	}

	// Validate input
	if err := validateRegisterRequest(req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: err.Error()})
		return
	}

	// Call the auth service to register the user
	user, err := s.Register(req.Email, req.Username, req.Password)
	if err != nil {
		switch {
		case errors.Is(err, ErrEmailInUse):
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusConflict)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "Email already exists"})
			return
		case errors.Is(err, ErrUsernameInUse):
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusConflict)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "Username already exists"})
			return
		default:
			http.Error(w, "Error creating user", http.StatusInternalServerError)
		}
		return
	}

	// Return the created user (without sensitive data)
	response := RegisterResponse{
		User: UserResponse{
			ID:       user.ID,
			Email:    user.Email,
			Username: user.Username,
		},
		Message: "User registered successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// LoginHandler handles user login
func (s *AuthService) LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse the request body
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Invalid JSON"})
		return
	}

	// Validate input
	if req.Email == "" || req.Password == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Email and password are required"})
		return
	}

	// Attempt to login
	token, user, err := s.Login(req.Email, req.Password)
	if err != nil {
		if errors.Is(err, ErrInvalidCredentials) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "Invalid credentials"})
		} else {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "Internal server error"})
		}
		return
	}

	// Create a clean user object for response (without password hash)
	userResponse := db.User{
		ID:        user.ID,
		Email:     user.Email,
		Username:  user.Username,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	// Return the token and user info
	response := LoginResponse{
		Token:   token,
		User:    userResponse,
		Message: "Login successful",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// ProfileHandler returns the authenticated user's profile
func (s *AuthService) ProfileHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get user from request context (set by auth middleware)
	user := GetUserFromContext(r.Context())
	if user == nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "User not found in context"})
		return
	}

	// Return user profile (excluding sensitive data)
	response := ProfileResponse{
		User: UserResponse{
			ID:       user.ID,
			Email:    user.Email,
			Username: user.Username,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// ValidateTokenHandler checks if a token is valid
func (s *AuthService) ValidateTokenHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract token from header
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		response := ValidateResponse{Valid: false}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Extract token from "Bearer TOKEN" format
	if !strings.HasPrefix(authHeader, "Bearer ") {
		response := ValidateResponse{Valid: false}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(response)
		return
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	user, err := s.GetUserFromToken(tokenString)
	if err != nil {
		response := ValidateResponse{Valid: false}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(response)
		return
	}

	userResponse := &UserResponse{
		ID:       user.ID,
		Email:    user.Email,
		Username: user.Username,
	}

	response := ValidateResponse{
		Valid: true,
		User:  userResponse,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// validateRegisterRequest validates the registration request
func validateRegisterRequest(req RegisterRequest) error {
	if req.Email == "" {
		return errors.New("email is required")
	}
	if req.Username == "" {
		return errors.New("username is required")
	}
	if req.Password == "" {
		return errors.New("password is required")
	}

	// Basic email format validation
	if !strings.Contains(req.Email, "@") || !strings.Contains(req.Email, ".") {
		return errors.New("invalid email format")
	}

	// Basic password strength validation
	if len(req.Password) < 6 {
		return errors.New("password must be at least 6 characters long")
	}

	// Basic username validation
	if len(req.Username) < 3 {
		return errors.New("username must be at least 3 characters long")
	}

	return nil
}

// Standalone handler functions for testing and easier integration

// RegisterHandler creates a handler function for user registration
func RegisterHandler(authService *AuthService) http.HandlerFunc {
	return http.HandlerFunc(authService.RegisterHandler)
}

// LoginHandler creates a handler function for user login
func LoginHandler(authService *AuthService) http.HandlerFunc {
	return http.HandlerFunc(authService.LoginHandler)
}

// ProfileHandler creates a handler function for user profile
func ProfileHandler(authService *AuthService) http.HandlerFunc {
	return http.HandlerFunc(authService.ProfileHandler)
}

// ValidateHandler creates a handler function for token validation
func ValidateHandler(authService *AuthService) http.HandlerFunc {
	return http.HandlerFunc(authService.ValidateTokenHandler)
}
