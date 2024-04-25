package auth

import (
	"github.com/gorilla/mux"
)

func NewRouter() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/api/v1/register", Register).Methods("POST")
	router.HandleFunc("/api/v1/login", Login).Methods("POST")
	return router
}
