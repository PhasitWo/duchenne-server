package mobile

import (
	"net/http"

	"github.com/PhasitWo/duchenne-server/config"
	"github.com/gin-gonic/gin"
)

func (m *MobileHandler) GetRequireMobileVersion(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"version":          config.AppConfig.REQUIRE_MOBILE_VERSION,
		"androidStoreLink": config.AppConfig.ANDROID_STORE_LINK,
		"iosStoreLink":     config.AppConfig.IOS_STORE_LINK,
	})
}
