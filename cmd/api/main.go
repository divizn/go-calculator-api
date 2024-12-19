package main

import (
	"fmt"
	"net/http"
	"strconv"
	"sync"

	"github.com/divizn/echo-calculator/internal/models"
	"github.com/divizn/echo-calculator/internal/services"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

var (
	seq  = 1
	lock = sync.Mutex{}
	db   = map[int]*models.Calculation{}
)

var validate *validator.Validate

/*
* CRUD API that takes 2 numbers and an operand, and stores it with the result in a database.
* Example:
* 4 + 4 would be stored as 4 + 4 = 8.
 */
func main() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	validate = models.RegisterValidations()

	e.GET("/calculations", getAllCalculations)
	e.POST("/calculations", createCalculation)
	e.GET("/calculations/:id", getCalculation)
	e.PUT("/calculations/:id", updateCalculation)
	e.DELETE("/calculations/:id", deleteCalculation)

	e.Logger.Fatal(e.Start(":1323"))
}

func getAllCalculations(c echo.Context) error {
	lock.Lock()
	defer lock.Unlock()
	return c.JSON(http.StatusOK, db)
}

func getCalculation(c echo.Context) error {
	lock.Lock()
	defer lock.Unlock()
	id, _ := strconv.Atoi(c.Param("id"))

	_, ok := db[id]
	if !ok {
		return c.JSON(http.StatusNotFound, models.CalcError{Message: fmt.Sprintf("Could not find calculation for ID %v", id)})
	}
	return c.JSON(http.StatusOK, db[id])
}

func createCalculation(c echo.Context) error {
	lock.Lock()
	defer lock.Unlock()
	calc := &models.Calculation{
		ID: seq,
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

	db[calc.ID] = calc
	seq++

	return c.JSON(http.StatusCreated, calc)
}

func updateCalculation(c echo.Context) error {
	lock.Lock()
	defer lock.Unlock()
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
	err := services.UpdateCalculation(calc, db, id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, db[id])
}

func deleteCalculation(c echo.Context) error {
	lock.Lock()
	defer lock.Unlock()

	id, _ := strconv.Atoi(c.Param("id"))
	delete(db, id)
	return c.NoContent(http.StatusNoContent)
}
