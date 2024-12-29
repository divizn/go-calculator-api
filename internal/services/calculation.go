package services

import (
	"context"
	"errors"
	"fmt"
	"math"

	"github.com/divizn/echo-calculator/internal/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
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

func UpdateCalculation(db *pgxpool.Pool, id int, calc *models.UpdateCalculationRequest) (*models.Calculation, error) {
	query := "SELECT num1, num2, operator FROM calculations WHERE id = $1"
	existingCalc := &models.Calculation{ID: id}
	err := db.QueryRow(context.Background(), query, id).Scan(&existingCalc.Num1, &existingCalc.Num2, &existingCalc.Operator)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("calculation not found")
		}
		return nil, fmt.Errorf("failed to fetch calculation: %v", err)
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

	if err := CalculateResult(existingCalc); err != nil {
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

func GetAllCalculations(db *pgxpool.Pool) ([]*models.Calculation, error) {
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

func GetCalculationByID(db *pgxpool.Pool, id int) (*models.Calculation, error) {
	query := "SELECT id, num1, num2, operator, result FROM calculations WHERE id = $1"

	// TODO: check if number postive int
	calc := &models.Calculation{}
	err := db.QueryRow(context.Background(), query, id).Scan(&calc.ID, &calc.Num1, &calc.Num2, &calc.Operator, &calc.Result)
	if err != nil {
		// TODO return error instead and check error there
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("calculation not found")
		}
		return nil, fmt.Errorf("failed to fetch calculation: %v", err)
	}

	return calc, nil
}

func DeleteCalculation(db *pgxpool.Pool, id int) error {
	// TODO: call getcalbyid to check whether this id exists, and then modify, since delete is costly just like updating so we need to check if it exists with a low cost select first
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
