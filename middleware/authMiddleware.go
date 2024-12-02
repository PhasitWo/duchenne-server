package middleware

import (
	"net/http"

	"github.com/PhasitWo/duchenne-server/auth"
	"github.com/PhasitWo/duchenne-server/config"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

func AuthMiddleware(c *gin.Context) {
	tokenString := c.GetHeader("Authorization")
	if tokenString == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No authorization header provided"})
		c.Abort() // stop the chain of middlewares
		return
	}
	// parse token
	claims := &auth.Claims{UserId: -1}
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
	c.Set("userId", claims.UserId)
	c.Next()
}
