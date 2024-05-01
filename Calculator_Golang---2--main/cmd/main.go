package main

import (
	"database/sql"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"sync"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

const (
	dbPath = "sql/my.db"
	port   = ":8081"
)

var (
	db         *sql.DB // Global variable for SQLite database
	numWorkers int
)

func main() {
	var err error
	// Initialize SQLite database
	db, err = initializeSQLiteDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	http.HandleFunc("/receive", receiveData)

	fmt.Println("Receiver server listening on port", port)
	http.ListenAndServe(port, nil)
}

func initializeSQLiteDB() (*sql.DB, error) {
	_, err := os.Stat(dbPath)
	if os.IsNotExist(err) {
		file, err := os.Create(dbPath)
		if err != nil {
			log.Fatalf("Error creating database file: %v\n", err)
			return nil, err
		}
		file.Close()
	}

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatalf("Error opening database: %v\n", err)
		return nil, err
	}

	// Create tables if they do not exist
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
		log.Fatalf("Error creating tables: %v\n", err)
		return nil, err
	}

	return db, nil
}

func receiveData(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}

	fmt.Printf("Received data: %s\n", string(body))
	start(string(body))
}

func start(expression string) {
	fmt.Println("Agent started")
	// Check expression validity
	re := regexp.MustCompile(`(\d+(\.\d+)?)|([+\-*\/])`)
	if !re.MatchString(expression) {
		fmt.Println("Invalid expression:", expression)
		plusResult(0, false)
	} else {
		fmt.Println("Calculating expression")
		// Solve expression
		result := startCount(expression)
		fmt.Println(result)
		// Write result to database
		plusResult(result, true)
	}
}

// Workers
func startCount(infixExpression string) float64 {
	postfixExpression := infixToPostfix(infixExpression)
	// Use workers for parallel processing of the expression
	if numWorkers <= 0 {
		numWorkers = 2
	}
	tasks := make(chan string)
	results := make(chan float64)
	var wg sync.WaitGroup

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go worker(&wg, tasks, results)
	}

	go func() {
		tasks <- postfixExpression
		close(tasks)
	}()

	go func() {
		wg.Wait()
		close(results)
	}()
	var result float64
	for result = range results {
		fmt.Println("Результат вычисления:", result)
	}
	return result
}

func worker(wg *sync.WaitGroup, tasks <-chan string, results chan<- float64) {
	defer wg.Done()
	for task := range tasks {
		result := evaluatePostfix(task)
		results <- result
	}
}

func precedence(op rune) int {
	switch op {
	case '+', '-':
		return 1
	case '*', '/':
		return 2
	}
	return 0
}

func infixToPostfix(infix string) string {
	var postfix []rune
	var stack []rune

	for _, char := range infix {
		if char >= '0' && char <= '9' {
			postfix = append(postfix, char)
		} else if char == '(' {
			stack = append(stack, char)
		} else if char == ')' {
			for len(stack) > 0 && stack[len(stack)-1] != '(' {
				postfix = append(postfix, stack[len(stack)-1])
				stack = stack[:len(stack)-1]
			}
			stack = stack[:len(stack)-1]
		} else {
			for len(stack) > 0 && precedence(char) <= precedence(stack[len(stack)-1]) {
				postfix = append(postfix, stack[len(stack)-1])
				stack = stack[:len(stack)-1]
			}
			stack = append(stack, char)
		}
	}

	for len(stack) > 0 {
		postfix = append(postfix, stack[len(stack)-1])
		stack = stack[:len(stack)-1]
	}

	return string(postfix)
}

func evaluatePostfix(postfix string) float64 {
	var stack []float64

	for _, char := range postfix {
		if char >= '0' && char <= '9' {
			num, _ := strconv.ParseFloat(string(char), 64)
			stack = append(stack, num)
		} else {
			b := stack[len(stack)-1]
			a := stack[len(stack)-2]
			stack = stack[:len(stack)-2]
			timeout := time.After(1 * time.Second) // Set a 1-second limit for operation execution
			var result float64
			switch char {
			case '+':
				select {
				case <-timeout:
					fmt.Println("Время выполнения операции + истекло")
					return -1
				default:
					result = a + b
				}

			case '-':
				select {
				case <-timeout:
					fmt.Println("Время выполнения операции - истекло")
					return -1
				default:
					result = a - b
				}
			case '*':
				select {
				case <-timeout:
					fmt.Println("Время выполнения операции * истекло")
					return -1
				default:
					result = a * b
				}
			case '/':
				select {
				case <-timeout:
					fmt.Println("Время выполнения операции / истекло")
					return -1
				default:
					result = a / b
				}
			}

			stack = append(stack, result)
		}
	}

	return stack[0]
}

func plusResult(res float64, invalid bool) {
	// Get the current ID value
	var id int
	err := db.QueryRow("SELECT id FROM id_counter").Scan(&id)
	if err != nil {
		log.Fatalf("Ошибка при чтении значения ID: %v\n", err)
		return
	}

	// Increase the ID value by 1
	_, err = db.Exec("UPDATE id_counter SET id = id + 1")
	if err != nil {
		log.Fatalf("Ошибка при обновлении значения ID: %v\n", err)
		return
	}

	// Check that the new ID value has indeed increased
	var newId int
	err = db.QueryRow("SELECT id FROM id_counter").Scan(&newId)
	if err != nil {
		log.Fatalf("Ошибка при чтении обновленного значения ID: %v\n", err)
		return
	}

	// Insert data into the expressions table with the updated ID
	_, err = db.Exec("INSERT INTO expressions (id, expression, responses, user) VALUES (?,?,?,?)",
		newId, "expression_placeholder", res, "user_placeholder")
	if err != nil {
		log.Fatalf("Ошибка при вставке данных в таблицу expressions: %v\n", err)
		return
	}
}
