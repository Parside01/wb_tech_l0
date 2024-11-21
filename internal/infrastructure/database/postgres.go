package database

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	sqlxmock "github.com/zhashkevych/go-sqlxmock"
	"os"
)

func NewPostgresDB() (*sqlx.DB, error) {
	dsn := os.Getenv("POSTGRES_DSN")
	if dsn == "" {
		dsn = "host=localhost user=postgres password=postgres dbname=postgres port=5432 sslmode=disable"
	}
	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func NewPostgresDBMock() (*sqlx.DB, error) {
	db, mock, err := sqlxmock.New()
	mock
}
