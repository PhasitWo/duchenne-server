package repository

import (
	"database/sql"
	"errors"
	"fmt"
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

// ERROR
var ErrDuplicateEntry = errors.New("duplicate entry")


// CRITERIA
type Criteria struct {
	QueryCriteria ColumnCriteria
	Value any
}

func (c *Criteria) ToString() string {
	return fmt.Sprintf(" %s %v ", c.QueryCriteria, c.Value)
}

type ColumnCriteria string

const (
	PATIENTID ColumnCriteria = "patient_id = "
	DOCTORID  ColumnCriteria = "doctor_id = "
	NONE      ColumnCriteria = ""
)

func attachCriteria(queryString string, criteria ...Criteria) string {
	if len(criteria) == 0 {
		return  queryString
	}
	for index, c := range criteria {
		if index == 0 {
			queryString += " WHERE" + c.ToString()
			continue
		}
		queryString += " AND" + c.ToString()
	}
	return queryString
}