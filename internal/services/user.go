package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/divizn/echo-calculator/internal/db"
	"github.com/divizn/echo-calculator/internal/models"
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
	// check if user exists by selecting first to lower costs of error
	id, err := db.UserIDInDB(&req.Username)
	if err != nil {
		if id != -2 {
			return nil, err
		}
	}

	passwordHash, err := s.GenerateHash(req.Password)
	if err != nil {
		return nil, err
	}

	user, err := db.RegisterUser(req, passwordHash)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// authenticate user and return JWT token if successful
func (s *Service) LoginUser(db *db.Database, req *models.LoginUserRequest) (string, error) {
	user, err := db.GetUserFromUsername(req.Username)
	if err != nil {
		return "", err
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
