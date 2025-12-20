package web

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/PhasitWo/duchenne-server/model"
	"github.com/PhasitWo/duchenne-server/repository"
	"github.com/PhasitWo/duchenne-server/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func (w *WebHandler) GetPatient(c *gin.Context) {
	id := c.Param("id")
	patient, err := w.Repo.GetPatientById(id)
	if err != nil {
		if errors.Unwrap(err) == gorm.ErrRecordNotFound { // no rows found
			c.Status(http.StatusNotFound)
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, patient)
}

func (w *WebHandler) GetAllPatient(c *gin.Context) {
	criteriaList := []repository.Criteria{}
	limit, offset, err := utils.Paging(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if search, exist := c.GetQuery("search"); exist {
		if search != "" {
			criteriaList = append(criteriaList, repository.Criteria{QueryCriteria: repository.PATIENT_SEARCH, Value: search})
		}
	}
	patients, err := w.Repo.GetAllPatient(limit, offset, criteriaList...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, patients)
}

// deprecated
// func (w *WebHandler) CreatePatient(c *gin.Context) {
// 	var input model.CreatePatientRequest
// 	if err := c.ShouldBindJSON(&input); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}
// 	insertedId, err := w.Repo.CreatePatient(model.Patient{
// 		Hn:         input.Hn,
// 		FirstName:  input.FirstName,
// 		MiddleName: input.MiddleName,
// 		LastName:   input.LastName,
// 		Email:      input.Email,
// 		Phone:      input.Phone,
// 		Verified:   input.Verified,
// 		Weight:     input.Weight,
// 		Height:     input.Height,
// 		BirthDate:  input.BirthDate,
// 	})
// 	if err != nil {
// 		if errors.Unwrap(err) == repository.ErrDuplicateEntry {
// 			c.JSON(http.StatusConflict, gin.H{"error": "duplicate hn"})
// 			return
// 		}
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}
// 	c.JSON(http.StatusCreated, gin.H{"id": insertedId})
// }

func (w *WebHandler) UpdatePatient(c *gin.Context) {
	var input model.UpdatePatientRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	i := c.Param("id")
	id, err := strconv.Atoi(i)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	_, err = w.Repo.GetPatientById(id) // check if this id exist
	if err != nil {
		if errors.Unwrap(err) == gorm.ErrRecordNotFound { // no rows found
			c.Status(http.StatusNotFound)
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	err = w.Repo.UpdatePatient(model.Patient{
		ID:         id,
		NID:        input.NID,
		Hn:         input.Hn,
		FirstName:  input.FirstName,
		MiddleName: input.MiddleName,
		LastName:   input.LastName,
		Email:      input.Email,
		Phone:      input.Phone,
		Verified:   input.Verified,
		Weight:     input.Weight,
		Height:     input.Height,
		BirthDate:  input.BirthDate,
	})
	if err != nil {
		if errors.Unwrap(err) == repository.ErrDuplicateEntry {
			c.JSON(http.StatusConflict, gin.H{"error": "duplicate HN or NID"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusOK)
}

// DeleteDoctor is idempotent
func (w *WebHandler) DeletePatient(c *gin.Context) {
	i := c.Param("id")
	id, err := strconv.Atoi(i)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err = w.Repo.DeletePatientById(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

func (w *WebHandler) UpdatePatientVaccineHistory(c *gin.Context) {
	var input model.UpdateVaccineHistoryRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	i := c.Param("id")
	id, err := strconv.Atoi(i)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	_, err = w.Repo.GetPatientById(id) // check if this id exist
	if err != nil {
		if errors.Unwrap(err) == gorm.ErrRecordNotFound { // no rows found
			c.Status(http.StatusNotFound)
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	err = w.Repo.UpdatePatientVaccineHistory(id, input.Data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusOK)
}

func (w *WebHandler) UpdatePatientMedicine(c *gin.Context) {
	var input model.UpdateMedicineRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	i := c.Param("id")
	id, err := strconv.Atoi(i)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	_, err = w.Repo.GetPatientById(id) // check if this id exist
	if err != nil {
		if errors.Unwrap(err) == gorm.ErrRecordNotFound { // no rows found
			c.Status(http.StatusNotFound)
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	err = w.Repo.UpdatePatientMedicine(id, input.Data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusOK)
}
