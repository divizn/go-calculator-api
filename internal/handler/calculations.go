package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/divizn/echo-calculator/internal/models"
	"github.com/divizn/echo-calculator/internal/services"
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
	return c.JSON(http.StatusOK, h.Db_OLD)
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
//	@Failure		404	{object}	models.CalcError			"not found"
//
//	@Router			/calculations/{id} [get]
func (h *Handler) GetCalculation(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	_, ok := h.Db_OLD[id]
	if !ok {
		return c.JSON(http.StatusNotFound, models.CalcError{Message: fmt.Sprintf("Could not find calculation for ID %v", id)})
	}
	return c.JSON(http.StatusOK, h.Db_OLD[id])
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
//	@Failure		400	{object}	models.CalcError			"bad request"
//
//	@Router			/calculations [post]
func (h *Handler) CreateCalculation(c echo.Context) error {
	calc := &models.Calculation{
		ID: h.seq,
	}

	req := new(models.CreateCalculationRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	if err := validate.Struct(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}

	calc.Num1 = req.Num1
	calc.Num2 = req.Num2
	calc.Operator = req.Operator

	err := services.CalculateResult(calc)
	if err != nil { // should be unreachable
		return c.JSON(http.StatusBadRequest, err)
	}

	h.Db_OLD[calc.ID] = calc
	h.seq++

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
//	@Failure		400	{object}	models.CalcError			"bad request"
//
//	@Router			/calculations/{id} [put]
func (h *Handler) UpdateCalculation(c echo.Context) error {
	calc := new(models.UpdateCalculationRequest)
	if err := c.Bind(calc); err != nil {
		return err
	}
	if err := validate.Struct(calc); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}
	id, _ := strconv.Atoi(c.Param("id"))
	err := services.UpdateCalculation(calc, h.Db_OLD, id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, h.Db_OLD[id])
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
	id, _ := strconv.Atoi(c.Param("id"))
	delete(h.Db_OLD, id)
	return c.NoContent(http.StatusNoContent)
}
