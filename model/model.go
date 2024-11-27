package model

import (
	// "database/sql"
)

type Patient struct {
	Id         int     `json:"id"`
	Hn         string  `json:"hn"`
	FirstName  string  `json:"firstName"`
	MiddleName *string `json:"middleName"` // nullable
	LastName   string  `json:"lastName"`
	Email      *string `json:"email"` // nullable
	Phone      *string `json:"phone"` // nullable
	Verified   bool    `json:"verified"`
}
