package main

import (
	"encoding/json"
	"log"
	"net/http"
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

var habits = []Habit{
	{ID: "1", Name: "Habit 1", Description: "Description 1", Frequency: "Daily", StartDate: "2021-01-01"},
	{ID: "2", Name: "Habit 2", Description: "Description 2", Frequency: "Weekly", StartDate: "2021-01-01"},
}

var trackings = []TrackingEntry{
	{ID: "1", HabitID: "1", Timestamp: "2021-01-01", Note: "Note 1"},
	{ID: "2", HabitID: "2", Timestamp: "2021-01-01", Note: "Note 2"},
}

func getHabits(w http.ResponseWriter, r *http.Request, params map[string]string) {
	// get habits stub
	json.NewEncoder(w).Encode(habits)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Habits fetched successfully"))
	return
}

func createHabit(w http.ResponseWriter, r *http.Request, params map[string]string) {
	// create habit stub
	habit := Habit{ID: "3", Name: "Habit 3", Description: "Description 3", Frequency: "Monthly", StartDate: "2021-01-01"}
	json.NewEncoder(w).Encode(habit)
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Habit created successfully"))
	return
}

func getHabit(w http.ResponseWriter, r *http.Request, params map[string]string) {
	// get habit stub
	habit := habits[0]
	json.NewEncoder(w).Encode(habit)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Habit fetched successfully"))
	return
}

func updateHabit(w http.ResponseWriter, r *http.Request, params map[string]string) {
	// update habit stub
	habit := habits[0]
	json.NewEncoder(w).Encode(habit)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Habit updated successfully"))
	return
}

func deleteHabit(w http.ResponseWriter, r *http.Request, params map[string]string) {
	// delete habit stub
	json.NewEncoder(w).Encode(nil)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Habit deleted successfully"))
	return
}

func createTracking(w http.ResponseWriter, r *http.Request, params map[string]string) {
	// create tracking stub
	tracking := trackings[0]
	json.NewEncoder(w).Encode(tracking)
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Tracking created successfully"))
	return
}

func getTracking(w http.ResponseWriter, r *http.Request, params map[string]string) {
	// get tracking stub
	tracking := trackings[0]
	json.NewEncoder(w).Encode(tracking)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Tracking fetched successfully"))
	return
}

func main() {
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
