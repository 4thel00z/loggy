package main

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"os"
	"path/filepath"
)

func applyMigrations(db *sqlx.DB) error {
	// Creating the table with created_at and updated_at columns
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS logs (
			id INTEGER PRIMARY KEY, 
			key TEXT, 
			message TEXT, 
			environment TEXT, 
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

func insertLog(db *sqlx.DB, logEntry LogEntry) error {
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

func ensureDatabasePath(config *Config) error {
	if err := os.MkdirAll(filepath.Dir(config.DatabasePath), 0755); err != nil {
		return err
	}
	return nil
}
