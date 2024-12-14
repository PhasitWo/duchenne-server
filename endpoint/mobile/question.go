package mobile

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	// "github.com/PhasitWo/duchenne-server/repository"
	"net/http"

	"github.com/PhasitWo/duchenne-server/repository"
	"github.com/gin-gonic/gin"
)

func (m *MobileHandler) GetAllPatientQuestion(c *gin.Context) {
	i, exists := c.Get("patientId")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "no 'patientId' from auth middleware"})
		return
	}
	id := i.(int)
	// _, queryExists := c.GetQuery("onlytopic")
	qs, err := m.Repo.GetAllQuestion(id, repository.PATIENTID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, qs)
}

func (m *MobileHandler) GetQuestion(c *gin.Context) {
	i, exists := c.Get("patientId")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "no 'patientId' from auth middleware"})
		return
	}
	patientId := i.(int)
	id := c.Param("id")
	q, err := m.Repo.GetQuestion(id)
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
}

const MAX_TOPIC_LENGTH = 50
const MAX_QUESTION_LENGTH = 700

func (m *MobileHandler) CreateQuestion(c *gin.Context) {
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
	// validate
	if len(input.Topic) > MAX_TOPIC_LENGTH {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": fmt.Sprintf("topic input is exceeding %d maximum of characters", MAX_TOPIC_LENGTH)})
		return
	}
	if len(input.Question) > MAX_QUESTION_LENGTH {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": fmt.Sprintf("question input is exceeding %d maximum of characters", MAX_QUESTION_LENGTH)})
		return
	}
	insertedId, err := m.Repo.CreateQuestion(patientId, input.Topic, input.Question, int(time.Now().Unix()))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"id": insertedId})
}

func (m *MobileHandler) DeleteQuestion(c *gin.Context) {
	// prepare param from url and auth middleware
	id := c.Param("id")
	i, exists := c.Get("patientId")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "no 'patientId' from auth middleware"})
		return
	}
	patientId := i.(int)
	// check if this question belongs to the patient
	q, err := m.Repo.GetQuestion(id)
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
	err = m.Repo.DeleteQuestion(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
