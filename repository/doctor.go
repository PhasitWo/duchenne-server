package repository

import (
	"errors"
	"fmt"

	"github.com/PhasitWo/duchenne-server/model"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

func (r *Repo) GetDoctorByUsername(username string) (model.Doctor, error) {
	var d model.Doctor
	err := r.db.Where("username = ?", username).First(&d).Error
	if err != nil {
		return d, fmt.Errorf("query : %w", err)
	}
	return d, nil
}

func (r *Repo) GetDoctorById(id any) (model.Doctor, error) {
	var d model.Doctor
	err := r.db.Where("id = ?", id).First(&d).Error
	if err != nil {
		return d, fmt.Errorf("query : %w", err)
	}
	return d, nil
}

func (r *Repo) GetAllDoctor(limit int, offset int, criteria ...Criteria) ([]model.TrimDoctor, error) {
	var res []model.TrimDoctor
	db := attachCriteria(r.db, criteria...)
	err := db.Model(&model.Doctor{}).Limit(limit).Offset(offset).Find(&res).Error
	if err != nil {
		return res, fmt.Errorf("query : %w", err)
	}
	return res, nil
}

// return last inserted id
func (r *Repo) CreateDoctor(doctor model.Doctor) (int, error) {
	err := r.db.Create(&doctor).Error
	if err != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
			return -1, fmt.Errorf("exec : %w", ErrDuplicateEntry)
		}
		return -1, fmt.Errorf("exec : %w", err)
	}
	return doctor.ID, nil
}

func (r *Repo) UpdateDoctor(doctor model.Doctor) error {
	err := r.db.Select("*").Updates(&doctor).Error
	if err != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
			return fmt.Errorf("exec : %w", ErrDuplicateEntry)
		}
		return fmt.Errorf("exec : %w", err)
	}
	return nil
}

func (r *Repo) DeleteDoctorById(id any) error {
	err := r.db.Transaction(func(tx *gorm.DB) error {
		// soft delete appointment
		err := tx.Where("doctor_id = ?", id).Delete(&model.Appointment{}).Error
		if err != nil {
			return err
		}
		// set question's doctorId as null
		err = tx.Model(&model.Question{}).Where("doctor_id = ?", id).Update("doctor_id", nil).Error
		if err != nil {
			return err
		}
		// soft delete doctor
		err = r.db.Where("id = ?", id).Delete(&model.Doctor{}).Error
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("exec : %w", err)
	}
	return err
}
