package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"habit-tracker/server/db"

	"github.com/google/uuid"
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

var database db.Database

func checkParams(w http.ResponseWriter, params map[string]string, requiredParams []string) bool {
	for _, param := range requiredParams {
		if _, ok := params[param]; !ok {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Missing required parameters"))
			return false
		}
	}
	return true
}

func getHabits(w http.ResponseWriter, r *http.Request, params map[string]string) {
	habits, err := database.GetAllHabits()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to retrieve habits"))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(habits)
}

func createHabit(w http.ResponseWriter, r *http.Request, params map[string]string) {
	var habit db.Habit
	if err := json.NewDecoder(r.Body).Decode(&habit); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid JSON"))
		return
	}

	if habit.ID == "" {
		habit.ID = uuid.New().String()
	}

	if err := database.CreateHabit(&habit); err != nil {
		if err == db.ErrDuplicate {
			w.WriteHeader(http.StatusConflict)
			w.Write([]byte("Habit already exists"))
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Failed to create habit"))
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(habit)
}

func getHabit(w http.ResponseWriter, r *http.Request, params map[string]string) {
	if !checkParams(w, params, []string{"id"}) {
		return
	}

	habit, err := database.GetHabit(params["id"])
	if err != nil {
		if err == db.ErrNotFound {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("Habit not found"))
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Failed to retrieve habit"))
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(habit)
}

func updateHabit(w http.ResponseWriter, r *http.Request, params map[string]string) {
	if !checkParams(w, params, []string{"id"}) {
		return
	}

	var habit db.Habit
	if err := json.NewDecoder(r.Body).Decode(&habit); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid JSON"))
		return
	}

	habit.ID = params["id"]

	if err := database.UpdateHabit(&habit); err != nil {
		if err == db.ErrNotFound {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("Habit not found"))
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Failed to update habit"))
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(habit)
}

func deleteHabit(w http.ResponseWriter, r *http.Request, params map[string]string) {
	if !checkParams(w, params, []string{"id"}) {
		return
	}

	if err := database.DeleteHabit(params["id"]); err != nil {
		if err == db.ErrNotFound {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("Habit not found"))
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Failed to delete habit"))
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func createTracking(w http.ResponseWriter, r *http.Request, params map[string]string) {
	if !checkParams(w, params, []string{"id"}) {
		return
	}

	var entry db.TrackingEntry
	if err := json.NewDecoder(r.Body).Decode(&entry); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid JSON"))
		return
	}

	if entry.ID == "" {
		entry.ID = uuid.New().String()
	}

	entry.HabitID = params["id"]

	if entry.Timestamp == "" {
		entry.Timestamp = time.Now().Format(time.RFC3339)
	}

	if err := database.CreateTrackingEntry(&entry); err != nil {
		if err == db.ErrDuplicate {
			w.WriteHeader(http.StatusConflict)
			w.Write([]byte("Tracking entry already exists"))
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Failed to create tracking entry"))
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(entry)
}

func getTracking(w http.ResponseWriter, r *http.Request, params map[string]string) {
	if !checkParams(w, params, []string{"id"}) {
		return
	}

	entries, err := database.GetTrackingEntriesByHabitID(params["id"])
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to retrieve tracking entries"))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(entries)
}

func main() {
	// Initialize database from configuration
	var err error
	database, err = db.NewDatabaseFromConfig()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Check database connection health
	if err := database.Ping(); err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}
	log.Println("Database connection successful")

	router := CreateRouter()

	router.Handle("GET", "/habits", getHabits)
	router.Handle("POST", "/habits", createHabit)
	router.Handle("GET", "/habits/:id", getHabit)
	router.Handle("PUT", "/habits/:id", updateHabit)
	router.Handle("DELETE", "/habits/:id", deleteHabit)
	router.Handle("POST", "/habits/:id/tracking", createTracking)
	router.Handle("GET", "/habits/:id/tracking", getTracking)

	log.Println("Server is running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
