// This file contains the repository implementation layer.
package repository

import (
	"database/sql"

	_ "github.com/lib/pq"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(dsn string) *Repository {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		panic(err)
	}

	return &Repository{
		db: db,
	}
}
