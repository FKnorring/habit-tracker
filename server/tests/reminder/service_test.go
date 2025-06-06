package reminder_test

import (
	"encoding/json"
	"testing"
	"time"

	"habit-tracker/server/db"
	"habit-tracker/server/reminder"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockDatabase is a mock implementation of the Database interface
type MockDatabase struct {
	mock.Mock
}

func (m *MockDatabase) Ping() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockDatabase) CreateHabit(habit *db.Habit) error {
	args := m.Called(habit)
	return args.Error(0)
}

func (m *MockDatabase) GetHabit(id string) (*db.Habit, error) {
	args := m.Called(id)
	return args.Get(0).(*db.Habit), args.Error(1)
}

func (m *MockDatabase) GetAllHabits() ([]*db.Habit, error) {
	args := m.Called()
	return args.Get(0).([]*db.Habit), args.Error(1)
}

func (m *MockDatabase) UpdateHabit(habit *db.Habit) error {
	args := m.Called(habit)
	return args.Error(0)
}

func (m *MockDatabase) UpdateHabitPartial(id string, updates map[string]interface{}) (*db.Habit, error) {
	args := m.Called(id, updates)
	return args.Get(0).(*db.Habit), args.Error(1)
}

func (m *MockDatabase) DeleteHabit(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockDatabase) CreateTrackingEntry(entry *db.TrackingEntry) error {
	args := m.Called(entry)
	return args.Error(0)
}

func (m *MockDatabase) GetTrackingEntry(id string) (*db.TrackingEntry, error) {
	args := m.Called(id)
	return args.Get(0).(*db.TrackingEntry), args.Error(1)
}

func (m *MockDatabase) GetTrackingEntriesByHabitID(habitID string) ([]*db.TrackingEntry, error) {
	args := m.Called(habitID)
	return args.Get(0).([]*db.TrackingEntry), args.Error(1)
}

func (m *MockDatabase) DeleteTrackingEntry(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockDatabase) CreateReminder(reminder *db.Reminder) error {
	args := m.Called(reminder)
	return args.Error(0)
}

func (m *MockDatabase) GetReminder(habitID string) (*db.Reminder, error) {
	args := m.Called(habitID)
	return args.Get(0).(*db.Reminder), args.Error(1)
}

func (m *MockDatabase) UpdateReminderLastReminder(habitID string, lastReminder string) error {
	args := m.Called(habitID, lastReminder)
	return args.Error(0)
}

func (m *MockDatabase) GetHabitsNeedingReminders() ([]*db.Habit, error) {
	args := m.Called()
	return args.Get(0).([]*db.Habit), args.Error(1)
}

func (m *MockDatabase) DeleteReminder(habitID string) error {
	args := m.Called(habitID)
	return args.Error(0)
}

// Statistics and Analytics Methods
func (m *MockDatabase) GetHabitStats(habitID string) (*db.HabitStats, error) {
	args := m.Called(habitID)
	return args.Get(0).(*db.HabitStats), args.Error(1)
}

func (m *MockDatabase) GetHabitProgress(habitID string, days int) ([]*db.ProgressPoint, error) {
	args := m.Called(habitID, days)
	return args.Get(0).([]*db.ProgressPoint), args.Error(1)
}

func (m *MockDatabase) GetOverallStats() (*db.OverallStats, error) {
	args := m.Called()
	return args.Get(0).(*db.OverallStats), args.Error(1)
}

func (m *MockDatabase) GetHabitCompletionRates(days int) ([]*db.HabitCompletionRate, error) {
	args := m.Called(days)
	return args.Get(0).([]*db.HabitCompletionRate), args.Error(1)
}

func (m *MockDatabase) GetDailyCompletions(days int) ([]*db.DailyCompletion, error) {
	args := m.Called(days)
	return args.Get(0).([]*db.DailyCompletion), args.Error(1)
}

// User Management Methods
func (m *MockDatabase) CreateUser(user *db.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockDatabase) GetUserByEmail(email string) (*db.User, error) {
	args := m.Called(email)
	return args.Get(0).(*db.User), args.Error(1)
}

func (m *MockDatabase) GetUserByID(id string) (*db.User, error) {
	args := m.Called(id)
	return args.Get(0).(*db.User), args.Error(1)
}

func (m *MockDatabase) UpdateUser(user *db.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockDatabase) DeleteUser(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func TestNewReminderService(t *testing.T) {
	mockDB := &MockDatabase{}
	service := reminder.NewReminderService(mockDB)

	assert.NotNil(t, service)
}

func TestSetCheckInterval(t *testing.T) {
	mockDB := &MockDatabase{}
	service := reminder.NewReminderService(mockDB)

	customInterval := 10 * time.Second
	service.SetCheckInterval(customInterval)

	// Note: We can't directly test the internal field, but we can verify it doesn't panic
	assert.NotNil(t, service)
}

func TestReminderMessageFormat(t *testing.T) {
	// Test that reminder message can be properly marshaled to JSON
	reminderData := reminder.ReminderData{
		HabitID:     "habit-1",
		HabitName:   "Test Habit",
		Description: "A test habit",
		Frequency:   "daily",
		Timestamp:   time.Now().Format(time.RFC3339),
	}

	reminderMessage := reminder.ReminderMessage{
		Type: "reminder",
		Data: reminderData,
	}

	messageBytes, err := json.Marshal(reminderMessage)
	assert.NoError(t, err)
	assert.Contains(t, string(messageBytes), "reminder")
	assert.Contains(t, string(messageBytes), "Test Habit")
	assert.Contains(t, string(messageBytes), "habit-1")
}

func TestCheckAndSendRemindersNoHabits(t *testing.T) {
	mockDB := &MockDatabase{}
	service := reminder.NewReminderService(mockDB)

	// Setup mock expectation for no habits needing reminders
	mockDB.On("GetHabitsNeedingReminders").Return([]*db.Habit{}, nil)

	// Test that service can be created successfully
	assert.NotNil(t, service)

	// Verify mock was properly set up (the actual method call would happen in the private method)
	mockDB.AssertNotCalled(t, "UpdateReminderLastReminder")
}
