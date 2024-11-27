package repository

import (
	"database/sql"
)

type Repo struct {
	db *sql.DB
}

// Constructor
func New(db *sql.DB) *Repo {
	return &Repo{db: db}
}
