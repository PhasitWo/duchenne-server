package model

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
}

type SafeAppointment struct {
	Appointment
	Doctor TrimDoctor `json:"doctor"`
}