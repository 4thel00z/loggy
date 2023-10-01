package main

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"log"

	"github.com/kataras/iris/v12"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	configPath, err := getConfigPath()
	if err != nil {
		log.Fatalf("Could not retrieve config path: %s", err.Error())
	}

	config, err := ensureConfig(configPath)
	if err != nil {
		log.Fatalf("Error ensuring config: %s", err.Error())
	}

	err = ensureDatabasePath(config)
	if err != nil {
		log.Fatalf("Error ensuring database path: %s", err.Error())
	}

	db, err := sqlx.Open("sqlite3", config.DatabasePath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := applyMigrations(db); err != nil {
		log.Fatal(err)
	}

	app := iris.Default()
	registerRoutes(app, db, config)
	host := EnvOrDefault("HOST", "127.0.0.1")
	port := EnvOrDefault("PORT", "12345")
	if err := app.Listen(fmt.Sprintf("%s:%s", host, port)); err != nil {
		log.Fatal(err)
	}
}
