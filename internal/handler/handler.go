package handler

import (
	"github.com/divizn/echo-calculator/internal/models"
	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

type (
	Handler struct {
		Db       map[int]*models.Calculation
		seq      int
		validate *validator.Validate
	}
)

func New(db map[int]*models.Calculation) *Handler {
	validate = models.RegisterValidations()

	handler := &Handler{
		Db:       db,
		seq:      1,
		validate: validate,
	}

	return handler
}
