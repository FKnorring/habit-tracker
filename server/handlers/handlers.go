package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"habit-tracker/server/db"

	"github.com/google/uuid"
)

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

	if habits == nil {
		habits = []*db.Habit{}
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

	if err := db.ValidateFrequency(string(habit.Frequency)); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid frequency: must be one of hourly, daily, weekly, biweekly, monthly, quarterly, yearly"))
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

	if err := db.ValidateFrequency(string(habit.Frequency)); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid frequency: must be one of hourly, daily, weekly, biweekly, monthly, quarterly, yearly"))
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

	if err := Database.UpdateReminderLastReminder(entry.HabitID, entry.Timestamp); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to update reminder"))
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

	if entries == nil {
		entries = []*db.TrackingEntry{}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(entries)
}

func UpdateReminder(w http.ResponseWriter, r *http.Request, params map[string]string) {
	if !checkParams(w, params, []string{"id"}) {
		return
	}

	var reminder db.Reminder
	if err := json.NewDecoder(r.Body).Decode(&reminder); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid JSON"))
		return
	}

	reminder.ID = params["id"]

	if err := Database.UpdateReminderLastReminder(reminder.HabitID, reminder.LastReminder); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to update reminder"))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(reminder)
}
