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
		return string(c.QueryCriteria)
	}
	return fmt.Sprintf(string(c.QueryCriteria), c.Value)
}

type ColumnCriteria string

const (
	PATIENTID            ColumnCriteria = "patient_id = %v"
	DOCTORID             ColumnCriteria = "doctor_id = %v"
	ANSWERAT_ISNULL      ColumnCriteria = "answer_at IS NULL"
	ANSWERAT_ISNOTNULL   ColumnCriteria = "answer_at IS NOT NULL"
	DATE_GREATERTHAN     ColumnCriteria = "date > %v"
	DATE_LESSTHAN        ColumnCriteria = "date < %v"
	CREATEAT_GREATERTHAN ColumnCriteria = "create_at > %v"
	IS_PUBLISHED         ColumnCriteria = "is_published = %v"
	CAN_BE_APPOINTED     ColumnCriteria = "can_be_appointed = %v"
	DOCTOR_SEARCH        ColumnCriteria = "first_name ILIKE '%%%[1]v%%' OR middle_name ILIKE '%%%[1]v%%' OR last_name ILIKE '%%%[1]v%%'"
	PATIENT_SEARCH       ColumnCriteria = "first_name ILIKE '%%%[1]v%%' OR middle_name ILIKE '%%%[1]v%%' OR last_name ILIKE '%%%[1]v%%' OR hn ILIKE '%%%[1]v%%'"
	QUESTION_SEARCH      ColumnCriteria = "topic ILIKE '%%%v%%'"
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
