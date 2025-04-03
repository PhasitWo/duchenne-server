package model

import "gorm.io/plugin/soft_delete"

type Appointment struct {
	ID        int     `json:"id"`
	CreateAt  int     `json:"createAt" gorm:"not null"`
	UpdateAt  int     `json:"updateAt" gorm:"autoUpdateTime;not null"`
	Date      int     `json:"date" gorm:"not null"`
	PatientID int     `json:"-" gorm:"not null"`
	Patient   Patient `json:"patient"`
	DoctorID  int     `json:"-" gorm:"not null"`
	Doctor    Doctor  `json:"doctor"`
	ApproveAt *int    `json:"approveAt"` // nullable
	DeletedAt soft_delete.DeletedAt
}

type SafeAppointment struct {
	Appointment
	Doctor TrimDoctor `json:"doctor"`
}