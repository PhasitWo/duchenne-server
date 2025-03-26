package repository

import (
	// "database/sql"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

// repository.New can accept both sql.DB or sql.Tx
// type DBTX interface {
// 	Exec(query string, args ...any) (sql.Result, error)
// 	Query(query string, args ...any) (*sql.Rows, error)
// 	QueryRow(query string, args ...any) *sql.Row
// }

type Repo struct {
	db *gorm.DB
}

// Constructor
func New(db *gorm.DB) *Repo {
	return &Repo{db: db}
}

// ERROR
var ErrDuplicateEntry = errors.New("duplicate entry")
var ErrForeignKeyFail = errors.New("foreign key error")

// CRITERIA
type Criteria struct {
	QueryCriteria ColumnCriteria
	Value         any
}

func (c *Criteria) ToString() string {
	if c.Value == nil {
		return fmt.Sprintf(" %s ", c.QueryCriteria)
	}
	return fmt.Sprintf(" %s %v ", c.QueryCriteria, c.Value)
}

type ColumnCriteria string

const (
	PATIENTID            ColumnCriteria = "patient_id = "
	DOCTORID             ColumnCriteria = "doctor_id = "
	ANSWERAT_ISNULL      ColumnCriteria = "answer_at IS NULL"
	ANSWERAT_ISNOTNULL   ColumnCriteria = "answer_at IS NOT NULL"
	DATE_GREATERTHAN     ColumnCriteria = "date > "
	DATE_LESSTHAN        ColumnCriteria = "date < "
	CREATEAT_GREATERTHAN ColumnCriteria = "create_at > "
)

func attachCriteria(db *gorm.DB, criteria ...Criteria) *gorm.DB {
	if len(criteria) == 0 {
		return db
	}
	for _, c := range criteria {
		db = db.Where(c.ToString())
	}
	return db
}
