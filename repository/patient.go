package repository

import (
	"errors"
	"fmt"

	"github.com/PhasitWo/duchenne-server/model"
	"github.com/go-sql-driver/mysql"
)

var GetPatientByIdQuery = `SELECT id, hn, first_name, middle_name, last_name, email, phone, verified FROM patient WHERE id=?`

func (r *Repo) GetPatientById(id any) (model.Patient, error) {
	var p model.Patient
	row := r.db.QueryRow(GetPatientByIdQuery, id)
	if err := row.Scan(&p.Id, &p.Hn, &p.FirstName, &p.MiddleName, &p.LastName, &p.Email, &p.Phone, &p.Verified); err != nil {
		return p, fmt.Errorf("query : %w", err)
	}
	return p, nil
}

var GetPatientByHNQuery = `SELECT id, hn, first_name, middle_name, last_name, email, phone, verified FROM patient WHERE hn=?`

func (r *Repo) GetPatientByHN(hn string) (model.Patient, error) {
	var p model.Patient
	row := r.db.QueryRow(GetPatientByHNQuery, hn)
	if err := row.Scan(&p.Id, &p.Hn, &p.FirstName, &p.MiddleName, &p.LastName, &p.Email, &p.Phone, &p.Verified); err != nil {
		return p, fmt.Errorf("query : %w", err)
	}
	return p, nil
}

func (r *Repo) GetAllPatient() ([]model.Patient, error) {
	rows, err := r.db.Query("SELECT id, hn, first_name, middle_name, last_name, email, phone, verified FROM patient")
	if err != nil {
		return nil, fmt.Errorf("query : %w", err)
	}
	defer rows.Close()

	var res []model.Patient
	for rows.Next() {
		var p model.Patient
		if err := rows.Scan(&p.Id, &p.Hn, &p.FirstName, &p.MiddleName, &p.LastName, &p.Email, &p.Phone, &p.Verified); err != nil {
			return nil, fmt.Errorf("query : %w", err)
		}
		res = append(res, p)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("query : %w", err)
	}
	return res, nil
}

const createPatientQuery = `
INSERT INTO patient (hn ,first_name, middle_name, last_name, email, phone, verified)
VALUES (?, ?, ?, ?, ?, ?, ?)
`

// return last inserted id
func (r *Repo) CreatePatient(patient model.Patient) (int, error) {
	result, err := r.db.Exec(
		createPatientQuery,
		patient.Hn,
		patient.FirstName,
		patient.MiddleName,
		patient.LastName,
		patient.Email,
		patient.Phone,
		patient.Verified)
	if err != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
			return -1, fmt.Errorf("exec : %w", ErrDuplicateEntry)
		}
		return -1, fmt.Errorf("exec : %w", err)
	}
	i, err := result.LastInsertId()
	if err != nil {
		return -1, fmt.Errorf("exec : %w", err)
	}
	lastId := int(i)
	return lastId, nil
}

const updatePatientQuery = `
UPDATE patient SET hn=? ,first_name = ?, middle_name=?, last_name=?, email=?, phone=?, verified=?
WHERE id = ?
`

func (r *Repo) UpdatePatient(patient model.Patient) error {
	// update should be idempotent -> error occur when this handler is called consecutively with same input -> err no affected row
	_, err := r.db.Exec(
		updatePatientQuery,
		patient.Hn,
		patient.FirstName,
		patient.MiddleName,
		patient.LastName,
		patient.Email,
		patient.Phone,
		patient.Verified,
		patient.Id)
	if err != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
			return fmt.Errorf("exec : %w", ErrDuplicateEntry)
		}
		return fmt.Errorf("exec : %w", err)
	}
	return nil
}

const deletePatientrQuery =  `
DELETE FROM patient
WHERE id = ?;
`

func (r *Repo) DeletePatientById(id any) error {
	_, err := r.db.Exec(deletePatientrQuery, id)
	if err != nil {
		return fmt.Errorf("exec : %w", err)
	}
	return nil
}