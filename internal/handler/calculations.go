package handler

import (
	"net/http"
	"strconv"

	"github.com/divizn/echo-calculator/internal/models"
	"github.com/labstack/echo/v4"
)

// GetAllCalculations example
//
//	@Summary		Shows all calculations in the database
//	@Description	Get all calculations
//	@ID				get-all-calculations
//	@Accept			json
//	@Produce		json
//	@Success		200	{array}	models.Calculation	"ok"
//	@Router			/calculations [get]
func (h *Handler) GetAllCalculations(c echo.Context) error {
	calculations, err := h.Service.GetAllCalculations(h.Db)
	if err != nil {
		return models.Return500InternalServerError(c, err)
	}

	return c.JSON(http.StatusOK, calculations)
}

// GetCalculation example
//
//	@Summary		Gets a calculation from the given ID
//	@Description	Get calculation by ID
//	@ID				get-calculation
//	@Param			id	path	int	true	"Some ID"
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	models.Calculation	"ok"
//	@Failure		404	{object}	models.CalcError	"not found"
//
//	@Router			/calculations/{id} [get]
func (h *Handler) GetCalculation(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return models.Return400BadRequest(c)
	}

	calc, err := h.Service.GetCalculationByID(h.Db, id, c)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}
	return c.JSON(http.StatusOK, calc)
}

// CreateCalculation example
//
//	@Summary		Creates a calculation
//	@Description	Createc calculation
//	@ID				create-calculation
//	@Param			request	body	models.CreateCalculationRequest	true	"request body"
//	@Accept			json
//	@Produce		json
//	@Success		201	{object}	models.Calculation	"created"
//	@Failure		400	{object}	models.CalcError	"bad request"
//
//	@Router			/calculations [post]
func (h *Handler) CreateCalculation(c echo.Context) error {
	req := new(models.CreateCalculationRequest)
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

	calc, err := h.Service.CreateCalculation(h.Db, req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to create calculation	",
		})
	}

	// Return the created calculation
	return c.JSON(http.StatusCreated, calc)
}

// UpdateCalculation example
//
//	@Summary		Updates a calculation from a given ID
//	@Description	Update calculation from given ID
//	@ID				update-calculation
//	@Param			id		path	int								true	"Some ID"
//	@Param			request	body	models.UpdateCalculationRequest	true	"request body"
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	models.Calculation	"ok"
//	@Failure		400	{object}	models.CalcError	"bad request"
//
//	@Router			/calculations/{id} [put]
func (h *Handler) UpdateCalculation(c echo.Context) error {
	calc := new(models.UpdateCalculationRequest)
	if err := c.Bind(calc); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to bind request",
		})
	}
	if err := h.validate.Struct(calc); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid ID",
		})
	}

	updatedCalc, err := h.Service.UpdateCalculation(h.Db, id, calc, c)
	if err != nil {
		if err.Error() == "calculation not found" {
			return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, updatedCalc)
}

// DeleteCalculation example
//
//	@Summary		Deletes a calculation from a given ID
//	@Description	Update calculation from given ID
//	@ID				delete-calculation
//	@Param			id	path	int	true	"Some ID"
//	@Accept			json
//	@Produce		json
//	@Success		204	"no content"
//
//	@Router			/calculations/{id} [delete]
func (h *Handler) DeleteCalculation(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return models.Return400BadRequest(c)
	}
	if err := h.Service.DeleteCalculation(h.Db, id, c); err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": err.Error(),
		})
	}
	return c.NoContent(http.StatusNoContent)
}
