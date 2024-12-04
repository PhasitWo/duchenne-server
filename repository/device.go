package repository

import (
	"fmt"
	"strconv"

	"github.com/PhasitWo/duchenne-server/model"
)

var allDeviceQuery = `
SELECT
id,
login_at,
device_name,
expo_token,
patient_id
FROM device
`

// Get all appointments with following criteria
func (r *Repo) GetAllDevice(id int, criteria QueryCriteria) ([]model.Device, error) {
	var queryString string
	switch criteria {
	case PATIENTID:
		queryString = allDeviceQuery + " " + string(PATIENTID) + strconv.Itoa(id)
	case NONE:
		queryString = allDeviceQuery
	default:
		return nil, fmt.Errorf("query : invalid criteria")
	}
	rows, err := r.db.Query(queryString + " ORDER BY login_at ASC")
	if err != nil {
		return nil, fmt.Errorf("query : %w", err)
	}
	defer rows.Close()
	res := []model.Device{}
	for rows.Next() {
		var d model.Device
		if err := rows.Scan(
			&d.Id,
			&d.LoginAt,
			&d.DeviceName,
			&d.ExpoToken,
			&d.PatientId,
		); err != nil {
			return nil, fmt.Errorf("query : %w", err)
		}
		res = append(res, d)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("query : %w", err)
	}
	return res, nil
}

var updateDeviceQuery = "UPDATE device SET login_at=?, device_name=?, expo_token=?, patient_id=? WHERE id = ?"

func (r *Repo) UpdateDevice(d model.Device) error {
	result, err := r.db.Exec(updateDeviceQuery, d.LoginAt, d.DeviceName, d.ExpoToken, d.PatientId, d.Id)
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
