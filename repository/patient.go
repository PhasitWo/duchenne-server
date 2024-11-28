package repository

import (
	"fmt"

	"github.com/PhasitWo/duchenne-server/model"
)


// GetPatient accept hn string  or  patient_id int
func (r *Repo) GetPatient(criteria interface{}) (model.Patient, error) {
	var p model.Patient
	var query string
	var v any
	switch val := criteria.(type) {
	case int:
		query = "SELECT id, hn, first_name, middle_name, last_name, email, phone, verified FROM patient WHERE id=?"
		v = val
	case string:
		query = "SELECT id, hn, first_name, middle_name, last_name, email, phone, verified FROM patient WHERE hn=?"
		v = val
	default:
		return p, fmt.Errorf("query : wrong parameter type")
	}
	row := r.db.QueryRow(query, v)
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

func (r *Repo) UpdatePatient(p model.Patient) error {
	result, err := r.db.Exec("UPDATE patient SET hn=? ,first_name = ?, middle_name=?, last_name=?, email=?, phone=?, verified=?  WHERE id = ?", p.Hn, p.FirstName, p.MiddleName, p.LastName, p.Email, p.Phone, p.Verified, p.Id)
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

// TODO: CreatePatient, DeletePatient

// func (r *Repo) VerifyPatient(id int) error {
// 	result, err := r.db.Exec("UPDATE patient SET verified = 1 WHERE id = ?", id)
// 	if err != nil {
// 		return fmt.Errorf("exec : %w", err)
// 	}
// 	rows, err := result.RowsAffected()
// 	if err != nil {
// 		return fmt.Errorf("exec : %w", err)
// 	}
// 	if rows != 1 {
// 		return fmt.Errorf("exec : no affected row")
// 	}
// 	return nil
// }
