package main

import (
	"log"
	"net/http"

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
*/

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

	// Start websocket message handler
	go sockets.HandleMessages()

	// Initialize and start reminder service
	reminderService := reminder.NewReminderService(database)
	reminderService.Start()
	log.Println("Reminder service started")

	// Create custom router for API endpoints
	router := handlers.CreateRouter()

	router.Handle("GET", "/habits", handlers.GetHabits)
	router.Handle("POST", "/habits", handlers.CreateHabit)
	router.Handle("GET", "/habits/:id", handlers.GetHabit)
	router.Handle("PATCH", "/habits/:id", handlers.UpdateHabit)
	router.Handle("DELETE", "/habits/:id", handlers.DeleteHabit)
	router.Handle("POST", "/habits/:id/tracking", handlers.CreateTracking)
	router.Handle("GET", "/habits/:id/tracking", handlers.GetTracking)
	router.Handle("PATCH", "/reminders/:id", handlers.UpdateReminder)

	// Create a main mux that handles both websockets and API routes
	mux := http.NewServeMux()
	mux.HandleFunc("/ws", sockets.WSHandler)
	mux.Handle("/", router)

	log.Println("Server is running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
