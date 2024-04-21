package auth

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/AkimKachaliev/Calculator_Golang---2-/models"
)

func Register(w http.ResponseWriter, r *http.Request) {
	// Parse and validate the request body
	var input models.User
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Validate the input
	if input.Login == "" || input.Password == "" {
		http.Error(w, "Login and password are required", http.StatusBadRequest)
		return
	}

	// Create a new user
	user, err := models.NewUser(input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Save the user to the database
	if err := user.Save(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return a success response
	w.WriteHeader(http.StatusCreated)
