package app

import (
	"github.com/divizn/echo-calculator/internal/db"
	"github.com/divizn/echo-calculator/internal/handler"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
)

type App struct {
	Echo *echo.Echo
	Db   *db.Database
}

func NewApp(db *db.Database) *App {
	e := echo.New()
	h := handler.NewHandler(db)

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/swagger/*", echoSwagger.WrapHandler)

	calculations := e.Group("/calculations")
	calculations.GET("", h.GetAllCalculations)
	calculations.POST("", h.CreateCalculation)
	calculations.GET("/:id", h.GetCalculation)
	calculations.PUT("/:id", h.UpdateCalculation)
	calculations.DELETE("/:id", h.DeleteCalculation)

	return &App{Echo: e, Db: db}
}
