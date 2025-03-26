package repository

import (
	"errors"
	"fmt"

	"github.com/PhasitWo/duchenne-server/model"
	"github.com/go-sql-driver/mysql"
	"gorm.io/datatypes"
)

var GetPatientByIdQuery = `SELECT id, hn, first_name, middle_name, last_name, email, phone, verified FROM patient WHERE id=?`

func (r *Repo) GetPatientById(id any) (model.Patient, error) {
	var p model.Patient
	err := r.db.Where("id = ?", id).First(&p).Error
	if err != nil {
		return p, fmt.Errorf("query : %w", err)
	}
	return p, nil
}

var GetPatientByHNQuery = `SELECT id, hn, first_name, middle_name, last_name, email, phone, verified FROM patient WHERE hn=?`

func (r *Repo) GetPatientByHN(hn string) (model.Patient, error) {
	var p model.Patient
	err := r.db.Where("hn = ?", hn).First(&p).Error
	if err != nil {
		return p, fmt.Errorf("query : %w", err)
	}
	return p, nil
}

func (r *Repo) GetAllPatient() ([]model.Patient, error) {
	var res []model.Patient
	err := r.db.Model(&model.Patient{}).Find(&res).Error
	if err != nil {
		return res, fmt.Errorf("query : %w", err)
	}
	return res, nil
}

// return last inserted id
func (r *Repo) CreatePatient(patient model.Patient) (int, error) {
	err := r.db.Create(&patient).Error
	if err != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
			return -1, fmt.Errorf("exec : %w", ErrDuplicateEntry)
		}
		return -1, fmt.Errorf("exec : %w", err)
	}
	return patient.ID, nil
}

func (r *Repo) UpdatePatient(patient model.Patient) error {
	err := r.db.Select("*").Omit("vaccine_history", "medicine").Updates(&patient).Error
	if err != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
			return fmt.Errorf("exec : %w", ErrDuplicateEntry)
		}
		return fmt.Errorf("exec : %w", err)
	}
	return nil
}

func (r *Repo) UpdatePatientVaccineHistory(patientId int, vaccineHistory []model.VaccineHistory) error {
	err := r.db.Select("vaccine_history").Updates(&model.Patient{
		ID:             patientId,
		VaccineHistory: datatypes.NewJSONSlice(vaccineHistory),
	}).Error
	if err != nil {
		return fmt.Errorf("exec : %w", err)
	}
	return nil
}

func (r *Repo) UpdatePatientMedicine(patientId int, medicines []model.Medicine) error {
	err := r.db.Select("medicine").Updates(&model.Patient{
		ID:       patientId,
		Medicine: datatypes.NewJSONSlice(medicines),
	}).Error
	if err != nil {
		return fmt.Errorf("exec : %w", err)
	}
	return nil
}

func (r *Repo) DeletePatientById(id any) error {
	err := r.db.Where("id = ?", id).Delete(&model.Patient{}).Error
	if err != nil {
		return fmt.Errorf("exec : %w", err)
	}
	return nil
}
