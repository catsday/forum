package main

import (
	"log"
	"net/http"
	"forum/internal"
)

func main() {
	log.Println("Starting server on : http://localhost:4000")
	err := http.ListenAndServe(":4000", internal.Router())
	if err != nil {
		log.Printf("Server start error %s", err)
	}
	log.Fatal(err)
}
