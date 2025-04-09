package repository

import (
	"fmt"

	"gorm.io/gorm"
)

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
	IS_PUBLISHED         ColumnCriteria = "is_published = "
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
