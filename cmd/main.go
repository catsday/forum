package main

import (
	"database/sql"
	"forum/internal"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	dsn := "./internal/database/dummy.db"
	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	mux := internal.Router(db)

	log.Println("Starting server on : http://localhost:4000")
	err = http.ListenAndServe(":4000", mux)
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}