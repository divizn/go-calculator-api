package db

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/divizn/echo-calculator/internal/models"
	"github.com/jackc/pgx/v5"
)

var (
	REDIS_CACHE_TIMEOUT = time.Minute * 60
)

// checks if the given user is in the db and returns id if it is present
func (db *Database) UserIDInDB(username *string) (int, error) {
	var id int

	query := `SELECT id FROM users WHERE username = $1`
	err := db.Pool.QueryRow(*db.Ctx, query, username).Scan(id)
	if err != nil {
		if err == pgx.ErrNoRows {
			return -2, fmt.Errorf("no user found")
		}
		fmt.Println(err)
		return -1, err
	}

	return id, nil
}

// caches calculations for given user
func (db *Database) CacheUserCalc(user_id int, calc *[]models.Calculation) error {
	key := createUserCachePrefix(user_id)

	err := db.Cache.Set(*db.Ctx, key, calc, REDIS_CACHE_TIMEOUT).Err()
	if err != nil {
		return err
	}
	return nil
}

func createUserCachePrefix(user_id int) string {
	return fmt.Sprintf("user:%v:calculations", user_id)
}

// returns the id if valid, otherwise nil and error
func validateID(id *int) (*int, error) {
	if *id <= 0 {
		return nil, fmt.Errorf("id cannot be 0 or negative")
	}
	return id, nil
}

// fetch the user from the database by username
func (db *Database) GetUserFromUsername(username string) (*models.User, error) {
	var user models.User
	query := `SELECT id, username, user_role, password_hash, created_at FROM users WHERE username = $1`
	err := db.Pool.QueryRow(context.Background(), query, username).Scan(&user.ID, &user.Username, &user.UserRole, &user.PasswordHash, &user.CreatedAt)
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
	err := db.Pool.QueryRow(context.Background(), query, req.Username, req.UserRole, passwordHash).
		Scan(&user.ID, &user.Username, &user.UserRole, &user.PasswordHash, &user.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to register user: %v", err)
	}

	return &user, nil
}
