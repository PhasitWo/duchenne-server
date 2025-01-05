package web

import (
	"net/http"
	"strconv"
	"time"

	"github.com/PhasitWo/duchenne-server/repository"
	"github.com/gin-gonic/gin"
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
