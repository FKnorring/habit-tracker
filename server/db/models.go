package db

import (
	"errors"
	"fmt"
)

var (
	ErrNotFound         = errors.New("record not found")
	ErrDuplicate        = errors.New("record already exists")
	ErrInvalidFrequency = errors.New("invalid frequency")
)

type Frequency string

const (
	FrequencyHourly    Frequency = "hourly"
	FrequencyDaily     Frequency = "daily"
	FrequencyWeekly    Frequency = "weekly"
	FrequencyBiweekly  Frequency = "biweekly"
	FrequencyMonthly   Frequency = "monthly"
	FrequencyQuarterly Frequency = "quarterly"
	FrequencyYearly    Frequency = "yearly"
)

var validFrequencies = map[Frequency]bool{
	FrequencyHourly:    true,
	FrequencyDaily:     true,
	FrequencyWeekly:    true,
	FrequencyBiweekly:  true,
	FrequencyMonthly:   true,
	FrequencyQuarterly: true,
	FrequencyYearly:    true,
}

func (f Frequency) IsValid() bool {
	return validFrequencies[f]
}

func (f Frequency) String() string {
	return string(f)
}

func ValidateFrequency(freq string) error {
	if !Frequency(freq).IsValid() {
		return fmt.Errorf("%w: %s", ErrInvalidFrequency, freq)
	}
	return nil
}

type Habit struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Frequency   Frequency `json:"frequency"`
	StartDate   string    `json:"startDate"`
}

type TrackingEntry struct {
	ID        string `json:"id"`
	HabitID   string `json:"habitId"`
	Timestamp string `json:"timestamp"`
	Note      string `json:"note"`
}

type Reminder struct {
	ID           string `json:"id"`
	HabitID      string `json:"habitId"`
	LastReminder string `json:"lastReminder"`
}

// Statistics and Analytics Structures

type HabitStats struct {
	HabitID        string    `json:"habitId"`
	HabitName      string    `json:"habitName"`
	Frequency      Frequency `json:"frequency"`
	StartDate      string    `json:"startDate"`
	TotalEntries   int       `json:"totalEntries"`
	CurrentStreak  int       `json:"currentStreak"`
	LongestStreak  int       `json:"longestStreak"`
	CompletionRate float64   `json:"completionRate"`
	LastCompleted  string    `json:"lastCompleted"`
}

type ProgressPoint struct {
	Date  string `json:"date"`
	Count int    `json:"count"`
}

type OverallStats struct {
	TotalHabits      int     `json:"totalHabits"`
	TotalEntries     int     `json:"totalEntries"`
	EntriesToday     int     `json:"entriesToday"`
	EntriesThisWeek  int     `json:"entriesThisWeek"`
	AvgEntriesPerDay float64 `json:"avgEntriesPerDay"`
}

type HabitCompletionRate struct {
	HabitID             string    `json:"habitId"`
	HabitName           string    `json:"habitName"`
	Frequency           Frequency `json:"frequency"`
	StartDate           string    `json:"startDate"`
	ActualCompletions   int       `json:"actualCompletions"`
	ExpectedCompletions int       `json:"expectedCompletions"`
	CompletionRate      float64   `json:"completionRate"`
}

type DailyCompletion struct {
	Date        string `json:"date"`
	Completions int    `json:"completions"`
}

type Database interface {
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

	CreateReminder(reminder *Reminder) error
	GetReminder(habitID string) (*Reminder, error)
	UpdateReminderLastReminder(habitID string, lastReminder string) error
	GetHabitsNeedingReminders() ([]*Habit, error)
	DeleteReminder(habitID string) error

	// Statistics and Analytics Methods
	GetHabitStats(habitID string) (*HabitStats, error)
	GetHabitProgress(habitID string, days int) ([]*ProgressPoint, error)
	GetOverallStats() (*OverallStats, error)
	GetHabitCompletionRates(days int) ([]*HabitCompletionRate, error)
	GetDailyCompletions(days int) ([]*DailyCompletion, error)
}
