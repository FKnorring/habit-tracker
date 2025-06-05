package reminder

import (
	"encoding/json"
	"log"
	"time"

	"habit-tracker/server/db"
	"habit-tracker/server/sockets"
)

type ReminderService struct {
	database      db.Database
	ticker        *time.Ticker
	stopChan      chan bool
	checkInterval time.Duration
}

type ReminderMessage struct {
	Type string       `json:"type"`
	Data ReminderData `json:"data"`
}

type ReminderData struct {
	HabitID     string `json:"habitId"`
	HabitName   string `json:"habitName"`
	Description string `json:"description"`
	Frequency   string `json:"frequency"`
	Timestamp   string `json:"timestamp"`
}

const (
	// Default user ID until proper user management is implemented
	DefaultUserID = "user-123"

	DefaultCheckInterval = 5 * time.Minute
)

func NewReminderService(database db.Database) *ReminderService {
	return &ReminderService{
		database:      database,
		stopChan:      make(chan bool),
		checkInterval: DefaultCheckInterval,
	}
}

func (rs *ReminderService) SetCheckInterval(interval time.Duration) {
	rs.checkInterval = interval
}

func (rs *ReminderService) Start() {
	log.Printf("Starting reminder service with check interval: %v", rs.checkInterval)

	rs.ticker = time.NewTicker(rs.checkInterval)

	go rs.checkAndSendReminders()

	go func() {
		for {
			select {
			case <-rs.ticker.C:
				rs.checkAndSendReminders()
			case <-rs.stopChan:
				log.Println("Reminder service stopped")
				return
			}
		}
	}()
}

func (rs *ReminderService) Stop() {
	if rs.ticker != nil {
		rs.ticker.Stop()
	}
	rs.stopChan <- true
}

func (rs *ReminderService) checkAndSendReminders() {
	log.Println("Checking for habits needing reminders...")

	habits, err := rs.database.GetHabitsNeedingReminders()
	if err != nil {
		log.Printf("Error fetching habits needing reminders: %v", err)
		return
	}

	if len(habits) == 0 {
		log.Println("No habits need reminders at this time")
		return
	}

	for _, habit := range habits {
		if err := rs.sendReminderForHabit(habit); err != nil {
			log.Printf("Error sending reminder for habit %s (%s): %v", habit.ID, habit.Name, err)
			continue
		}
	}
}

func (rs *ReminderService) sendReminderForHabit(habit *db.Habit) error {
	reminderData := ReminderData{
		HabitID:     habit.ID,
		HabitName:   habit.Name,
		Description: habit.Description,
		Frequency:   string(habit.Frequency),
		Timestamp:   time.Now().Format(time.RFC3339),
	}

	reminderMessage := ReminderMessage{
		Type: "reminder",
		Data: reminderData,
	}

	messageBytes, err := json.Marshal(reminderMessage)
	if err != nil {
		return err
	}

	return sockets.MessageUser(DefaultUserID, messageBytes)
}
