package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/divizn/echo-calculator/internal/db"
	"github.com/divizn/echo-calculator/internal/models"
	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
)

// hashes using bcrypt
func (s *Service) GenerateHash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// compares plaintext password with hashed password
func (s *Service) CompareHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// generates JWT token for given user ID and username
func (s *Service) GenerateJWT(userID int, username string) (string, error) {
	secretKey := []byte(s.Config.JWT_SECRET)

	claims := jwt.MapClaims{
		"sub":      userID, // subject
		"username": username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(), // expires in 24 hours
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

// registers new user
func (s *Service) RegisterUser(db *db.Database, req *models.RegisterUserRequest) (*models.User, error) {
	passwordHash, err := s.GenerateHash(req.Password)
	if err != nil {
		return nil, err
	}

	// TODO check if user exists by selecting instead of erroring because duplicate key

	var user models.User
	query := `
	INSERT INTO users (username, user_role, password_hash, created_at)
	VALUES ($1, $2, $3, NOW())
	RETURNING id, username, user_role, password_hash, created_at
    `
	err = db.Pool.QueryRow(context.Background(), query, req.Username, req.UserRole, passwordHash).
		Scan(&user.ID, &user.Username, &user.UserRole, &user.PasswordHash, &user.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to register user: %v", err)
	}

	return &user, nil
}

// authenticate user and return JWT token if successful
func (s *Service) LoginUser(db *db.Database, req *models.LoginUserRequest) (string, error) {
	// Fetch the user from the database by username
	var user models.User
	query := `SELECT id, username, user_role, password_hash, created_at FROM users WHERE username = $1`
	err := db.Pool.QueryRow(context.Background(), query, req.Username).Scan(&user.ID, &user.Username, &user.UserRole, &user.PasswordHash, &user.CreatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return "", errors.New("invalid credentials")
		}
		return "", fmt.Errorf("failed to fetch user: %v", err)
	}

	if !s.CompareHash(req.Password, user.PasswordHash) {
		return "", errors.New("invalid credentials")
	}

	token, err := s.GenerateJWT(user.ID, user.Username)
	if err != nil {
		return "", fmt.Errorf("failed to generate JWT: %v", err)
	}

	return token, nil
}
