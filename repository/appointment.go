package repository

import (
	"fmt"
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

func (r *Repo) GetAppointment(appointmentId any) (model.Appointment, error) {
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

// Get all appointments with following criteria
func (r *Repo) GetAllAppointment(criteria ...Criteria) ([]model.Appointment, error) {
	queryString := attachCriteria(allAppointmentQuery, criteria...)
	rows, err := r.db.Query(queryString  + " ORDER BY date ASC LIMIT 15")
	if err != nil {
		return nil, fmt.Errorf("query : %w", err)
	}
	defer rows.Close()
	res := []model.Appointment{}
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

func (r *Repo) CreateAppointment(createAt int, date int, patientId int, doctorId int) (int, error) {
	result, err := r.db.Exec(createAppointmentQuery, createAt, date, patientId, doctorId)
	if err != nil {
		return -1, fmt.Errorf("exec : %w", err)
	}
	i, err := result.LastInsertId()
	if err != nil {
		return -1, fmt.Errorf("exec : %w", err)
	}
	lastId := int(i)
	return lastId, nil
}

var deleteAppointmentQuery = `
DELETE FROM appointment
WHERE id = ?;
`

func (r *Repo) DeleteAppointment(appointmentId any) error {
	_, err := r.db.Exec(deleteAppointmentQuery, appointmentId)
	if err != nil {
		return fmt.Errorf("exec : %w", err)
	}
	return nil
}
