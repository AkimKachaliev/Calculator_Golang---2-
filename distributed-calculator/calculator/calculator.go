package main

import (
	"encoding/json"
	"net/http"

	"github.com/AkimKachaliev/distributed-calculator/distributed-calculator/models"
	"github.com/golang-jwt/jwt"
	"github.com/gorilla/mux"
)

func Calculate(w http.ResponseWriter, r *http.Request) {
	var req models.СalculateRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	token, err := jwt.Parse(r.Header.Get("Authorization"), func(token *jwt.Token) (interface{}, error) {
		return []byte("secret"), nil
	})
	if err != nil || !token.Valid {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	claims := token.Claims.(jwt.MapClaims)
	userID := int64(claims["user_id"].(float64))

	if userID != req.UserID {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	result, err := Evaluate(req.Expression)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	expression := &models.Expression{
		UserID:     req.UserID,
		Expression: req.Expression,
		Result:     result,
	}

	err = expression.Create()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(models.CalculateResponse{
		Result: result,
	})
}

func Evaluate(expression string) (string, error) {
	// Implement expression evaluation logic here
	return "42", nil
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/api/v1/calculate", Calculate).Methods("POST")
	http.ListenAndServe(":8080", r)
}
