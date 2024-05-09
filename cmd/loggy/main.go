package main

import (
	"context"
	"fmt"
	"github.com/4thel00z/loggy"
	"log"
	"strconv"
	"time"

	"github.com/kataras/iris/v12"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	configPath, err := loggy.GetConfigPath()
	if err != nil {
		log.Fatalf("Could not retrieve config path: %s", err.Error())
	}

	config, err := loggy.EnsureConfig(configPath)
	if err != nil {
		log.Fatalf("Error ensuring config: %s", err.Error())
	}

	err = loggy.EnsureDatabasePath(config)
	if err != nil {
		log.Fatalf("Error ensuring database path: %s", err.Error())
	}

	db, err := loggy.OpenDatabase(err, config)
	if err != nil {
		log.Fatal(err.Error())
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer func() {

		cancel()
		_, err := db.Exec(`PRAGMA wal_checkpoint(FULL);`)
		if err != nil {
			fmt.Println("Error during checkpoint:", err)
		} else {
			fmt.Println("Checkpoint successful.")
		}
		if err := db.Close(); err != nil {
			fmt.Println("Error closing database:", err.Error())
		}
	}()

	interval, err := strconv.ParseInt(loggy.EnvOrDefault("INTERVAL", "15"), 10, 64)
	if err != nil {
		interval = 15
	}

	go loggy.StartWALCheckpointer(ctx, db, time.Duration(interval))

	if err := loggy.ApplyMigrations(db); err != nil {
		log.Fatal(err)
	}

	app := iris.Default()
	loggy.RegisterRoutes(app, db, config)
	host := loggy.EnvOrDefault("HOST", "127.0.0.1")
	port := loggy.EnvOrDefault("PORT", "12345")
	if err := app.Listen(fmt.Sprintf("%s:%s", host, port)); err != nil {
		log.Fatal(err)
	}
}
