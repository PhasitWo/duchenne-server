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
	Hn         string  `json:"hn" gorm:"unique"`
	FirstName  string  `json:"firstName"`
	MiddleName *string `json:"middleName"` // nullable
	LastName   string  `json:"lastName"`
	Email      *string `json:"email"` // nullable
	Phone      *string `json:"phone"` // nullable
	Verified   bool    `json:"verified"`
}

type Doctor struct {
	ID         int     `json:"id"`
	FirstName  string  `json:"firstName"`
	MiddleName *string `json:"middleName"` // nullable
	LastName   string  `json:"lastName"`
	Username   string  `json:"username" gorm:"unique"`
	Password   string  `json:"password"`
	Role       Role    `json:"role"`
}

type TrimDoctor struct {
	ID         int     `json:"id"`
	FirstName  string  `json:"firstName"`
	MiddleName *string `json:"middleName"` // nullable
	LastName   string  `json:"lastName"`
	Role       Role    `json:"role"`
}

type Appointment struct {
	ID        int     `json:"id"`
	CreateAt  int     `json:"createAt"`
	Date      int     `json:"date"`
	PatientID int     `json:"-"`
	Patient   Patient `json:"patient"`
	DoctorID  int     `json:"-"`
	Doctor    Doctor  `json:"doctor"`
}

type SafeAppointment struct {
	Appointment
	Doctor TrimDoctor `json:"doctor"`
}

type Question struct {
	ID        int     `json:"id"`
	Topic     string  `json:"topic"`
	Question  string  `json:"question"`
	CreateAt  int     `json:"createAt"`
	Answer    *string `json:"answer"`   // nullable
	AnswerAt  *int    `json:"answerAt"` // nullable
	PatientID int     `json:"-"`
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
