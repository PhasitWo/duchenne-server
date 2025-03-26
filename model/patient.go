package model

import (
	"gorm.io/datatypes"
)

type VaccineHistory struct {
	Id              string  `json:"id" binding:"required"`
	VaccineName     string  `json:"vaccineName" binding:"required"`
	VaccineLocation *string `json:"vaccineLocation"`
	VaccineAt       int     `json:"vaccineAt" binding:"required"`
	Description     *string `json:"description"` // nullable
}

type Medicine struct {
	Id           string  `json:"id" binding:"required"`
	MedicineName string  `json:"medicineName" binding:"required"`
	Description  *string `json:"description"` // nullable
}

type Patient struct {
	ID             int                                 `json:"id"`
	Hn             string                              `json:"hn" gorm:"unique;not null"`
	FirstName      string                              `json:"firstName" gorm:"not null"`
	MiddleName     *string                             `json:"middleName"` // nullable
	LastName       string                              `json:"lastName" gorm:"not null"`
	Email          *string                             `json:"email"` // nullable
	Phone          *string                             `json:"phone"` // nullable
	Verified       bool                                `json:"verified" gorm:"not null;default:0"`
	VaccineHistory datatypes.JSONSlice[VaccineHistory] `json:"vaccineHistory"` // nullable
	Medicine       datatypes.JSONSlice[Medicine]       `json:"medicine"`       // nullable
}
