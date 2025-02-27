package models

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

type Calculation struct {
	Num1     float32 `json:"number1" validate:"required" example:"1"`
	Num2     float32 `json:"number2" validate:"required" example:"1"`
	Operator string  `json:"operator" validate:"required,operator" example:"+"`
	Result   float32 `json:"result" example:"2"`
	ID       int     `json:"id" example:"1"`
}

type CreateCalculationRequest struct {
	Num1     float32 `json:"number1" validate:"required" example:"1"`
	Num2     float32 `json:"number2" validate:"required" example:"1"`
	Operator string  `json:"operator" validate:"required,operator" example:"+"`
}

type UpdateCalculationRequest struct {
	Num1     *float32 `json:"number1,omitempty" example:"1"`
	Num2     *float32 `json:"number2,omitempty" example:"1"`
	Operator *string  `json:"operator,omitempty" validate:"operator" example:"+"`
}

type CalcError struct {
	Message string `json:"message" example:"calculation error occured"`
}

var validate *validator.Validate

func RegisterValidations() *validator.Validate {
	validate = validator.New()

	validate.RegisterValidation("operator", validateOperator)

	fmt.Println("Registered all validators")
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
