package db

import (
	"context"
	"fmt"
	"os"

	"github.com/divizn/echo-calculator/internal/utils"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type Database struct {
	Pool  *pgxpool.Pool
	Cache *redis.Client
	Ctx   *context.Context
}

// initialise the db and return a pointer to the connection pool to be utilised in handler.Handler
func InitDB(config *utils.IConfig) (*Database, error) {
	dbpool, err := pgxpool.New(context.Background(), config.DB_URL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create connection pool: %v\n", err)
		os.Exit(1)
	}

	var testConnection string
	err = dbpool.QueryRow(context.Background(), "select 'Test DB query successful'").Scan(&testConnection)
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		return nil, err
	}
	fmt.Println(testConnection)

	ctx := context.TODO()

	rdb := redis.NewClient(&redis.Options{
		Addr:     config.REDIS_ADDR,
		Username: "",
		Password: "",
	})

	return &Database{Pool: dbpool, Cache: rdb, Ctx: &ctx}, nil
}

func (db *Database) Close() {
	db.Pool.Close()

	if err := db.Cache.Close(); err != nil {
		panic(err)
	}
}
