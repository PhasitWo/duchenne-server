package mobile

import (
	"net/http"

	"github.com/PhasitWo/duchenne-server/repository"
	"github.com/PhasitWo/duchenne-server/utils"
	"github.com/gin-gonic/gin"
)

func (m *MobileHandler) GetAllDoctor(c *gin.Context) {
	criteriaList := []repository.Criteria{}
	var err error
	// get url query param
	limit, offset, err := utils.Paging(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if canBeAppointedStr, exist := c.GetQuery("canBeAppointed"); exist {
		var canBeAppointed bool
		switch canBeAppointedStr {
		case "true":
			canBeAppointed = true
		case "false":
			canBeAppointed = false
		default:
			c.JSON(http.StatusBadRequest, gin.H{"error": "canBeAppointed can either be 'true' or 'false'"})
			return
		}
		criteriaList = append(criteriaList, repository.Criteria{QueryCriteria: repository.CAN_BE_APPOINTED, Value: canBeAppointed})
	}
	doctors, err := m.Repo.GetAllDoctor(limit, offset, criteriaList...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, doctors)
}
