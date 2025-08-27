package repository

import (
	"errors"
	"fmt"

	"github.com/PhasitWo/duchenne-server/model"
	"github.com/go-sql-driver/mysql"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

func (r *Repo) GetPatientById(id any) (model.Patient, error) {
	var p model.Patient
	err := r.db.Where("id = ?", id).First(&p).Error
	if err != nil {
		return p, fmt.Errorf("query : %w", err)
	}
	return p, nil
}

func (r *Repo) GetPatientByHN(hn string) (model.Patient, error) {
	var p model.Patient
	err := r.db.Where("hn = ?", hn).First(&p).Error
	if err != nil {
		return p, fmt.Errorf("query : %w", err)
	}
	return p, nil
}

func (r *Repo) GetAllPatient(limit int, offset int, criteria ...Criteria) ([]model.Patient, error) {
	var res []model.Patient
	db := attachCriteria(r.db, criteria...)
	err := db.Model(&model.Patient{}).Limit(limit).Offset(offset).Find(&res).Error
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
	err := r.db.Transaction(func(tx *gorm.DB) error {
		// soft delete appointment
		err := tx.Where("patient_id = ?", id).Delete(&model.Appointment{}).Error
		if err != nil {
			return err
		}
		// soft delete question
		err = tx.Where("patient_id = ?", id).Delete(&model.Question{}).Error
		if err != nil {
			return err
		}
		// change hn
		var storePatient model.Patient
		err = tx.Where("id = ?", id).First(&storePatient).Error
		if err != nil {
			return err
		}
		err = tx.Updates(model.Patient{ID: storePatient.ID, Hn: "DEL" + storePatient.Hn}).Error
		if err != nil {
			return err
		}
		// soft delete patient
		err = tx.Where("id = ?", id).Delete(&model.Patient{}).Error
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("exec : %w", err)
	}
	return nil
}
