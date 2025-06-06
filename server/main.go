package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"habit-tracker/server/auth"
	"habit-tracker/server/db"
	"habit-tracker/server/handlers"
	"habit-tracker/server/reminder"
	"habit-tracker/server/sockets"
)

/*
	Endpoints:
		GET /habits
		GET /habits/:id
		POST /habits
		PATCH /habits/:id
		DELETE /habits/:id
		POST /habits/:id/tracking
		GET /habits/:id/tracking
		WS /ws

	Authentication Endpoints:
		POST /auth/register
		POST /auth/login
		GET /auth/profile
		GET /auth/validate

	Statistics Endpoints:
		GET /habits/:id/stats
		GET /habits/:id/progress
		GET /stats/overview
		GET /stats/completion-rates
		GET /stats/daily-completions
*/

// Auth handler wrappers to adapt from http.HandlerFunc to handlers.HandlerFunc
func wrapAuthHandler(handler func(http.ResponseWriter, *http.Request)) handlers.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, params map[string]string) {
		handler(w, r)
	}
}

func wrapAuthMiddleware(authService *auth.AuthService, handler func(http.ResponseWriter, *http.Request)) handlers.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, params map[string]string) {
		// Create a wrapper handler that calls the auth middleware
		middlewareHandler := authService.AuthMiddleware(http.HandlerFunc(handler))
		middlewareHandler.ServeHTTP(w, r)
	}
}

func main() {
	database, err := db.NewDatabaseFromConfig()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	handlers.Database = database

	if err := database.Ping(); err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}
	log.Println("Database connection successful")

	// Initialize auth service
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "your-secret-key-change-in-production"
		log.Println("Warning: Using default JWT secret. Set JWT_SECRET environment variable for production.")
	}

	authService := auth.NewAuthService(database, jwtSecret, 24*time.Hour)

	go sockets.HandleMessages()

	reminderService := reminder.NewReminderService(database)
	reminderService.Start()
	log.Println("Reminder service started")

	router := handlers.CreateRouter()

	// Authentication routes (public)
	router.Handle("POST", "/auth/register", wrapAuthHandler(authService.RegisterHandler))
	router.Handle("POST", "/auth/login", wrapAuthHandler(authService.LoginHandler))

	// Authentication routes (protected)
	router.Handle("GET", "/auth/profile", wrapAuthMiddleware(authService, authService.ProfileHandler))
	router.Handle("GET", "/auth/validate", wrapAuthHandler(authService.ValidateTokenHandler))

	// Habit routes
	router.Handle("GET", "/habits", handlers.GetHabits)
	router.Handle("POST", "/habits", handlers.CreateHabit)
	router.Handle("GET", "/habits/:id", handlers.GetHabit)
	router.Handle("PATCH", "/habits/:id", handlers.UpdateHabit)
	router.Handle("DELETE", "/habits/:id", handlers.DeleteHabit)

	// Tracking routes
	router.Handle("POST", "/habits/:id/tracking", handlers.CreateTracking)
	router.Handle("GET", "/habits/:id/tracking", handlers.GetTracking)

	// Reminder routes
	router.Handle("PATCH", "/reminders/:id", handlers.UpdateReminder)

	// Statistics routes
	router.Handle("GET", "/habits/:id/stats", handlers.GetHabitStats)
	router.Handle("GET", "/habits/:id/progress", handlers.GetHabitProgress)
	router.Handle("GET", "/stats/overview", handlers.GetOverallStats)
	router.Handle("GET", "/stats/completion-rates", handlers.GetHabitCompletionRates)
	router.Handle("GET", "/stats/daily-completions", handlers.GetDailyCompletions)

	mux := http.NewServeMux()
	mux.HandleFunc("/ws", sockets.WSHandler)
	mux.Handle("/", router)

	log.Println("Server is running on port 8080")
	log.Println("Auth endpoints available:")
	log.Println("POST /auth/register - Register a new user")
	log.Println("POST /auth/login - Login and get JWT token")
	log.Println("GET /auth/profile - Get user profile (requires Bearer token)")
	log.Println("GET /auth/validate - Validate JWT token")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
