package auth

import (
	"log"
	"net/http"
	"os"
	"time"

	"habit-tracker/server/db"
)

// Example shows how to integrate authentication with your habit tracker
func Example() {
	// Initialize your database (SQLite or in-memory)
	database, err := db.NewDatabaseFromConfig()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Get JWT secret from environment variable
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "your-secret-key-change-in-production"
		log.Println("Warning: Using default JWT secret. Set JWT_SECRET environment variable for production.")
	}

	// Create auth service with 24-hour token expiry
	authService := NewAuthService(database, jwtSecret, 24*time.Hour)

	// Create HTTP mux
	mux := http.NewServeMux()

	// Public routes (no authentication required)
	mux.HandleFunc("/auth/register", authService.RegisterHandler)
	mux.HandleFunc("/auth/login", authService.LoginHandler)

	// Protected routes (authentication required)
	protectedMux := http.NewServeMux()
	protectedMux.HandleFunc("/auth/profile", authService.ProfileHandler)
	protectedMux.HandleFunc("/auth/validate", authService.ValidateTokenHandler)

	// Apply auth middleware to protected routes
	mux.Handle("/auth/profile", authService.AuthMiddleware(protectedMux))
	mux.Handle("/auth/validate", authService.AuthMiddleware(protectedMux))

	// Example of how to protect your existing habit routes
	// mux.Handle("/habits", authService.AuthMiddleware(http.HandlerFunc(yourHabitHandler)))

	log.Println("Auth service example running on :8080")
	log.Println("Try these endpoints:")
	log.Println("POST /auth/register - Register a new user")
	log.Println("POST /auth/login - Login and get JWT token")
	log.Println("GET /auth/profile - Get user profile (requires Bearer token)")
	log.Println("GET /auth/validate - Validate JWT token")

	log.Fatal(http.ListenAndServe(":8080", mux))
}

// Integration example showing how to get the current user in your handlers
func ExampleProtectedHandler(authService *AuthService) http.HandlerFunc {
	return authService.AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get user from context
		user, ok := GetUserFromRequest(r)
		if !ok {
			http.Error(w, "User not found in context", http.StatusInternalServerError)
			return
		}

		// Use user information
		log.Printf("Request from user: %s (%s)", user.Username, user.Email)

		// Your handler logic here
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"message": "Success", "user": "` + user.Username + `"}`))
	})).ServeHTTP
}
