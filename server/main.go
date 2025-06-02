package main

import (
	"log"
	"net/http"

	"habit-tracker/server/db"
	"habit-tracker/server/handlers"
)

/*
	Endpoints:
		GET /habits
		GET /habits/:id
		POST /habits
		PUT /habits/:id
		DELETE /habits/:id
		POST /habits/:id/tracking
		GET /habits/:id/tracking
*/

func main() {
	// Initialize database from configuration
	database, err := db.NewDatabaseFromConfig()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Set the database in the handlers package
	handlers.Database = database

	// Check database connection health
	if err := database.Ping(); err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}
	log.Println("Database connection successful")

	router := handlers.CreateRouter()

	router.Handle("GET", "/habits", handlers.GetHabits)
	router.Handle("POST", "/habits", handlers.CreateHabit)
	router.Handle("GET", "/habits/:id", handlers.GetHabit)
	router.Handle("PUT", "/habits/:id", handlers.UpdateHabit)
	router.Handle("DELETE", "/habits/:id", handlers.DeleteHabit)
	router.Handle("POST", "/habits/:id/tracking", handlers.CreateTracking)
	router.Handle("GET", "/habits/:id/tracking", handlers.GetTracking)

	log.Println("Server is running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
