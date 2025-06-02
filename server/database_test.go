package main

import (
	"testing"
)

func TestMapDatabase_HabitOperations(t *testing.T) {
	db := NewMapDatabase()

	habit := &Habit{
		ID:          "habit1",
		Name:        "Exercise",
		Description: "Daily exercise routine",
		Frequency:   "daily",
		StartDate:   "2024-01-01",
	}

	err := db.CreateHabit(habit)
	if err != nil {
		t.Fatalf("Failed to create habit: %v", err)
	}

	// Test duplicate creation
	err = db.CreateHabit(habit)
	if err != ErrDuplicate {
		t.Fatalf("Expected ErrDuplicate, got: %v", err)
	}

	// Test getting a habit
	retrieved, err := db.GetHabit("habit1")
	if err != nil {
		t.Fatalf("Failed to get habit: %v", err)
	}

	if retrieved.Name != habit.Name {
		t.Fatalf("Expected name %s, got %s", habit.Name, retrieved.Name)
	}

	// Test getting non-existent habit
	_, err = db.GetHabit("nonexistent")
	if err != ErrNotFound {
		t.Fatalf("Expected ErrNotFound, got: %v", err)
	}

	// Test updating a habit
	habit.Name = "Updated Exercise"
	err = db.UpdateHabit(habit)
	if err != nil {
		t.Fatalf("Failed to update habit: %v", err)
	}

	retrieved, err = db.GetHabit("habit1")
	if err != nil {
		t.Fatalf("Failed to get updated habit: %v", err)
	}

	if retrieved.Name != "Updated Exercise" {
		t.Fatalf("Expected updated name, got %s", retrieved.Name)
	}

	// Test getting all habits
	habits, err := db.GetAllHabits()
	if err != nil {
		t.Fatalf("Failed to get all habits: %v", err)
	}

	if len(habits) != 1 {
		t.Fatalf("Expected 1 habit, got %d", len(habits))
	}

	// Test deleting a habit
	err = db.DeleteHabit("habit1")
	if err != nil {
		t.Fatalf("Failed to delete habit: %v", err)
	}

	// Verify deletion
	_, err = db.GetHabit("habit1")
	if err != ErrNotFound {
		t.Fatalf("Expected ErrNotFound after deletion, got: %v", err)
	}
}

func TestMapDatabase_TrackingOperations(t *testing.T) {
	db := NewMapDatabase()

	// Test creating a tracking entry
	entry := &TrackingEntry{
		ID:        "track1",
		HabitID:   "habit1",
		Timestamp: "2024-01-01T10:00:00Z",
		Note:      "First workout",
	}

	err := db.CreateTrackingEntry(entry)
	if err != nil {
		t.Fatalf("Failed to create tracking entry: %v", err)
	}

	// Test getting a tracking entry
	retrieved, err := db.GetTrackingEntry("track1")
	if err != nil {
		t.Fatalf("Failed to get tracking entry: %v", err)
	}

	if retrieved.Note != entry.Note {
		t.Fatalf("Expected note %s, got %s", entry.Note, retrieved.Note)
	}

	// Test getting tracking entries by habit ID
	entries, err := db.GetTrackingEntriesByHabitID("habit1")
	if err != nil {
		t.Fatalf("Failed to get tracking entries by habit ID: %v", err)
	}

	if len(entries) != 1 {
		t.Fatalf("Expected 1 tracking entry, got %d", len(entries))
	}

	// Test deleting a tracking entry
	err = db.DeleteTrackingEntry("track1")
	if err != nil {
		t.Fatalf("Failed to delete tracking entry: %v", err)
	}

	// Verify deletion
	_, err = db.GetTrackingEntry("track1")
	if err != ErrNotFound {
		t.Fatalf("Expected ErrNotFound after deletion, got: %v", err)
	}
}

// Example of how you would use this in your main application
func TestDatabaseInjection(t *testing.T) {
	// For testing, use the map database
	testDB := NewMapDatabase()

	// For production, you would use something like:
	// prodDB, err := NewSQLDatabase("connection_string")

	// Your application code works with either implementation
	testService := &HabitService{DB: testDB}

	habit := &Habit{
		ID:   "test-habit",
		Name: "Test Habit",
	}

	err := testService.CreateHabit(habit)
	if err != nil {
		t.Fatalf("Service failed to create habit: %v", err)
	}
}

// Example service that uses dependency injection
type HabitService struct {
	DB Database // This can be any implementation of the Database interface
}

func (s *HabitService) CreateHabit(habit *Habit) error {
	return s.DB.CreateHabit(habit)
}

func (s *HabitService) GetHabit(id string) (*Habit, error) {
	return s.DB.GetHabit(id)
}

// ... other service methods
