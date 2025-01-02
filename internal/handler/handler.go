package handler

import (
	"github.com/divizn/echo-calculator/internal/db"
	"github.com/divizn/echo-calculator/internal/models"
	"github.com/divizn/echo-calculator/internal/services"
	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

type (
	Handler struct {
		Db       *db.Database
		validate *validator.Validate
		Service  *services.Service
	}
)

func NewHandler(db *db.Database) *Handler {
	validate = models.RegisterValidations()

	handler := &Handler{
		Db:       db,
		validate: validate,
		Service:  services.NewService(),
	}

	return handler
}
