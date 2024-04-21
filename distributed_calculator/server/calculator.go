package server

import (
	"context"
	"github.com/AkimKachaliev/Calculator_Golang---2-/server/grpc"
	"github.com/AkimKachaliev/Calculator_Golang---2-/server/http"
)

func (c *Calculator) Save() error {
    // Подключение к базе данных
    db, err := InitDB()
    if err != nil {
        return err
    }
    defer db.Close()

    // Подготовка SQL запроса
    stmt, err := db.Prepare("INSERT INTO calculations (user_id, expression, result) VALUES (?, ?, ?)")
    if err != nil {
        return err
    }

    // Выполнение SQL запроса
    result, err := stmt.Exec(c.UserID, c.Expression, c.Result)
    if err != nil {
        return err
    }

    // Получение ID новой записи
    id, _ := result.LastInsertId()
    c.ID = int(id)

    return nil
}

func CalculatorGRPC(ctx context.Context, req *grpc.CalculatorRequest) (*grpc.CalculatorResponse, error) {
    // Получение идентификатора пользователя из токена авторизации
    userID, err := GetUserIDFromToken(req.Token)
    if err != nil {
        return nil, err
    }

    // Выполнение вычислений
    result := Calculate(req.Expression)

    // Сохранение результата в базе данных
    err = SaveCalculationResult(result, userID)
    if err != nil {
        return nil, err
    }

    // Создание ответа gRPC
    response := &grpc.CalculatorResponse{
        Result: result,
    }

    return response, nil
}
