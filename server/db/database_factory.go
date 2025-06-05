package db

import (
	"flag"
	"fmt"
	"log"
	"os"
)

type DatabaseConfig struct {
	Driver string
	// SQLite specific
	SQLitePath string
}

func NewDatabaseFromConfig() (Database, error) {
	config := parseFlags()

	switch config.Driver {
	case "memory":
		log.Println("Using in-memory database")
		return NewMapDatabase(), nil
	case "sqlite":
		log.Printf("Using SQLite database at: %s", config.SQLitePath)
		return NewSQLiteDatabase(config.SQLitePath)
	default:
		return nil, fmt.Errorf("unsupported database driver: %s", config.Driver)
	}
}

func parseFlags() DatabaseConfig {
	var config DatabaseConfig

	// Check for environment variables first
	dbPath := os.Getenv("DB_PATH")

	// Set defaults based on environment variables
	defaultDriver := "memory"
	defaultSQLitePath := "./habits.db"

	if dbPath != "" {
		defaultDriver = "sqlite"
		defaultSQLitePath = dbPath
	}

	flag.StringVar(&config.Driver, "db-driver", defaultDriver, "Database driver to use (memory, sqlite)")
	flag.StringVar(&config.SQLitePath, "sqlite-path", defaultSQLitePath, "Path to SQLite database file")

	flag.Parse()

	return config
}
