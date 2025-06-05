package db

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type SQLiteDatabase struct {
	db *sql.DB
}

func NewSQLiteDatabase(dbPath string) (*SQLiteDatabase, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	sqliteDB := &SQLiteDatabase{db: db}

	if err := sqliteDB.createTables(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to create tables: %w", err)
	}

	return sqliteDB, nil
}

func (db *SQLiteDatabase) createTables() error {
	createHabitsTable := `
		CREATE TABLE IF NOT EXISTS habits (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL,
			description TEXT,
			frequency TEXT,
			start_date TEXT
		);
	`

	createTrackingTable := `
		CREATE TABLE IF NOT EXISTS tracking_entries (
			id TEXT PRIMARY KEY,
			habit_id TEXT NOT NULL,
			timestamp TEXT NOT NULL,
			note TEXT,
			FOREIGN KEY (habit_id) REFERENCES habits(id) ON DELETE CASCADE
		);
	`

	createRemindersTable := `
		CREATE TABLE IF NOT EXISTS reminders (
			id TEXT PRIMARY KEY,
			habit_id TEXT NOT NULL UNIQUE,
			last_reminder TEXT NOT NULL,
			FOREIGN KEY (habit_id) REFERENCES habits(id) ON DELETE CASCADE
		);
	`

	if _, err := db.db.Exec(createHabitsTable); err != nil {
		return fmt.Errorf("failed to create habits table: %w", err)
	}

	if _, err := db.db.Exec(createTrackingTable); err != nil {
		return fmt.Errorf("failed to create tracking_entries table: %w", err)
	}

	if _, err := db.db.Exec(createRemindersTable); err != nil {
		return fmt.Errorf("failed to create reminders table: %w", err)
	}

	return nil
}

func (db *SQLiteDatabase) Ping() error {
	return db.db.Ping()
}

func (db *SQLiteDatabase) CreateHabit(habit *Habit) error {
	tx, err := db.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	habitQuery := `
		INSERT INTO habits (id, name, description, frequency, start_date)
		VALUES (?, ?, ?, ?, ?)
	`

	_, err = tx.Exec(habitQuery, habit.ID, habit.Name, habit.Description, habit.Frequency, habit.StartDate)
	if err != nil {
		if sqliteError, ok := err.(interface{ Error() string }); ok {
			if ContainsString(sqliteError.Error(), "UNIQUE constraint failed") {
				return ErrDuplicate
			}
		}
		return fmt.Errorf("failed to create habit: %w", err)
	}

	reminderQuery := `
		INSERT INTO reminders (id, habit_id, last_reminder)
		VALUES (?, ?, ?)
	`

	now := time.Now().Format(time.RFC3339)
	_, err = tx.Exec(reminderQuery, habit.ID+"-reminder", habit.ID, now)
	if err != nil {
		return fmt.Errorf("failed to create reminder: %w", err)
	}

	return tx.Commit()
}

func (db *SQLiteDatabase) GetHabit(id string) (*Habit, error) {
	query := `SELECT id, name, description, frequency, start_date FROM habits WHERE id = ?`

	habit := &Habit{}
	var frequencyStr string
	err := db.db.QueryRow(query, id).Scan(
		&habit.ID, &habit.Name, &habit.Description, &frequencyStr, &habit.StartDate,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to get habit: %w", err)
	}

	habit.Frequency = Frequency(frequencyStr)
	return habit, nil
}

func (db *SQLiteDatabase) GetAllHabits() ([]*Habit, error) {
	query := `SELECT id, name, description, frequency, start_date FROM habits`

	rows, err := db.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query habits: %w", err)
	}
	defer rows.Close()

	var habits []*Habit
	for rows.Next() {
		habit := &Habit{}
		var frequencyStr string
		err := rows.Scan(&habit.ID, &habit.Name, &habit.Description, &frequencyStr, &habit.StartDate)
		if err != nil {
			return nil, fmt.Errorf("failed to scan habit: %w", err)
		}
		habit.Frequency = Frequency(frequencyStr)
		habits = append(habits, habit)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating habits: %w", err)
	}

	return habits, nil
}

func (db *SQLiteDatabase) UpdateHabit(habit *Habit) error {
	query := `
		UPDATE habits 
		SET name = ?, description = ?, frequency = ?, start_date = ?
		WHERE id = ?
	`

	result, err := db.db.Exec(query, habit.Name, habit.Description, habit.Frequency, habit.StartDate, habit.ID)
	if err != nil {
		return fmt.Errorf("failed to update habit: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}

func (db *SQLiteDatabase) UpdateHabitPartial(id string, updates map[string]interface{}) (*Habit, error) {
	// First check if habit exists
	existing, err := db.GetHabit(id)
	if err != nil {
		return nil, err
	}

	// Build dynamic update query
	setParts := []string{}
	args := []interface{}{}

	// Map JSON field names to database column names
	fieldMap := map[string]string{
		"name":        "name",
		"description": "description",
		"frequency":   "frequency",
		"startDate":   "start_date",
	}

	for jsonField, value := range updates {
		if dbField, ok := fieldMap[jsonField]; ok {
			// Validate frequency if it's being updated
			if jsonField == "frequency" {
				if freqStr, ok := value.(string); ok {
					if err := ValidateFrequency(freqStr); err != nil {
						return nil, err
					}
				}
			}
			setParts = append(setParts, dbField+" = ?")
			args = append(args, value)
		}
	}

	if len(setParts) == 0 {
		// No valid fields to update, return existing habit
		return existing, nil
	}

	// Add the ID parameter for the WHERE clause
	args = append(args, id)

	// Build the query by joining the SET parts
	setClause := ""
	for i, part := range setParts {
		if i > 0 {
			setClause += ", "
		}
		setClause += part
	}

	query := fmt.Sprintf("UPDATE habits SET %s WHERE id = ?", setClause)

	result, err := db.db.Exec(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to update habit: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return nil, ErrNotFound
	}

	// Return the updated habit
	return db.GetHabit(id)
}

func (db *SQLiteDatabase) DeleteHabit(id string) error {
	query := `DELETE FROM habits WHERE id = ?`

	result, err := db.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete habit: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}

func (db *SQLiteDatabase) CreateTrackingEntry(entry *TrackingEntry) error {
	query := `
		INSERT INTO tracking_entries (id, habit_id, timestamp, note)
		VALUES (?, ?, ?, ?)
	`

	_, err := db.db.Exec(query, entry.ID, entry.HabitID, entry.Timestamp, entry.Note)
	if err != nil {
		if sqliteError, ok := err.(interface{ Error() string }); ok {
			if ContainsString(sqliteError.Error(), "UNIQUE constraint failed") {
				return ErrDuplicate
			}
		}
		return fmt.Errorf("failed to create tracking entry: %w", err)
	}

	return nil
}

func (db *SQLiteDatabase) GetTrackingEntry(id string) (*TrackingEntry, error) {
	query := `SELECT id, habit_id, timestamp, note FROM tracking_entries WHERE id = ?`

	entry := &TrackingEntry{}
	err := db.db.QueryRow(query, id).Scan(
		&entry.ID, &entry.HabitID, &entry.Timestamp, &entry.Note,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to get tracking entry: %w", err)
	}

	return entry, nil
}

func (db *SQLiteDatabase) GetTrackingEntriesByHabitID(habitID string) ([]*TrackingEntry, error) {
	query := `SELECT id, habit_id, timestamp, note FROM tracking_entries WHERE habit_id = ? ORDER BY timestamp DESC`

	rows, err := db.db.Query(query, habitID)
	if err != nil {
		return nil, fmt.Errorf("failed to query tracking entries: %w", err)
	}
	defer rows.Close()

	var entries []*TrackingEntry
	for rows.Next() {
		entry := &TrackingEntry{}
		err := rows.Scan(&entry.ID, &entry.HabitID, &entry.Timestamp, &entry.Note)
		if err != nil {
			return nil, fmt.Errorf("failed to scan tracking entry: %w", err)
		}
		entries = append(entries, entry)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating tracking entries: %w", err)
	}

	return entries, nil
}

func (db *SQLiteDatabase) DeleteTrackingEntry(id string) error {
	query := `DELETE FROM tracking_entries WHERE id = ?`

	result, err := db.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete tracking entry: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}

func (db *SQLiteDatabase) CreateReminder(reminder *Reminder) error {
	query := `
		INSERT INTO reminders (id, habit_id, last_reminder)
		VALUES (?, ?, ?)
	`

	_, err := db.db.Exec(query, reminder.ID, reminder.HabitID, reminder.LastReminder)
	if err != nil {
		if sqliteError, ok := err.(interface{ Error() string }); ok {
			if ContainsString(sqliteError.Error(), "UNIQUE constraint failed") {
				return ErrDuplicate
			}
		}
		return fmt.Errorf("failed to create reminder: %w", err)
	}

	return nil
}

func (db *SQLiteDatabase) GetReminder(habitID string) (*Reminder, error) {
	query := `SELECT id, habit_id, last_reminder FROM reminders WHERE habit_id = ?`

	reminder := &Reminder{}
	err := db.db.QueryRow(query, habitID).Scan(
		&reminder.ID, &reminder.HabitID, &reminder.LastReminder,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to get reminder: %w", err)
	}

	return reminder, nil
}

func (db *SQLiteDatabase) UpdateReminderLastReminder(habitID string, lastReminder string) error {
	query := `UPDATE reminders SET last_reminder = ? WHERE habit_id = ?`

	result, err := db.db.Exec(query, lastReminder, habitID)
	if err != nil {
		return fmt.Errorf("failed to update reminder: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}

func (db *SQLiteDatabase) GetHabitsNeedingReminders() ([]*Habit, error) {
	query := `
		SELECT h.id, h.name, h.description, h.frequency, h.start_date, r.last_reminder
		FROM habits h
		JOIN reminders r ON h.id = r.habit_id
	`

	rows, err := db.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query habits with reminders: %w", err)
	}
	defer rows.Close()

	var needingReminders []*Habit
	now := time.Now()

	for rows.Next() {
		habit := &Habit{}
		var frequencyStr, lastReminderStr string
		err := rows.Scan(&habit.ID, &habit.Name, &habit.Description, &frequencyStr, &habit.StartDate, &lastReminderStr)
		if err != nil {
			return nil, fmt.Errorf("failed to scan habit: %w", err)
		}

		habit.Frequency = Frequency(frequencyStr)

		lastReminder, err := time.Parse(time.RFC3339, lastReminderStr)
		if err != nil {
			continue
		}

		nextReminderTime := CalculateNextReminderTime(lastReminder, habit.Frequency)
		log.Println("nextReminderTime", fmt.Sprintf("%v", nextReminderTime))
		if now.After(nextReminderTime) {
			needingReminders = append(needingReminders, habit)
		}
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating habits: %w", err)
	}

	return needingReminders, nil
}

func (db *SQLiteDatabase) DeleteReminder(habitID string) error {
	query := `DELETE FROM reminders WHERE habit_id = ?`

	result, err := db.db.Exec(query, habitID)
	if err != nil {
		return fmt.Errorf("failed to delete reminder: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}

// Statistics and Analytics Methods

func (db *SQLiteDatabase) GetHabitStats(habitID string) (*HabitStats, error) {
	// Get basic habit info
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

	// Get total entries count
	countQuery := `SELECT COUNT(*) FROM tracking_entries WHERE habit_id = ?`
	err = db.db.QueryRow(countQuery, habitID).Scan(&stats.TotalEntries)
	if err != nil {
		return nil, fmt.Errorf("failed to get total entries: %w", err)
	}

	// Get current streak
	stats.CurrentStreak = db.calculateCurrentStreak(habitID, habit.Frequency)

	// Get longest streak
	stats.LongestStreak = db.calculateLongestStreak(habitID, habit.Frequency)

	// Get completion rate
	stats.CompletionRate = db.calculateCompletionRate(habitID, habit.Frequency, habit.StartDate)

	// Get last completed date
	lastQuery := `SELECT MAX(timestamp) FROM tracking_entries WHERE habit_id = ?`
	var lastCompleted sql.NullString
	err = db.db.QueryRow(lastQuery, habitID).Scan(&lastCompleted)
	if err == nil && lastCompleted.Valid {
		stats.LastCompleted = lastCompleted.String
	}

	return stats, nil
}

func (db *SQLiteDatabase) GetHabitProgress(habitID string, days int) ([]*ProgressPoint, error) {
	startDate := time.Now().AddDate(0, 0, -days).Format("2006-01-02")

	query := `
		SELECT DATE(timestamp) as date, COUNT(*) as count
		FROM tracking_entries 
		WHERE habit_id = ? AND DATE(timestamp) >= ?
		GROUP BY DATE(timestamp)
		ORDER BY DATE(timestamp)
	`

	rows, err := db.db.Query(query, habitID, startDate)
	if err != nil {
		return nil, fmt.Errorf("failed to query progress: %w", err)
	}
	defer rows.Close()

	var progress []*ProgressPoint
	for rows.Next() {
		point := &ProgressPoint{}
		err := rows.Scan(&point.Date, &point.Count)
		if err != nil {
			return nil, fmt.Errorf("failed to scan progress point: %w", err)
		}
		progress = append(progress, point)
	}

	return progress, nil
}

func (db *SQLiteDatabase) GetOverallStats() (*OverallStats, error) {
	stats := &OverallStats{}

	// Total habits
	habitsQuery := `SELECT COUNT(*) FROM habits`
	err := db.db.QueryRow(habitsQuery).Scan(&stats.TotalHabits)
	if err != nil {
		return nil, fmt.Errorf("failed to get total habits: %w", err)
	}

	// Total entries
	entriesQuery := `SELECT COUNT(*) FROM tracking_entries`
	err = db.db.QueryRow(entriesQuery).Scan(&stats.TotalEntries)
	if err != nil {
		return nil, fmt.Errorf("failed to get total entries: %w", err)
	}

	// Entries today
	todayQuery := `SELECT COUNT(*) FROM tracking_entries WHERE DATE(timestamp) = DATE('now')`
	err = db.db.QueryRow(todayQuery).Scan(&stats.EntriesToday)
	if err != nil {
		return nil, fmt.Errorf("failed to get today's entries: %w", err)
	}

	// Entries this week
	weekQuery := `
		SELECT COUNT(*) FROM tracking_entries 
		WHERE DATE(timestamp) >= DATE('now', '-6 days')
	`
	err = db.db.QueryRow(weekQuery).Scan(&stats.EntriesThisWeek)
	if err != nil {
		return nil, fmt.Errorf("failed to get this week's entries: %w", err)
	}

	// Average entries per day (last 30 days)
	avgQuery := `
		SELECT CAST(COUNT(*) AS FLOAT) / 30 FROM tracking_entries 
		WHERE DATE(timestamp) >= DATE('now', '-30 days')
	`
	err = db.db.QueryRow(avgQuery).Scan(&stats.AvgEntriesPerDay)
	if err != nil {
		return nil, fmt.Errorf("failed to get average entries: %w", err)
	}

	return stats, nil
}

func (db *SQLiteDatabase) GetHabitCompletionRates(days int) ([]*HabitCompletionRate, error) {
	startDate := time.Now().AddDate(0, 0, -days).Format("2006-01-02")

	query := `
		SELECT h.id, h.name, h.frequency, h.start_date,
			   COUNT(te.id) as actual_completions
		FROM habits h
		LEFT JOIN tracking_entries te ON h.id = te.habit_id 
			AND DATE(te.timestamp) >= ?
		GROUP BY h.id, h.name, h.frequency, h.start_date
		ORDER BY h.name
	`

	rows, err := db.db.Query(query, startDate)
	if err != nil {
		return nil, fmt.Errorf("failed to query completion rates: %w", err)
	}
	defer rows.Close()

	var rates []*HabitCompletionRate
	for rows.Next() {
		rate := &HabitCompletionRate{}
		var frequencyStr string
		err := rows.Scan(&rate.HabitID, &rate.HabitName, &frequencyStr, &rate.StartDate, &rate.ActualCompletions)
		if err != nil {
			return nil, fmt.Errorf("failed to scan completion rate: %w", err)
		}

		rate.Frequency = Frequency(frequencyStr)
		rate.ExpectedCompletions = db.calculateExpectedCompletions(Frequency(frequencyStr), days)

		if rate.ExpectedCompletions > 0 {
			rate.CompletionRate = float64(rate.ActualCompletions) / float64(rate.ExpectedCompletions)
		}

		rates = append(rates, rate)
	}

	return rates, nil
}

func (db *SQLiteDatabase) GetDailyCompletions(days int) ([]*DailyCompletion, error) {
	startDate := time.Now().AddDate(0, 0, -days).Format("2006-01-02")

	query := `
		SELECT DATE(timestamp) as date, COUNT(*) as completions
		FROM tracking_entries 
		WHERE DATE(timestamp) >= ?
		GROUP BY DATE(timestamp)
		ORDER BY DATE(timestamp)
	`

	rows, err := db.db.Query(query, startDate)
	if err != nil {
		return nil, fmt.Errorf("failed to query daily completions: %w", err)
	}
	defer rows.Close()

	var completions []*DailyCompletion
	for rows.Next() {
		completion := &DailyCompletion{}
		err := rows.Scan(&completion.Date, &completion.Completions)
		if err != nil {
			return nil, fmt.Errorf("failed to scan daily completion: %w", err)
		}
		completions = append(completions, completion)
	}

	return completions, nil
}

// Helper methods for calculations

func (db *SQLiteDatabase) calculateCurrentStreak(habitID string, frequency Frequency) int {
	// Implementation depends on frequency - for now, let's do daily streaks
	query := `
		SELECT DATE(timestamp) FROM tracking_entries 
		WHERE habit_id = ? 
		ORDER BY timestamp DESC
	`

	rows, err := db.db.Query(query, habitID)
	if err != nil {
		return 0
	}
	defer rows.Close()

	var dates []string
	for rows.Next() {
		var date string
		if err := rows.Scan(&date); err == nil {
			dates = append(dates, date)
		}
	}

	if len(dates) == 0 {
		return 0
	}

	// Remove duplicates and sort
	uniqueDates := make(map[string]bool)
	for _, date := range dates {
		uniqueDates[date] = true
	}

	today := time.Now().Format("2006-01-02")
	currentDate, _ := time.Parse("2006-01-02", today)
	streak := 0

	for {
		dateStr := currentDate.Format("2006-01-02")
		if uniqueDates[dateStr] {
			streak++
			currentDate = currentDate.AddDate(0, 0, -1)
		} else {
			break
		}
	}

	return streak
}

func (db *SQLiteDatabase) calculateLongestStreak(habitID string, frequency Frequency) int {
	// Simplified implementation - can be enhanced based on frequency
	query := `
		SELECT DISTINCT DATE(timestamp) FROM tracking_entries 
		WHERE habit_id = ? 
		ORDER BY DATE(timestamp)
	`

	rows, err := db.db.Query(query, habitID)
	if err != nil {
		return 0
	}
	defer rows.Close()

	var dates []time.Time
	for rows.Next() {
		var dateStr string
		if err := rows.Scan(&dateStr); err == nil {
			if date, err := time.Parse("2006-01-02", dateStr); err == nil {
				dates = append(dates, date)
			}
		}
	}

	if len(dates) == 0 {
		return 0
	}

	maxStreak := 1
	currentStreak := 1

	for i := 1; i < len(dates); i++ {
		if dates[i].Sub(dates[i-1]).Hours() == 24 {
			currentStreak++
			if currentStreak > maxStreak {
				maxStreak = currentStreak
			}
		} else {
			currentStreak = 1
		}
	}

	return maxStreak
}

func (db *SQLiteDatabase) calculateCompletionRate(habitID string, frequency Frequency, startDate string) float64 {
	start, err := time.Parse("2006-01-02", startDate[:10])
	if err != nil {
		return 0.0
	}

	daysSinceStart := int(time.Since(start).Hours() / 24)
	if daysSinceStart <= 0 {
		return 0.0
	}

	query := `SELECT COUNT(*) FROM tracking_entries WHERE habit_id = ?`
	var actualCompletions int
	if err := db.db.QueryRow(query, habitID).Scan(&actualCompletions); err != nil {
		return 0.0
	}

	expectedCompletions := db.calculateExpectedCompletions(frequency, daysSinceStart)
	if expectedCompletions == 0 {
		return 0.0
	}

	return float64(actualCompletions) / float64(expectedCompletions)
}

func (db *SQLiteDatabase) calculateExpectedCompletions(frequency Frequency, days int) int {
	switch frequency {
	case FrequencyDaily:
		return days
	case FrequencyWeekly:
		return days / 7
	case FrequencyBiweekly:
		return days / 14
	case FrequencyMonthly:
		return days / 30
	case FrequencyQuarterly:
		return days / 90
	case FrequencyYearly:
		return days / 365
	case FrequencyHourly:
		return days * 24 // Assuming once per hour per day
	default:
		return days // Default to daily
	}
}
