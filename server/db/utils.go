package db

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

func CalculateNextReminderTime(lastReminder time.Time, frequency Frequency) time.Time {
	switch frequency {
	case FrequencyHourly:
		return lastReminder.Add(time.Hour)
	case FrequencyDaily:
		return lastReminder.AddDate(0, 0, 1)
	case FrequencyWeekly:
		return lastReminder.AddDate(0, 0, 7)
	case FrequencyBiweekly:
		return lastReminder.AddDate(0, 0, 14)
	case FrequencyMonthly:
		return lastReminder.AddDate(0, 1, 0)
	case FrequencyQuarterly:
		return lastReminder.AddDate(0, 3, 0)
	case FrequencyYearly:
		return lastReminder.AddDate(1, 0, 0)
	default:
		return lastReminder.AddDate(0, 0, 1) // Default to daily
	}
}

func ContainsString(str, substr string) bool {
	return strings.Contains(strings.ToLower(str), strings.ToLower(substr))
}

// ParseCSV parses a comma-separated string into a slice
func ParseCSV(csv string) []string {
	if csv == "" {
		return []string{}
	}
	parts := strings.Split(csv, ",")
	result := make([]string, len(parts))
	for i, part := range parts {
		result[i] = strings.TrimSpace(part)
	}
	return result
}

// generateUUID generates a new UUID string
func generateUUID() string {
	return uuid.New().String()
}
