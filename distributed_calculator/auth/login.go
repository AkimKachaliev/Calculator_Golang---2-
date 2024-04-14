package auth

import (
	"encoding/json"
	"net/http"

	"github.com/yourusername/yourproject/models"
	"github.com/yourusername/yourproject/utils"
)

func Login(w http.ResponseWriter, r *http.Request) {
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

	// Find the user by login
	user, err := models.FindUserByLogin(input.Login)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	// Compare the provided password with the stored hash
	if !user.CheckPassword(input.Password) {
		http.Error(w, "Invalid password", http.StatusUnauthorized)
		return
	}

	// Generate a JWT token
	tokenString, err := utils.GenerateJWT(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return the token in the response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": tokenString})
}
