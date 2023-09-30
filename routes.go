package main

import (
	"database/sql"
	"fmt"
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

func handlePostLogs(db *sql.DB) context.Handler {
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
		log.Printf("Received log: %s => %s", logEntry.Key, logEntry.Message)

		err := insertLog(db, logEntry)
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
		ctx.JSON(iris.Map{"status": fmt.Sprintf("Log entry received %s => %s", logEntry.Key, logEntry.Message)})
	}
}

func registerRoutes(app *iris.Application, db *sql.DB, config *Config) {
	app.Post("/logs", handlePostLogs(db))
	app.Get("/errors", handleGetErrors(config))
}
