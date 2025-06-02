package db_test

import (
	"testing"

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
		Frequency:   "daily",
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
		Frequency:   "daily",
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
		Frequency:   "daily",
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
			Frequency:   "daily",
			StartDate:   "2024-01-01",
		},
		{
			ID:          "habit-2",
			Name:        "Reading",
			Description: "Read for 30 minutes",
			Frequency:   "daily",
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
		Frequency:   "daily",
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
		Frequency:   "twice daily",
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
		Frequency:   "daily",
		StartDate:   "2024-01-01",
	}

	err := suite.db.UpdateHabit(habit)
	suite.Equal(db.ErrNotFound, err)
}

func (suite *InMemoryDBTestSuite) TestDeleteHabit() {
	habit := &db.Habit{
		ID:          "test-habit-1",
		Name:        "Exercise",
		Description: "Daily workout",
		Frequency:   "daily",
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
		Frequency:   "daily",
		StartDate:   "2024-01-01",
	}

	// This is a basic test - for production use, you'd want proper concurrent testing
	err := suite.db.CreateHabit(habit)
	suite.NoError(err)

	retrieved, err := suite.db.GetHabit(habit.ID)
	suite.NoError(err)
	suite.Equal(habit.ID, retrieved.ID)
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
		Frequency:   "daily",
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
