package models

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

// TODO Fix so it actually sends errors and terminates, rather than sending error and continuing

func Return500InternalServerError(ctx echo.Context, err error) error {
	fmt.Printf("internal Server Error: %v\n", err.Error())

	return ctx.JSON(http.StatusInternalServerError, map[string]string{
		"error": "500 Internal Server Error",
	})
}

func Return404NotFound(ctx echo.Context) error {
	return ctx.JSON(http.StatusNotFound, map[string]string{
		"error": "404 Not Found",
	})
}

func Return400BadRequest(ctx echo.Context) error {
	return ctx.JSON(http.StatusBadRequest, map[string]string{
		"error": "400 Bad Request",
	})
}
