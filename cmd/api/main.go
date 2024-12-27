package main

import (
	"log"

	"github.com/divizn/echo-calculator/internal/handler"
	"github.com/divizn/echo-calculator/internal/utils"

	_ "github.com/divizn/echo-calculator/docs"

	"github.com/joho/godotenv"
	_ "github.com/joho/godotenv/autoload"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
)

// @title			Calculator API
// @version		1.0
// @description	CRUD API that takes 2 numbers and an operand, and stores it with the result in a database.
// @contact.name	Repository
// @contact.url	http://github.com/divizn/go-calculator-api
func main() {
	godotenv.Load()

	var env utils.IConfig
	err := env.New()
	if err != nil {
		log.Fatal("Could not load environment variables (check if they are present)")
	}

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Use(echojwt.JWT([]byte(env.JWT_SECRET)))

	h := handler.New()

	defer h.Db.Close() // close pool on server shut down todo: not graceful

	e.GET("/swagger/*", echoSwagger.WrapHandler)

	e.GET("/calculations", h.GetAllCalculations)
	e.POST("/calculations", h.CreateCalculation)
	e.GET("/calculations/:id", h.GetCalculation)
	e.PUT("/calculations/:id", h.UpdateCalculation)
	e.DELETE("/calculations/:id", h.DeleteCalculation)

	e.Logger.Fatal(e.Start(":1323"))
}
