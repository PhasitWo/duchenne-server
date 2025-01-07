package repository

import (
	"errors"
	"fmt"

	"github.com/PhasitWo/duchenne-server/model"
	"github.com/go-sql-driver/mysql"
)

const getDoctorByUsernameQuery = "SELECT id, first_name, middle_name, last_name, username, password, role FROM doctor WHERE username = ?"

func (r *Repo) GetDoctorByUsername(username string) (model.Doctor, error) {
	var d model.Doctor
	row := r.db.QueryRow(getDoctorByUsernameQuery, username)
	if err := row.Scan(&d.Id, &d.FirstName, &d.MiddleName, &d.LastName, &d.Username, &d.Password, &d.Role); err != nil {
		return d, fmt.Errorf("query : %w", err)
	}
	return d, nil
}

const getDoctorByIdQuery = "SELECT id, first_name, middle_name, last_name, username, password, role FROM doctor WHERE id = ?"

func (r *Repo) GetDoctorById(id any) (model.Doctor, error) {
	var d model.Doctor
	row := r.db.QueryRow(getDoctorByIdQuery, id)
	if err := row.Scan(&d.Id, &d.FirstName, &d.MiddleName, &d.LastName, &d.Username, &d.Password, &d.Role); err != nil {
		return d, fmt.Errorf("query : %w", err)
	}
	return d, nil
}

var allDoctorQuery = "SELECT id, first_name, middle_name, last_name, role FROM doctor"

func (r *Repo) GetAllDoctor() ([]model.TrimDoctor, error) {
	rows, err := r.db.Query(allDoctorQuery)
	if err != nil {
		return nil, fmt.Errorf("query : %w", err)
	}
	defer rows.Close()

	var res []model.TrimDoctor
	for rows.Next() {
		var d model.TrimDoctor
		if err := rows.Scan(&d.Id, &d.FirstName, &d.MiddleName, &d.LastName, &d.Role); err != nil {
			return nil, fmt.Errorf("query : %w", err)
		}
		res = append(res, d)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("query : %w", err)
	}
	return res, nil
}

const createDoctorQuery = `
INSERT INTO doctor (first_name, middle_name, last_name, username, password, role)
VALUES (?, ?, ?, ?, ?, ?)
`

// return last inserted id
func (r *Repo) CreateDoctor(doctor model.Doctor) (int, error) {
	result, err := r.db.Exec(createDoctorQuery, doctor.FirstName, doctor.MiddleName, doctor.LastName, doctor.Username, doctor.Password, doctor.Role)
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

const updateDoctorQuery = `
UPDATE doctor SET first_name = ?, middle_name=?, last_name=?, username=?, password=?, role=?
WHERE id = ?`

func (r *Repo) UpdateDoctor(doctor model.Doctor) error {
	// update should be idempotent -> error occur when this handler is called consecutively with same input -> err no affected row
	_, err := r.db.Exec(updateDoctorQuery, doctor.FirstName, doctor.MiddleName, doctor.LastName, doctor.Username, doctor.Password, doctor.Role, doctor.Id)
	if err != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
			return fmt.Errorf("exec : %w", ErrDuplicateEntry)
		}
		return fmt.Errorf("exec : %w", err)
	}
	return nil
}

const deleteDoctorQuery =  `
DELETE FROM doctor
WHERE id = ?;
`

func (r *Repo) DeleteDoctorById(id any) error {
	_, err := r.db.Exec(deleteDoctorQuery, id)
	if err != nil {
		return fmt.Errorf("exec : %w", err)
	}
	return nil
}