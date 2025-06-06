package db

import (
	"time"
)

type MapDatabase struct {
	habits    map[string]*Habit
	tracking  map[string]*TrackingEntry
	reminders map[string]*Reminder
	users     map[string]*User
}

func NewMapDatabase() *MapDatabase {
	return &MapDatabase{
		habits:    make(map[string]*Habit),
		tracking:  make(map[string]*TrackingEntry),
		reminders: make(map[string]*Reminder),
		users:     make(map[string]*User),
	}
}

func (db *MapDatabase) Ping() error {
	// In-memory database is always available
	return nil
}

func (db *MapDatabase) CreateHabit(habit *Habit) error {
	if _, exists := db.habits[habit.ID]; exists {
		return ErrDuplicate
	}

	habitCopy := *habit
	db.habits[habit.ID] = &habitCopy

	reminder := &Reminder{
		ID:           habit.ID + "-reminder",
		HabitID:      habit.ID,
		LastReminder: time.Now().Format(time.RFC3339),
	}
	db.reminders[habit.ID] = reminder

	return nil
}

func (db *MapDatabase) GetHabit(id string) (*Habit, error) {
	habit, exists := db.habits[id]
	if !exists {
		return nil, ErrNotFound
	}

	habitCopy := *habit
	return &habitCopy, nil
}

func (db *MapDatabase) GetAllHabits() ([]*Habit, error) {
	habits := make([]*Habit, 0, len(db.habits))
	for _, habit := range db.habits {
		habitCopy := *habit
		habits = append(habits, &habitCopy)
	}
	return habits, nil
}

func (db *MapDatabase) UpdateHabit(habit *Habit) error {
	if _, exists := db.habits[habit.ID]; !exists {
		return ErrNotFound
	}

	habitCopy := *habit
	db.habits[habit.ID] = &habitCopy
	return nil
}

func (db *MapDatabase) UpdateHabitPartial(id string, updates map[string]interface{}) (*Habit, error) {
	existing, exists := db.habits[id]
	if !exists {
		return nil, ErrNotFound
	}

	// Create a copy of the existing habit
	updated := *existing

	// Apply updates
	for field, value := range updates {
		switch field {
		case "name":
			if name, ok := value.(string); ok {
				updated.Name = name
			}
		case "description":
			if desc, ok := value.(string); ok {
				updated.Description = desc
			}
		case "frequency":
			if freqStr, ok := value.(string); ok {
				if err := ValidateFrequency(freqStr); err != nil {
					return nil, err
				}
				updated.Frequency = Frequency(freqStr)
			}
		case "startDate":
			if startDate, ok := value.(string); ok {
				updated.StartDate = startDate
			}
		}
	}

	// Store the updated habit
	db.habits[id] = &updated

	// Return a copy
	result := updated
	return &result, nil
}

func (db *MapDatabase) DeleteHabit(id string) error {
	if _, exists := db.habits[id]; !exists {
		return ErrNotFound
	}

	delete(db.habits, id)
	delete(db.reminders, id)
	return nil
}

func (db *MapDatabase) CreateTrackingEntry(entry *TrackingEntry) error {
	if _, exists := db.tracking[entry.ID]; exists {
		return ErrDuplicate
	}

	entryCopy := *entry
	db.tracking[entry.ID] = &entryCopy
	return nil
}

func (db *MapDatabase) GetTrackingEntry(id string) (*TrackingEntry, error) {
	entry, exists := db.tracking[id]
	if !exists {
		return nil, ErrNotFound
	}

	entryCopy := *entry
	return &entryCopy, nil
}

func (db *MapDatabase) GetTrackingEntriesByHabitID(habitID string) ([]*TrackingEntry, error) {
	var entries []*TrackingEntry
	for _, entry := range db.tracking {
		if entry.HabitID == habitID {
			entryCopy := *entry
			entries = append(entries, &entryCopy)
		}
	}
	return entries, nil
}

func (db *MapDatabase) DeleteTrackingEntry(id string) error {
	if _, exists := db.tracking[id]; !exists {
		return ErrNotFound
	}

	delete(db.tracking, id)
	return nil
}

func (db *MapDatabase) CreateReminder(reminder *Reminder) error {
	if _, exists := db.reminders[reminder.HabitID]; exists {
		return ErrDuplicate
	}

	reminderCopy := *reminder
	db.reminders[reminder.HabitID] = &reminderCopy
	return nil
}

func (db *MapDatabase) GetReminder(habitID string) (*Reminder, error) {
	reminder, exists := db.reminders[habitID]
	if !exists {
		return nil, ErrNotFound
	}

	reminderCopy := *reminder
	return &reminderCopy, nil
}

func (db *MapDatabase) UpdateReminderLastReminder(habitID string, lastReminder string) error {
	reminder, exists := db.reminders[habitID]
	if !exists {
		return ErrNotFound
	}

	reminder.LastReminder = lastReminder
	return nil
}

func (db *MapDatabase) GetHabitsNeedingReminders() ([]*Habit, error) {
	var needingReminders []*Habit
	now := time.Now()

	for habitID, reminder := range db.reminders {
		habit, exists := db.habits[habitID]
		if !exists {
			continue
		}

		lastReminder, err := time.Parse(time.RFC3339, reminder.LastReminder)
		if err != nil {
			continue
		}

		nextReminderTime := CalculateNextReminderTime(lastReminder, habit.Frequency)
		if now.After(nextReminderTime) {
			habitCopy := *habit
			needingReminders = append(needingReminders, &habitCopy)
		}
	}

	return needingReminders, nil
}

func (db *MapDatabase) DeleteReminder(habitID string) error {
	if _, exists := db.reminders[habitID]; !exists {
		return ErrNotFound
	}

	delete(db.reminders, habitID)
	return nil
}

// Statistics and Analytics Methods for MapDatabase

func (db *MapDatabase) GetHabitStats(habitID string) (*HabitStats, error) {
	habit, err := db.GetHabit(habitID)
	if err != nil {
		return nil, err
	}

	stats := &HabitStats{
		HabitID:   habitID,
		HabitName: habit.Name,
		Frequency: habit.Frequency,
		StartDate: habit.StartDate,
	}

	// Count total entries for this habit
	totalEntries := 0
	var lastCompleted string
	for _, entry := range db.tracking {
		if entry.HabitID == habitID {
			totalEntries++
			if entry.Timestamp > lastCompleted {
				lastCompleted = entry.Timestamp
			}
		}
	}

	stats.TotalEntries = totalEntries
	stats.LastCompleted = lastCompleted

	// For simplicity in in-memory implementation, set basic values
	stats.CurrentStreak = 0    // Would need more complex logic
	stats.LongestStreak = 0    // Would need more complex logic
	stats.CompletionRate = 0.0 // Would need more complex logic

	return stats, nil
}

func (db *MapDatabase) GetHabitProgress(habitID string, days int) ([]*ProgressPoint, error) {
	// Simple implementation - in a real scenario would need date parsing and filtering
	progress := []*ProgressPoint{}

	dateCount := make(map[string]int)
	for _, entry := range db.tracking {
		if entry.HabitID == habitID {
			// Extract date part from timestamp (simplified)
			date := entry.Timestamp[:10] // Assumes RFC3339 format
			dateCount[date]++
		}
	}

	for date, count := range dateCount {
		progress = append(progress, &ProgressPoint{
			Date:  date,
			Count: count,
		})
	}

	return progress, nil
}

func (db *MapDatabase) GetOverallStats() (*OverallStats, error) {
	stats := &OverallStats{
		TotalHabits:      len(db.habits),
		TotalEntries:     len(db.tracking),
		EntriesToday:     0,   // Would need date filtering
		EntriesThisWeek:  0,   // Would need date filtering
		AvgEntriesPerDay: 0.0, // Would need date calculations
	}

	return stats, nil
}

func (db *MapDatabase) GetHabitCompletionRates(days int) ([]*HabitCompletionRate, error) {
	var rates []*HabitCompletionRate

	for _, habit := range db.habits {
		rate := &HabitCompletionRate{
			HabitID:             habit.ID,
			HabitName:           habit.Name,
			Frequency:           habit.Frequency,
			StartDate:           habit.StartDate,
			ActualCompletions:   0,
			ExpectedCompletions: 0,
			CompletionRate:      0.0,
		}

		// Count actual completions for this habit
		for _, entry := range db.tracking {
			if entry.HabitID == habit.ID {
				rate.ActualCompletions++
			}
		}

		rates = append(rates, rate)
	}

	return rates, nil
}

func (db *MapDatabase) GetDailyCompletions(days int) ([]*DailyCompletion, error) {
	dateCount := make(map[string]int)

	for _, entry := range db.tracking {
		// Extract date part from timestamp (simplified)
		date := entry.Timestamp[:10] // Assumes RFC3339 format
		dateCount[date]++
	}

	var completions []*DailyCompletion
	for date, count := range dateCount {
		completions = append(completions, &DailyCompletion{
			Date:        date,
			Completions: count,
		})
	}

	return completions, nil
}

// User Management Methods for MapDatabase

func (db *MapDatabase) CreateUser(user *User) error {
	// Generate UUID for user if not provided
	if user.ID == "" {
		user.ID = generateUUID()
	}

	// Check if email already exists
	for _, existingUser := range db.users {
		if existingUser.Email == user.Email {
			return ErrDuplicate
		}
		if existingUser.Username == user.Username {
			return ErrDuplicate
		}
	}

	// Create a copy to store
	userCopy := *user
	db.users[user.ID] = &userCopy
	return nil
}

func (db *MapDatabase) GetUserByEmail(email string) (*User, error) {
	for _, user := range db.users {
		if user.Email == email {
			userCopy := *user
			return &userCopy, nil
		}
	}
	return nil, ErrNotFound
}

func (db *MapDatabase) GetUserByID(id string) (*User, error) {
	user, exists := db.users[id]
	if !exists {
		return nil, ErrNotFound
	}

	userCopy := *user
	return &userCopy, nil
}

func (db *MapDatabase) UpdateUser(user *User) error {
	if _, exists := db.users[user.ID]; !exists {
		return ErrNotFound
	}

	// Create a copy to store
	userCopy := *user
	db.users[user.ID] = &userCopy
	return nil
}

func (db *MapDatabase) DeleteUser(id string) error {
	if _, exists := db.users[id]; !exists {
		return ErrNotFound
	}

	delete(db.users, id)
	return nil
}
