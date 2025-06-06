package auth

import (
	"errors"
	"time"

	"habit-tracker/server/db"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidToken       = errors.New("invalid token")
	ErrExpiredToken       = errors.New("token is expired")
	ErrEmailInUse         = errors.New("email already in use")
	ErrUsernameInUse      = errors.New("username already in use")
)

// AuthService provides authentication functionality
type AuthService struct {
	database    db.Database
	jwtSecret   []byte
	tokenExpiry time.Duration
}

// NewAuthService creates a new authentication service
func NewAuthService(database db.Database, jwtSecret string, tokenExpiry time.Duration) *AuthService {
	if database == nil {
		panic("database cannot be nil")
	}
	if jwtSecret == "" {
		panic("JWT secret cannot be empty")
	}
	if tokenExpiry <= 0 {
		panic("token expiry must be positive")
	}

	return &AuthService{
		database:    database,
		jwtSecret:   []byte(jwtSecret),
		tokenExpiry: tokenExpiry,
	}
}

// HashPassword creates a bcrypt hash from a plain-text password
func (s *AuthService) HashPassword(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedBytes), nil
}

// VerifyPassword checks if the provided password matches the stored hash
func (s *AuthService) VerifyPassword(hashedPassword, providedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(providedPassword))
}

// Register creates a new user with the provided credentials
func (s *AuthService) Register(email, username, password string) (*db.User, error) {
	// Validate input
	if email == "" {
		return nil, errors.New("email is required")
	}
	if username == "" {
		return nil, errors.New("username is required")
	}
	if password == "" {
		return nil, errors.New("password is required")
	}

	// Check if user already exists by email
	_, err := s.database.GetUserByEmail(email)
	if err == nil {
		return nil, ErrEmailInUse
	}
	if !errors.Is(err, db.ErrNotFound) {
		return nil, err
	}

	// Hash the password
	hashedPassword, err := s.HashPassword(password)
	if err != nil {
		return nil, err
	}

	// Create the user
	user := &db.User{
		Email:        email,
		Username:     username,
		PasswordHash: hashedPassword,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	err = s.database.CreateUser(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// Login authenticates a user and returns a JWT token
func (s *AuthService) Login(email, password string) (string, *db.User, error) {
	// Get the user from the database
	user, err := s.database.GetUserByEmail(email)
	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			return "", nil, ErrInvalidCredentials
		}
		return "", nil, err
	}

	// Verify the password
	if err := s.VerifyPassword(user.PasswordHash, password); err != nil {
		return "", nil, ErrInvalidCredentials
	}

	// Generate a JWT token
	token, err := s.generateToken(user)
	if err != nil {
		return "", nil, err
	}

	return token, user, nil
}

// generateToken creates a new JWT token for a user
func (s *AuthService) generateToken(user *db.User) (string, error) {
	// Set the expiration time
	expirationTime := time.Now().Add(s.tokenExpiry)

	// Create the JWT claims
	claims := jwt.MapClaims{
		"sub":      user.ID,               // subject (user ID)
		"email":    user.Email,            // custom claim
		"username": user.Username,         // custom claim
		"exp":      expirationTime.Unix(), // expiration time
		"iat":      time.Now().Unix(),     // issued at time
	}

	// Create the token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with our secret key
	tokenString, err := token.SignedString(s.jwtSecret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ValidateToken verifies a JWT token and returns the claims
func (s *AuthService) ValidateToken(tokenString string) (jwt.MapClaims, error) {
	// Parse the token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return s.jwtSecret, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	// Extract and validate claims
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, ErrInvalidToken
}

// GetUserFromToken extracts user information from a valid token
func (s *AuthService) GetUserFromToken(tokenString string) (*db.User, error) {
	claims, err := s.ValidateToken(tokenString)
	if err != nil {
		return nil, err
	}

	userID, ok := claims["sub"].(string)
	if !ok {
		return nil, ErrInvalidToken
	}

	user, err := s.database.GetUserByID(userID)
	if err != nil {
		return nil, err
	}

	return user, nil
}
