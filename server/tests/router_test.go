package tests

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"habit-tracker/server/handlers"

	"github.com/stretchr/testify/assert"
)

func TestCreateRouter(t *testing.T) {
	router := handlers.CreateRouter()
	assert.NotNil(t, router)
	// Note: we can't directly access router.routes anymore since it's unexported
	// But we can test the functionality instead
}

func TestRouterHandle(t *testing.T) {
	router := handlers.CreateRouter()
	testHandler := func(w http.ResponseWriter, r *http.Request, params map[string]string) {
		w.WriteHeader(http.StatusOK)
	}

	router.Handle("GET", "/test", testHandler)
	// We can't directly check routes slice, but we can test by making a request
	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestMatchFunction(t *testing.T) {
	tests := []struct {
		name           string
		pattern        string
		path           string
		expectedMatch  bool
		expectedParams map[string]string
	}{
		{
			name:           "exact match",
			pattern:        "/habits",
			path:           "/habits",
			expectedMatch:  true,
			expectedParams: map[string]string{},
		},
		{
			name:           "no match - different paths",
			pattern:        "/habits",
			path:           "/users",
			expectedMatch:  false,
			expectedParams: nil,
		},
		{
			name:           "single parameter match",
			pattern:        "/habits/:id",
			path:           "/habits/123",
			expectedMatch:  true,
			expectedParams: map[string]string{"id": "123"},
		},
		{
			name:           "multiple parameters match",
			pattern:        "/habits/:id/tracking/:entryId",
			path:           "/habits/123/tracking/456",
			expectedMatch:  true,
			expectedParams: map[string]string{"id": "123", "entryId": "456"},
		},
		{
			name:           "parameter mismatch - wrong length",
			pattern:        "/habits/:id",
			path:           "/habits/123/extra",
			expectedMatch:  false,
			expectedParams: nil,
		},
		{
			name:           "parameter mismatch - wrong path",
			pattern:        "/habits/:id/tracking",
			path:           "/habits/123/invalid",
			expectedMatch:  false,
			expectedParams: nil,
		},
		{
			name:           "root path",
			pattern:        "/",
			path:           "/",
			expectedMatch:  true,
			expectedParams: map[string]string{},
		},
		{
			name:           "empty paths",
			pattern:        "",
			path:           "",
			expectedMatch:  true,
			expectedParams: map[string]string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			params, matched := handlers.Match(tt.pattern, tt.path)
			assert.Equal(t, tt.expectedMatch, matched)
			if tt.expectedMatch {
				assert.Equal(t, tt.expectedParams, params)
			} else {
				assert.Nil(t, params)
			}
		})
	}
}

func TestRouterServeHTTP(t *testing.T) {
	router := handlers.CreateRouter()

	// Register test handlers
	router.Handle("GET", "/test", func(w http.ResponseWriter, r *http.Request, params map[string]string) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("GET test"))
	})

	router.Handle("POST", "/test", func(w http.ResponseWriter, r *http.Request, params map[string]string) {
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("POST test"))
	})

	router.Handle("GET", "/test/:id", func(w http.ResponseWriter, r *http.Request, params map[string]string) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ID: " + params["id"]))
	})

	tests := []struct {
		name           string
		method         string
		path           string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "GET request to /test",
			method:         "GET",
			path:           "/test",
			expectedStatus: http.StatusOK,
			expectedBody:   "GET test",
		},
		{
			name:           "POST request to /test",
			method:         "POST",
			path:           "/test",
			expectedStatus: http.StatusCreated,
			expectedBody:   "POST test",
		},
		{
			name:           "GET request with parameter",
			method:         "GET",
			path:           "/test/123",
			expectedStatus: http.StatusOK,
			expectedBody:   "ID: 123",
		},
		{
			name:           "Not found",
			method:         "GET",
			path:           "/nonexistent",
			expectedStatus: http.StatusNotFound,
			expectedBody:   "404 page not found\n",
		},
		{
			name:           "Method not allowed",
			method:         "DELETE",
			path:           "/test",
			expectedStatus: http.StatusNotFound,
			expectedBody:   "404 page not found\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.path, nil)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			assert.Equal(t, tt.expectedBody, w.Body.String())
		})
	}
}

func TestCORSHeaders(t *testing.T) {
	router := handlers.CreateRouter()
	router.Handle("GET", "/test", func(w http.ResponseWriter, r *http.Request, params map[string]string) {
		w.WriteHeader(http.StatusOK)
	})

	// Test CORS headers on regular request
	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, "*", w.Header().Get("Access-Control-Allow-Origin"))
	assert.Equal(t, "GET, POST, PUT, DELETE, OPTIONS", w.Header().Get("Access-Control-Allow-Methods"))
	assert.Equal(t, "Content-Type, Authorization", w.Header().Get("Access-Control-Allow-Headers"))
	assert.Equal(t, "86400", w.Header().Get("Access-Control-Max-Age"))
}

func TestOPTIONSRequest(t *testing.T) {
	router := handlers.CreateRouter()
	router.Handle("GET", "/test", func(w http.ResponseWriter, r *http.Request, params map[string]string) {
		w.WriteHeader(http.StatusOK)
	})

	// Test preflight OPTIONS request
	req := httptest.NewRequest("OPTIONS", "/test", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "*", w.Header().Get("Access-Control-Allow-Origin"))
	assert.Equal(t, "GET, POST, PUT, DELETE, OPTIONS", w.Header().Get("Access-Control-Allow-Methods"))
	assert.Equal(t, "Content-Type, Authorization", w.Header().Get("Access-Control-Allow-Headers"))
	assert.Equal(t, "86400", w.Header().Get("Access-Control-Max-Age"))
}
