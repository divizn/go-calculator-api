package handler

import (
	"net/http"
	"strings"

	"github.com/divizn/echo-calculator/internal/models"
	"github.com/labstack/echo/v4"
)

// RegisterUser example
//
//	@Summary		Register a new user
//	@Description	Register a new user in the system
//	@ID				register-user
//	@Param			request	body	models.RegisterUserRequest	true	"request body"
//	@Accept			json
//	@Produce		json
//	@Success		201	{object}	models.User		"created"
//	@Failure		400	{object}	models.Error	"bad request"
//	@Failure		500	{object}	models.Error	"internal server error"
//	@Router			/users/register [post]
func (h *Handler) RegisterUser(c echo.Context) error {
	req := new(models.RegisterUserRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to bind request",
		})
	}

	if err := h.validate.Struct(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}

	user, err := h.Service.RegisterUser(h.Db, req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to register user",
		})
	}

	return c.JSON(http.StatusCreated, user)
}

// LoginUser example
//
//	@Summary		Login a user
//	@Description	Authenticate user and return a JWT token
//	@ID				login-user
//	@Param			request	body	models.LoginUserRequest	true	"request body"
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	models.LoginResponse	"logged in"
//	@Failure		400	{object}	models.Error			"bad request"
//	@Failure		401	{object}	models.Error			"unauthorized"
//	@Failure		500	{object}	models.Error			"internal server error"
//	@Router			/users/login [post]
func (h *Handler) LoginUser(c echo.Context) error {
	req := new(models.LoginUserRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to bind request",
		})
	}

	if err := h.validate.Struct(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}

	token, err := h.Service.LoginUser(h.Db, req)
	if err != nil {
		if strings.Contains(err.Error(), "invalid credentials") {
			return c.JSON(http.StatusUnauthorized, map[string]string{
				"error": "username or password incorrect, please try again",
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to login",
		})
	}

	// Return JWT
	return c.JSON(http.StatusOK, models.LoginResponse{
		Token: token,
	})
}
