package main

import (
	"database/sql"
	"log"
	"net/url"
	"os"
	"path/filepath"

	"plex-query/controller"
	"plex-query/database"
	"plex-query/server"

	_ "modernc.org/sqlite"
)

const DefaultPlexDBFile = "/var/lib/plexmediaserver/Library/Application Support/Plex Media Server/Plug-in Support/Databases/com.plexapp.plugins.library.db"

func openPlexDB(path string) (*sql.DB, error) {
	abs, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}
	dsn := url.URL{Scheme: "file", Path: abs, RawQuery: "cache=shared&mode=ro"}
	return sql.Open("sqlite", dsn.String())
}

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

	db, err := openPlexDB(pdp)
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
