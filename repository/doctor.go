package repository

import (
	"fmt"

	"github.com/PhasitWo/duchenne-server/model"
)

var allDoctorQuery = "SELECT id, first_name, middle_name, last_name FROM doctor"

func (r *Repo) GetAllDoctor() ([]model.TrimDoctor, error) {
	rows, err := r.db.Query(allDoctorQuery)
	if err != nil {
		return nil, fmt.Errorf("query : %w", err)
	}
	defer rows.Close()

	var res []model.TrimDoctor
	for rows.Next() {
		var d model.TrimDoctor
		if err := rows.Scan(&d.Id, &d.FirstName, &d.MiddleName, &d.LastName); err != nil {
			return nil, fmt.Errorf("query : %w", err)
		}
		res = append(res, d)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("query : %w", err)
	}
	return res, nil
}
