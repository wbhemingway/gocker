package main

import (
	"database/sql"
	"log"

	"github.com/wbhemingway/gocker/internal/cli"
	"github.com/wbhemingway/gocker/internal/db"
	"github.com/wbhemingway/gocker/internal/engine"
	_ "modernc.org/sqlite"
)

func main() {
	database, err := sql.Open("sqlite", "gocker.db?_pragma=foreign_keys(1)&_pragma=journal_mode(WAL)")
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}
	defer database.Close()

	queries := db.New(database)
	appEngine := engine.NewEngine(queries, database)
	cli.Init(appEngine)
	cli.Execute()
}
