package tests

import (
	"github.com/yourusername/distributed-calculator/config"
	"github.com/yourusername/distributed-calculator/models"
	"testing"
)

func TestUserRegistration(t *testing.T) {
	db, err := config.InitDB()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	user := models.User{
		Login:    "test",
		Password: "test",
	}

	err = user.Create()
	if err != nil {
		t.Error(err)
	}

	var foundUser models.User
	err = db.QueryRow("SELECT id, login, password FROM users WHERE login = ?", user.Login).Scan(&foundUser.ID, &foundUser.Login, &foundUser.Password)
	if err != nil {
		t.Error(err)
	}

	if foundUser.ID == 0 {
		t.Error("user not found")
	}

	if foundUser.Login != user.Login {
		t.Error("wrong login")
	}

	if foundUser.Password != user.Password {
		t.Error("wrong password")
	}
}

func TestExpressionCalculation(t *testing.T) {
	db, err := config.InitDB()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	expression := models.Expression{
		UserID:     1,
		Expression: "2+2",
	}

	err = expression.Calculate()
	if err != nil {
		t.Error(err)
	}

	var foundExpression models.Expression
	err = db.QueryRow("SELECT id, user_id, expression, result FROM expressions WHERE user_id = ?", expression.UserID).Scan(&foundExpression.ID, &foundExpression.UserID, &foundExpression.Expression, &foundExpression.Result)
	if err != nil {
		t.Error(err)
	}

	if foundExpression.ID == 0 {
		t.Error("expression not found")
	}

	if foundExpression.UserID != expression.UserID {
		t.Error("wrong user_id")
	}

	if foundExpression.Expression != expression.Expression {
		t.Error("wrong expression")
	}

	if foundExpression.Result != expression.Result {
		t.Error("wrong result")
	}
}
