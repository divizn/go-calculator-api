package services

import (
	"errors"
	"math"

	"github.com/divizn/echo-calculator/internal/models"
)

func CalculateResult(calc *models.Calculation) error {

	switch calc.Operator {
	case "+":
		calc.Result = calc.Num1 + calc.Num2
	case "-":
		calc.Result = calc.Num1 - calc.Num2
	case "*":
		calc.Result = calc.Num1 * calc.Num2
	case "/":
		// 0 check is not necessary because 0 is float32s 0 value and "required" in validator library sees is as nothing there https://github.com/go-playground/validator/issues/290
		if calc.Num2 == 0 {
			return errors.New("division by zero is not allowed")
		}
		calc.Result = calc.Num1 / calc.Num2
	case "^":
		calc.Result = float32(math.Pow(float64(calc.Num1), float64(calc.Num2)))
	case "%":
		if calc.Num2 == 0 {
			return errors.New("division by zero is not allowed")
		}
		calc.Result = float32(int(calc.Num1) % int(calc.Num2))
	}

	return nil
}

func UpdateCalculation(calc *models.UpdateCalculationRequest, db map[int]*models.Calculation, id int) error {

	if (calc.Num1) != nil {
		db[id].Num1 = *calc.Num1
	}

	if (calc.Num2) != nil {
		db[id].Num2 = *calc.Num2
	}

	if (calc.Operator) != nil {
		db[id].Operator = *calc.Operator
	}
	CalculateResult(db[id])
	return nil
}
