package handler

import (
	"log"

	"github.com/divizn/echo-calculator/internal/db"
	"github.com/divizn/echo-calculator/internal/models"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgxpool"
)

var validate *validator.Validate

type (
	Handler struct {
		Db       *pgxpool.Pool
		seq      int
		validate *validator.Validate
	}
)

func New() *Handler {
	validate = models.RegisterValidations()

	db, err := db.InitDB()
	if err != nil {
		log.Fatal("Database was not successfully initialised, exiting...")
	}

	handler := &Handler{
		Db:       db,
		seq:      1,
		validate: validate,
	}

	return handler
}
