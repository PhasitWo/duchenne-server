package repository

import (
	"database/sql"
	"fmt"

	"github.com/PhasitWo/duchenne-server/model"
	_ "github.com/go-sql-driver/mysql"
)

type Repo struct {
	db *sql.DB
}

// Constructor
func New(db *sql.DB) *Repo {
	return &Repo{db: db}
}

// Patient
func (r *Repo) GetPatient(hn string) (model.Patient, error) {
	row := r.db.QueryRow("SELECT id, hn, first_name, middle_name, last_name, email, phone, verified FROM patient WHERE hn=?", hn)

	var p model.Patient
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

func (r *Repo) VerifyPatient(id int) error {
	result, err := r.db.Exec("UPDATE patient SET verified = 1 WHERE id = ?", id)
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
