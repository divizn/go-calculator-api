package main

import (
	"log"

	"github.com/divizn/echo-calculator/internal/app"
	"github.com/divizn/echo-calculator/internal/db"
	"github.com/divizn/echo-calculator/internal/utils"

	_ "github.com/divizn/echo-calculator/docs"

	_ "github.com/joho/godotenv/autoload"
)

// @title			Calculator API
// @version		1.0
// @description	CRUD API that takes 2 numbers and an operand, and stores it with the result in a database.
// @contact.name	Repository
// @contact.url	http://github.com/divizn/go-calculator-api
func main() {
	config, err := utils.NewConfig()
	if err != nil {
		log.Fatal("Could not load environment variables (check if they are present)", err)
	}

	db, err := db.InitDB(config)
	if err != nil {
		log.Fatal("Database was not successfully initialised, exiting...")
	}

	app := app.NewApp(db)
	defer app.Db.Close() // close pool on server shut down todo: not graceful

	app.Echo.Logger.Fatal(app.Echo.Start(config.SERVER_ADDR))
}
