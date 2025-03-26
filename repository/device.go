package repository

import (
	"fmt"

	"github.com/PhasitWo/duchenne-server/model"
)

func (r *Repo) GetAllDevice(criteria ...Criteria) ([]model.Device, error) {
	res := []model.Device{}
	db := attachCriteria(r.db, criteria...)
	err := db.Find(&res).Error
	if err != nil {
		return nil, fmt.Errorf("query : %w", err)
	}
	return res, nil
}

func (r *Repo) UpdateDevice(d model.Device) error {
	err := r.db.Updates(&d).Error
	if err != nil {
		return fmt.Errorf("exec : %w", err)
	}
	return nil
}

func (r *Repo) CreateDevice(d model.Device) (int, error) {
	err := r.db.Create(&d).Error
	if err != nil {
		return -1, fmt.Errorf("exec : %w", err)
	}
	return d.ID, nil
}

func (r *Repo) DeleteDevice(deviceId any) error {
	err := r.db.Where("id = ?", deviceId).Delete(&model.Device{}).Error
	if err != nil {
		return fmt.Errorf("exec : %w", err)
	}
	return nil
}
