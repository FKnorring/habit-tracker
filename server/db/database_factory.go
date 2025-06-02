package db

import (
	"flag"
	"fmt"
	"log"
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

	flag.StringVar(&config.Driver, "db-driver", "memory", "Database driver to use (memory, sqlite)")
	flag.StringVar(&config.SQLitePath, "sqlite-path", "./habits.db", "Path to SQLite database file")

	flag.Parse()

	return config
}
