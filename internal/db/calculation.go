package db

import (
	"fmt"

	"github.com/divizn/echo-calculator/internal/models"
	"github.com/jackc/pgx/v5"
)

func (db *Database) CreateCalculation(calc *models.Calculation) error {
	query := `
        INSERT INTO calculations (num1, num2, operator, result)
        VALUES ($1, $2, $3, $4)
        RETURNING id
    `
	row := db.Pool.QueryRow(*db.Ctx, query, calc.Num1, calc.Num2, calc.Operator, calc.Result)

	// Get the newly created ID
	if err := row.Scan(&calc.ID); err != nil {
		return fmt.Errorf("failed to scan calculation: %v", err)
	}
	return nil
}

func (db *Database) GetCalculationByID(id int) (*models.Calculation, error) {
	calc := &models.Calculation{}
	query := "SELECT id, num1, num2, operator, result FROM calculations WHERE id = $1"

	err := db.Pool.QueryRow(*db.Ctx, query, id).Scan(&calc.ID, &calc.Num1, &calc.Num2, &calc.Operator, &calc.Result)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("calculation not found")
		}
		return nil, fmt.Errorf("failed to fetch calculation: %v", err)
	}
	return calc, nil
}

func (db *Database) DeleteCalculation(id int) error {
	query := "DELETE FROM calculations WHERE id = $1"

	cmdTag, err := db.Pool.Exec(*db.Ctx, query, id)
	if err != nil {
		return err
	}

	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("no calculation found with id %d", id)
	}

	return nil
}

func (db *Database) GetAllCalculations() ([]*models.Calculation, error) {
	query := "SELECT id, num1, num2, operator, result FROM calculations"
	rows, err := db.Pool.Query(*db.Ctx, query)
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

func (db *Database) UpdateCalculation(id int, calc *models.Calculation) (*models.Calculation, error) {
	updateQuery := `
        UPDATE calculations
        SET num1 = $1, num2 = $2, operator = $3, result = $4
        WHERE id = $5
    `
	_, err := db.Pool.Exec(*db.Ctx, updateQuery, calc.Num1, calc.Num2, calc.Operator, calc.Result, id)
	if err != nil {
		return nil, fmt.Errorf("failed to update calculation: %v", err)
	}

	return calc, nil
}
