package auth_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"habit-tracker/server/auth"
	"habit-tracker/server/db"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type MiddlewareTestSuite struct {
	suite.Suite
	authService *auth.AuthService
	database    *db.MapDatabase
}

func (suite *MiddlewareTestSuite) SetupTest() {
	suite.database = db.NewMapDatabase()
	suite.authService = auth.NewAuthService(suite.database, "test-secret", time.Hour)
}

func (suite *MiddlewareTestSuite) TestAuthMiddlewareWithValidToken() {
	// Register and login user
	email := "test@example.com"
	username := "testuser"
	password := "password123"

	_, err := suite.authService.Register(email, username, password)
	suite.NoError(err)

	token, _, err := suite.authService.Login(email, password)
	suite.NoError(err)

	// Create test handler
	called := false
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		user := auth.GetUserFromContext(r.Context())
		suite.NotNil(user)
		suite.Equal(email, user.Email)
		suite.Equal(username, user.Username)
		w.WriteHeader(http.StatusOK)
	})

	// Create middleware
	middleware := auth.AuthMiddleware(suite.authService)
	protectedHandler := middleware(testHandler)

	// Create request with valid token
	req := httptest.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rr := httptest.NewRecorder()

	protectedHandler.ServeHTTP(rr, req)

	suite.Equal(http.StatusOK, rr.Code)
	suite.True(called)
}

func (suite *MiddlewareTestSuite) TestAuthMiddlewareWithInvalidToken() {
	// Create test handler
	called := false
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
	})

	// Create middleware
	middleware := auth.AuthMiddleware(suite.authService)
	protectedHandler := middleware(testHandler)

	// Create request with invalid token
	req := httptest.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer invalid.token.here")
	rr := httptest.NewRecorder()

	protectedHandler.ServeHTTP(rr, req)

	suite.Equal(http.StatusUnauthorized, rr.Code)
	suite.False(called)
}

func (suite *MiddlewareTestSuite) TestAuthMiddlewareWithMissingToken() {
	// Create test handler
	called := false
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
	})

	// Create middleware
	middleware := auth.AuthMiddleware(suite.authService)
	protectedHandler := middleware(testHandler)

	// Create request without token
	req := httptest.NewRequest("GET", "/protected", nil)
	rr := httptest.NewRecorder()

	protectedHandler.ServeHTTP(rr, req)

	suite.Equal(http.StatusUnauthorized, rr.Code)
	suite.False(called)
}

func (suite *MiddlewareTestSuite) TestAuthMiddlewareWithMalformedHeader() {
	// Create test handler
	called := false
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
	})

	// Create middleware
	middleware := auth.AuthMiddleware(suite.authService)
	protectedHandler := middleware(testHandler)

	// Create request with malformed authorization header
	req := httptest.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "InvalidFormat")
	rr := httptest.NewRecorder()

	protectedHandler.ServeHTTP(rr, req)

	suite.Equal(http.StatusUnauthorized, rr.Code)
	suite.False(called)
}

func (suite *MiddlewareTestSuite) TestOptionalAuthMiddlewareWithValidToken() {
	// Register and login user
	email := "test@example.com"
	username := "testuser"
	password := "password123"

	_, err := suite.authService.Register(email, username, password)
	suite.NoError(err)

	token, _, err := suite.authService.Login(email, password)
	suite.NoError(err)

	// Create test handler
	called := false
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		user := auth.GetUserFromContext(r.Context())
		suite.NotNil(user)
		suite.Equal(email, user.Email)
		suite.Equal(username, user.Username)
		w.WriteHeader(http.StatusOK)
	})

	// Create middleware
	middleware := auth.OptionalAuthMiddleware(suite.authService)
	handler := middleware(testHandler)

	// Create request with valid token
	req := httptest.NewRequest("GET", "/optional", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	suite.Equal(http.StatusOK, rr.Code)
	suite.True(called)
}

func (suite *MiddlewareTestSuite) TestOptionalAuthMiddlewareWithInvalidToken() {
	// Create test handler
	called := false
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		user := auth.GetUserFromContext(r.Context())
		suite.Nil(user) // Should be nil for invalid token
		w.WriteHeader(http.StatusOK)
	})

	// Create middleware
	middleware := auth.OptionalAuthMiddleware(suite.authService)
	handler := middleware(testHandler)

	// Create request with invalid token
	req := httptest.NewRequest("GET", "/optional", nil)
	req.Header.Set("Authorization", "Bearer invalid.token.here")
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	suite.Equal(http.StatusOK, rr.Code)
	suite.True(called)
}

func (suite *MiddlewareTestSuite) TestOptionalAuthMiddlewareWithoutToken() {
	// Create test handler
	called := false
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		user := auth.GetUserFromContext(r.Context())
		suite.Nil(user) // Should be nil when no token provided
		w.WriteHeader(http.StatusOK)
	})

	// Create middleware
	middleware := auth.OptionalAuthMiddleware(suite.authService)
	handler := middleware(testHandler)

	// Create request without token
	req := httptest.NewRequest("GET", "/optional", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	suite.Equal(http.StatusOK, rr.Code)
	suite.True(called)
}

func (suite *MiddlewareTestSuite) TestGetUserFromContextWithUser() {
	email := "test@example.com"
	username := "testuser"
	password := "password123"

	user, err := suite.authService.Register(email, username, password)
	suite.NoError(err)

	// Create context with user
	req := httptest.NewRequest("GET", "/test", nil)
	ctx := auth.SetUserInContext(req.Context(), user)
	req = req.WithContext(ctx)

	// Get user from context
	retrievedUser := auth.GetUserFromContext(req.Context())
	suite.NotNil(retrievedUser)
	suite.Equal(user.ID, retrievedUser.ID)
	suite.Equal(user.Email, retrievedUser.Email)
	suite.Equal(user.Username, retrievedUser.Username)
}

func (suite *MiddlewareTestSuite) TestGetUserFromContextWithoutUser() {
	req := httptest.NewRequest("GET", "/test", nil)

	// Get user from context without setting one
	user := auth.GetUserFromContext(req.Context())
	suite.Nil(user)
}

func (suite *MiddlewareTestSuite) TestGetUserIDFromContextWithUser() {
	email := "test@example.com"
	username := "testuser"
	password := "password123"

	user, err := suite.authService.Register(email, username, password)
	suite.NoError(err)

	// Create context with user
	req := httptest.NewRequest("GET", "/test", nil)
	ctx := auth.SetUserInContext(req.Context(), user)
	req = req.WithContext(ctx)

	// Get user ID from context
	userID := auth.GetUserIDFromContext(req.Context())
	suite.Equal(user.ID, userID)
}

func (suite *MiddlewareTestSuite) TestGetUserIDFromContextWithoutUser() {
	req := httptest.NewRequest("GET", "/test", nil)

	// Get user ID from context without setting one
	userID := auth.GetUserIDFromContext(req.Context())
	suite.Empty(userID)
}

func (suite *MiddlewareTestSuite) TestAuthMiddlewareWithExpiredToken() {
	// Create auth service with very short expiry
	shortExpiryService := auth.NewAuthService(suite.database, "test-secret", time.Millisecond)

	email := "test@example.com"
	username := "testuser"
	password := "password123"

	_, err := shortExpiryService.Register(email, username, password)
	suite.NoError(err)

	token, _, err := shortExpiryService.Login(email, password)
	suite.NoError(err)

	// Wait for token to expire
	time.Sleep(10 * time.Millisecond)

	// Create test handler
	called := false
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
	})

	// Create middleware
	middleware := auth.AuthMiddleware(shortExpiryService)
	protectedHandler := middleware(testHandler)

	// Create request with expired token
	req := httptest.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rr := httptest.NewRecorder()

	protectedHandler.ServeHTTP(rr, req)

	suite.Equal(http.StatusUnauthorized, rr.Code)
	suite.False(called)
}

// Run the test suite
func TestMiddlewareTestSuite(t *testing.T) {
	suite.Run(t, new(MiddlewareTestSuite))
}

// Additional standalone tests
func TestGetUserFromContextNilContext(t *testing.T) {
	user := auth.GetUserFromContext(nil)
	assert.Nil(t, user)
}

func TestGetUserIDFromContextNilContext(t *testing.T) {
	userID := auth.GetUserIDFromContext(nil)
	assert.Empty(t, userID)
}
