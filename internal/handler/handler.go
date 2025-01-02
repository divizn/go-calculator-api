package handler

import (
	"github.com/divizn/echo-calculator/internal/models"
	"github.com/divizn/echo-calculator/internal/services"
	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

type (
	Handler struct {
		validate *validator.Validate
		Service  *services.Service
	}
)

func NewHandler() *Handler {
	validate = models.RegisterValidations()

	handler := &Handler{
		validate: validate,
		Service:  services.NewService(),
	}

	return handler
}
