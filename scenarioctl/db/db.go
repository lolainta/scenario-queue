package db

import (
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

func Connect() *pgxpool.Pool {

	url := os.Getenv("DB_URL")
	if url == "" {
		log.Fatal("DB_URL not set")
	}

	pool, err := pgxpool.New(context.Background(), url)
	if err != nil {
		log.Fatal(err)
	}

	return pool
}
