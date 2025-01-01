package repository

import (
	"fmt"
	"github.com/PhasitWo/duchenne-server/model"
)

var allDeviceQuery = `SELECT id, login_at, device_name, expo_token, patient_id FROM device`

func (r *Repo) GetAllDevice(criteria ...Criteria) ([]model.Device, error) {
	queryString := attachCriteria(allDeviceQuery, criteria...)
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

var createDeviceQuery = `
INSERT INTO device (login_at, device_name, expo_token, patient_id)
VALUES (?, ?, ?, ?)
`

func (r *Repo) CreateDevice(d model.Device) (int, error) {
	result, err := r.db.Exec(createDeviceQuery, d.LoginAt, d.DeviceName, d.ExpoToken, d.PatientId)
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

var deleteDeviceQuery = `
DELETE FROM device
WHERE id = ?;
`

func (r *Repo) DeleteDevice(deviceId any) error {
	_, err := r.db.Exec(deleteDeviceQuery, deviceId)
	if err != nil {
		return fmt.Errorf("exec : %w", err)
	}
	return nil
}
