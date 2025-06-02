package db

import (
	"database/sql"
	"fmt"
	"strings"

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

	// Create tables if they don't exist
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

	if _, err := db.db.Exec(createHabitsTable); err != nil {
		return fmt.Errorf("failed to create habits table: %w", err)
	}

	if _, err := db.db.Exec(createTrackingTable); err != nil {
		return fmt.Errorf("failed to create tracking_entries table: %w", err)
	}

	return nil
}

func (db *SQLiteDatabase) Ping() error {
	return db.db.Ping()
}

func (db *SQLiteDatabase) CreateHabit(habit *Habit) error {
	query := `
		INSERT INTO habits (id, name, description, frequency, start_date)
		VALUES (?, ?, ?, ?, ?)
	`

	_, err := db.db.Exec(query, habit.ID, habit.Name, habit.Description, habit.Frequency, habit.StartDate)
	if err != nil {
		// Check for duplicate key error (UNIQUE constraint failed)
		if sqliteError, ok := err.(interface{ Error() string }); ok {
			if containsString(sqliteError.Error(), "UNIQUE constraint failed") {
				return ErrDuplicate
			}
		}
		return fmt.Errorf("failed to create habit: %w", err)
	}

	return nil
}

func (db *SQLiteDatabase) GetHabit(id string) (*Habit, error) {
	query := `SELECT id, name, description, frequency, start_date FROM habits WHERE id = ?`

	habit := &Habit{}
	err := db.db.QueryRow(query, id).Scan(
		&habit.ID, &habit.Name, &habit.Description, &habit.Frequency, &habit.StartDate,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to get habit: %w", err)
	}

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
		err := rows.Scan(&habit.ID, &habit.Name, &habit.Description, &habit.Frequency, &habit.StartDate)
		if err != nil {
			return nil, fmt.Errorf("failed to scan habit: %w", err)
		}
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
		// Check for duplicate key error
		if sqliteError, ok := err.(interface{ Error() string }); ok {
			if containsString(sqliteError.Error(), "UNIQUE constraint failed") {
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

// Helper function to check if a string contains a substring
func containsString(s, substr string) bool {
	return strings.Contains(s, substr)
}
