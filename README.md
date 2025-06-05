# Habit Tracker

A full-stack habit tracking application with a Go backend and Next.js frontend.

## Quick Start with Docker

### Prerequisites

- Docker
- Docker Compose

### Running the Application

1. **Clone the repository and navigate to the project directory**

2. **Start the application:**
   ```bash
   docker-compose up --build
   ```

3. **Access the application:**
   - **Frontend:** http://localhost:3000
   - **API:** http://localhost:8080

4. **Stop the application:**
   ```bash
   docker-compose down
   ```

5. **Stop and remove all data (including database):**
   ```bash
   docker-compose down -v
   ```

## Architecture

- **Backend:** Go server serving the data model RESTfully
   - **Custom Router:** Pattern-based HTTP router with parameter extraction
   - **Injectable Database:** Interface-based database abstraction for interchangeability
   - **WebSocket Service:** Real-time communication for notifications and updates
   - **Reminder Service:** Automated habit reminders with configurable frequency-based scheduling
- **Frontend:** Next.js with Tailwind CSS

## Injectable Database System

The application uses an interface-based database abstraction that allows for easy swapping of database implementations:

### Database Interface

All database operations are defined through the `Database` interface in `models.go`, providing methods for:
- **Habit Management:** CRUD operations for habits
- **Tracking Entries:** CRUD operations for habit tracking entries
- **Reminders:** Management of habit reminder schedules
- **Statistics & Analytics:** Comprehensive habit tracking analytics
- **Health Checks:** Database connectivity verification

### Current Implementations

- **SQLite Database:** Persistent storage with file-based SQLite
- **In-Memory Database:** Fast, temporary storage for testing and development

## WebSocket Service

The WebSocket service provides real-time communication between the server and frontend clients:

- **Real-time Notifications:** Instant delivery of habit reminders and updates
- **User Authentication:** Client authentication with user ID mapping
- **Message Broadcasting:** Support for both broadcast and targeted user messaging
- **Connection Management:** Automatic client lifecycle management with cleanup

## Reminder Service

The reminder service automatically monitors habits and sends timely notifications:

- **Frequency-based Scheduling:** Intelligent reminder timing based on habit frequency (hourly, daily, weekly, etc.)
- **Background Processing:** Runs continuously with configurable check intervals
- **WebSocket Integration:** Delivers reminders via real-time WebSocket connections
- **Automatic Updates:** Updates reminder timestamps when habits are completed

## API Endpoints

### Core Habit Management
- `GET /habits` - Get all habits
- `GET /habits/:id` - Get a specific habit
- `POST /habits` - Create a new habit
- `PATCH /habits/:id` - Update a habit
- `DELETE /habits/:id` - Delete a habit

### Habit Tracking
- `POST /habits/:id/tracking` - Add tracking entry
- `GET /habits/:id/tracking` - Get tracking entries for a habit

### Reminders
- `PATCH /reminders/:id` - Update reminder last reminder timestamp

### Statistics & Analytics
- `GET /habits/:id/stats` - Get comprehensive statistics for a specific habit (streaks, completion rate, total entries)
- `GET /habits/:id/progress` - Get daily progress data for a habit over specified time period (supports `?days=N` query parameter)
- `GET /stats/overview` - Get overall application statistics (total habits, entries, daily averages)
- `GET /stats/completion-rates` - Get habit completion rates compared to expected frequency (supports `?days=N` query parameter)
- `GET /stats/daily-completions` - Get daily completion counts across all habits (supports `?days=N` query parameter)

### Real-time Communication
- `WS /ws` - WebSocket endpoint for real-time notifications and updates

## Data Models

### Habit
- `id`: *string* (UUID) - Unique identifier for the habit
- `name`: *string* - Name of the habit
- `description`: *string* - Detailed description of the habit
- `frequency`: *string* - How often the habit should be performed (hourly, daily, weekly, biweekly, monthly, quarterly, yearly)
- `startDate`: *datetime* - When the habit tracking started

### TrackingEntry
- `id`: *string* (UUID) - Unique identifier for the tracking entry
- `habitId`: *string* (UUID) - Reference to the associated habit
- `timestamp`: *datetime* - When the habit was completed
- `note`: *string* - Optional note about the completion

### Reminder
- `id`: *string* (UUID) - Unique identifier for the reminder
- `habitId`: *string* (UUID) - Reference to the associated habit
- `lastReminder`: *datetime* - Timestamp of the last reminder sent for this habit

