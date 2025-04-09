package web

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/PhasitWo/duchenne-server/model"
	"github.com/PhasitWo/duchenne-server/repository"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func (w *WebHandler) GetAllAppointment(c *gin.Context) {
	criteriaList := []repository.Criteria{}
	limit := 9999
	offset := 0
	var err error
	// get url query param
	if l, exist := c.GetQuery("limit"); exist {
		limit, err = strconv.Atoi(l)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "cannot parse limit value"})
			return
		}
	}
	if of, exist := c.GetQuery("offset"); exist {
		offset, err = strconv.Atoi(of)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "cannot parse offset value"})
			return
		}
	}
	if d, exist := c.GetQuery("doctorId"); exist {
		doctorId, err := strconv.Atoi(d)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "cannot parse doctorId value"})
			return
		}
		criteriaList = append(criteriaList, repository.Criteria{QueryCriteria: repository.DOCTORID, Value: doctorId})
	}
	if p, exist := c.GetQuery("patientId"); exist {
		patientId, err := strconv.Atoi(p)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "cannot parse patientId value"})
			return
		}
		criteriaList = append(criteriaList, repository.Criteria{QueryCriteria: repository.PATIENTID, Value: patientId})
	}
	if t, exist := c.GetQuery("type"); exist {
		if t == "incoming" {
			criteriaList = append(criteriaList, repository.Criteria{QueryCriteria: repository.DATE_GREATERTHAN, Value: int(time.Now().Unix())})
		} else if t == "history" {
			criteriaList = append(criteriaList, repository.Criteria{QueryCriteria: repository.DATE_LESSTHAN, Value: int(time.Now().Unix())})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid type value"})
			return
		}
	}
	// query
	aps, err := w.Repo.GetAllAppointment(limit, offset, criteriaList...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, aps)
}

func (w *WebHandler) GetAppointment(c *gin.Context) {
	id := c.Param("id")
	apm, err := w.Repo.GetAppointment(id)
	if err != nil {
		if errors.Unwrap(err) == gorm.ErrRecordNotFound { // no rows found
			c.Status(http.StatusNotFound)
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, apm)
}


func (w *WebHandler) CreateAppointment(c *gin.Context) {
	// binding request body
	var input model.CreateAppointmentRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// validate input.date
	now := int(time.Now().Unix())
	if input.Date < now {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "'Date' is before current time"})
		return
	}
	// create new appointment
	var approveAt *int = nil
	if input.Approve {
		approveAt = &now
	}
	insertedId, err := w.Repo.CreateAppointment(model.Appointment{
		Date:      input.Date,
		PatientID: input.PatientId,
		DoctorID:  input.DoctorId,
		ApproveAt: approveAt,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	go w.NotiService.SendNotiByPatientId(input.PatientId, "คุณมีนัดหมายใหม่!", "ดูข้อมูลในแอปพลิเคชัน")
	c.JSON(http.StatusCreated, gin.H{"id": insertedId})
}

func (w *WebHandler) UpdateAppointment(c *gin.Context) {
	var input model.CreateAppointmentRequest
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
	// validate input.date
	now := int(time.Now().Unix())
	if input.Date < now {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "'Date' is before current time"})
		return
	}
	// update
	var approveAt *int = nil
	if input.Approve {
		approveAt = &now
	}
	err = w.Repo.UpdateAppointment(model.Appointment{
		ID:        id,
		Date:      input.Date,
		PatientID: input.PatientId,
		DoctorID:  input.DoctorId,
		ApproveAt: approveAt,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	go w.NotiService.SendNotiByPatientId(input.PatientId, "นัดหมายของคุณมีการเปลี่ยนแปลง!", "เช็คสถานะในแอปพลิเคชัน")
	c.Status(http.StatusOK)
}

func (w *WebHandler) DeleteAppointment(c *gin.Context) {
	id := c.Param("id")
	apm, err := w.Repo.GetAppointment(id)
	if err != nil {
		if errors.Unwrap(err) == gorm.ErrRecordNotFound { // no rows found
			c.Status(http.StatusNotFound)
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// delete appointment
	err = w.Repo.DeleteAppointment(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	go w.NotiService.SendNotiByPatientId(apm.PatientID, "นัดหมายของคุณถูกลบ!", "คุณหมอลบนัดหมายของคุณ")
	c.Status(http.StatusNoContent)
}
