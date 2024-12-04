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
	Id         int     `json:"id"`
	Hn         string  `json:"hn"`
	FirstName  string  `json:"firstName"`
	MiddleName *string `json:"middleName"` // nullable
	LastName   string  `json:"lastName"`
	Email      *string `json:"email"` // nullable
	Phone      *string `json:"phone"` // nullable
	Verified   bool    `json:"verified"`
}

type Doctor struct {
	Id         int     `json:"id"`
	FirstName  string  `json:"firstName"`
	MiddleName *string `json:"middleName"` // nullable
	LastName   string  `json:"lastName"`
	Username   string  `json:"username"`
	Password   string  `json:"password"`
	Role       Role    `json:"role"`
}

type TrimDoctor struct {
	Id         int     `json:"id"`
	FirstName  string  `json:"firstName"`
	MiddleName *string `json:"middleName"` // nullable
	LastName   string  `json:"lastName"`
}

type Appointment struct {
	Id       int        `json:"id"`
	CreateAt int        `json:"createAt"`
	Date     int        `json:"date"`
	Patient  Patient    `json:"patient"`
	Doctor   TrimDoctor `json:"doctor"`
}

type Question struct {
	Id       int         `json:"id"`
	Topic    string      `json:"topic"`
	Question string      `json:"question"`
	CreateAt int         `json:"createAt"`
	Answer   *string     `json:"answer"`   // nullable
	AnswerAt *int        `json:"answerAt"` // nullable
	Patient  Patient     `json:"patient"`
	Doctor   *TrimDoctor `json:"doctor"` // nullable
}

type Device struct {
	Id         int    `json:"id"`
	LoginAt    int    `json:"loginAt"`
	DeviceName string `json:"deviceName"`
	ExpoToken  string `json:"expoToken"`
	PatientId  int    `json:"patientId"`
}
