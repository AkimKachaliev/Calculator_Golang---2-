package auth

import (
	"html/template"
	"net/http"

	"github.com/AkimKachaliev/distributed-calculator/distributed-calculator/models"
	"github.com/golang-jwt/jwt"
)

func Login(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		t, err := template.ParseFiles("auth/html/login.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		t.Execute(w, nil)
		return
	}

	var user models.User
	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	user.Login = r.FormValue("login")
	user.Password = r.FormValue("password")

	foundUser, err := models.GetUser(user.Login, user.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": foundUser.ID,
	})

	tokenString, err := token.SignedString([]byte("secret"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:  "token",
		Value: tokenString,
	})

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
