package repository

import (
	"errors"
	"fmt"
	"time"

	// "strconv"

	"github.com/PhasitWo/duchenne-server/model"
	"github.com/go-sql-driver/mysql"
)

func (r *Repo) GetAppointment(appointmentId any) (model.SafeAppointment, error) {
	var ap model.SafeAppointment
	err := r.db.Model(&model.Appointment{}).Joins("Doctor").Preload("Patient").Where("Appointments.id = ?", appointmentId).First(&ap).Error
	if err != nil {
		return ap, fmt.Errorf("exec : %w", err)
	}
	return ap, nil
}

// Get all appointments with following criteria
func (r *Repo) GetAllAppointment(limit int, offset int, criteria ...Criteria) ([]model.SafeAppointment, error) {
	res := []model.SafeAppointment{}
	db := attachCriteria(r.db, criteria...)
	err := db.Model(&model.Appointment{}).Joins("Doctor").Preload("Patient").Limit(limit).Offset(offset).Order("date ASC").Find(&res).Error
	if err != nil {
		return res, fmt.Errorf("exec : %w", err)
	}
	return res, nil
}

func (r *Repo) CreateAppointment(appointment model.Appointment) (int, error) {
	now := int(time.Now().Unix())
	appointment.CreateAt = now
	appointment.UpdateAt = now
	err := r.db.Create(&appointment).Error
	if err != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1452 {
			return -1, fmt.Errorf("exec : %w", ErrForeignKeyFail)
		}
		return -1, fmt.Errorf("exec : %w", err)
	}
	return appointment.ID, nil
}

func (r *Repo) UpdateAppointment(appointment model.Appointment) error {
	result := r.db.Omit("create_at").Updates(&appointment)
	err := result.Error
	if err != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1452 {
			return fmt.Errorf("exec : %w", ErrForeignKeyFail)
		}
		return fmt.Errorf("exec : %w", err)
	}
	return nil
}

func (r *Repo) DeleteAppointment(appointmentId any) error {
	err := r.db.Where("id = ?", appointmentId).Delete(&model.Appointment{}).Error
	if err != nil {
		return fmt.Errorf("exec : %w", err)
	}
	return nil
}
