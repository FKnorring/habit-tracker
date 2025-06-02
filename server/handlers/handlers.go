package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"habit-tracker/server/db"

	"github.com/google/uuid"
)

type HandlerFunc func(http.ResponseWriter, *http.Request, map[string]string)

type route struct {
	method  string
	pattern string
	handler HandlerFunc
}

type Router struct {
	routes []route
}

var Database db.Database

func CreateRouter() *Router {
	return &Router{
		routes: []route{},
	}
}

func (r *Router) Handle(method, pattern string, handler HandlerFunc) {
	r.routes = append(r.routes, route{method: method, pattern: pattern, handler: handler})
}

func addCORSHeaders(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*") // TODO: Change to only allow requests from the frontend
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	w.Header().Set("Access-Control-Max-Age", "86400")
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	addCORSHeaders(w)

	if req.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	for _, route := range r.routes {
		if req.Method != route.method {
			continue
		}

		params, matched := Match(route.pattern, req.URL.Path)

		if !matched {
			continue
		}

		log.Println(route.method, req.URL.Path)

		route.handler(w, req, params)
		return
	}

	http.NotFound(w, req)
}

func Match(pattern, path string) (map[string]string, bool) {
	patternParts := strings.Split(strings.Trim(pattern, "/"), "/")
	pathParts := strings.Split(strings.Trim(path, "/"), "/")

	if len(patternParts) != len(pathParts) {
		return nil, false
	}

	params := make(map[string]string)

	for i := range patternParts {
		if patternParts[i] == pathParts[i] {
			continue
		}

		if strings.HasPrefix(patternParts[i], ":") {
			params[strings.TrimPrefix(patternParts[i], ":")] = pathParts[i]
		} else {
			return nil, false
		}
	}

	return params, true
}

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

func GetHabits(w http.ResponseWriter, r *http.Request, params map[string]string) {
	habits, err := Database.GetAllHabits()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to retrieve habits"))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(habits)
}

func CreateHabit(w http.ResponseWriter, r *http.Request, params map[string]string) {
	var habit db.Habit
	if err := json.NewDecoder(r.Body).Decode(&habit); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid JSON"))
		return
	}

	if habit.ID == "" {
		habit.ID = uuid.New().String()
	}

	if err := Database.CreateHabit(&habit); err != nil {
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

func GetHabit(w http.ResponseWriter, r *http.Request, params map[string]string) {
	if !checkParams(w, params, []string{"id"}) {
		return
	}

	habit, err := Database.GetHabit(params["id"])
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

func UpdateHabit(w http.ResponseWriter, r *http.Request, params map[string]string) {
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

	if err := Database.UpdateHabit(&habit); err != nil {
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

func DeleteHabit(w http.ResponseWriter, r *http.Request, params map[string]string) {
	if !checkParams(w, params, []string{"id"}) {
		return
	}

	if err := Database.DeleteHabit(params["id"]); err != nil {
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

func CreateTracking(w http.ResponseWriter, r *http.Request, params map[string]string) {
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

	if err := Database.CreateTrackingEntry(&entry); err != nil {
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

func GetTracking(w http.ResponseWriter, r *http.Request, params map[string]string) {
	if !checkParams(w, params, []string{"id"}) {
		return
	}

	entries, err := Database.GetTrackingEntriesByHabitID(params["id"])
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to retrieve tracking entries"))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(entries)
}
