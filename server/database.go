package main

import (
	"errors"
)

var (
	ErrNotFound  = errors.New("record not found")
	ErrDuplicate = errors.New("record already exists")
)

type Database interface {
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
