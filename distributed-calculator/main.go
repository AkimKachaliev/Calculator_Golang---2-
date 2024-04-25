package main

import (
	"log"
	"net/http"

	"github.com/AkimKachaliev/distributed-calculator/auth"
	"github.com/AkimKachaliev/distributed-calculator/config"
	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()

	// Initialize the database
	db, err := config.InitDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Register the routes
	r.HandleFunc("/api/v1/register", auth.Register).Methods("POST")
	r.HandleFunc("/api/v1/login", auth.Login).Methods("POST")

	// Start the HTTP server
	log.Println("Starting HTTP server on :8080...")
	log.Fatal(http.ListenAndServe(":8080", r))
}
