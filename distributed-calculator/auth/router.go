package auth

import (
	"github.com/gorilla/mux"
)

func NewRouter() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/register", Register).Methods("GET", "POST")
	router.HandleFunc("/login", Login).Methods("GET", "POST")
	return router
}
