package middleware

import (
	"net/http"

	"github.com/PhasitWo/duchenne-server/auth"
	"github.com/PhasitWo/duchenne-server/config"
	"github.com/PhasitWo/duchenne-server/model"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

func WebAuthMiddleware(c *gin.Context) {
	cookie, err := c.Cookie("web_auth_token")
	if err != nil {
		// assume that the err is http.ErrNoCookie or expired
		c.JSON(http.StatusUnauthorized, gin.H{"error": "cannot get cookie from request"})
		c.Abort()
		return
	}

	// parse token
	claims := &auth.DoctorClaims{DoctorId: -1}
	token, err := jwt.ParseWithClaims(cookie, claims, func(token *jwt.Token) (interface{}, error) {
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

	if claims.DoctorId == -1 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		c.Abort()
		return
	}
	c.Set("doctorId", claims.DoctorId)
	c.Set("doctorRole", claims.Role)
	c.Next()
}

type permission string

const (
	CreateDoctorPermission  permission = "createDoctorPermission"
	UpdateDoctorPermission  permission = "updateDoctorPermission"
	DeleteDoctorPermission  permission = "deleteDoctorPermission"
	CreatePatientPermission permission = "createPatientPermission"
	UpdatePatientPermission permission = "updatePatientPermission"
	DeletePatientPermission permission = "deletePatientPermission"
)

var rolePermissionsMap = map[model.Role][]permission{
	model.USER:  {},
	model.ADMIN: {CreatePatientPermission, UpdatePatientPermission, DeletePatientPermission},
	model.ROOT:  {CreatePatientPermission, UpdatePatientPermission, DeletePatientPermission, CreateDoctorPermission, UpdateDoctorPermission, DeleteDoctorPermission},
}

func WebRBACMiddleware(requiredPermission permission) gin.HandlerFunc {
	return func(c *gin.Context) {
		r, exists := c.Get("doctorRole")
		if !exists {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "no 'doctorRole' from auth middleware"})
			c.Abort()
			return
		}
		role := r.(model.Role)
		permissions, exists := rolePermissionsMap[role]
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{"error": "Role not found"})
			c.Abort()
			return
		}
		// check if this role has requiredPermission
		hasPermission := false
		for _, permission := range permissions {
			if permission == requiredPermission {
				hasPermission = true
				break
			}
		}
		if !hasPermission {
			c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
			c.Abort()
			return
		}
		c.Next()
	}
}
