package main

import (
	_ "database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	"log"
	"strings"
)

func handleGetErrors(config *Config) context.Handler {
	return func(ctx iris.Context) {
		var builder strings.Builder
		for _, desc := range config.ErrorDescriptions {
			builder.WriteString(fmt.Sprintf("%s\n", desc))
		}
		ctx.Header("Content-Type", "text/markdown")
		ctx.Text(builder.String())
	}
}

func handleGetLogs(db *sqlx.DB) context.Handler {
	return func(ctx *context.Context) {
		limit := ctx.URLParamIntDefault("limit", 10)  // default limit is 10
		offset := ctx.URLParamIntDefault("offset", 0) // default offset is 0
		// Initialize a slice to hold the retrieved log entries.
		logEntries := []LogEntry{}

		// Query the database to retrieve log entries, with the given limit and offset.
		rows, err := db.Queryx("SELECT key, message, environment, app_version FROM logs LIMIT ? OFFSET ?", limit, offset)
		if err != nil {
			problem := ProblemDetails{
				Type:   "/errors#internal-server-error",
				Title:  "Internal Server Error",
				Status: iris.StatusInternalServerError,
				Detail: "Error retrieving the log entries",
			}
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.JSON(problem)
			return
		}
		defer rows.Close()

		// Iterate over the result set and append each row to the logEntries slice.
		for rows.Next() {
			var entry LogEntry
			if err := rows.StructScan(&entry); err != nil {
				log.Printf("Error scanning row: %v", err)
				continue
			}
			logEntries = append(logEntries, entry)
		}

		// Check for errors from iterating over rows.
		if err := rows.Err(); err != nil {
			problem := ProblemDetails{
				Type:   "/errors#internal-server-error",
				Title:  "Internal Server Error",
				Status: iris.StatusInternalServerError,
				Detail: "Error reading the log entries",
			}
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.JSON(problem)
			return
		}

		// Return the retrieved log entries as JSON.
		ctx.JSON(logEntries)
	}
}

func handlePostLogs(db *sqlx.DB) context.Handler {
	return func(ctx iris.Context) {
		var logEntry LogEntry
		if err := ctx.ReadJSON(&logEntry); err != nil {
			problem := ProblemDetails{
				Type:   "/errors#bad-request",
				Title:  "Bad Request",
				Status: iris.StatusBadRequest,
				Detail: "Invalid JSON format",
			}
			ctx.StatusCode(iris.StatusBadRequest)
			ctx.JSON(problem)
			return
		}
		log.Printf("Received log: %s", logEntry.String())

		err := InsertLog(db, logEntry)
		if err != nil {
			problem := ProblemDetails{
				Type:   "/errors#internal-server-error",
				Title:  "Internal Server Error",
				Status: iris.StatusInternalServerError,
				Detail: "Error saving the log entry",
			}
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.JSON(problem)
			return
		}

		ctx.StatusCode(iris.StatusOK)
		ctx.JSON(iris.Map{"status": fmt.Sprintf("Log entry received: %s", logEntry.String())})
	}
}

func RegisterRoutes(app *iris.Application, db *sqlx.DB, config *Config) {
	app.Get("/logs", handleGetLogs(db))
	app.Post("/logs", handlePostLogs(db))
	app.Get("/errors", handleGetErrors(config))
}
