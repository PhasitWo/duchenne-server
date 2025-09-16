package model

import (
	"gorm.io/datatypes"
	"gorm.io/plugin/soft_delete"
)

type VaccineHistory struct {
	Id              string  `json:"id" binding:"required"`
	VaccineName     string  `json:"vaccineName" binding:"required"`
	VaccineLocation *string `json:"vaccineLocation"`
	VaccineAt       int     `json:"vaccineAt" binding:"required"`
	Complication    *string `json:"complication"` // nullable
}

type Medicine struct {
	Id              string  `json:"id" binding:"required"`
	MedicineName    string  `json:"medicineName" binding:"required"`
	Dose            *string `json:"dose"`            // nullable
	FrequencyPerDay *string `json:"frequencyPerDay"` // nullable
	Instruction     *string `json:"instruction"`     // nullable
	Quantity        *string `json:"quantity"`        // nullable
}

type Patient struct {
	ID             int                                 `json:"id"`
	NID            string                              `json:"nid" gorm:"unique;not null;column:nid"`
	Password       string                              `gorm:"not null"`
	Hn             string                              `json:"hn" gorm:"unique;not null"`
	Pin            string                              `gorm:"not null"`
	FirstName      string                              `json:"firstName" gorm:"not null"`
	MiddleName     *string                             `json:"middleName"` // nullable
	LastName       string                              `json:"lastName" gorm:"not null"`
	Email          *string                             `json:"email"` // nullable
	Phone          *string                             `json:"phone"` // nullable
	Verified       bool                                `json:"verified" gorm:"not null;default:0"`
	Weight         *float32                            `json:"weight"` // nullable
	Height         *float32                            `json:"height"` // nullable
	BirthDate      int                                 `json:"birthDate" gorm:"not null"`
	VaccineHistory datatypes.JSONSlice[VaccineHistory] `json:"vaccineHistory"` // nullable
	Medicine       datatypes.JSONSlice[Medicine]       `json:"medicine"`       // nullable
	DeletedAt      soft_delete.DeletedAt               `gorm:"default:0"`
}

// type CreatePatientRequest struct {
// 	Hn         string   `json:"hn" binding:"required"`
// 	FirstName  string   `json:"firstName" binding:"required"`
// 	MiddleName *string  `json:"middleName"`
// 	LastName   string   `json:"lastName" binding:"required"`
// 	Email      *string  `json:"email"`
// 	Phone      *string  `json:"phone"`
// 	Verified   bool     `json:"verified"`
// 	Weight     *float32 `json:"weight"`
// 	Height     *float32 `json:"height"`
// 	BirthDate  int      `json:"birthDate"`
// }

type UpdatePatientRequest struct {
	NID        string   `json:"nid" binding:"required,min=13"`
	Hn         string   `json:"hn" binding:"required"`
	FirstName  string   `json:"firstName" binding:"required"`
	MiddleName *string  `json:"middleName"`
	LastName   string   `json:"lastName" binding:"required"`
	Email      *string  `json:"email"`
	Phone      *string  `json:"phone"`
	Verified   bool     `json:"verified"`
	Weight     *float32 `json:"weight"`
	Height     *float32 `json:"height"`
	BirthDate  int      `json:"birthDate"`
}

type UpdateVaccineHistoryRequest struct {
	Data []VaccineHistory `json:"data" binding:"dive"`
}

type UpdateMedicineRequest struct {
	Data []Medicine `json:"data" binding:"dive"`
}
