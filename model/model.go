package model

import (
	"database/sql"
)

type Patient struct {
	Id         int            `json:"id"`
	Hn         string         `json:"hn"`
	FirstName  string         `json:"firstName"`
	MiddleName sql.NullString `json:"middleName"` // nullable
	LastName   string         `json:"lastName"`
	Email      sql.NullString `json:"email"` // nullable
	Phone      sql.NullString `json:"phone"` // nullable
	Verified   bool           `json:"verified"`
}
