package db

import (
	"context"
	"fmt"
	"os"

	"github.com/divizn/echo-calculator/internal/utils"
	"github.com/jackc/pgx/v5/pgxpool"
)

// initialise the db and return a pointer to the connection pool to be utilised in handler.Handler
func InitDB() (*pgxpool.Pool, error) {
	var env utils.IConfig
	env.New()
	dbpool, err := pgxpool.New(context.Background(), env.DB_URL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create connection pool: %v\n", err)
		os.Exit(1)
	}

	var testConnection string
	err = dbpool.QueryRow(context.Background(), "select 'Test DB query successful\n\n'").Scan(&testConnection)
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		return nil, err
	}

	fmt.Println(testConnection)

	return dbpool, nil
}
