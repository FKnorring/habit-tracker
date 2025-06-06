package auth_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"habit-tracker/server/auth"
	"habit-tracker/server/db"

	"github.com/stretchr/testify/suite"
)

type HandlersTestSuite struct {
	suite.Suite
	authService *auth.AuthService
	database    *db.MapDatabase
}

func (suite *HandlersTestSuite) SetupTest() {
	suite.database = db.NewMapDatabase()
	suite.authService = auth.NewAuthService(suite.database, "test-secret", time.Hour)
}

func (suite *HandlersTestSuite) TestRegisterHandlerSuccess() {
	registerData := auth.RegisterRequest{
		Email:    "test@example.com",
		Username: "testuser",
		Password: "password123",
	}

	body, err := json.Marshal(registerData)
	suite.NoError(err)

	req := httptest.NewRequest("POST", "/auth/register", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	handler := auth.RegisterHandler(suite.authService)
	handler.ServeHTTP(rr, req)

	suite.Equal(http.StatusCreated, rr.Code)

	var response auth.RegisterResponse
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	suite.NoError(err)
	suite.Equal("User registered successfully", response.Message)
	suite.NotEmpty(response.User.ID)
	suite.Equal(registerData.Email, response.User.Email)
	suite.Equal(registerData.Username, response.User.Username)
}

func (suite *HandlersTestSuite) TestRegisterHandlerInvalidJSON() {
	req := httptest.NewRequest("POST", "/auth/register", bytes.NewReader([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	handler := auth.RegisterHandler(suite.authService)
	handler.ServeHTTP(rr, req)

	suite.Equal(http.StatusBadRequest, rr.Code)

	var response auth.ErrorResponse
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	suite.NoError(err)
	suite.Contains(response.Error, "Invalid JSON")
}

func (suite *HandlersTestSuite) TestRegisterHandlerMissingFields() {
	registerData := auth.RegisterRequest{
		Email: "test@example.com",
		// Missing Username and Password
	}

	body, err := json.Marshal(registerData)
	suite.NoError(err)

	req := httptest.NewRequest("POST", "/auth/register", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	handler := auth.RegisterHandler(suite.authService)
	handler.ServeHTTP(rr, req)

	suite.Equal(http.StatusBadRequest, rr.Code)

	var response auth.ErrorResponse
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	suite.NoError(err)
	suite.Contains(response.Error, "required")
}

func (suite *HandlersTestSuite) TestRegisterHandlerDuplicateEmail() {
	// Register first user
	_, err := suite.authService.Register("test@example.com", "user1", "password123")
	suite.NoError(err)

	// Try to register another user with same email
	registerData := auth.RegisterRequest{
		Email:    "test@example.com",
		Username: "user2",
		Password: "password123",
	}

	body, err := json.Marshal(registerData)
	suite.NoError(err)

	req := httptest.NewRequest("POST", "/auth/register", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	handler := auth.RegisterHandler(suite.authService)
	handler.ServeHTTP(rr, req)

	suite.Equal(http.StatusConflict, rr.Code)

	var response auth.ErrorResponse
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	suite.NoError(err)
	suite.Contains(response.Error, "already exists")
}

func (suite *HandlersTestSuite) TestLoginHandlerSuccess() {
	// Register user first
	email := "test@example.com"
	username := "testuser"
	password := "password123"

	_, err := suite.authService.Register(email, username, password)
	suite.NoError(err)

	loginData := auth.LoginRequest{
		Email:    email,
		Password: password,
	}

	body, err := json.Marshal(loginData)
	suite.NoError(err)

	req := httptest.NewRequest("POST", "/auth/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	handler := auth.LoginHandler(suite.authService)
	handler.ServeHTTP(rr, req)

	suite.Equal(http.StatusOK, rr.Code)

	var response auth.LoginResponse
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	suite.NoError(err)
	suite.Equal("Login successful", response.Message)
	suite.NotEmpty(response.Token)
	suite.Equal(email, response.User.Email)
	suite.Equal(username, response.User.Username)
}

func (suite *HandlersTestSuite) TestLoginHandlerInvalidCredentials() {
	loginData := auth.LoginRequest{
		Email:    "nonexistent@example.com",
		Password: "password123",
	}

	body, err := json.Marshal(loginData)
	suite.NoError(err)

	req := httptest.NewRequest("POST", "/auth/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	handler := auth.LoginHandler(suite.authService)
	handler.ServeHTTP(rr, req)

	suite.Equal(http.StatusUnauthorized, rr.Code)

	var response auth.ErrorResponse
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	suite.NoError(err)
	suite.Contains(response.Error, "Invalid credentials")
}

func (suite *HandlersTestSuite) TestLoginHandlerInvalidJSON() {
	req := httptest.NewRequest("POST", "/auth/login", bytes.NewReader([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	handler := auth.LoginHandler(suite.authService)
	handler.ServeHTTP(rr, req)

	suite.Equal(http.StatusBadRequest, rr.Code)

	var response auth.ErrorResponse
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	suite.NoError(err)
	suite.Contains(response.Error, "Invalid JSON")
}

func (suite *HandlersTestSuite) TestLoginHandlerMissingFields() {
	loginData := auth.LoginRequest{
		Email: "test@example.com",
		// Missing Password
	}

	body, err := json.Marshal(loginData)
	suite.NoError(err)

	req := httptest.NewRequest("POST", "/auth/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	handler := auth.LoginHandler(suite.authService)
	handler.ServeHTTP(rr, req)

	suite.Equal(http.StatusBadRequest, rr.Code)

	var response auth.ErrorResponse
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	suite.NoError(err)
	suite.Contains(response.Error, "required")
}

func (suite *HandlersTestSuite) TestProfileHandlerSuccess() {
	// Register and login user
	email := "test@example.com"
	username := "testuser"
	password := "password123"

	user, err := suite.authService.Register(email, username, password)
	suite.NoError(err)

	req := httptest.NewRequest("GET", "/auth/profile", nil)

	// Set user in context (simulating middleware)
	ctx := auth.SetUserInContext(req.Context(), user)
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()

	handler := auth.ProfileHandler(suite.authService)
	handler.ServeHTTP(rr, req)

	suite.Equal(http.StatusOK, rr.Code)

	var response auth.ProfileResponse
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	suite.NoError(err)
	suite.Equal(email, response.User.Email)
	suite.Equal(username, response.User.Username)
}

func (suite *HandlersTestSuite) TestProfileHandlerNoUser() {
	req := httptest.NewRequest("GET", "/auth/profile", nil)
	rr := httptest.NewRecorder()

	handler := auth.ProfileHandler(suite.authService)
	handler.ServeHTTP(rr, req)

	suite.Equal(http.StatusUnauthorized, rr.Code)

	var response auth.ErrorResponse
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	suite.NoError(err)
	suite.Contains(response.Error, "User not found")
}

func (suite *HandlersTestSuite) TestValidateHandlerSuccess() {
	// Register and login user
	email := "test@example.com"
	username := "testuser"
	password := "password123"

	_, err := suite.authService.Register(email, username, password)
	suite.NoError(err)

	token, _, err := suite.authService.Login(email, password)
	suite.NoError(err)

	req := httptest.NewRequest("GET", "/auth/validate", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rr := httptest.NewRecorder()

	handler := auth.ValidateHandler(suite.authService)
	handler.ServeHTTP(rr, req)

	suite.Equal(http.StatusOK, rr.Code)

	var response auth.ValidateResponse
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	suite.NoError(err)
	suite.True(response.Valid)
	suite.Equal(email, response.User.Email)
	suite.Equal(username, response.User.Username)
}

func (suite *HandlersTestSuite) TestValidateHandlerInvalidToken() {
	req := httptest.NewRequest("GET", "/auth/validate", nil)
	req.Header.Set("Authorization", "Bearer invalid.token.here")
	rr := httptest.NewRecorder()

	handler := auth.ValidateHandler(suite.authService)
	handler.ServeHTTP(rr, req)

	suite.Equal(http.StatusUnauthorized, rr.Code)

	var response auth.ValidateResponse
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	suite.NoError(err)
	suite.False(response.Valid)
	suite.Nil(response.User)
}

func (suite *HandlersTestSuite) TestValidateHandlerMissingToken() {
	req := httptest.NewRequest("GET", "/auth/validate", nil)
	rr := httptest.NewRecorder()

	handler := auth.ValidateHandler(suite.authService)
	handler.ServeHTTP(rr, req)

	suite.Equal(http.StatusUnauthorized, rr.Code)

	var response auth.ValidateResponse
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	suite.NoError(err)
	suite.False(response.Valid)
	suite.Nil(response.User)
}

func (suite *HandlersTestSuite) TestValidateHandlerExpiredToken() {
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

	req := httptest.NewRequest("GET", "/auth/validate", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rr := httptest.NewRecorder()

	handler := auth.ValidateHandler(shortExpiryService)
	handler.ServeHTTP(rr, req)

	suite.Equal(http.StatusUnauthorized, rr.Code)

	var response auth.ValidateResponse
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	suite.NoError(err)
	suite.False(response.Valid)
	suite.Nil(response.User)
}

func (suite *HandlersTestSuite) TestHandlersWithWrongHTTPMethods() {
	// Test GET request to register endpoint
	req := httptest.NewRequest("GET", "/auth/register", nil)
	rr := httptest.NewRecorder()

	handler := auth.RegisterHandler(suite.authService)
	handler.ServeHTTP(rr, req)

	suite.Equal(http.StatusMethodNotAllowed, rr.Code)

	// Test POST request to profile endpoint
	req = httptest.NewRequest("POST", "/auth/profile", nil)
	rr = httptest.NewRecorder()

	handler = auth.ProfileHandler(suite.authService)
	handler.ServeHTTP(rr, req)

	suite.Equal(http.StatusMethodNotAllowed, rr.Code)

	// Test POST request to validate endpoint
	req = httptest.NewRequest("POST", "/auth/validate", nil)
	rr = httptest.NewRecorder()

	handler = auth.ValidateHandler(suite.authService)
	handler.ServeHTTP(rr, req)

	suite.Equal(http.StatusMethodNotAllowed, rr.Code)
}

// Run the test suite
func TestHandlersTestSuite(t *testing.T) {
	suite.Run(t, new(HandlersTestSuite))
}
