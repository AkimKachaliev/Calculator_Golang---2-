package models

import (
	"errors"
	"github.com/AkimKachaliev/distributed-calculator/config"
)

type Expression struct {
	ID         int64  `json:"id"`
	UserID     int64  `json:"user_id"`
	Expression string `json:"expression"`
	Result     string `json:"result"`
}

func (e *Expression) Create() error {
	db, err := config.InitDB()
	if err != nil {
		return err
	}
	defer db.Close()

	statement, err := db.Prepare("INSERT INTO expressions (user_id, expression) VALUES (?, ?)")
	if err != nil {
		return err
	}
	defer statement.Close()

	result, err := statement.Exec(e.UserID, e.Expression)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	e.ID = id
	return nil
}

func (e *Expression) GetResult() error {
	db, err := config.InitDB()
	if err != nil {
		return err
	}
	defer db.Close()

	statement, err := db.Prepare("UPDATE expressions SET result = ? WHERE id = ?")
	if err != nil {
		return err
	}
	defer statement.Close()

	result, err := statement.Exec(e.Result, e.ID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("expression not found")
	}

	return nil
}
