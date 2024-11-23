package main

import (
	"database/sql"
	"fmt"
	"forum/internal"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func main() {

	dsn := "./internal/database/dummy.db"
	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	err = initializeDatabase(db)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	mux := internal.Router(db)

	staticPath := "ui/static"
	fs := http.FileServer(http.Dir(staticPath))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Starting server on : http://%s:%s", "localhost", port)
	err = http.ListenAndServe("localhost:"+port, mux)
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

func initializeDatabase(db *sql.DB) error {
	sqlFile := "./internal/database/init.sql"
	initSQL, err := ioutil.ReadFile(sqlFile)
	if err != nil {
		return fmt.Errorf("failed to read init.sql: %v", err)
	}

	_, err = db.Exec(string(initSQL))
	if err != nil {
		return fmt.Errorf("failed to execute init.sql: %v", err)
	}

	log.Println("Database initialized successfully.")
	return nil
}
