package repository

import (
	"database/sql"
)

// repository.New can accept both sql.DB or sql.Tx
type DBTX interface {
	Exec(query string, args ...any) (sql.Result, error)
	Query(query string, args ...any) (*sql.Rows, error)
	QueryRow(query string, args ...any) *sql.Row
}

type Repo struct {
	db DBTX
}

// Constructor
func New(db DBTX) *Repo {
	return &Repo{db: db}
}

// CRITERIA
type QueryCriteria string

const (
	PATIENTID QueryCriteria = "WHERE patient_id = "
	DOCTORID  QueryCriteria = "WHERE doctor_id = "
	NONE      QueryCriteria = ""
)
