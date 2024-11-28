package repository

import (
	"fmt"
	"strconv"

	"github.com/PhasitWo/duchenne-server/model"
)

var appointmentQuery = `
SELECT
appointment.id,
create_at,
date,
patient_id,
patient.hn,
patient.first_name,
patient.middle_name,
patient.last_name,
patient.email,
patient.phone,
patient.verified,
doctor_id,
doctor.first_name,
doctor.middle_name,
doctor.last_name
FROM appointment
INNER JOIN patient ON appointment.patient_id = patient.id 
INNER JOIN doctor ON appointment.doctor_id = doctor.id
WHERE appointment.id = ?
`

func (r *Repo) GetAppointment(appointmentId int) (model.Appointment, error) {
	var ap model.Appointment
	row := r.db.QueryRow(appointmentQuery, appointmentId)
	if err := row.Scan(
		&ap.Id,
		&ap.CreateAt,
		&ap.Date,
		&ap.Patient.Id,
		&ap.Patient.Hn,
		&ap.Patient.FirstName,
		&ap.Patient.MiddleName,
		&ap.Patient.LastName,
		&ap.Patient.Email,
		&ap.Patient.Phone,
		&ap.Patient.Verified,
		&ap.Doctor.Id,
		&ap.Doctor.FirstName,
		&ap.Doctor.MiddleName,
		&ap.Doctor.LastName,
	); err != nil {
		return ap, fmt.Errorf("query : %w", err)
	}
	return ap, nil
}

type AppointmentCriteria string

const (
	PATIENTID AppointmentCriteria = "WHERE patient_id = "
	DOCTORID  AppointmentCriteria = "WHERE doctor_id = "
	NONE      AppointmentCriteria = ""
)

var allAppointmentQuery = `
SELECT
appointment.id,
create_at,
date,
patient_id,
patient.hn,
patient.first_name,
patient.middle_name,
patient.last_name,
patient.email,
patient.phone,
patient.verified,
doctor_id,
doctor.first_name,
doctor.middle_name,
doctor.last_name
FROM appointment
INNER JOIN patient ON appointment.patient_id = patient.id 
INNER JOIN doctor ON appointment.doctor_id = doctor.id
`

// Get all appointment with following criteria
func (r *Repo) GetAllAppointment(id int, criteria AppointmentCriteria) ([]model.Appointment, error) {
	var queryString string
	switch criteria {
	case PATIENTID:
		queryString = allAppointmentQuery + " " + string(PATIENTID) + strconv.Itoa(id)
	case DOCTORID:
		queryString = allAppointmentQuery + " " + string(DOCTORID) + strconv.Itoa(id)
	case NONE:
		queryString = allAppointmentQuery
	default:
		return nil, fmt.Errorf("query : invalid criteria")
	}
	fmt.Println(queryString)
	rows, err := r.db.Query(queryString)
	if err != nil {
		return nil, fmt.Errorf("query : %w", err)
	}
	defer rows.Close()
	var res []model.Appointment
	for rows.Next() {
		var ap model.Appointment
		if err := rows.Scan(
			&ap.Id,
			&ap.CreateAt,
			&ap.Date,
			&ap.Patient.Id,
			&ap.Patient.Hn,
			&ap.Patient.FirstName,
			&ap.Patient.MiddleName,
			&ap.Patient.LastName,
			&ap.Patient.Email,
			&ap.Patient.Phone,
			&ap.Patient.Verified,
			&ap.Doctor.Id,
			&ap.Doctor.FirstName,
			&ap.Doctor.MiddleName,
			&ap.Doctor.LastName,
		); err != nil {
			return nil, fmt.Errorf("query : %w", err)
		}
		res = append(res, ap)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("query : %w", err)
	}
	return res, nil
}

var createAppointmentQuery = `
INSERT INTO appointment (create_at, date, patient_id, doctor_id)
VALUES (?, ?, ?, ?)
`

func (r *Repo) CreateAppointment(create_at int, date int, patient_id int, doctor_id int) error {
	result, err := r.db.Exec(createAppointmentQuery, create_at, date, patient_id, doctor_id)
	if err != nil {
		return fmt.Errorf("exec : %w", err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("exec : %w", err)
	}
	if rows != 1 {
		return fmt.Errorf("exec : no affected row")
	}
	return nil
}
