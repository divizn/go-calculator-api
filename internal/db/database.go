package db

import (
	"context"
	"fmt"
	"os"

	"github.com/divizn/echo-calculator/internal/utils"
	"github.com/jackc/pgx/v5/pgxpool"
)

// TODO:
// database
// functions for creating, updating by id, retrieving by id, and getting all entries (maybe in chunks research pagination)

// initialise the db and return a pointer to the connection pool to be utilised in handler.Handler
func InitDB() (*pgxpool.Pool, error) {
	var env utils.IConfig
	env.New()
	dbpool, err := pgxpool.New(context.Background(), env.DB_URL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create connection pool: %v\n", err)
		os.Exit(1)
	}
	defer dbpool.Close()

	var greeting string
	err = dbpool.QueryRow(context.Background(), "select 'Hello, world!'").Scan(&greeting)
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		return nil, err
	}

	fmt.Println(greeting)

	return dbpool, nil
}
