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

	// Create habit
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

	// Create associated reminder
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
			continue // Skip if we can't parse the time
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
