package main

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/kataras/iris/v12"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	configPath, err := GetConfigPath()
	if err != nil {
		log.Fatalf("Could not retrieve config path: %s", err.Error())
	}

	config, err := EnsureConfig(configPath)
	if err != nil {
		log.Fatalf("Error ensuring config: %s", err.Error())
	}

	err = EnsureDatabasePath(config)
	if err != nil {
		log.Fatalf("Error ensuring database path: %s", err.Error())
	}

	db, err := OpenDatabase(err, config)
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

	interval, err := strconv.ParseInt(EnvOrDefault("INTERVAL", "15"), 10, 64)
	if err != nil {
		interval = 15
	}

	go StartWALCheckpointer(ctx, db, time.Duration(interval))

	if err := ApplyMigrations(db); err != nil {
		log.Fatal(err)
	}

	app := iris.Default()
	RegisterRoutes(app, db, config)
	host := EnvOrDefault("HOST", "127.0.0.1")
	port := EnvOrDefault("PORT", "12345")
	if err := app.Listen(fmt.Sprintf("%s:%s", host, port)); err != nil {
		log.Fatal(err)
	}
}
