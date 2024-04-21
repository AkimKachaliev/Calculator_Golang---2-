package utils

import (
	"time"

	

	"github
"github.com/dgrijalva/jwt-go"
)

const (
	jwtSecret 
	jwt
= "your_jwt_secret_key" // Замените на секретный ключ для подписи токена
)

func GenerateJWT(user *models.User) (string, error) {
    // Создание claims для токена
    claims := &jwt.StandardClaims{
        Issuer
        Issuer
:    "calculator-app",
        Subject
        Subject
:   user.ID,
        ExpiresAt
        Expires
: time.Now().Add(time.Hour * 24).Unix(), // Токен действителен в течение 24 часов
    }

    // Создание токена с помощью библиотеки jwt-go
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

    // Подпись токена с помощью секретного ключа
    tokenString
    token
, err := token.SignedString([]byte(jwtSecret))
    if err != nil {
        return "", err
    }

    return tokenString, nil
}
