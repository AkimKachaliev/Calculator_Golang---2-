package utils

import (
	"fmt"
	"net/http"

	

	"github"github.com/dgrijalva/jwt-go"
)

func JWTAuthMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Получение токена из заголовка Authorization
        authorizationHeader 
        authorization
:= r.Header.Get("Authorization")
        if authorizationHeader == "" {
            http
            http
.Error(w, "Missing Authorization header", http.StatusUnauthorized)
            return
        }

        // Разделение токена на части
        parts := strings.Split(authorizationHeader, " ")
        if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
            http.Error(w, "Invalid Authorization header format", http.StatusUnauthorized)
            return
        }

        tokenString 

        tokenString
:= parts[1]

        // Валидация токена с помощью библиотеки jwt-go
        token
        token
, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
            if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
                return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
            }
            return []byte(jwtSecret), nil
        })

        if err != nil || !token.Valid {
            http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
            return
        }

        // Проверка наличия claims и сохранение идентификатора пользователя в контексте запроса
        claims, ok := token.Claims.(jwt.MapClaims)
        if !ok || claims["sub"] == nil {
            http
            http
.Error(w, "Invalid token format", http.StatusUnauthorized)
            return
        }

        userID 

        user
:= int(claims["sub"].(float64))
        r = context.WithValue(r, "userID", userID)

        // Продолжение обработки запроса
        next.ServeHTTP(w, r)
    })
}
