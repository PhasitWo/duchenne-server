package web

import (
	"database/sql"
	"errors"
	"net/http"
	"strconv"

	"github.com/PhasitWo/duchenne-server/repository"
	"github.com/gin-gonic/gin"
)

func (w *WebHandler) GetAllQuestion(c *gin.Context) {
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
			c.JSON(http.StatusBadRequest, gin.H{"error": "cannot parse offset value"})
			return
		}
		criteriaList = append(criteriaList, repository.Criteria{QueryCriteria: repository.DOCTORID, Value: doctorId})
	}
	if p, exist := c.GetQuery("patientId"); exist {
		patientId, err := strconv.Atoi(p)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "cannot parse offset value"})
			return
		}
		criteriaList = append(criteriaList, repository.Criteria{QueryCriteria: repository.PATIENTID, Value: patientId})
	}
	if t, exist := c.GetQuery("type"); exist {
		if t == "replied" {
			criteriaList = append(criteriaList, repository.Criteria{QueryCriteria: repository.ANSWERAT_ISNOTNULL})
		} else if t == "unreplied" {
			criteriaList = append(criteriaList, repository.Criteria{QueryCriteria: repository.ANSWERAT_ISNULL})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid type value"})
			return
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
		if errors.Unwrap(err) == sql.ErrNoRows { // no rows found
			c.Status(http.StatusNotFound)
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, q)
}

type questionAnswer struct {
	Answer string `json:"answer" binding:"required,max=500"`
}

func (w *WebHandler) AnswerQuestion(c *gin.Context) {
	id := c.Param("id")
	q, err := w.Repo.GetQuestion(id)
	if err != nil {
		if errors.Unwrap(err) == sql.ErrNoRows { // no rows found
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
	var input questionAnswer
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	doctorId, exists := c.Get("doctorId")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "no 'doctorId' from auth middleware"})
		return
	}
	// query
	err = w.Repo.UpdateQuestionAnswer(id, input.Answer, doctorId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusOK)
}
