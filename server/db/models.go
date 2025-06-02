package db

import (
	"errors"
)

var (
	ErrNotFound  = errors.New("record not found")
	ErrDuplicate = errors.New("record already exists")
)

type Habit struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Frequency   string `json:"frequency"`
	StartDate   string `json:"startDate"`
}

type TrackingEntry struct {
	ID        string `json:"id"`
	HabitID   string `json:"habitId"`
	Timestamp string `json:"timestamp"`
	Note      string `json:"note"`
}

type Database interface {
	// Connection health check
	Ping() error

	CreateHabit(habit *Habit) error
	GetHabit(id string) (*Habit, error)
	GetAllHabits() ([]*Habit, error)
	UpdateHabit(habit *Habit) error
	DeleteHabit(id string) error

	CreateTrackingEntry(entry *TrackingEntry) error
	GetTrackingEntry(id string) (*TrackingEntry, error)
	GetTrackingEntriesByHabitID(habitID string) ([]*TrackingEntry, error)
	DeleteTrackingEntry(id string) error
}
