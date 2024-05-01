package main

import (
	"bytes"
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/golang-jwt/jwt"
	_ "github.com/mattn/go-sqlite3"
)

var (
	newData     string
	tokenString string
	currentUser string
)

const (
	hmacSampleSecret = "super_secret_signature"
)

var (
	id                 int
	numWorkers         int
	additionTime       time.Duration = 1
	subtractionTime    time.Duration = 1
	multiplicationTime time.Duration = 1
	divisionTime       time.Duration = 1
	mu                 sync.Mutex
)

// CheckID проверяет и возвращает текущий ID из базы данных.
func CheckID(db *sql.DB) int {
	var id int
	err := db.QueryRow("SELECT id FROM id_counter").Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			// Если таблица пуста, инициализируем ID
			_, err = db.Exec("INSERT INTO id_counter (id) VALUES (1)")
			if err != nil {
				log.Fatalf("Ошибка при инициализации ID: %v\n", err)
			}
			return 1
		} else {
			log.Fatalf("Ошибка при чтении значения ID: %v\n", err)
		}
	}
	return id
}

func main() {
	db, err := InitializeSQLiteDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	http.HandleFunc("/s", func(w http.ResponseWriter, r *http.Request) {
		FormHandler(w, r, db)
	})

	http.HandleFunc("/n", func(w http.ResponseWriter, r *http.Request) {
		FormHandlerAuth(w, r, db)
	})

	http.HandleFunc("/e", func(w http.ResponseWriter, r *http.Request) {
		FormHandlerReg(w, r, db)
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		listPage(w, r)
	})

	http.HandleFunc("/add_page", func(w http.ResponseWriter, r *http.Request) {
		AddPage(w, r, db, currentUser)
	})

	http.HandleFunc("/auth_page", func(w http.ResponseWriter, r *http.Request) {
		AuthPage(w, r, db)
	})

	http.HandleFunc("/reg_page", func(w http.ResponseWriter, r *http.Request) {
		regPage(w, r)
	})

	http.ListenAndServe(":8080", nil)
}

// InitializeSQLiteDB инициализирует базу данных SQLite.
func InitializeSQLiteDB() (*sql.DB, error) {
	dbPath := "sql/my.db"
	_, err := os.Stat(dbPath)
	if os.IsNotExist(err) {
		file, err := os.Create(dbPath)
		if err != nil {
			log.Fatalf("Ошибка создания файла базы данных: %v\n", err)
			return nil, err
		}
		file.Close()
	}

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatalf("Ошибка открытия базы данных: %v\n", err)
		return nil, err
	}

	// Создаем таблицы, если они еще не существуют
	createTableQuery := `
    CREATE TABLE IF NOT EXISTS id_counter (
        id INTEGER PRIMARY KEY
    );
    CREATE TABLE IF NOT EXISTS expressions (
        id INTEGER PRIMARY KEY,
        expression TEXT,
        responses TEXT,
        user TEXT
    );
    `
	_, err = db.Exec(createTableQuery)
	if err != nil {
		log.Fatalf("Ошибка создания таблиц: %v\n", err)
		return nil, err
	}

	return db, nil
}

// listPage отображает главную страницу.
func listPage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("html/main.html", "html/head.html")
	if err != nil {
		log.Fatalf("Ошибка при инициализации: %v\n", err)
	}
	tmpl.ExecuteTemplate(w, "list_page", nil)
}

// regPage отображает страницу регистрации.
func regPage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("html/registration_jwt.html", "html/head.html")
	if err != nil {
		log.Fatalf("Ошибка при инициализации: %v\n", err)
	}
	tmpl.ExecuteTemplate(w, "reg_page", nil)
}

// AddPage отображает страницу добавления нового выражения.
func AddPage(w http.ResponseWriter, r *http.Request, db *sql.DB, currentUser string) {
	tmpl, err := template.ParseFiles("html/calculate.html", "html/head.html")
	if err != nil {
		log.Fatalf("Ошибка при инициализации: %v\n", err)
	}

	rows, err := db.Query("SELECT expression FROM expressions WHERE user =?", currentUser)
	if err != nil {
		log.Fatalf("Ошибка при выполнении запроса: %v\n", err)
	}
	defer rows.Close()

	var expressions []string
	for rows.Next() {
		var expression string
		err := rows.Scan(&expression)
		if err != nil {
			log.Fatalf("Ошибка при сканировании строки: %v\n", err)
		}
		expressions = append(expressions, expression)
	}

	data := struct {
		Expressions []string
	}{
		Expressions: expressions,
	}

	tmpl.ExecuteTemplate(w, "add_page", data)
}

// AuthPage отображает страницу аутентификации.
func AuthPage(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	tmpl, err := template.ParseFiles("html/auth.html", "html/head.html")
	if err != nil {
		log.Fatalf("Ошибка при инициализации: %v\n", err)
	}

	rows, err := db.Query("SELECT expression FROM expressions")
	if err != nil {
		log.Fatalf("Ошибка при выполнении запроса: %v\n", err)
	}
	defer rows.Close()

	var expressions []string
	for rows.Next() {
		var expression string
		err := rows.Scan(&expression)
		if err != nil {
			log.Fatalf("Ошибка при сканировании строки: %v\n", err)
		}
		expressions = append(expressions, expression)
	}

	data := struct {
		Expressions []string
	}{
		Expressions: expressions,
	}

	tmpl.ExecuteTemplate(w, "auth_page", data)
}

// FormHandler обрабатывает форму для добавления выражения.
func FormHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	data := r.FormValue("data")
	fmt.Printf("Данные: %s\n", data)
	WriteToDatabase(currentUser, data, db)

	http.Redirect(w, r, "/add_page", http.StatusMovedPermanently)
}

// FormHandlerReg обрабатывает форму регистрации.
func FormHandlerReg(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	user := r.FormValue("user")
	fmt.Printf("Получено: %s\n", user)
	now := time.Now()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"name": user,
		"nbf":  now.Unix(),
		"exp":  now.Add(10 * time.Minute).Unix(),
		"iat":  now.Unix(),
	})
	TokenString, _ := token.SignedString([]byte(hmacSampleSecret))
	fmt.Println("Token string:", TokenString)
	WriteToDatabase(TokenString, "", db)

	tmpl, err := template.ParseFiles("html/registration_jwt.html", "html/head.html")
	if err != nil {
		log.Fatalf("Ошибка при инициализации: %v\n", err)
	}

	data := struct {
		Variable string
	}{
		Variable: TokenString,
	}

	tmpl.ExecuteTemplate(w, "reg_page", data)
	http.Redirect(w, r, "/reg_page", http.StatusMovedPermanently)
}

// FormHandlerAuth обрабатывает форму аутентификации.
func FormHandlerAuth(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	ifErr := ""
	user := ""
	tokenString := r.FormValue("user_token")

	fmt.Printf("Получен токен: %s\n", tokenString)
	found := FindTokenInDatabase(tokenString, db)

	if !found {
		ifErr = "Неверный токен"
	} else {
		tokenFromString, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				fmt.Printf("неожиданный метод подписи: %v", token.Header["alg"])
			}
			return []byte(hmacSampleSecret), nil
		})

		if err != nil {
			log.Printf("Ошибка: %v\n", err)
			ifErr = "Неверный токен"
		} else {
			if claims, ok := tokenFromString.Claims.(jwt.MapClaims); ok {
				fmt.Println("Имя пользователя: ", claims["name"])
				user = "Добрый день, " + claims["name"].(string)
				currentUser = claims["name"].(string)
			} else {
				log.Printf("Ошибка: %v\n", err)
				ifErr = "Неверный токен"
			}
		}
	}

	if ifErr != "" {
		user = ifErr
	}

	tmpl, err := template.ParseFiles("html/auth.html", "html/head.html")
	if err != nil {
		log.Printf("Ошибка при инициализации: %v\n", err)
	}

	data := struct {
		Var string
	}{
		Var: user,
	}

	tmpl.ExecuteTemplate(w, "auth_page", data)
	http.Redirect(w, r, "/auth_page", http.StatusMovedPermanently)
}

// FindTokenInDatabase проверяет наличие токена в базе данных.
func FindTokenInDatabase(tokenCheck string, db *sql.DB) bool {
	id := CheckID(db)
	for i := 1; i <= id; i++ {
		var user string
		err := db.QueryRow("SELECT user FROM expressions WHERE id =? AND user =?", i, tokenCheck).Scan(&user)
		if err == nil {
			return true
		}
	}
	return false
}

// WriteToDatabase записывает данные в базу данных.
func WriteToDatabase(user string, data string, db *sql.DB) {
	id := CheckID(db)
	id++

	_, err := db.Exec("UPDATE id_counter SET id =?", id)
	if err != nil {
		log.Fatalf("Ошибка при обновлении значения ID: %v\n", err)
	}

	_, err = db.Exec("INSERT INTO expressions (id, expression, responses, user) VALUES (?,?, 'new',?)", id, data, user)
	if err != nil {
		log.Fatalf("Ошибка при записи данных в базу данных: %v\n", err)
	}

	log.Println("Данные успешно записаны в базу данных!")
	if data != "" {
		// Преобразование строки в байты для отправки
		dataBytes := []byte(data)

		// Отправка данных на сервер
		response, err := http.Post("http://localhost:8081/receive", "text/plain", bytes.NewBuffer(dataBytes))
		if err != nil {
			// Вывод ошибки при отправке данных
			fmt.Printf("Ошибка при отправке данных: %v\n", err)
			return
		}

		// Закрытие тела ответа после использования
		defer response.Body.Close()

		// Вывод сообщения об успешной отправке данных
		fmt.Println("Данные успешно отправлены")
	}
}
