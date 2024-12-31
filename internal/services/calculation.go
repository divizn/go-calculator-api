package services

import (
	"context"
	"errors"
	"fmt"
	"math"

	"github.com/divizn/echo-calculator/internal/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
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

func (s *Service) UpdateCalculation(db *pgxpool.Pool, id int, calc *models.UpdateCalculationRequest, ctx echo.Context) (*models.Calculation, error) {
	existingCalc, err := s.GetCalculationByID(db, id, ctx)
	if err != nil {
		return nil, err
	}

	if calc.Num1 != nil {
		existingCalc.Num1 = *calc.Num1
	}
	if calc.Num2 != nil {
		existingCalc.Num2 = *calc.Num2
	}
	if calc.Operator != nil {
		existingCalc.Operator = *calc.Operator
	}

	if err := s.calculateResult(existingCalc); err != nil {
		return nil, fmt.Errorf("failed to calculate result: %v", err)
	}

	// Update query
	updateQuery := `
        UPDATE calculations
        SET num1 = $1, num2 = $2, operator = $3, result = $4
        WHERE id = $5
    `
	_, err = db.Exec(context.Background(), updateQuery, existingCalc.Num1, existingCalc.Num2, existingCalc.Operator, existingCalc.Result, id)
	if err != nil {
		return nil, fmt.Errorf("failed to update calculation: %v", err)
	}

	return existingCalc, nil
}

func (s *Service) GetAllCalculations(db *pgxpool.Pool) ([]*models.Calculation, error) {
	query := "SELECT id, num1, num2, operator, result FROM calculations"
	rows, err := db.Query(context.Background(), query)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch calculations: %v", err)
	}
	defer rows.Close()

	var calculations []*models.Calculation
	for rows.Next() {
		calc := &models.Calculation{}
		err := rows.Scan(&calc.ID, &calc.Num1, &calc.Num2, &calc.Operator, &calc.Result)
		if err != nil {
			return nil, fmt.Errorf("failed to scan calculation: %v", err)
		}
		calculations = append(calculations, calc)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("error iterating rows: %v", rows.Err())
	}

	if calculations == nil {
		return nil, fmt.Errorf("no calculations in db")
	}

	return calculations, nil
}

func (s *Service) GetCalculationByID(db *pgxpool.Pool, id int, ctx echo.Context) (*models.Calculation, error) {
	if id <= 0 {
		return nil, fmt.Errorf("id cannot be 0 or less")
	}

	query := "SELECT id, num1, num2, operator, result FROM calculations WHERE id = $1"
	calc := &models.Calculation{}
	err := db.QueryRow(context.Background(), query, id).Scan(&calc.ID, &calc.Num1, &calc.Num2, &calc.Operator, &calc.Result)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("calculation not found")
		}
		return nil, fmt.Errorf("failed to fetch calculation: %v", err)
	}

	return calc, nil
}

func (s *Service) DeleteCalculation(db *pgxpool.Pool, id int, ctx echo.Context) error {
	_, err := s.GetCalculationByID(db, id, ctx) // get calculation first since deleting is costly unlike select, so first use select to check if the id is valid to save db costs
	if err != nil {
		return err
	}
	query := "DELETE FROM calculations WHERE id = $1"

	cmdTag, err := db.Exec(context.Background(), query, id)
	if err != nil {
		return err
	}

	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("no calculation found with id %d", id)
	}

	return nil
}

func (s *Service) CreateCalculation(db *pgxpool.Pool, req *models.CreateCalculationRequest) (*models.Calculation, error) {
	calc := &models.Calculation{
		Num1:     req.Num1,
		Num2:     req.Num2,
		Operator: req.Operator,
	}
	if err := s.calculateResult(calc); err != nil {
		return nil, fmt.Errorf("failed to calculate result: %v", err)
	}

	// Create query
	query := `
        INSERT INTO calculations (num1, num2, operator, result)
        VALUES ($1, $2, $3, $4)
        RETURNING id
    `
	row := db.QueryRow(context.Background(), query, calc.Num1, calc.Num2, calc.Operator, calc.Result)

	// Get the newly created ID
	if err := row.Scan(&calc.ID); err != nil {
		return nil, fmt.Errorf("failed to scan calculation: %v", err)
	}
	return calc, nil
}
