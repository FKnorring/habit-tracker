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
- **Frontend:** Next.js with Tailwind CSS

## Injectable Database System

The application uses an interface-based database abstraction that allows for easy swapping of database implementations:

### Database Interface

All database operations are defined through the `Database` interface in `models.go`, providing methods for:
- **Habit Management:** CRUD operations for habits
- **Tracking Entries:** CRUD operations for habit tracking entries
- **Health Checks:** Database connectivity verification

### Current Implementations

- **SQLite Database:** Persistent storage with file-based SQLite
- **In-Memory Database:** Fast, temporary storage for testing and development


## API Endpoints

- `GET /habits` - Get all habits
- `GET /habits/:id` - Get a specific habit
- `POST /habits` - Create a new habit
- `PUT /habits/:id` - Update a habit
- `DELETE /habits/:id` - Delete a habit
- `POST /habits/:id/tracking` - Add tracking entry
- `GET /habits/:id/tracking` - Get tracking entries

## Data Models

### Habit
- `id`: `string` (UUID) - Unique identifier for the habit
- `name`: `string` - Name of the habit
- `description`: `string` - Detailed description of the habit
- `frequency`: `string` - How often the habit should be performed
- `startDate`: `string` - When the habit tracking started

### TrackingEntry
- `id`: `string` (UUID) - Unique identifier for the tracking entry
- `habitId`: `string` (UUID) - Reference to the associated habit
- `timestamp`: `string` - When the habit was completed
- `note`: `string` - Optional note about the completion

