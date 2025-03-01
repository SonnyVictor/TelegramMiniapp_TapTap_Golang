package database

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // Import driver PostgreSQL
)

func NewPostgresDB() (*sqlx.DB, error) {
	connStr := "host=localhost port=5432 user=root dbname=simple_telegram sslmode=disable password=secret"
	db, err := sqlx.Connect("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	return db, nil
}
