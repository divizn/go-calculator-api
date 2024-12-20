package main

import (
	"fmt"
	"log"

	"github.com/divizn/echo-calculator/internal/handler"
	"github.com/divizn/echo-calculator/internal/models"
	"github.com/divizn/echo-calculator/internal/utils"
	"github.com/joho/godotenv"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

var (
	db = map[int]*models.Calculation{}
)

/*
* CRUD API that takes 2 numbers and an operand, and stores it with the result in a database.
* Example:
* 4 + 4 would be stored as 4 + 4 = 8.
 */
func main() {
	godotenv.Load()

	var env utils.IConfig

	err := env.New()
	if err != nil {
		log.Fatal("Could not load environment variables (check if they are present)")
	}

	fmt.Println("db: ", env.PORT)

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	h := handler.New(db)

	e.GET("/calculations", h.GetAllCalculations)
	e.POST("/calculations", h.CreateCalculation)
	e.GET("/calculations/:id", h.GetCalculation)
	e.PUT("/calculations/:id", h.UpdateCalculation)
	e.DELETE("/calculations/:id", h.DeleteCalculation)

	e.Logger.Fatal(e.Start(":1323"))
}
