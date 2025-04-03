package repository

import (
	// "database/sql"
	"errors"

	"github.com/PhasitWo/duchenne-server/model"
	"gorm.io/gorm"
)


type Repo struct {
	db *gorm.DB
}

// Constructor
func New(db *gorm.DB) *Repo {
	return &Repo{db: db}
}

// ERROR
var ErrDuplicateEntry = errors.New("duplicate entry")
var ErrForeignKeyFail = errors.New("foreign key error")

type IRepo interface {
	GetAppointment(appointmentId any) (model.SafeAppointment, error)
	GetAllAppointment(limit int, offset int, criteria ...Criteria) ([]model.SafeAppointment, error)
	CreateAppointment(appointment model.Appointment) (int, error)
	UpdateAppointment(appointment model.Appointment) error
	DeleteAppointment(appointmentId any) error
	GetAllDevice(criteria ...Criteria) ([]model.Device, error)
	UpdateDevice(d model.Device) error
	CreateDevice(d model.Device) (int, error)
	DeleteDevice(deviceId any) error
	GetDoctorByUsername(username string) (model.Doctor, error)
	GetDoctorById(id any) (model.Doctor, error)
	GetAllDoctor() ([]model.TrimDoctor, error)
	CreateDoctor(doctor model.Doctor) (int, error)
	UpdateDoctor(doctor model.Doctor) error
	DeleteDoctorById(id any) error
	GetPatientById(id any) (model.Patient, error)
	GetPatientByHN(hn string) (model.Patient, error)
	GetAllPatient() ([]model.Patient, error)
	CreatePatient(patient model.Patient) (int, error)
	UpdatePatient(patient model.Patient) error
	UpdatePatientVaccineHistory(patientId int, vaccineHistory []model.VaccineHistory) error
	UpdatePatientMedicine(patientId int, medicines []model.Medicine) error
	DeletePatientById(id any) error
	GetQuestion(questionId any) (model.SafeQuestion, error)
	GetAllQuestion(limit int, offset int, criteria ...Criteria) ([]model.QuestionTopic, error)
	CreateQuestion(patientId int, topic string, question string, createAt int) (int, error)
	UpdateQuestionAnswer(questionId int, answer string, doctorId int) error
	DeleteQuestion(questionId any) error
}