package main

import (
	"fmt"
	"math"
	"net/http"
	"strconv"
	"sync"

	"github.com/divizn/echo-calculator/internal/models"
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
	// validate := models.RegisterValidations()
	if err := validate.Struct(calc); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}

	switch calc.Operator {
	case "+":
		calc.Result = calc.Num1 + calc.Num2
	case "-":
		calc.Result = calc.Num1 - calc.Num2
	case "*":
		calc.Result = calc.Num1 * calc.Num2
	case "/":
		// 0 check is not necessary because 0 is float32s 0 value and "required" in validator library sees is as nothing there https://github.com/go-playground/validator/issues/290
		if calc.Num2 == 0 {
			return c.JSON(http.StatusBadRequest, models.CalcError{Message: "Division by zero is not allowed"})
		}
		calc.Result = calc.Num1 / calc.Num2
	case "^":
		calc.Result = float32(math.Pow(float64(calc.Num1), float64(calc.Num2)))
	case "%":
		if calc.Num2 == 0 {
			return c.JSON(http.StatusBadRequest, models.CalcError{Message: "Modulo by zero is not allowed"})
		}
		calc.Result = float32(int(calc.Num1) % int(calc.Num2))
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
	id, _ := strconv.Atoi(c.Param("id"))
	db[id].Num1 = calc.Num1
	db[id].Num1 = calc.Num1
	db[id].Result = calc.Num1 + calc.Num2
	return c.JSON(http.StatusOK, db[id])
}

func deleteCalculation(c echo.Context) error {
	lock.Lock()
	defer lock.Unlock()

	id, _ := strconv.Atoi(c.Param("id"))
	delete(db, id)
	return c.NoContent(http.StatusNoContent)
}
