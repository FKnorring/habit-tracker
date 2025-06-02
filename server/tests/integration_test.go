package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"habit-tracker/server/db"

	"github.com/stretchr/testify/suite"
)

type IntegrationTestSuite struct {
	suite.Suite
	router *Router
	server *httptest.Server
	origDB db.Database
}

func (suite *IntegrationTestSuite) SetupSuite() {
	// Store original database
	suite.origDB = database

	// Create in-memory database for testing
	database = db.NewMapDatabase()

	// Create router and register handlers
	suite.router = CreateRouter()
	suite.router.Handle("GET", "/habits", getHabits)
	suite.router.Handle("POST", "/habits", createHabit)
	suite.router.Handle("GET", "/habits/:id", getHabit)
	suite.router.Handle("PUT", "/habits/:id", updateHabit)
	suite.router.Handle("DELETE", "/habits/:id", deleteHabit)
	suite.router.Handle("POST", "/habits/:id/tracking", createTracking)
	suite.router.Handle("GET", "/habits/:id/tracking", getTracking)

	// Create test server
	suite.server = httptest.NewServer(suite.router)
}

func (suite *IntegrationTestSuite) TearDownSuite() {
	suite.server.Close()
	// Restore original database
	database = suite.origDB
}

func (suite *IntegrationTestSuite) SetupTest() {
	// Reset database state before each test
	database = db.NewMapDatabase()
}

func (suite *IntegrationTestSuite) TestGetHabitsEmpty() {
	resp, err := http.Get(suite.server.URL + "/habits")
	suite.NoError(err)
	defer resp.Body.Close()

	suite.Equal(http.StatusOK, resp.StatusCode)

	var habits []db.Habit
	err = json.NewDecoder(resp.Body).Decode(&habits)
	suite.NoError(err)
	suite.Empty(habits)
}

func (suite *IntegrationTestSuite) TestCreateHabit() {
	habitData := db.Habit{
		Name:        "Exercise",
		Description: "Daily workout routine",
		Frequency:   "daily",
		StartDate:   "2024-01-01",
	}

	jsonData, err := json.Marshal(habitData)
	suite.NoError(err)

	resp, err := http.Post(suite.server.URL+"/habits", "application/json", bytes.NewBuffer(jsonData))
	suite.NoError(err)
	defer resp.Body.Close()

	suite.Equal(http.StatusCreated, resp.StatusCode)

	var createdHabit db.Habit
	err = json.NewDecoder(resp.Body).Decode(&createdHabit)
	suite.NoError(err)

	suite.NotEmpty(createdHabit.ID)
	suite.Equal(habitData.Name, createdHabit.Name)
	suite.Equal(habitData.Description, createdHabit.Description)
	suite.Equal(habitData.Frequency, createdHabit.Frequency)
	suite.Equal(habitData.StartDate, createdHabit.StartDate)
}

func (suite *IntegrationTestSuite) TestCreateHabitInvalidJSON() {
	resp, err := http.Post(suite.server.URL+"/habits", "application/json", bytes.NewBuffer([]byte("invalid json")))
	suite.NoError(err)
	defer resp.Body.Close()

	suite.Equal(http.StatusBadRequest, resp.StatusCode)
}

func (suite *IntegrationTestSuite) TestGetHabitById() {
	// First create a habit
	habit := &db.Habit{
		ID:          "test-habit-1",
		Name:        "Reading",
		Description: "Read for 30 minutes",
		Frequency:   "daily",
		StartDate:   "2024-01-01",
	}
	err := database.CreateHabit(habit)
	suite.NoError(err)

	// Get the habit
	resp, err := http.Get(suite.server.URL + "/habits/test-habit-1")
	suite.NoError(err)
	defer resp.Body.Close()

	suite.Equal(http.StatusOK, resp.StatusCode)

	var retrievedHabit db.Habit
	err = json.NewDecoder(resp.Body).Decode(&retrievedHabit)
	suite.NoError(err)
	suite.Equal(habit.ID, retrievedHabit.ID)
	suite.Equal(habit.Name, retrievedHabit.Name)
}

func (suite *IntegrationTestSuite) TestGetHabitNotFound() {
	resp, err := http.Get(suite.server.URL + "/habits/nonexistent")
	suite.NoError(err)
	defer resp.Body.Close()

	suite.Equal(http.StatusNotFound, resp.StatusCode)
}

func (suite *IntegrationTestSuite) TestUpdateHabit() {
	// First create a habit
	habit := &db.Habit{
		ID:          "test-habit-2",
		Name:        "Meditation",
		Description: "Daily meditation",
		Frequency:   "daily",
		StartDate:   "2024-01-01",
	}
	err := database.CreateHabit(habit)
	suite.NoError(err)

	// Update the habit
	updatedHabit := db.Habit{
		Name:        "Updated Meditation",
		Description: "Updated daily meditation practice",
		Frequency:   "twice daily",
		StartDate:   "2024-01-02",
	}

	jsonData, err := json.Marshal(updatedHabit)
	suite.NoError(err)

	req, err := http.NewRequest("PUT", suite.server.URL+"/habits/test-habit-2", bytes.NewBuffer(jsonData))
	suite.NoError(err)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	suite.NoError(err)
	defer resp.Body.Close()

	suite.Equal(http.StatusOK, resp.StatusCode)

	var responseHabit db.Habit
	err = json.NewDecoder(resp.Body).Decode(&responseHabit)
	suite.NoError(err)
	suite.Equal("test-habit-2", responseHabit.ID)
	suite.Equal(updatedHabit.Name, responseHabit.Name)
}

func (suite *IntegrationTestSuite) TestUpdateHabitNotFound() {
	updatedHabit := db.Habit{
		Name:        "Nonexistent",
		Description: "This habit doesn't exist",
		Frequency:   "never",
		StartDate:   "2024-01-01",
	}

	jsonData, err := json.Marshal(updatedHabit)
	suite.NoError(err)

	req, err := http.NewRequest("PUT", suite.server.URL+"/habits/nonexistent", bytes.NewBuffer(jsonData))
	suite.NoError(err)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	suite.NoError(err)
	defer resp.Body.Close()

	suite.Equal(http.StatusNotFound, resp.StatusCode)
}

func (suite *IntegrationTestSuite) TestDeleteHabit() {
	// First create a habit
	habit := &db.Habit{
		ID:          "test-habit-3",
		Name:        "Journaling",
		Description: "Daily journaling",
		Frequency:   "daily",
		StartDate:   "2024-01-01",
	}
	err := database.CreateHabit(habit)
	suite.NoError(err)

	// Delete the habit
	req, err := http.NewRequest("DELETE", suite.server.URL+"/habits/test-habit-3", nil)
	suite.NoError(err)

	client := &http.Client{}
	resp, err := client.Do(req)
	suite.NoError(err)
	defer resp.Body.Close()

	suite.Equal(http.StatusNoContent, resp.StatusCode)

	// Verify it's deleted
	_, err = database.GetHabit("test-habit-3")
	suite.Equal(db.ErrNotFound, err)
}

func (suite *IntegrationTestSuite) TestDeleteHabitNotFound() {
	req, err := http.NewRequest("DELETE", suite.server.URL+"/habits/nonexistent", nil)
	suite.NoError(err)

	client := &http.Client{}
	resp, err := client.Do(req)
	suite.NoError(err)
	defer resp.Body.Close()

	suite.Equal(http.StatusNotFound, resp.StatusCode)
}

func (suite *IntegrationTestSuite) TestCreateTracking() {
	// First create a habit
	habit := &db.Habit{
		ID:          "test-habit-4",
		Name:        "Running",
		Description: "Daily run",
		Frequency:   "daily",
		StartDate:   "2024-01-01",
	}
	err := database.CreateHabit(habit)
	suite.NoError(err)

	// Create tracking entry
	trackingData := db.TrackingEntry{
		Note: "Great run today!",
	}

	jsonData, err := json.Marshal(trackingData)
	suite.NoError(err)

	resp, err := http.Post(suite.server.URL+"/habits/test-habit-4/tracking", "application/json", bytes.NewBuffer(jsonData))
	suite.NoError(err)
	defer resp.Body.Close()

	suite.Equal(http.StatusCreated, resp.StatusCode)

	var createdEntry db.TrackingEntry
	err = json.NewDecoder(resp.Body).Decode(&createdEntry)
	suite.NoError(err)

	suite.NotEmpty(createdEntry.ID)
	suite.Equal("test-habit-4", createdEntry.HabitID)
	suite.Equal(trackingData.Note, createdEntry.Note)
	suite.NotEmpty(createdEntry.Timestamp)
}

func (suite *IntegrationTestSuite) TestCreateTrackingWithCustomTimestamp() {
	// First create a habit
	habit := &db.Habit{
		ID:          "test-habit-5",
		Name:        "Water",
		Description: "Drink water",
		Frequency:   "hourly",
		StartDate:   "2024-01-01",
	}
	err := database.CreateHabit(habit)
	suite.NoError(err)

	// Create tracking entry with custom timestamp
	customTime := time.Now().Add(-2 * time.Hour).Format(time.RFC3339)
	trackingData := db.TrackingEntry{
		Note:      "Drank 500ml",
		Timestamp: customTime,
	}

	jsonData, err := json.Marshal(trackingData)
	suite.NoError(err)

	resp, err := http.Post(suite.server.URL+"/habits/test-habit-5/tracking", "application/json", bytes.NewBuffer(jsonData))
	suite.NoError(err)
	defer resp.Body.Close()

	suite.Equal(http.StatusCreated, resp.StatusCode)

	var createdEntry db.TrackingEntry
	err = json.NewDecoder(resp.Body).Decode(&createdEntry)
	suite.NoError(err)
	suite.Equal(customTime, createdEntry.Timestamp)
}

func (suite *IntegrationTestSuite) TestGetTracking() {
	// First create a habit
	habit := &db.Habit{
		ID:          "test-habit-6",
		Name:        "Stretching",
		Description: "Daily stretching",
		Frequency:   "daily",
		StartDate:   "2024-01-01",
	}
	err := database.CreateHabit(habit)
	suite.NoError(err)

	// Create some tracking entries
	entries := []*db.TrackingEntry{
		{
			ID:        "entry-1",
			HabitID:   "test-habit-6",
			Timestamp: time.Now().Format(time.RFC3339),
			Note:      "Morning stretch",
		},
		{
			ID:        "entry-2",
			HabitID:   "test-habit-6",
			Timestamp: time.Now().Add(-24 * time.Hour).Format(time.RFC3339),
			Note:      "Evening stretch",
		},
	}

	for _, entry := range entries {
		err := database.CreateTrackingEntry(entry)
		suite.NoError(err)
	}

	// Get tracking entries
	resp, err := http.Get(suite.server.URL + "/habits/test-habit-6/tracking")
	suite.NoError(err)
	defer resp.Body.Close()

	suite.Equal(http.StatusOK, resp.StatusCode)

	var retrievedEntries []db.TrackingEntry
	err = json.NewDecoder(resp.Body).Decode(&retrievedEntries)
	suite.NoError(err)
	suite.Len(retrievedEntries, 2)
}

func (suite *IntegrationTestSuite) TestGetTrackingEmpty() {
	// First create a habit
	habit := &db.Habit{
		ID:          "test-habit-7",
		Name:        "Empty Habit",
		Description: "No tracking entries",
		Frequency:   "daily",
		StartDate:   "2024-01-01",
	}
	err := database.CreateHabit(habit)
	suite.NoError(err)

	// Get tracking entries (should be empty)
	resp, err := http.Get(suite.server.URL + "/habits/test-habit-7/tracking")
	suite.NoError(err)
	defer resp.Body.Close()

	suite.Equal(http.StatusOK, resp.StatusCode)

	var retrievedEntries []db.TrackingEntry
	err = json.NewDecoder(resp.Body).Decode(&retrievedEntries)
	suite.NoError(err)
	suite.Empty(retrievedEntries)
}

func (suite *IntegrationTestSuite) TestFullWorkflow() {
	// Test complete workflow: Create habit -> Create tracking -> Get all -> Update -> Delete

	// 1. Create a habit
	habitData := db.Habit{
		Name:        "Full Workflow Test",
		Description: "Testing complete workflow",
		Frequency:   "daily",
		StartDate:   "2024-01-01",
	}

	jsonData, err := json.Marshal(habitData)
	suite.NoError(err)

	resp, err := http.Post(suite.server.URL+"/habits", "application/json", bytes.NewBuffer(jsonData))
	suite.NoError(err)
	defer resp.Body.Close()

	var createdHabit db.Habit
	err = json.NewDecoder(resp.Body).Decode(&createdHabit)
	suite.NoError(err)
	habitID := createdHabit.ID

	// 2. Create tracking entry
	trackingData := db.TrackingEntry{
		Note: "Workflow test tracking",
	}

	jsonData, err = json.Marshal(trackingData)
	suite.NoError(err)

	resp, err = http.Post(suite.server.URL+"/habits/"+habitID+"/tracking", "application/json", bytes.NewBuffer(jsonData))
	suite.NoError(err)
	resp.Body.Close()
	suite.Equal(http.StatusCreated, resp.StatusCode)

	// 3. Get all habits
	resp, err = http.Get(suite.server.URL + "/habits")
	suite.NoError(err)
	defer resp.Body.Close()

	var habits []db.Habit
	err = json.NewDecoder(resp.Body).Decode(&habits)
	suite.NoError(err)
	suite.Len(habits, 1)

	// 4. Update habit
	updatedHabit := db.Habit{
		Name:        "Updated Workflow Test",
		Description: "Updated workflow description",
		Frequency:   "weekly",
		StartDate:   "2024-01-02",
	}

	jsonData, err = json.Marshal(updatedHabit)
	suite.NoError(err)

	req, err := http.NewRequest("PUT", suite.server.URL+"/habits/"+habitID, bytes.NewBuffer(jsonData))
	suite.NoError(err)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err = client.Do(req)
	suite.NoError(err)
	resp.Body.Close()
	suite.Equal(http.StatusOK, resp.StatusCode)

	// 5. Get tracking entries
	resp, err = http.Get(suite.server.URL + "/habits/" + habitID + "/tracking")
	suite.NoError(err)
	defer resp.Body.Close()

	var trackingEntries []db.TrackingEntry
	err = json.NewDecoder(resp.Body).Decode(&trackingEntries)
	suite.NoError(err)
	suite.Len(trackingEntries, 1)

	// 6. Delete habit
	req, err = http.NewRequest("DELETE", suite.server.URL+"/habits/"+habitID, nil)
	suite.NoError(err)

	resp, err = client.Do(req)
	suite.NoError(err)
	resp.Body.Close()
	suite.Equal(http.StatusNoContent, resp.StatusCode)

	// 7. Verify deletion - should return empty list
	resp, err = http.Get(suite.server.URL + "/habits")
	suite.NoError(err)
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&habits)
	suite.NoError(err)
	suite.Empty(habits)
}

// Test helper functions
func (suite *IntegrationTestSuite) TestCheckParams() {
	tests := []struct {
		name           string
		params         map[string]string
		requiredParams []string
		expectedResult bool
	}{
		{
			name:           "all params present",
			params:         map[string]string{"id": "123", "name": "test"},
			requiredParams: []string{"id", "name"},
			expectedResult: true,
		},
		{
			name:           "missing param",
			params:         map[string]string{"id": "123"},
			requiredParams: []string{"id", "name"},
			expectedResult: false,
		},
		{
			name:           "no required params",
			params:         map[string]string{"id": "123"},
			requiredParams: []string{},
			expectedResult: true,
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			w := httptest.NewRecorder()
			result := checkParams(w, tt.params, tt.requiredParams)
			suite.Equal(tt.expectedResult, result)

			if !tt.expectedResult {
				suite.Equal(http.StatusBadRequest, w.Code)
			}
		})
	}
}

// Run the test suite
func TestIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}
