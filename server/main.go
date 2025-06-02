package main

import (
	"encoding/json"
	"log"
	"net/http"
)

/*
	Habit:
		id
		name
		description
		frequency
		startDate

	Tracking:
		id
		habitId
		timestamp
		note?
*/

type Habit struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Frequency   string `json:"frequency"`
	StartDate   string `json:"startDate"`
}

type Tracking struct {
	ID        string `json:"id"`
	HabitID   string `json:"habitId"`
	Timestamp string `json:"timestamp"`
	Note      string `json:"note"`
}

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

var trackings = []Tracking{
	{ID: "1", HabitID: "1", Timestamp: "2021-01-01", Note: "Note 1"},
	{ID: "2", HabitID: "2", Timestamp: "2021-01-01", Note: "Note 2"},
}

func getHabits(w http.ResponseWriter, r *http.Request) {
	// get habits stub
	json.NewEncoder(w).Encode(habits)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Habits fetched successfully"))
	return
}

func createHabit(w http.ResponseWriter, r *http.Request) {
	// create habit stub
	habit := Habit{ID: "3", Name: "Habit 3", Description: "Description 3", Frequency: "Monthly", StartDate: "2021-01-01"}
	json.NewEncoder(w).Encode(habit)
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Habit created successfully"))
	return
}

func habitsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		getHabits(w, r)
	case "POST":
		createHabit(w, r)
	}
}

func getHabit(w http.ResponseWriter, r *http.Request) {
	// get habit stub
	habit := habits[0]
	json.NewEncoder(w).Encode(habit)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Habit fetched successfully"))
	return
}

func updateHabit(w http.ResponseWriter, r *http.Request) {
	// update habit stub
	habit := habits[0]
	json.NewEncoder(w).Encode(habit)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Habit updated successfully"))
	return
}

func deleteHabit(w http.ResponseWriter, r *http.Request) {
	// delete habit stub
	json.NewEncoder(w).Encode(nil)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Habit deleted successfully"))
	return
}

func habitHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		getHabit(w, r)
	case "PUT":
		updateHabit(w, r)
	case "DELETE":
		deleteHabit(w, r)
	}
}

func createTracking(w http.ResponseWriter, r *http.Request) {
	// create tracking stub
	tracking := trackings[0]
	json.NewEncoder(w).Encode(tracking)
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Tracking created successfully"))
	return
}

func getTracking(w http.ResponseWriter, r *http.Request) {
	// get tracking stub
	tracking := trackings[0]
	json.NewEncoder(w).Encode(tracking)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Tracking fetched successfully"))
	return
}

func trackingHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		createTracking(w, r)
	case "GET":
		getTracking(w, r)
	}
}

func main() {
	http.HandleFunc("/habits", habitsHandler)
	http.HandleFunc("/habits/:id", habitHandler)
	http.HandleFunc("/habits/:id/tracking", trackingHandler)
	log.Println("Server is running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
