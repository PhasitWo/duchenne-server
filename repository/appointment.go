package repository

import (
	"fmt"
	// "strconv"

	"github.com/PhasitWo/duchenne-server/model"
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

func (r *Repo) CreateAppointment(createAt int, date int, patientId int, doctorId int) (int, error) {
	ap := model.Appointment{CreateAt: createAt, Date: date, PatientID: patientId, DoctorID: doctorId}
	err := r.db.Create(&ap).Error
	if err != nil {
		return -1, fmt.Errorf("exec : %w", err)
	}
	return ap.ID, nil
}

func (r *Repo) DeleteAppointment(appointmentId any) error {
	err := r.db.Where("id = ?", appointmentId).Delete(&model.Appointment{}).Error
	if err != nil {
		return fmt.Errorf("exec : %w", err)
	}
	return nil
}
