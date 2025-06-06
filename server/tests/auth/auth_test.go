package auth_test

import (
	"testing"
	"time"

	"habit-tracker/server/auth"
	"habit-tracker/server/db"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type AuthTestSuite struct {
	suite.Suite
	authService *auth.AuthService
	database    *db.MapDatabase
}

func (suite *AuthTestSuite) SetupTest() {
	suite.database = db.NewMapDatabase()
	suite.authService = auth.NewAuthService(suite.database, "test-secret", time.Hour)
}

func (suite *AuthTestSuite) TestNewAuthService() {
	authService := auth.NewAuthService(suite.database, "test-secret", time.Hour)
	suite.NotNil(authService)
}

func (suite *AuthTestSuite) TestRegister() {
	email := "test@example.com"
	username := "testuser"
	password := "password123"

	user, err := suite.authService.Register(email, username, password)
	suite.NoError(err)
	suite.NotNil(user)
	suite.Equal(email, user.Email)
	suite.Equal(username, user.Username)
	suite.NotEmpty(user.ID)
	suite.NotEqual(password, user.PasswordHash) // Password should be hashed
	suite.NotZero(user.CreatedAt)
	suite.NotZero(user.UpdatedAt)
}

func (suite *AuthTestSuite) TestRegisterDuplicateEmail() {
	email := "test@example.com"
	username1 := "testuser1"
	username2 := "testuser2"
	password := "password123"

	// Register first user
	_, err := suite.authService.Register(email, username1, password)
	suite.NoError(err)

	// Try to register another user with same email
	_, err = suite.authService.Register(email, username2, password)
	suite.Error(err)
	suite.Contains(err.Error(), "email already in use")
}

func (suite *AuthTestSuite) TestLoginSuccess() {
	email := "test@example.com"
	username := "testuser"
	password := "password123"

	// Register user first
	_, err := suite.authService.Register(email, username, password)
	suite.NoError(err)

	// Login with correct credentials
	token, user, err := suite.authService.Login(email, password)
	suite.NoError(err)
	suite.NotEmpty(token)
	suite.NotNil(user)
	suite.Equal(email, user.Email)
	suite.Equal(username, user.Username)
}

func (suite *AuthTestSuite) TestLoginInvalidEmail() {
	email := "nonexistent@example.com"
	password := "password123"

	token, user, err := suite.authService.Login(email, password)
	suite.Error(err)
	suite.Empty(token)
	suite.Nil(user)
	suite.Contains(err.Error(), "invalid credentials")
}

func (suite *AuthTestSuite) TestLoginInvalidPassword() {
	email := "test@example.com"
	username := "testuser"
	password := "password123"
	wrongPassword := "wrongpassword"

	// Register user first
	_, err := suite.authService.Register(email, username, password)
	suite.NoError(err)

	// Login with wrong password
	token, user, err := suite.authService.Login(email, wrongPassword)
	suite.Error(err)
	suite.Empty(token)
	suite.Nil(user)
	suite.Contains(err.Error(), "invalid credentials")
}

func (suite *AuthTestSuite) TestValidateTokenSuccess() {
	email := "test@example.com"
	username := "testuser"
	password := "password123"

	// Register and login user
	_, err := suite.authService.Register(email, username, password)
	suite.NoError(err)

	token, _, err := suite.authService.Login(email, password)
	suite.NoError(err)

	// Validate token
	user, err := suite.authService.GetUserFromToken(token)
	suite.NoError(err)
	suite.NotNil(user)
	suite.Equal(email, user.Email)
	suite.Equal(username, user.Username)
}

func (suite *AuthTestSuite) TestValidateTokenInvalid() {
	invalidToken := "invalid.token.here"

	user, err := suite.authService.GetUserFromToken(invalidToken)
	suite.Error(err)
	suite.Nil(user)
}

func (suite *AuthTestSuite) TestValidateTokenExpired() {
	// Create auth service with very short expiry
	shortExpiryService := auth.NewAuthService(suite.database, "test-secret", time.Millisecond)

	email := "test@example.com"
	username := "testuser"
	password := "password123"

	// Register and login user
	_, err := shortExpiryService.Register(email, username, password)
	suite.NoError(err)

	token, _, err := shortExpiryService.Login(email, password)
	suite.NoError(err)

	// Wait for token to expire
	time.Sleep(10 * time.Millisecond)

	// Validate expired token
	user, err := shortExpiryService.GetUserFromToken(token)
	suite.Error(err)
	suite.Nil(user)
	suite.Contains(err.Error(), "token is expired")
}

func (suite *AuthTestSuite) TestHashPasswordAndVerify() {
	password := "mySecretPassword123"

	hashedPassword, err := suite.authService.HashPassword(password)
	suite.NoError(err)
	suite.NotEmpty(hashedPassword)
	suite.NotEqual(password, hashedPassword)

	// Verify correct password
	err = suite.authService.VerifyPassword(hashedPassword, password)
	suite.NoError(err)

	// Verify incorrect password
	err = suite.authService.VerifyPassword(hashedPassword, "wrongPassword")
	suite.Error(err)
}

func (suite *AuthTestSuite) TestPasswordHashing() {
	password := "testPassword123"

	// Hash the same password multiple times
	hash1, err := suite.authService.HashPassword(password)
	suite.NoError(err)

	hash2, err := suite.authService.HashPassword(password)
	suite.NoError(err)

	// Hashes should be different (bcrypt uses salt)
	suite.NotEqual(hash1, hash2)

	// But both should verify against the original password
	suite.NoError(suite.authService.VerifyPassword(hash1, password))
	suite.NoError(suite.authService.VerifyPassword(hash2, password))
}

func (suite *AuthTestSuite) TestTokenContainsClaims() {
	email := "test@example.com"
	username := "testuser"
	password := "password123"

	// Register and login user
	user, err := suite.authService.Register(email, username, password)
	suite.NoError(err)

	token, _, err := suite.authService.Login(email, password)
	suite.NoError(err)

	// Validate token and check if user data matches
	validatedUser, err := suite.authService.GetUserFromToken(token)
	suite.NoError(err)
	suite.Equal(user.ID, validatedUser.ID)
	suite.Equal(user.Email, validatedUser.Email)
	suite.Equal(user.Username, validatedUser.Username)
}

func (suite *AuthTestSuite) TestRegisterEmptyFields() {
	// Test empty email
	_, err := suite.authService.Register("", "username", "password")
	suite.Error(err)

	// Test empty username
	_, err = suite.authService.Register("test@example.com", "", "password")
	suite.Error(err)

	// Test empty password
	_, err = suite.authService.Register("test@example.com", "username", "")
	suite.Error(err)
}

func (suite *AuthTestSuite) TestLoginEmptyFields() {
	// Test empty email
	_, _, err := suite.authService.Login("", "password")
	suite.Error(err)

	// Test empty password
	_, _, err = suite.authService.Login("test@example.com", "")
	suite.Error(err)
}

// Run the test suite
func TestAuthTestSuite(t *testing.T) {
	suite.Run(t, new(AuthTestSuite))
}

// Additional standalone tests
func TestAuthServiceWithNilDatabase(t *testing.T) {
	assert.Panics(t, func() {
		auth.NewAuthService(nil, "secret", time.Hour)
	})
}

func TestAuthServiceWithEmptySecret(t *testing.T) {
	db := db.NewMapDatabase()
	assert.Panics(t, func() {
		auth.NewAuthService(db, "", time.Hour)
	})
}

func TestAuthServiceWithZeroExpiry(t *testing.T) {
	db := db.NewMapDatabase()
	assert.Panics(t, func() {
		auth.NewAuthService(db, "secret", 0)
	})
}
