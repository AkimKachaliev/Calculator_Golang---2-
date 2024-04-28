package config

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

func InitDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "./distributed-calculator.db")
	if err != nil {
		return nil, err
	}

	statement, err := db.Prepare(`CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		login TEXT UNIQUE NOT NULL,
		password TEXT NOT NULL
	)`)
	if err != nil {
		return nil, err
	}
	statement.Exec()

	statement, err = db.Prepare(`CREATE TABLE IF NOT EXISTS expressions (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL,
		expression TEXT NOT NULL,
		result TEXT,
		FOREIGN KEY(user_id) REFERENCES users(id)
	)`)
	if err != nil {
		return nil, err
	}
	statement.Exec()

	return db, nil
}
