package database

import (
	"github.com/jmoiron/sqlx"
	"os"
)

func NewPostgresDB() (*sqlx.DB, error) {
	dsn := os.Getenv("POSTGRES_CONNECTION")
	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		return nil, err
	}
	return db, nil
}
