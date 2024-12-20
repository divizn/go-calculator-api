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
		Db_OLD   map[int]*models.Calculation
		Db       *pgxpool.Pool // maybe instead make a wrapper over db so i can use the functions for creating/updating/getting and have access to db model too
		seq      int
		validate *validator.Validate
	}
)

func New(db_old map[int]*models.Calculation) *Handler {
	validate = models.RegisterValidations()

	db, err := db.InitDB()
	if err != nil {
		log.Fatal("Database was not successfully initialised, exiting...")
	}

	handler := &Handler{
		Db_OLD:   db_old,
		Db:       db,
		seq:      1,
		validate: validate,
	}

	return handler
}
