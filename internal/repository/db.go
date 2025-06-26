package repository

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

func NewDB() (*sql.DB, error) {
connStr := "postgres://postgres:postgres@localhost:5432/url_shortener?sslmode=disable"

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("could not connect to DB: %w", err)
	}
	return db, db.Ping()
}
