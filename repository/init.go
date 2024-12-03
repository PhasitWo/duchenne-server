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


type QueryCriteria string

const (
	PATIENTID QueryCriteria = "WHERE patient_id = "
	DOCTORID  QueryCriteria = "WHERE doctor_id = "
	NONE      QueryCriteria = ""
)