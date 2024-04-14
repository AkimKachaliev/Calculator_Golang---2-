package main

import (
	"github.com/gorilla/mux"
	"github.com/yourusername/yourproject/auth"
	"github.com/yourusername/yourproject/server"
	"net/http"
)

func main() {
	r := mux.NewRouter()

	// Register the routes
	r.HandleFunc("/api/v1/register", auth.Register).Methods("POST")
	r.HandleFunc("/api/v1/login", auth.Login).Methods("POST")
	r.HandleFunc("/api/v1/calculator", server.Calculator).Methods("POST")

	// Start the HTTP server
	http.ListenAndServe(":8080", r)
}
