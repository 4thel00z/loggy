package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/jmoiron/sqlx"
)

func OpenDatabase(err error, config *Config) (*sqlx.DB, error) {
	db, err := sqlx.Open("sqlite3", config.DatabasePath)
	if err != nil {
		return nil, err
	}

	// Applying SQLite pragmas when opening the DB
	_, err = db.Exec(`PRAGMA journal_mode = wal;`)
	if err != nil {
		return nil, fmt.Errorf("failed to set pragma journal_mode: %w", err)
	}

	_, err = db.Exec(`PRAGMA synchronous = normal;`)
	if err != nil {
		return nil, fmt.Errorf("failed to set pragma synchronous: %w", err)
	}

	_, err = db.Exec(`PRAGMA foreign_keys = on;`)
	if err != nil {
		return nil, fmt.Errorf("failed to set pragma foreign_keys: %w", err)
	}

	return db, nil
}

// StartWALCheckpointer start the checkpointing goroutine
func StartWALCheckpointer(ctx context.Context, db *sqlx.DB, interval time.Duration) {
	ticker := time.NewTicker(interval * time.Second)
	defer ticker.Stop()
	log.Println("Starting WALCheckpointer")
	for {
		select {
		case <-ctx.Done():
			// Context has been canceled, so stop the coroutine
			fmt.Println("Stopping checkpointing coroutine.")
			return
		case <-ticker.C:

			_, err := db.Exec(`PRAGMA wal_checkpoint(FULL);`)
			if err != nil {
				fmt.Println("Error during checkpoint:", err)
			} else {
				fmt.Println("Checkpoint successful.")
			}
		}
	}
}

func ApplyMigrations(db *sqlx.DB) error {
	// Creating the table with created_at and updated_at columns
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS logs (
			id INTEGER PRIMARY KEY, 
			key TEXT, 
			message TEXT, 
			environment TEXT, 
			device_name TEXT, 
			app_version TEXT, 
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP NOT NULL, 
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP NOT NULL
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create table: %w", err)
	}

	// Creating a trigger to update the updated_at column on each update operation
	_, err = db.Exec(`
		CREATE TRIGGER IF NOT EXISTS update_updated_at
		AFTER UPDATE
			ON logs
		FOR EACH ROW
		BEGIN
			UPDATE logs SET updated_at = CURRENT_TIMESTAMP WHERE id = OLD.id;
		END;
	`)
	if err != nil {
		return fmt.Errorf("failed to create trigger: %w", err)
	}

	return nil
}

func InsertLog(db *sqlx.DB, logEntry LogEntry) error {
	_, err := db.Exec(
		`INSERT INTO logs (key, message, environment, app_version, device_name ) VALUES (?, ?, ?, ?, ?)`,
		logEntry.Key,
		logEntry.Message,
		logEntry.Environment,
		logEntry.AppVersion,
		logEntry.DeviceName,
	)
	return err
}

func EnsureDatabasePath(config *Config) error {
	if err := os.MkdirAll(filepath.Dir(config.DatabasePath), 0755); err != nil {
		return err
	}
	return nil
}
