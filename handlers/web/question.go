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

func (w *WebHandler) GetAllQuestion(c *gin.Context) {
	criteriaList := []repository.Criteria{}
	var err error
	// get url query param
	limit, offset, err := utils.Paging(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
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
		switch t {
		case "replied":
			criteriaList = append(criteriaList, repository.Criteria{QueryCriteria: repository.ANSWERAT_ISNOTNULL})
		case "unreplied":
			criteriaList = append(criteriaList, repository.Criteria{QueryCriteria: repository.ANSWERAT_ISNULL})
		default:
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid type value"})
			return
		}
	}
	if search, exist := c.GetQuery("search"); exist {
		if search != "" {
			criteriaList = append(criteriaList, repository.Criteria{QueryCriteria: repository.QUESTION_SEARCH, Value: search})
		}
	}
	// query
	qs, err := w.Repo.GetAllQuestion(limit, offset, criteriaList...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, qs)
}

func (w *WebHandler) GetQuestion(c *gin.Context) {
	id := c.Param("id")
	q, err := w.Repo.GetQuestion(id)
	if err != nil {
		if errors.Unwrap(err) == gorm.ErrRecordNotFound { // no rows found
			c.Status(http.StatusNotFound)
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, q)
}

func (w *WebHandler) AnswerQuestion(c *gin.Context) {
	id := c.Param("id")
	questionId, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	q, err := w.Repo.GetQuestion(id)
	if err != nil {
		if errors.Unwrap(err) == gorm.ErrRecordNotFound { // no rows found
			c.Status(http.StatusNotFound)
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// check question status
	if q.AnswerAt != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "this question has been replied"})
		return
	}
	// input
	var input model.QuestionAnswerRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	dId, exists := c.Get("doctorId")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "no 'doctorId' from auth middleware"})
		return
	}
	doctorId := dId.(int)
	// query
	err = w.Repo.UpdateQuestionAnswer(questionId, input.Answer, doctorId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	go w.NotiService.SendNotiByPatientId(q.PatientID, "แพทย์ตอบคำถามของคุณแล้ว!", "ดูคำตอบในแอปพลิเคชัน")
	c.Status(http.StatusOK)
}
