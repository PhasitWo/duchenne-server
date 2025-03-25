package model

// "database/sql"

// Doctor roles
type Role string

const (
	ROOT  Role = "root"
	ADMIN Role = "admin"
	USER  Role = "user"
)

type Patient struct {
	ID         int     `json:"id"`
	Hn         string  `json:"hn" gorm:"unique;not null"`
	FirstName  string  `json:"firstName" gorm:"not null"`
	MiddleName *string `json:"middleName"` // nullable
	LastName   string  `json:"lastName" gorm:"not null"`
	Email      *string `json:"email"` // nullable
	Phone      *string `json:"phone"` // nullable
	Verified   bool    `json:"verified" gorm:"not null;default:0"`
}

type Doctor struct {
	ID         int     `json:"id"`
	FirstName  string  `json:"firstName" gorm:"not null"`
	MiddleName *string `json:"middleName"` // nullable
	LastName   string  `json:"lastName" gorm:"not null"`
	Username   string  `json:"username" gorm:"unique;not null"`
	Password   string  `json:"password" gorm:"not null"`
	Role       Role    `json:"role" gorm:"not null"`
}

type TrimDoctor struct {
	Doctor
	Username *string `json:"username" gorm:"-"`
	Password *string `json:"password" gorm:"-"`
}

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

type Question struct {
	ID        int     `json:"id"`
	Topic     string  `json:"topic" gorm:"not null"`
	Question  string  `json:"question" gorm:"not null"`
	CreateAt  int     `json:"createAt" gorm:"not null"`
	Answer    *string `json:"answer"`   // nullable
	AnswerAt  *int    `json:"answerAt"` // nullable
	PatientID int     `json:"-" gorm:"not null"`
	Patient   Patient `json:"patient"`
	DoctorID  *int    `json:"-"`
	Doctor    *Doctor `json:"doctor"` // nullable
}

type SafeQuestion struct {
	Question
	Doctor TrimDoctor `json:"doctor"`
}

type QuestionTopic struct {
	ID        int         `json:"id"`
	Topic     string      `json:"topic"`
	CreateAt  int         `json:"createAt"`
	AnswerAt  *int        `json:"answerAt"` // nullable
	PatientID int         `json:"-"`
	Patient   Patient     `json:"patient"`
	DoctorID  int         `json:"-"`
	Doctor    *TrimDoctor `json:"doctor"` // nullable
}

type Device struct {
	ID         int    `json:"id"`
	LoginAt    int    `json:"loginAt"`
	DeviceName string `json:"deviceName"`
	ExpoToken  string `json:"expoToken"`
	PatientId  int    `json:"patientId"`
}

type AppointmentDevice struct {
	AppointmentId int    `json:"appointment_id"`
	Date          int    `json:"date"`
	DeviceId      int    `json:"device_id"`
	DeviceName    string `json:"device_name"`
	ExpoToken     string `json:"expoToken"`
	PatientId     int    `json:"patient_id"`
}
