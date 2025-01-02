package services

import (
	"errors"
	"fmt"
	"math"

	"github.com/divizn/echo-calculator/internal/models"
	"github.com/labstack/echo/v4"
)

func (s *Service) calculateResult(calc *models.Calculation) error {

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

func (s *Service) UpdateCalculation(id int, req *models.UpdateCalculationRequest, ctx echo.Context) (*models.Calculation, error) {
	calc, err := s.GetCalculationByID(id)
	if err != nil {
		return nil, err
	}

	if req.Num1 != nil {
		calc.Num1 = *req.Num1
	}
	if req.Num2 != nil {
		calc.Num2 = *req.Num2
	}
	if req.Operator != nil {
		calc.Operator = *req.Operator
	}

	if err := s.calculateResult(calc); err != nil {
		return nil, fmt.Errorf("failed to calculate result: %v", err)
	}

	s.Db.UpdateCalculation(id, calc)
	return calc, nil
}

func (s *Service) GetAllCalculations() ([]*models.Calculation, error) {
	calculations, err := s.Db.GetAllCalculations()
	if err != nil {
		return nil, err
	}

	return calculations, nil
}

func (s *Service) GetCalculationByID(id int) (*models.Calculation, error) {
	if err := s.validateID(id); err != nil {
		return nil, err
	}

	calc, err := s.Db.GetCalculationByID(id)
	if err != nil {
		return nil, err
	}

	return calc, nil
}

func (s *Service) DeleteCalculation(id int) error {
	_, err := s.GetCalculationByID(id) // get calculation first since deleting is costly unlike select, so first use select to check if the id is valid to save db costs, also validates id
	if err != nil {
		return err
	}

	if err = s.Db.DeleteCalculation(id); err != nil {
		return err
	}

	return nil
}

func (s *Service) CreateCalculation(req *models.CreateCalculationRequest) (*models.Calculation, error) {
	calc := &models.Calculation{
		Num1:     req.Num1,
		Num2:     req.Num2,
		Operator: req.Operator,
	}
	// modifies calc directly
	if err := s.calculateResult(calc); err != nil {
		return nil, fmt.Errorf("failed to calculate result: %v", err)
	}

	if err := s.Db.CreateCalculation(calc); err != nil {
		return nil, err
	}

	return calc, nil
}
