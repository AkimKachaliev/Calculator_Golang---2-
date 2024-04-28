package tests

import (
	"testing"

	"github.com/AkimKachaliev/distributed-calculator/distributed-calculator/models"
)

func TestUserCreate(t *testing.T) {
	user := models.User{
		Login:    "test",
		Password: "test",
	}

	err := user.Create()
	if err != nil {
		t.Error(err)
	}
}

func TestExpressionCalculate(t *testing.T) {
	expression := models.Expression{
		UserID:     1,
		Expression: "2+2",
	}

	err := expression.Calculate()
	if err != nil {
		t.Error(err)
	}

	if expression.Result != "4" {
		t.Error("wrong result")
	}
}
