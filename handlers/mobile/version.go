package mobile

import (
	"net/http"

	"github.com/PhasitWo/duchenne-server/config"
	"github.com/gin-gonic/gin"
)

func (m *MobileHandler) GetRequireMobileVersion(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"version": config.AppConfig.REQUIRE_MOBILE_VERSION})
}
