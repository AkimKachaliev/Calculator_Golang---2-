package models

import "github.com/AkimKachaliev/distributed-calculator/distributed-calculator/config"

type Expression struct {
	ID         int64  `json:"id"`
	UserID     int64  `json:"user_id"`
	Expression string `json:"expression"`
	Result     string `json:"result"`
}

type CalculateRequest struct {
	UserID     int64  `json:"user_id"`
	Expression string `json:"expression"`
}

type CalculateResponse struct {
	Result string `json:"result"`
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

	_, err = statement.Exec(e.Result, e.ID)
	if err != nil {
		return err
	}

	return nil
}
