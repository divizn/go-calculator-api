package db

import "github.com/divizn/echo-calculator/internal/models"

// returns the id if valid, otherwise nil and error
func ValidateID(id string) (int, error) {
	return 0, nil
}

// checks if the given, valid ID is in the db
func IDInDB(id int, db *map[int]*models.Calculation) bool {
	return false
}

// func EnvVars -> loads env variables and types them
