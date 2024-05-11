package main

import (
	"fmt"
	"github.com/AkimKachaliev/distributed-calculator/distributed-calculator/auth"
	"github.com/AkimKachaliev/distributed-calculator/distributed-calculator/config"
	"github.com/AkimKachaliev/distributed-calculator/distributed-calculator/models"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)

func main() {
	// Инициализация базы данных
	_, err := config.InitDB()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Создание маршрутизатора
	r := mux.NewRouter()

	// Регистрация маршрутов для аутентификации
	authRouter := auth.NewRouter()
	r.PathPrefix("/auth").Handler(http.StripPrefix("/auth", authRouter))

	// Регистрация маршрута для вычислений
	r.HandleFunc("/calculator", CalculateHandler).Methods("POST")

	// Запуск сервера
	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}

func CalculateHandler(w http.ResponseWriter, r *http.Request) {
	// Получение данных из запроса
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Failed to parse form data", http.StatusBadRequest)
		return
	}

	// Извлечение параметров для вычислений
	expression := r.FormValue("expression")
	userID, err := strconv.ParseInt(r.FormValue("user_id"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid user_id", http.StatusBadRequest)
		return
	}

	// Выполнение вычислений
	result, err := performCalculation(expression)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Сохранение результата в базе данных
	expr := models.Expression{
		UserID:     userID,
		Expression: expression,
		Result:     result,
	}
	err = expr.Create()
	if err != nil {
		http.Error(w, "Failed to save expression", http.StatusInternalServerError)
		return
	}

	// Отправка результата
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%s", result)
}

func performCalculation(expression string) (string, error) {
	// Здесь вы должны реализовать логику вычислений
	// Например, можно использовать пакет "math/big" для работы с большими числами

	// Пример простой реализации:
	expr, err := govaluate.NewEvaluableExpression(expression)
	if err != nil {
		return "", err
	}
	result, err := expr.Evaluate(nil)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%.2f", result), nil
}
