package db

import (
	"errors"
	"fmt"

	"github.com/divizn/echo-calculator/internal/models"
	"github.com/jackc/pgx/v5"
)

// checks if the given user is in the db and returns id if it is present
func (db *Database) UserIDInDB(username *string) (int, error) {
	var id int

	query := `SELECT id FROM users WHERE username = $1`
	err := db.Pool.QueryRow(*db.Ctx, query, username).Scan(&id)
	if err != nil {
		if err == pgx.ErrNoRows {
			return -2, fmt.Errorf("no user found")
		}
		return -1, err
	}
	return id, nil
}

// fetch the user from the database by username
func (db *Database) GetUserFromUsername(username string) (*models.User, error) {
	var user models.User
	query := `SELECT id, username, user_role, password_hash, created_at FROM users WHERE username = $1`
	err := db.Pool.QueryRow(*db.Ctx, query, username).Scan(&user.ID, &user.Username, &user.UserRole, &user.PasswordHash, &user.CreatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, errors.New("invalid credentials")
		}
		return nil, fmt.Errorf("failed to fetch user: %v", err)
	}

	return &user, nil
}

func (db *Database) RegisterUser(req *models.RegisterUserRequest, passwordHash string) (*models.User, error) {
	var user models.User
	query := `
	INSERT INTO users (username, user_role, password_hash, created_at)
	VALUES ($1, $2, $3, NOW())
	RETURNING id, username, user_role, password_hash, created_at
    `
	err := db.Pool.QueryRow(*db.Ctx, query, req.Username, req.UserRole, passwordHash).
		Scan(&user.ID, &user.Username, &user.UserRole, &user.PasswordHash, &user.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to register user: %v", err)
	}

	return &user, nil
}
