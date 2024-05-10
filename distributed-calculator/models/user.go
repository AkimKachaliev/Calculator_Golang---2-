package models

import "github.com/AkimKachaliev/distributed-calculator/distributed-calculator/config"

type User struct {
	ID       int64  `json:"id"`
	Login    string `json:"login"`
	Password string `json:"password"`
}

func (u *User) Create() error {
	db, err := config.InitDB()
	if err != nil {
		return err
	}
	defer db.Close()

	statement, err := db.Prepare("INSERT INTO users (login, password) VALUES (?, ?)")
	if err != nil {
		return err
	}
	defer statement.Close()

	result, err := statement.Exec(u.Login, u.Password)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	u.ID = id
	return nil
}

func GetUser(login, password string) (*User, error) {
	db, err := config.InitDB()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	statement, err := db.Prepare("SELECT id, login, password FROM users WHERE login = ? AND password = ?")
	if err != nil {
		return nil, err
	}
	defer statement.Close()

	row := statement.QueryRow(login, password)

	var user User
	err = row.Scan(&user.ID, &user.Login, &user.Password)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
