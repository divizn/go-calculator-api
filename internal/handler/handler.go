package handler

import (
	"github.com/divizn/echo-calculator/internal/db"
	"github.com/divizn/echo-calculator/internal/models"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgxpool"
)

var validate *validator.Validate

type (
	Handler struct {
		Db       *pgxpool.Pool
		validate *validator.Validate
	}
)

func NewHandler(db *db.Database) *Handler {
	validate = models.RegisterValidations()

	handler := &Handler{
		Db:       db.Pool,
		validate: validate,
	}

	return handler
}
