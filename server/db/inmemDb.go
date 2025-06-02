package db

type MapDatabase struct {
	habits   map[string]*Habit
	tracking map[string]*TrackingEntry
}

func NewMapDatabase() *MapDatabase {
	return &MapDatabase{
		habits:   make(map[string]*Habit),
		tracking: make(map[string]*TrackingEntry),
	}
}

func (db *MapDatabase) Ping() error {
	// In-memory database is always available
	return nil
}

func (db *MapDatabase) CreateHabit(habit *Habit) error {
	if _, exists := db.habits[habit.ID]; exists {
		return ErrDuplicate
	}

	habitCopy := *habit
	db.habits[habit.ID] = &habitCopy
	return nil
}

func (db *MapDatabase) GetHabit(id string) (*Habit, error) {
	habit, exists := db.habits[id]
	if !exists {
		return nil, ErrNotFound
	}

	habitCopy := *habit
	return &habitCopy, nil
}

func (db *MapDatabase) GetAllHabits() ([]*Habit, error) {
	habits := make([]*Habit, 0, len(db.habits))
	for _, habit := range db.habits {
		habitCopy := *habit
		habits = append(habits, &habitCopy)
	}
	return habits, nil
}

func (db *MapDatabase) UpdateHabit(habit *Habit) error {
	if _, exists := db.habits[habit.ID]; !exists {
		return ErrNotFound
	}

	habitCopy := *habit
	db.habits[habit.ID] = &habitCopy
	return nil
}

func (db *MapDatabase) DeleteHabit(id string) error {
	if _, exists := db.habits[id]; !exists {
		return ErrNotFound
	}

	delete(db.habits, id)
	return nil
}

func (db *MapDatabase) CreateTrackingEntry(entry *TrackingEntry) error {
	if _, exists := db.tracking[entry.ID]; exists {
		return ErrDuplicate
	}

	entryCopy := *entry
	db.tracking[entry.ID] = &entryCopy
	return nil
}

func (db *MapDatabase) GetTrackingEntry(id string) (*TrackingEntry, error) {
	entry, exists := db.tracking[id]
	if !exists {
		return nil, ErrNotFound
	}

	entryCopy := *entry
	return &entryCopy, nil
}

func (db *MapDatabase) GetTrackingEntriesByHabitID(habitID string) ([]*TrackingEntry, error) {
	var entries []*TrackingEntry
	for _, entry := range db.tracking {
		if entry.HabitID == habitID {
			entryCopy := *entry
			entries = append(entries, &entryCopy)
		}
	}
	return entries, nil
}

func (db *MapDatabase) DeleteTrackingEntry(id string) error {
	if _, exists := db.tracking[id]; !exists {
		return ErrNotFound
	}

	delete(db.tracking, id)
	return nil
}
