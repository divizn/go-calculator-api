package models

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

type Calculation struct {
	Num1     float32 `json:"number1" validate:"required"`
	Num2     float32 `json:"number2" validate:"required"`
	Operator string  `json:"operator" validate:"required,operator"`
	Result   float32 `json:"result"`
	ID       int     `json:"id"`
}

type CalcError struct {
	Message string `json:"message"`
}

var validate *validator.Validate

func RegisterValidations() *validator.Validate {
	validate = validator.New()
	defer fmt.Println("Registered all validators")

	validate.RegisterValidation("operator", validateOperator)

	return validate
}

func validateOperator(fl validator.FieldLevel) bool {
	validOperators := "+-*/^%"
	operator := fl.Field().String()

	for _, op := range validOperators {
		if string(op) == operator {
			return true
		}
	}
	return false
}
