package mobile

import (
	"database/sql"
	"errors"
	// "github.com/PhasitWo/duchenne-server/repository"
	"net/http"

	"github.com/PhasitWo/duchenne-server/repository"
	"github.com/gin-gonic/gin"
)

func (m *mobileHandler) GetAllPatientQuestion(c *gin.Context) {
	i, exists := c.Get("patientId")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "no 'patientId' from auth middleware"})
		return
	}
	id := i.(int)
	_, queryExists := c.GetQuery("onlytopic")
	var qs any
	var err error
	if queryExists {
		qs, err = m.repo.GetAllQuestionTopic(id, repository.PATIENTID)
	} else {
		qs, err = m.repo.GetAllQuestion(id, repository.PATIENTID)
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, qs)
}

func (m *mobileHandler) GetQuestion(c *gin.Context) {
	i, exists := c.Get("patientId")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "no 'patientId' from auth middleware"})
		return
	}
	patientId := i.(int)
	id := c.Param("id")
	q, err := m.repo.GetQuestion(id)
	if err != nil {
		if errors.Unwrap(err) == sql.ErrNoRows { // no rows found
			c.Status(http.StatusNotFound)
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// check if this question belongs to the patient
	if patientId != q.Patient.Id {
		c.Status(http.StatusUnauthorized)
		return
	}
	c.JSON(http.StatusOK, q)
}

type questionInput struct {
	Topic    string `json:"topic" binding:"required"`
	Question string `json:"question" binding:"required"`
	CreateAt int    `json:"createAt" binding:"required"`
}

func (m *mobileHandler) CreateQuestion(c *gin.Context) {
	// get patientId from auth header
	i, exists := c.Get("patientId")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "no 'patientId' from auth middleware"})
		return
	}
	patientId := i.(int)
	// binding request body
	var input questionInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	_, err := m.repo.CreateQuestion(patientId, input.Topic, input.Question, input.CreateAt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusCreated)
}

func (m *mobileHandler) DeleteQuestion(c *gin.Context) {
	// prepare param from url and auth middleware
	id := c.Param("id")
	i, exists := c.Get("patientId")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "no 'patientId' from auth middleware"})
		return
	}
	patientId := i.(int)
	// check if this question belongs to the patient
	q, err := m.repo.GetQuestion(id)
	if err != nil {
		if errors.Unwrap(err) == sql.ErrNoRows { // no row found
			c.Status(http.StatusNotFound)
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if patientId != q.Patient.Id {
		c.Status(http.StatusUnauthorized)
		return
	}
	err = m.repo.DeleteQuestion(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
