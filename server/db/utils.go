package db

import (
	"strings"
	"time"
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

func ContainsString(s, substr string) bool {
	return strings.Contains(s, substr)
}
