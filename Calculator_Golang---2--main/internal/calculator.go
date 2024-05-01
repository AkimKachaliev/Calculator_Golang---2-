package agent

import (
	"errors"

	"github.com/Knetic/govaluate"
)

// CalculatorAgent представляет агента, способного выполнять математические вычисления.
type CalculatorAgent struct{}

// Calculate выполняет математическое вычисление заданной выражения.
func (a *CalculatorAgent) Calculate(expression string) (float64, error) {
	// Создание новой вычислимой выражения из переданной строки.
	expressionEvaluator, err := govaluate.NewEvaluableExpression(expression)
	if err != nil {
		// В случае ошибки при создании выражения, возвращаем ошибку.
		return 0, err
	}

	// Вычисление значения выражения.
	result, err := expressionEvaluator.Evaluate(nil)
	if err != nil {
		// В случае ошибки при вычислении, возвращаем ошибку.
		return 0, err
	}

	// Проверка, является ли результат числом.
	if resultFloat, ok := result.(float64); ok {
		// Если результат является числом, возвращаем его.
		return resultFloat, nil
	}

	// В случае, если результат не является числом, возвращаем ошибку.
	return 0, errors.New("результат вычисления не является числом")
}
