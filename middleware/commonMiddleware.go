package middleware

import (
	"net/http"

	"github.com/PhasitWo/duchenne-server/config"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

func CommonAuthMiddleware(c *gin.Context) {
	var err error
	var tokenString string
	if tokenString = c.GetHeader("Authorization"); tokenString == "" {
		tokenString, err = c.Cookie(config.Constants.WEB_ACCESS_COOKIE_NAME)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "cannot get cookie from request"})
			c.Abort()
			return
		}
	}
	if tokenString == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "no token from header or cookie"})
		c.Abort()
	}
	// parse token
	claims := &jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.AppConfig.JWT_KEY), nil
	})
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		c.Abort()
		return
	}
	if !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		c.Abort()
		return
	}
	c.Next()
}
