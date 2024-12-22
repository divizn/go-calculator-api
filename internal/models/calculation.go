package models

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

// TODO: deprecate this struct and use the ones below
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
	Message string `json:"message" example:"error message goes here"`
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
