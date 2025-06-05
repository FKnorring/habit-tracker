package db_test

import (
	"testing"
	"time"

	"habit-tracker/server/db"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type InMemoryDBTestSuite struct {
	suite.Suite
	db *db.MapDatabase
}

func (suite *InMemoryDBTestSuite) SetupTest() {
	suite.db = db.NewMapDatabase()
}

func (suite *InMemoryDBTestSuite) TestNewMapDatabase() {
	database := db.NewMapDatabase()
	suite.NotNil(database)
	suite.NoError(database.Ping())
}

func (suite *InMemoryDBTestSuite) TestPing() {
	err := suite.db.Ping()
	suite.NoError(err)
}

func (suite *InMemoryDBTestSuite) TestCreateHabit() {
	habit := &db.Habit{
		ID:          "test-habit-1",
		Name:        "Exercise",
		Description: "Daily workout",
		Frequency:   db.FrequencyDaily,
		StartDate:   "2024-01-01",
	}

	err := suite.db.CreateHabit(habit)
	suite.NoError(err)

	// Verify habit is stored by retrieving it
	stored, err := suite.db.GetHabit(habit.ID)
	suite.NoError(err)
	suite.NotNil(stored)
	suite.Equal(habit.ID, stored.ID)
	suite.Equal(habit.Name, stored.Name)
	suite.Equal(habit.Description, stored.Description)
	suite.Equal(habit.Frequency, stored.Frequency)
	suite.Equal(habit.StartDate, stored.StartDate)

	// Verify it's a copy, not the same reference
	suite.NotSame(habit, stored)
}

func (suite *InMemoryDBTestSuite) TestCreateHabitDuplicate() {
	habit := &db.Habit{
		ID:          "test-habit-1",
		Name:        "Exercise",
		Description: "Daily workout",
		Frequency:   db.FrequencyDaily,
		StartDate:   "2024-01-01",
	}

	// Create habit first time
	err := suite.db.CreateHabit(habit)
	suite.NoError(err)

	// Try to create same habit again
	err = suite.db.CreateHabit(habit)
	suite.Equal(db.ErrDuplicate, err)
}

func (suite *InMemoryDBTestSuite) TestGetHabit() {
	habit := &db.Habit{
		ID:          "test-habit-1",
		Name:        "Exercise",
		Description: "Daily workout",
		Frequency:   db.FrequencyDaily,
		StartDate:   "2024-01-01",
	}

	// Create habit
	err := suite.db.CreateHabit(habit)
	suite.NoError(err)

	// Get habit
	retrieved, err := suite.db.GetHabit(habit.ID)
	suite.NoError(err)
	suite.NotNil(retrieved)
	suite.Equal(habit.ID, retrieved.ID)
	suite.Equal(habit.Name, retrieved.Name)
	suite.Equal(habit.Description, retrieved.Description)
	suite.Equal(habit.Frequency, retrieved.Frequency)
	suite.Equal(habit.StartDate, retrieved.StartDate)

	// Verify it's a copy, not the same reference
	suite.NotSame(habit, retrieved)
}

func (suite *InMemoryDBTestSuite) TestGetHabitNotFound() {
	retrieved, err := suite.db.GetHabit("nonexistent")
	suite.Nil(retrieved)
	suite.Equal(db.ErrNotFound, err)
}

func (suite *InMemoryDBTestSuite) TestGetAllHabitsEmpty() {
	habits, err := suite.db.GetAllHabits()
	suite.NoError(err)
	suite.Empty(habits)
}

func (suite *InMemoryDBTestSuite) TestGetAllHabits() {
	habits := []*db.Habit{
		{
			ID:          "habit-1",
			Name:        "Exercise",
			Description: "Daily workout",
			Frequency:   db.FrequencyDaily,
			StartDate:   "2024-01-01",
		},
		{
			ID:          "habit-2",
			Name:        "Reading",
			Description: "Read for 30 minutes",
			Frequency:   db.FrequencyDaily,
			StartDate:   "2024-01-01",
		},
	}

	// Create habits
	for _, habit := range habits {
		err := suite.db.CreateHabit(habit)
		suite.NoError(err)
	}

	// Get all habits
	retrieved, err := suite.db.GetAllHabits()
	suite.NoError(err)
	suite.Len(retrieved, 2)

	// Verify each habit
	idMap := make(map[string]*db.Habit)
	for _, h := range retrieved {
		idMap[h.ID] = h
	}

	for _, originalHabit := range habits {
		retrievedHabit := idMap[originalHabit.ID]
		suite.NotNil(retrievedHabit)
		suite.Equal(originalHabit.ID, retrievedHabit.ID)
		suite.Equal(originalHabit.Name, retrievedHabit.Name)
		suite.NotSame(originalHabit, retrievedHabit) // Verify it's a copy
	}
}

func (suite *InMemoryDBTestSuite) TestUpdateHabit() {
	habit := &db.Habit{
		ID:          "test-habit-1",
		Name:        "Exercise",
		Description: "Daily workout",
		Frequency:   db.FrequencyDaily,
		StartDate:   "2024-01-01",
	}

	// Create habit
	err := suite.db.CreateHabit(habit)
	suite.NoError(err)

	// Update habit
	updatedHabit := &db.Habit{
		ID:          "test-habit-1",
		Name:        "Updated Exercise",
		Description: "Updated daily workout",
		Frequency:   db.FrequencyWeekly,
		StartDate:   "2024-01-02",
	}

	err = suite.db.UpdateHabit(updatedHabit)
	suite.NoError(err)

	// Verify update
	retrieved, err := suite.db.GetHabit(habit.ID)
	suite.NoError(err)
	suite.Equal(updatedHabit.Name, retrieved.Name)
	suite.Equal(updatedHabit.Description, retrieved.Description)
	suite.Equal(updatedHabit.Frequency, retrieved.Frequency)
	suite.Equal(updatedHabit.StartDate, retrieved.StartDate)
}

func (suite *InMemoryDBTestSuite) TestUpdateHabitNotFound() {
	habit := &db.Habit{
		ID:          "nonexistent",
		Name:        "Exercise",
		Description: "Daily workout",
		Frequency:   db.FrequencyDaily,
		StartDate:   "2024-01-01",
	}

	err := suite.db.UpdateHabit(habit)
	suite.Equal(db.ErrNotFound, err)
}

func (suite *InMemoryDBTestSuite) TestUpdateHabitPartial() {
	// Create a habit first
	habit := &db.Habit{
		ID:          "test-habit-partial",
		Name:        "Original Name",
		Description: "Original Description",
		Frequency:   db.FrequencyDaily,
		StartDate:   "2024-01-01",
	}

	err := suite.db.CreateHabit(habit)
	suite.NoError(err)

	// Test partial update with only name
	updates := map[string]interface{}{
		"name": "Updated Name",
	}

	updatedHabit, err := suite.db.UpdateHabitPartial("test-habit-partial", updates)
	suite.NoError(err)
	suite.NotNil(updatedHabit)

	// Verify only name was updated
	suite.Equal("test-habit-partial", updatedHabit.ID)
	suite.Equal("Updated Name", updatedHabit.Name)
	suite.Equal("Original Description", updatedHabit.Description) // Should remain unchanged
	suite.Equal(db.FrequencyDaily, updatedHabit.Frequency)        // Should remain unchanged
	suite.Equal("2024-01-01", updatedHabit.StartDate)             // Should remain unchanged

	// Test partial update with multiple fields
	updates2 := map[string]interface{}{
		"description": "Updated Description",
		"frequency":   "weekly",
	}

	updatedHabit2, err := suite.db.UpdateHabitPartial("test-habit-partial", updates2)
	suite.NoError(err)
	suite.NotNil(updatedHabit2)

	// Verify multiple fields were updated while others remained
	suite.Equal("test-habit-partial", updatedHabit2.ID)
	suite.Equal("Updated Name", updatedHabit2.Name)               // Should remain from previous update
	suite.Equal("Updated Description", updatedHabit2.Description) // Should be updated
	suite.Equal(db.FrequencyWeekly, updatedHabit2.Frequency)      // Should be updated
	suite.Equal("2024-01-01", updatedHabit2.StartDate)            // Should remain unchanged

	// Test with invalid frequency
	invalidUpdates := map[string]interface{}{
		"frequency": "invalid_frequency",
	}

	_, err = suite.db.UpdateHabitPartial("test-habit-partial", invalidUpdates)
	suite.Error(err)
	suite.Contains(err.Error(), "invalid frequency")

	// Test with empty updates
	emptyUpdates := map[string]interface{}{}

	unchangedHabit, err := suite.db.UpdateHabitPartial("test-habit-partial", emptyUpdates)
	suite.NoError(err)
	suite.NotNil(unchangedHabit)

	// Should return the existing habit unchanged
	suite.Equal("test-habit-partial", unchangedHabit.ID)
	suite.Equal("Updated Name", unchangedHabit.Name)
	suite.Equal("Updated Description", unchangedHabit.Description)
	suite.Equal(db.FrequencyWeekly, unchangedHabit.Frequency)
	suite.Equal("2024-01-01", unchangedHabit.StartDate)
}

func (suite *InMemoryDBTestSuite) TestUpdateHabitPartialNotFound() {
	updates := map[string]interface{}{
		"name": "Updated Name",
	}

	_, err := suite.db.UpdateHabitPartial("nonexistent", updates)
	suite.Equal(db.ErrNotFound, err)
}

func (suite *InMemoryDBTestSuite) TestDeleteHabit() {
	habit := &db.Habit{
		ID:          "test-habit-1",
		Name:        "Exercise",
		Description: "Daily workout",
		Frequency:   db.FrequencyDaily,
		StartDate:   "2024-01-01",
	}

	// Create habit
	err := suite.db.CreateHabit(habit)
	suite.NoError(err)

	// Delete habit
	err = suite.db.DeleteHabit(habit.ID)
	suite.NoError(err)

	// Verify deletion
	_, err = suite.db.GetHabit(habit.ID)
	suite.Equal(db.ErrNotFound, err)
}

func (suite *InMemoryDBTestSuite) TestDeleteHabitNotFound() {
	err := suite.db.DeleteHabit("nonexistent")
	suite.Equal(db.ErrNotFound, err)
}

func (suite *InMemoryDBTestSuite) TestCreateTrackingEntry() {
	entry := &db.TrackingEntry{
		ID:        "entry-1",
		HabitID:   "habit-1",
		Timestamp: "2024-01-01T10:00:00Z",
		Note:      "Great workout!",
	}

	err := suite.db.CreateTrackingEntry(entry)
	suite.NoError(err)

	// Verify entry is stored by retrieving it
	stored, err := suite.db.GetTrackingEntry(entry.ID)
	suite.NoError(err)
	suite.NotNil(stored)
	suite.Equal(entry.ID, stored.ID)
	suite.Equal(entry.HabitID, stored.HabitID)
	suite.Equal(entry.Timestamp, stored.Timestamp)
	suite.Equal(entry.Note, stored.Note)

	// Verify it's a copy
	suite.NotSame(entry, stored)
}

func (suite *InMemoryDBTestSuite) TestCreateTrackingEntryDuplicate() {
	entry := &db.TrackingEntry{
		ID:        "entry-1",
		HabitID:   "habit-1",
		Timestamp: "2024-01-01T10:00:00Z",
		Note:      "Great workout!",
	}

	// Create entry first time
	err := suite.db.CreateTrackingEntry(entry)
	suite.NoError(err)

	// Try to create same entry again
	err = suite.db.CreateTrackingEntry(entry)
	suite.Equal(db.ErrDuplicate, err)
}

func (suite *InMemoryDBTestSuite) TestGetTrackingEntry() {
	entry := &db.TrackingEntry{
		ID:        "entry-1",
		HabitID:   "habit-1",
		Timestamp: "2024-01-01T10:00:00Z",
		Note:      "Great workout!",
	}

	// Create entry
	err := suite.db.CreateTrackingEntry(entry)
	suite.NoError(err)

	// Get entry
	retrieved, err := suite.db.GetTrackingEntry(entry.ID)
	suite.NoError(err)
	suite.NotNil(retrieved)
	suite.Equal(entry.ID, retrieved.ID)
	suite.Equal(entry.HabitID, retrieved.HabitID)
	suite.Equal(entry.Timestamp, retrieved.Timestamp)
	suite.Equal(entry.Note, retrieved.Note)

	// Verify it's a copy
	suite.NotSame(entry, retrieved)
}

func (suite *InMemoryDBTestSuite) TestGetTrackingEntryNotFound() {
	retrieved, err := suite.db.GetTrackingEntry("nonexistent")
	suite.Nil(retrieved)
	suite.Equal(db.ErrNotFound, err)
}

func (suite *InMemoryDBTestSuite) TestGetTrackingEntriesByHabitID() {
	entries := []*db.TrackingEntry{
		{
			ID:        "entry-1",
			HabitID:   "habit-1",
			Timestamp: "2024-01-01T10:00:00Z",
			Note:      "Morning workout",
		},
		{
			ID:        "entry-2",
			HabitID:   "habit-1",
			Timestamp: "2024-01-02T10:00:00Z",
			Note:      "Evening workout",
		},
		{
			ID:        "entry-3",
			HabitID:   "habit-2",
			Timestamp: "2024-01-01T20:00:00Z",
			Note:      "Reading session",
		},
	}

	// Create entries
	for _, entry := range entries {
		err := suite.db.CreateTrackingEntry(entry)
		suite.NoError(err)
	}

	// Get entries for habit-1
	retrieved, err := suite.db.GetTrackingEntriesByHabitID("habit-1")
	suite.NoError(err)
	suite.Len(retrieved, 2)

	// Verify entries belong to habit-1
	for _, entry := range retrieved {
		suite.Equal("habit-1", entry.HabitID)
	}

	// Get entries for habit-2
	retrieved, err = suite.db.GetTrackingEntriesByHabitID("habit-2")
	suite.NoError(err)
	suite.Len(retrieved, 1)
	suite.Equal("habit-2", retrieved[0].HabitID)

	// Get entries for nonexistent habit
	retrieved, err = suite.db.GetTrackingEntriesByHabitID("nonexistent")
	suite.NoError(err)
	suite.Empty(retrieved)
}

func (suite *InMemoryDBTestSuite) TestDeleteTrackingEntry() {
	entry := &db.TrackingEntry{
		ID:        "entry-1",
		HabitID:   "habit-1",
		Timestamp: "2024-01-01T10:00:00Z",
		Note:      "Great workout!",
	}

	// Create entry
	err := suite.db.CreateTrackingEntry(entry)
	suite.NoError(err)

	// Delete entry
	err = suite.db.DeleteTrackingEntry(entry.ID)
	suite.NoError(err)

	// Verify deletion
	_, err = suite.db.GetTrackingEntry(entry.ID)
	suite.Equal(db.ErrNotFound, err)
}

func (suite *InMemoryDBTestSuite) TestDeleteTrackingEntryNotFound() {
	err := suite.db.DeleteTrackingEntry("nonexistent")
	suite.Equal(db.ErrNotFound, err)
}

// Test concurrent access safety (basic test)
func (suite *InMemoryDBTestSuite) TestConcurrentAccess() {
	habit := &db.Habit{
		ID:          "concurrent-habit",
		Name:        "Concurrent Test",
		Description: "Testing concurrent access",
		Frequency:   db.FrequencyDaily,
		StartDate:   "2024-01-01",
	}

	// This is a basic test - for production use, you'd want proper concurrent testing
	err := suite.db.CreateHabit(habit)
	suite.NoError(err)

	retrieved, err := suite.db.GetHabit(habit.ID)
	suite.NoError(err)
	suite.Equal(habit.ID, retrieved.ID)
}

func (suite *InMemoryDBTestSuite) TestCreateReminder() {
	// First create a habit
	habit := &db.Habit{
		ID:          "test-habit-1",
		Name:        "Exercise",
		Description: "Daily workout",
		Frequency:   db.FrequencyDaily,
		StartDate:   "2024-01-01",
	}
	err := suite.db.CreateHabit(habit)
	suite.NoError(err)

	// Get the automatically created reminder
	stored, err := suite.db.GetReminder(habit.ID)
	suite.NoError(err)
	suite.NotNil(stored)
	suite.Equal(habit.ID+"-reminder", stored.ID)
	suite.Equal(habit.ID, stored.HabitID)
	suite.NotEmpty(stored.LastReminder)

	// Verify it's a copy, not the same reference
	suite.NotSame(habit, stored)
}

func (suite *InMemoryDBTestSuite) TestCreateReminderDuplicate() {
	// First create a habit
	habit := &db.Habit{
		ID:          "test-habit-1",
		Name:        "Exercise",
		Description: "Daily workout",
		Frequency:   db.FrequencyDaily,
		StartDate:   "2024-01-01",
	}
	err := suite.db.CreateHabit(habit)
	suite.NoError(err)

	// Try to create another reminder for the same habit
	duplicateReminder := &db.Reminder{
		ID:           "reminder-2",
		HabitID:      "test-habit-1",
		LastReminder: "2024-01-02T10:00:00Z",
	}

	err = suite.db.CreateReminder(duplicateReminder)
	suite.Equal(db.ErrDuplicate, err)
}

func (suite *InMemoryDBTestSuite) TestGetReminder() {
	// First create a habit
	habit := &db.Habit{
		ID:          "test-habit-1",
		Name:        "Exercise",
		Description: "Daily workout",
		Frequency:   db.FrequencyDaily,
		StartDate:   "2024-01-01",
	}
	err := suite.db.CreateHabit(habit)
	suite.NoError(err)

	// Get reminder
	retrieved, err := suite.db.GetReminder(habit.ID)
	suite.NoError(err)
	suite.NotNil(retrieved)
	suite.Equal(habit.ID+"-reminder", retrieved.ID)
	suite.Equal(habit.ID, retrieved.HabitID)
	suite.NotEmpty(retrieved.LastReminder)

	// Verify it's a copy, not the same reference
	suite.NotSame(habit, retrieved)
}

func (suite *InMemoryDBTestSuite) TestGetReminderNotFound() {
	retrieved, err := suite.db.GetReminder("nonexistent-habit")
	suite.Nil(retrieved)
	suite.Equal(db.ErrNotFound, err)
}

func (suite *InMemoryDBTestSuite) TestUpdateReminderLastReminder() {
	// First create a habit
	habit := &db.Habit{
		ID:          "test-habit-1",
		Name:        "Exercise",
		Description: "Daily workout",
		Frequency:   db.FrequencyDaily,
		StartDate:   "2024-01-01",
	}
	err := suite.db.CreateHabit(habit)
	suite.NoError(err)

	// Update last reminder time
	newLastReminder := "2024-01-02T10:00:00Z"
	err = suite.db.UpdateReminderLastReminder(habit.ID, newLastReminder)
	suite.NoError(err)

	// Verify update
	retrieved, err := suite.db.GetReminder(habit.ID)
	suite.NoError(err)
	suite.Equal(newLastReminder, retrieved.LastReminder)
}

func (suite *InMemoryDBTestSuite) TestUpdateReminderLastReminderNotFound() {
	err := suite.db.UpdateReminderLastReminder("nonexistent-habit", "2024-01-01T10:00:00Z")
	suite.Equal(db.ErrNotFound, err)
}

func (suite *InMemoryDBTestSuite) TestDeleteReminder() {
	// First create a habit
	habit := &db.Habit{
		ID:          "test-habit-1",
		Name:        "Exercise",
		Description: "Daily workout",
		Frequency:   db.FrequencyDaily,
		StartDate:   "2024-01-01",
	}
	err := suite.db.CreateHabit(habit)
	suite.NoError(err)

	// Verify reminder exists
	retrieved, err := suite.db.GetReminder(habit.ID)
	suite.NoError(err)
	suite.NotNil(retrieved)

	// Delete reminder
	err = suite.db.DeleteReminder(habit.ID)
	suite.NoError(err)

	// Verify reminder is deleted
	retrieved, err = suite.db.GetReminder(habit.ID)
	suite.Nil(retrieved)
	suite.Equal(db.ErrNotFound, err)
}

func (suite *InMemoryDBTestSuite) TestDeleteReminderNotFound() {
	err := suite.db.DeleteReminder("nonexistent-habit")
	suite.Equal(db.ErrNotFound, err)
}

func (suite *InMemoryDBTestSuite) TestGetHabitsNeedingReminders() {
	// Create habits with different frequencies
	habits := []*db.Habit{
		{
			ID:          "habit-daily",
			Name:        "Daily Exercise",
			Description: "Daily workout",
			Frequency:   db.FrequencyDaily,
			StartDate:   "2024-01-01",
		},
		{
			ID:          "habit-weekly",
			Name:        "Weekly Reading",
			Description: "Read a book",
			Frequency:   db.FrequencyWeekly,
			StartDate:   "2024-01-01",
		},
		{
			ID:          "habit-hourly",
			Name:        "Hourly Water",
			Description: "Drink water",
			Frequency:   db.FrequencyHourly,
			StartDate:   "2024-01-01",
		},
	}

	// Create habits
	for _, habit := range habits {
		err := suite.db.CreateHabit(habit)
		suite.NoError(err)
	}

	// Update reminders with old timestamps to trigger reminders
	oldTime := "2024-01-01T10:00:00Z"
	for _, habit := range habits {
		err := suite.db.UpdateReminderLastReminder(habit.ID, oldTime)
		suite.NoError(err)
	}

	// Get habits needing reminders
	needingReminders, err := suite.db.GetHabitsNeedingReminders()
	suite.NoError(err)
	suite.Len(needingReminders, 3) // All should need reminders due to old timestamp

	// Verify all expected habits are in the result
	habitIDs := make(map[string]bool)
	for _, habit := range needingReminders {
		habitIDs[habit.ID] = true
	}

	suite.True(habitIDs["habit-daily"])
	suite.True(habitIDs["habit-weekly"])
	suite.True(habitIDs["habit-hourly"])
}

func (suite *InMemoryDBTestSuite) TestGetHabitsNeedingRemindersNoReminders() {
	// Create a habit but no reminder
	habit := &db.Habit{
		ID:          "habit-no-reminder",
		Name:        "Exercise",
		Description: "Daily workout",
		Frequency:   db.FrequencyDaily,
		StartDate:   "2024-01-01",
	}

	err := suite.db.CreateHabit(habit)
	suite.NoError(err)

	// Delete the auto-created reminder to test empty case
	err = suite.db.DeleteReminder(habit.ID)
	suite.NoError(err)

	// Get habits needing reminders - should be empty
	needingReminders, err := suite.db.GetHabitsNeedingReminders()
	suite.NoError(err)
	suite.Empty(needingReminders)
}

func (suite *InMemoryDBTestSuite) TestGetHabitsNeedingRemindersRecentReminders() {
	// Create a habit
	habit := &db.Habit{
		ID:          "habit-recent",
		Name:        "Exercise",
		Description: "Daily workout",
		Frequency:   db.FrequencyDaily,
		StartDate:   "2024-01-01",
	}

	err := suite.db.CreateHabit(habit)
	suite.NoError(err)

	// Update reminder to very recent time (future)
	futureTime := "2099-01-01T10:00:00Z"
	err = suite.db.UpdateReminderLastReminder(habit.ID, futureTime)
	suite.NoError(err)

	// Get habits needing reminders - should be empty
	needingReminders, err := suite.db.GetHabitsNeedingReminders()
	suite.NoError(err)
	suite.Empty(needingReminders)
}

func (suite *InMemoryDBTestSuite) TestGetHabitsNeedingRemindersInvalidTimestamp() {
	// Create a habit
	habit := &db.Habit{
		ID:          "habit-invalid",
		Name:        "Exercise",
		Description: "Daily workout",
		Frequency:   db.FrequencyDaily,
		StartDate:   "2024-01-01",
	}

	err := suite.db.CreateHabit(habit)
	suite.NoError(err)

	// Update reminder to invalid timestamp
	err = suite.db.UpdateReminderLastReminder(habit.ID, "invalid-timestamp")
	suite.NoError(err)

	// Get habits needing reminders - should handle invalid timestamp gracefully
	needingReminders, err := suite.db.GetHabitsNeedingReminders()
	suite.NoError(err)
	// Should not include the habit with invalid timestamp
	for _, h := range needingReminders {
		suite.NotEqual("habit-invalid", h.ID)
	}
}

func (suite *InMemoryDBTestSuite) TestHabitReminderIntegration() {
	// Test that creating a habit automatically creates a reminder
	habit := &db.Habit{
		ID:          "test-habit-integration",
		Name:        "Exercise",
		Description: "Daily workout",
		Frequency:   db.FrequencyDaily,
		StartDate:   "2024-01-01",
	}

	err := suite.db.CreateHabit(habit)
	suite.NoError(err)

	// Verify reminder was automatically created
	reminder, err := suite.db.GetReminder(habit.ID)
	suite.NoError(err)
	suite.NotNil(reminder)
	suite.Equal(habit.ID+"-reminder", reminder.ID)
	suite.Equal(habit.ID, reminder.HabitID)
	suite.NotEmpty(reminder.LastReminder)

	// Test that deleting a habit also deletes the reminder
	err = suite.db.DeleteHabit(habit.ID)
	suite.NoError(err)

	// Verify reminder is also deleted
	reminder, err = suite.db.GetReminder(habit.ID)
	suite.Nil(reminder)
	suite.Equal(db.ErrNotFound, err)
}

func TestInMemoryDBTestSuite(t *testing.T) {
	suite.Run(t, new(InMemoryDBTestSuite))
}

// Additional unit tests for specific edge cases
func TestHabitCopyIntegrity(t *testing.T) {
	database := db.NewMapDatabase()
	original := &db.Habit{
		ID:          "test-habit",
		Name:        "Original Name",
		Description: "Original Description",
		Frequency:   db.FrequencyDaily,
		StartDate:   "2024-01-01",
	}

	// Create habit
	err := database.CreateHabit(original)
	assert.NoError(t, err)

	// Modify original after creation
	original.Name = "Modified Name"

	// Verify stored habit wasn't affected
	retrieved, err := database.GetHabit("test-habit")
	assert.NoError(t, err)
	assert.Equal(t, "Original Name", retrieved.Name)

	// Modify retrieved habit
	retrieved.Name = "Another Modification"

	// Verify stored habit still wasn't affected
	retrieved2, err := database.GetHabit("test-habit")
	assert.NoError(t, err)
	assert.Equal(t, "Original Name", retrieved2.Name)
}

func TestTrackingEntryCopyIntegrity(t *testing.T) {
	database := db.NewMapDatabase()
	original := &db.TrackingEntry{
		ID:        "test-entry",
		HabitID:   "test-habit",
		Timestamp: "2024-01-01T10:00:00Z",
		Note:      "Original Note",
	}

	// Create entry
	err := database.CreateTrackingEntry(original)
	assert.NoError(t, err)

	// Modify original after creation
	original.Note = "Modified Note"

	// Verify stored entry wasn't affected
	retrieved, err := database.GetTrackingEntry("test-entry")
	assert.NoError(t, err)
	assert.Equal(t, "Original Note", retrieved.Note)

	// Modify retrieved entry
	retrieved.Note = "Another Modification"

	// Verify stored entry still wasn't affected
	retrieved2, err := database.GetTrackingEntry("test-entry")
	assert.NoError(t, err)
	assert.Equal(t, "Original Note", retrieved2.Note)
}

func TestReminderCopyIntegrity(t *testing.T) {
	database := db.NewMapDatabase()

	// Create a habit first
	habit := &db.Habit{
		ID:          "test-habit",
		Name:        "Exercise",
		Description: "Daily workout",
		Frequency:   db.FrequencyDaily,
		StartDate:   "2024-01-01",
	}
	err := database.CreateHabit(habit)
	assert.NoError(t, err)

	// Get the automatically created reminder
	retrieved, err := database.GetReminder(habit.ID)
	assert.NoError(t, err)
	assert.NotNil(t, retrieved)

	// Store the original last reminder time
	originalLastReminder := retrieved.LastReminder

	// Modify the retrieved reminder
	retrieved.LastReminder = "2024-01-02T10:00:00Z"

	// Get the reminder again to ensure database stores copies
	retrieved2, err := database.GetReminder(habit.ID)
	assert.NoError(t, err)
	assert.NotSame(t, retrieved, retrieved2)                       // Should be different objects
	assert.Equal(t, originalLastReminder, retrieved2.LastReminder) // Should have original data
}

func TestCalculateNextReminderTime(t *testing.T) {
	baseTime := time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC)

	tests := []struct {
		name      string
		frequency db.Frequency
		expected  time.Time
	}{
		{
			name:      "Hourly frequency",
			frequency: db.FrequencyHourly,
			expected:  baseTime.Add(time.Hour),
		},
		{
			name:      "Daily frequency",
			frequency: db.FrequencyDaily,
			expected:  baseTime.AddDate(0, 0, 1),
		},
		{
			name:      "Weekly frequency",
			frequency: db.FrequencyWeekly,
			expected:  baseTime.AddDate(0, 0, 7),
		},
		{
			name:      "Biweekly frequency",
			frequency: db.FrequencyBiweekly,
			expected:  baseTime.AddDate(0, 0, 14),
		},
		{
			name:      "Monthly frequency",
			frequency: db.FrequencyMonthly,
			expected:  baseTime.AddDate(0, 1, 0),
		},
		{
			name:      "Quarterly frequency",
			frequency: db.FrequencyQuarterly,
			expected:  baseTime.AddDate(0, 3, 0),
		},
		{
			name:      "Yearly frequency",
			frequency: db.FrequencyYearly,
			expected:  baseTime.AddDate(1, 0, 0),
		},
		{
			name:      "Invalid frequency defaults to daily",
			frequency: db.Frequency("invalid"),
			expected:  baseTime.AddDate(0, 0, 1),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := db.CalculateNextReminderTime(baseTime, tt.frequency)
			assert.Equal(t, tt.expected, result)
		})
	}
}
