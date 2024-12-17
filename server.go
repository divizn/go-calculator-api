package main

import (
	"fmt"
	"net/http"
	"strconv"
	"sync"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

type calculation struct {
	Num1     int    `json:"number1"`
	Num2     int    `json:"number2"`
	Operator string `json:"operator"`
	ID       int    `json:"id"`
}

type calcError struct {
	message string
}

var (
	seq  = 1
	lock = sync.Mutex{}
	db   = map[int]*calculation{}
)

/*
* CRUD API that takes 2 numbers and an operand, and stores it with the result in a database.
* Example:
* 4 + 4 would be stored as 4 + 4 = 8.
 */
func main() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

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
		return c.JSON(http.StatusNotFound, calcError{message: fmt.Sprintf("Could not find calculation for ID %v", id)})
	}
	return c.JSON(http.StatusOK, db[id])
}

func createCalculation(c echo.Context) error {
	lock.Lock()
	defer lock.Unlock()

	calc := &calculation{
		ID: seq,
	}
	if err := c.Bind(calc); err != nil {
		return err
	}

	db[calc.ID] = calc
	seq++

	return c.JSON(http.StatusCreated, calc)
}

func updateCalculation(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}

func deleteCalculation(c echo.Context) error {
	return c.String(http.StatusOK, " delete")
}
