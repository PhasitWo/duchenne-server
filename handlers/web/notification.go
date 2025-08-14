package web

import (
	"net/http"
	"strconv"

	"github.com/PhasitWo/duchenne-server/config"
	"github.com/gin-gonic/gin"
)

func (w *WebHandler) SendDailyNotifications(c *gin.Context) {
	// get url query param
	secret, exist := c.GetQuery("secret")
	if !exist {
		c.JSON(http.StatusBadRequest, gin.H{"error": "require secret"})
		return
	}
	if secret != config.AppConfig.NOTIFY_SECRET {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid secret", "received": secret})
		return
	}
	var dayRange *int
	if l, exist := c.GetQuery("day_range"); exist {
		converted, err := strconv.Atoi(l)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "cannot parse day_range value"})
			return
		}
		dayRange = &converted
	}
	err := w.NotiService.SendDailyNotifications(dayRange)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	c.Status(200)
}
