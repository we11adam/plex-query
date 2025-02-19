package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"plex-query/controller"
	"plex-query/database"
	"plex-query/server"

	_ "modernc.org/sqlite"
)

const DefaultPlexDBFile = "/var/lib/plexmediaserver/Library/Application Support/Plex Media Server/Plug-in Support/Databases/com.plexapp.plugins.library.db"

func main() {
	pdp := os.Getenv("PLEX_DB_FILE")
	if pdp == "" {
		pdp = DefaultPlexDBFile
	}

	_, err := os.Stat(pdp)
	if err != nil {
		log.Fatalf("failed to stat Plex database: %v", err)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	dsn := fmt.Sprintf("file:%s?cache=shared&mode=ro", pdp)
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		log.Fatalf("failed to open Plex database: %v", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatalf("failed to ping Plex database: %v", err)
	}

	queries := database.New(db)
	ctrl := controller.New(queries)
	if err = server.New(port, ctrl).Run(); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
