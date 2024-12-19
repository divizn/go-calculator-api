package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/divizn/echo-calculator/internal/models"
	"github.com/divizn/echo-calculator/internal/services"
	"github.com/labstack/echo"
)

func (h *Handler) GetAllCalculations(c echo.Context) error {
	return c.JSON(http.StatusOK, h.Db)
}

func (h *Handler) GetCalculation(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	_, ok := h.Db[id]
	if !ok {
		return c.JSON(http.StatusNotFound, models.CalcError{Message: fmt.Sprintf("Could not find calculation for ID %v", id)})
	}
	return c.JSON(http.StatusOK, h.Db[id])
}

func (h *Handler) CreateCalculation(c echo.Context) error {
	calc := &models.Calculation{
		ID: h.seq,
	}
	if err := c.Bind(calc); err != nil {
		return err
	}
	if err := validate.Struct(calc); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}

	err := services.CalculateResult(calc)
	if err != nil { // should be unreachable
		return c.JSON(http.StatusBadRequest, err)
	}

	h.Db[calc.ID] = calc
	h.seq++

	return c.JSON(http.StatusCreated, calc)
}

func (h *Handler) UpdateCalculation(c echo.Context) error {
	calc := new(models.Calculation)
	if err := c.Bind(calc); err != nil {
		return err
	}
	if err := validate.Struct(calc); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}
	id, _ := strconv.Atoi(c.Param("id"))
	err := services.UpdateCalculation(calc, h.Db, id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, h.Db[id])
}

func (h *Handler) DeleteCalculation(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	delete(h.Db, id)
	return c.NoContent(http.StatusNoContent)
}
