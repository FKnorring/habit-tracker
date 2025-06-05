package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"habit-tracker/server/db"
	"habit-tracker/server/handlers"

	"github.com/stretchr/testify/suite"
)

type IntegrationTestSuite struct {
	suite.Suite
	router *handlers.Router
	server *httptest.Server
	origDB db.Database
}

func (suite *IntegrationTestSuite) SetupSuite() {
	// Store original database
	suite.origDB = handlers.Database

	// Create in-memory database for testing
	handlers.Database = db.NewMapDatabase()

	// Create router and register handlers
	suite.router = handlers.CreateRouter()
	suite.router.Handle("GET", "/habits", handlers.GetHabits)
	suite.router.Handle("POST", "/habits", handlers.CreateHabit)
	suite.router.Handle("GET", "/habits/:id", handlers.GetHabit)
	suite.router.Handle("PATCH", "/habits/:id", handlers.UpdateHabit)
	suite.router.Handle("DELETE", "/habits/:id", handlers.DeleteHabit)
	suite.router.Handle("POST", "/habits/:id/tracking", handlers.CreateTracking)
	suite.router.Handle("GET", "/habits/:id/tracking", handlers.GetTracking)

	// Create test server
	suite.server = httptest.NewServer(suite.router)
}

func (suite *IntegrationTestSuite) TearDownSuite() {
	suite.server.Close()
	// Restore original database
	handlers.Database = suite.origDB
}

func (suite *IntegrationTestSuite) SetupTest() {
	// Reset database state before each test
	handlers.Database = db.NewMapDatabase()
}

func (suite *IntegrationTestSuite) TestGetHabitsEmpty() {
	resp, err := http.Get(suite.server.URL + "/habits")
	suite.NoError(err)
	defer resp.Body.Close()

	suite.Equal(http.StatusOK, resp.StatusCode)

	var habits []db.Habit
	err = json.NewDecoder(resp.Body).Decode(&habits)
	suite.NoError(err)

	// Ensure we get an empty array, not null
	suite.NotNil(habits)
	suite.Len(habits, 0)
}

func (suite *IntegrationTestSuite) TestCreateHabit() {
	habitData := db.Habit{
		Name:        "Exercise",
		Description: "Daily workout routine",
		Frequency:   db.FrequencyDaily,
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

func (suite *IntegrationTestSuite) TestCreateHabitInvalidFrequency() {
	habitData := map[string]interface{}{
		"name":        "Exercise",
		"description": "Daily workout routine",
		"frequency":   "invalid_frequency",
		"startDate":   "2024-01-01",
	}

	jsonData, err := json.Marshal(habitData)
	suite.NoError(err)

	resp, err := http.Post(suite.server.URL+"/habits", "application/json", bytes.NewBuffer(jsonData))
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
		Frequency:   db.FrequencyDaily,
		StartDate:   "2024-01-01",
	}
	err := handlers.Database.CreateHabit(habit)
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
		Frequency:   db.FrequencyDaily,
		StartDate:   "2024-01-01",
	}
	err := handlers.Database.CreateHabit(habit)
	suite.NoError(err)

	// Update the habit
	updatedHabit := db.Habit{
		Name:        "Updated Meditation",
		Description: "Updated daily meditation practice",
		Frequency:   db.FrequencyWeekly,
		StartDate:   "2024-01-02",
	}

	jsonData, err := json.Marshal(updatedHabit)
	suite.NoError(err)

	req, err := http.NewRequest("PATCH", suite.server.URL+"/habits/test-habit-2", bytes.NewBuffer(jsonData))
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
		Frequency:   db.FrequencyHourly,
		StartDate:   "2024-01-01",
	}

	jsonData, err := json.Marshal(updatedHabit)
	suite.NoError(err)

	req, err := http.NewRequest("PATCH", suite.server.URL+"/habits/nonexistent", bytes.NewBuffer(jsonData))
	suite.NoError(err)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	suite.NoError(err)
	defer resp.Body.Close()

	suite.Equal(http.StatusNotFound, resp.StatusCode)
}

func (suite *IntegrationTestSuite) TestUpdateHabitInvalidFrequency() {
	// First create a habit
	habit := &db.Habit{
		ID:          "test-habit-invalid-freq",
		Name:        "Test Habit",
		Description: "Test description",
		Frequency:   db.FrequencyDaily,
		StartDate:   "2024-01-01",
	}
	err := handlers.Database.CreateHabit(habit)
	suite.NoError(err)

	// Try to update with invalid frequency
	updateData := map[string]interface{}{
		"name":        "Updated Test Habit",
		"description": "Updated description",
		"frequency":   "invalid_frequency",
		"startDate":   "2024-01-01",
	}

	jsonData, err := json.Marshal(updateData)
	suite.NoError(err)

	req, err := http.NewRequest("PATCH", suite.server.URL+"/habits/test-habit-invalid-freq", bytes.NewBuffer(jsonData))
	suite.NoError(err)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	suite.NoError(err)
	defer resp.Body.Close()

	suite.Equal(http.StatusBadRequest, resp.StatusCode)
}

func (suite *IntegrationTestSuite) TestUpdateHabitPartial() {
	// First create a habit
	originalHabit := &db.Habit{
		ID:          "test-habit-partial",
		Name:        "Original Name",
		Description: "Original Description",
		Frequency:   db.FrequencyDaily,
		StartDate:   "2024-01-01",
	}
	err := handlers.Database.CreateHabit(originalHabit)
	suite.NoError(err)

	// Test updating only the name field
	updateData := map[string]interface{}{
		"name": "Updated Name Only",
	}

	jsonData, err := json.Marshal(updateData)
	suite.NoError(err)

	req, err := http.NewRequest("PATCH", suite.server.URL+"/habits/test-habit-partial", bytes.NewBuffer(jsonData))
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

	// Verify only the name was updated, other fields should remain unchanged
	suite.Equal("test-habit-partial", responseHabit.ID)
	suite.Equal("Updated Name Only", responseHabit.Name)
	suite.Equal("Original Description", responseHabit.Description) // Should not change
	suite.Equal(db.FrequencyDaily, responseHabit.Frequency)        // Should not change
	suite.Equal("2024-01-01", responseHabit.StartDate)             // Should not change

	// Test updating only the frequency field
	updateData2 := map[string]interface{}{
		"frequency": "weekly",
	}

	jsonData2, err := json.Marshal(updateData2)
	suite.NoError(err)

	req2, err := http.NewRequest("PATCH", suite.server.URL+"/habits/test-habit-partial", bytes.NewBuffer(jsonData2))
	suite.NoError(err)
	req2.Header.Set("Content-Type", "application/json")

	resp2, err := client.Do(req2)
	suite.NoError(err)
	defer resp2.Body.Close()

	suite.Equal(http.StatusOK, resp2.StatusCode)

	var responseHabit2 db.Habit
	err = json.NewDecoder(resp2.Body).Decode(&responseHabit2)
	suite.NoError(err)

	// Verify only the frequency was updated
	suite.Equal("test-habit-partial", responseHabit2.ID)
	suite.Equal("Updated Name Only", responseHabit2.Name)           // Should remain from previous update
	suite.Equal("Original Description", responseHabit2.Description) // Should still not change
	suite.Equal(db.FrequencyWeekly, responseHabit2.Frequency)       // Should be updated
	suite.Equal("2024-01-01", responseHabit2.StartDate)             // Should still not change

	// Test updating with empty body (no fields should change)
	emptyData := map[string]interface{}{}
	jsonData3, err := json.Marshal(emptyData)
	suite.NoError(err)

	req3, err := http.NewRequest("PATCH", suite.server.URL+"/habits/test-habit-partial", bytes.NewBuffer(jsonData3))
	suite.NoError(err)
	req3.Header.Set("Content-Type", "application/json")

	resp3, err := client.Do(req3)
	suite.NoError(err)
	defer resp3.Body.Close()

	suite.Equal(http.StatusOK, resp3.StatusCode)

	var responseHabit3 db.Habit
	err = json.NewDecoder(resp3.Body).Decode(&responseHabit3)
	suite.NoError(err)

	// All fields should remain the same as the previous state
	suite.Equal("test-habit-partial", responseHabit3.ID)
	suite.Equal("Updated Name Only", responseHabit3.Name)
	suite.Equal("Original Description", responseHabit3.Description)
	suite.Equal(db.FrequencyWeekly, responseHabit3.Frequency)
	suite.Equal("2024-01-01", responseHabit3.StartDate)
}

func (suite *IntegrationTestSuite) TestDeleteHabit() {
	// First create a habit
	habit := &db.Habit{
		ID:          "test-habit-3",
		Name:        "Journaling",
		Description: "Daily journaling",
		Frequency:   db.FrequencyDaily,
		StartDate:   "2024-01-01",
	}
	err := handlers.Database.CreateHabit(habit)
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
	_, err = handlers.Database.GetHabit("test-habit-3")
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
		Frequency:   db.FrequencyDaily,
		StartDate:   "2024-01-01",
	}
	err := handlers.Database.CreateHabit(habit)
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
		Frequency:   db.FrequencyHourly,
		StartDate:   "2024-01-01",
	}
	err := handlers.Database.CreateHabit(habit)
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
		Frequency:   db.FrequencyDaily,
		StartDate:   "2024-01-01",
	}
	err := handlers.Database.CreateHabit(habit)
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
		err := handlers.Database.CreateTrackingEntry(entry)
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
		Frequency:   db.FrequencyDaily,
		StartDate:   "2024-01-01",
	}
	err := handlers.Database.CreateHabit(habit)
	suite.NoError(err)

	// Get tracking entries (should be empty)
	resp, err := http.Get(suite.server.URL + "/habits/test-habit-7/tracking")
	suite.NoError(err)
	defer resp.Body.Close()

	suite.Equal(http.StatusOK, resp.StatusCode)

	var retrievedEntries []db.TrackingEntry
	err = json.NewDecoder(resp.Body).Decode(&retrievedEntries)
	suite.NoError(err)

	// Ensure we get an empty array, not null
	suite.NotNil(retrievedEntries)
	suite.Len(retrievedEntries, 0)
}

func (suite *IntegrationTestSuite) TestFullWorkflow() {
	// Test complete workflow: Create habit -> Create tracking -> Get all -> Update -> Delete

	// 1. Create a habit
	habitData := db.Habit{
		Name:        "Full Workflow Test",
		Description: "Testing complete workflow",
		Frequency:   db.FrequencyDaily,
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
		Frequency:   db.FrequencyWeekly,
		StartDate:   "2024-01-02",
	}

	jsonData, err = json.Marshal(updatedHabit)
	suite.NoError(err)

	req, err := http.NewRequest("PATCH", suite.server.URL+"/habits/"+habitID, bytes.NewBuffer(jsonData))
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

	// Ensure we get an empty array, not null
	suite.NotNil(habits)
	suite.Len(habits, 0)
}

func (suite *IntegrationTestSuite) TestParameterValidation() {
	// Test that endpoints properly validate required path parameters
	// This effectively tests the internal checkParams functionality

	tests := []struct {
		name           string
		method         string
		url            string
		expectedStatus int
	}{
		{
			name:           "GET all habits (valid route)",
			method:         "GET",
			url:            "/habits",
			expectedStatus: http.StatusOK, // This should work - gets all habits
		},
		{
			name:           "GET habit with valid ID format",
			method:         "GET",
			url:            "/habits/valid-id",
			expectedStatus: http.StatusNotFound, // ID doesn't exist, but parameter is valid
		},
		{
			name:           "DELETE habit with valid ID format",
			method:         "DELETE",
			url:            "/habits/valid-id",
			expectedStatus: http.StatusNotFound, // ID doesn't exist, but parameter is valid
		},
		{
			name:           "GET tracking with valid habit ID format",
			method:         "GET",
			url:            "/habits/valid-id/tracking",
			expectedStatus: http.StatusOK, // Should return empty array for non-existent habit
		},
		{
			name:           "GET invalid route",
			method:         "GET",
			url:            "/invalid-route",
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			var resp *http.Response
			var err error

			switch tt.method {
			case "GET":
				resp, err = http.Get(suite.server.URL + tt.url)
			case "DELETE":
				req, reqErr := http.NewRequest("DELETE", suite.server.URL+tt.url, nil)
				suite.NoError(reqErr)
				client := &http.Client{}
				resp, err = client.Do(req)
			}

			suite.NoError(err)
			defer resp.Body.Close()
			suite.Equal(tt.expectedStatus, resp.StatusCode)
		})
	}
}

// Run the test suite
func TestIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}
