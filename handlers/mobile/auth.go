package mobile

import (
	"errors"
	"net/http"
	"time"

	"github.com/PhasitWo/duchenne-server/auth"
	"github.com/PhasitWo/duchenne-server/config"
	"github.com/PhasitWo/duchenne-server/repository"
	"github.com/golang-jwt/jwt/v4"
	"gorm.io/gorm"

	"github.com/PhasitWo/duchenne-server/model"

	"github.com/gin-gonic/gin"
)

type login struct {
	Hn         string `json:"hn" binding:"required"`
	FirstName  string `json:"firstName" binding:"required"`
	LastName   string `json:"lastName" binding:"required"`
	DeviceName string `json:"deviceName" binding:"required"`
	ExpoToken  string `json:"expoToken" binding:"required"`
}

func (m *MobileHandler) Login(c *gin.Context) {
	var input login
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// fetch patient from database
	storedPatient, err := m.Repo.GetPatientByHN(input.Hn)
	if err != nil {
		if errors.Unwrap(err) == gorm.ErrRecordNotFound { // no rows found
			c.Status(http.StatusNotFound)
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// checking
	if !storedPatient.Verified {
		c.JSON(http.StatusForbidden, gin.H{"error": "unverified account"})
		return
	}
	if input.FirstName != storedPatient.FirstName || input.LastName != storedPatient.LastName {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credential"})
		return
	}
	// save this device for notification stuff
	criteria := repository.Criteria{QueryCriteria: repository.PATIENTID, Value: storedPatient.ID}
	devices, err := m.Repo.GetAllDevice(criteria)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	newDevice := model.Device{
		LoginAt:    int(time.Now().Unix()),
		DeviceName: input.DeviceName,
		ExpoToken:  input.ExpoToken,
		PatientId:  storedPatient.ID,
	}
	tx := m.DBConn.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if err := tx.Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	repoWithTx := repository.New(tx)
	if len(devices) >= config.AppConfig.MAX_DEVICE {
		// remove the oldest login device
		toRemoveDeviceId := devices[0].ID
		err = repoWithTx.DeleteDevice(toRemoveDeviceId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}
	// insert new device
	deviceId, err := repoWithTx.CreateDevice(newDevice)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// generate token
	token, err := auth.GeneratePatientToken(storedPatient.ID, deviceId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// commit tx
	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Tx can't commit"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": token})
}

type signup struct {
	Hn         string `json:"hn" binding:"required"`
	FirstName  string `json:"firstName" binding:"required"`
	MiddleName string `json:"middleName"`
	LastName   string `json:"lastName" binding:"required"`
	Phone      string `json:"phone" binding:"required"`
	Email      string `json:"email" binding:"required"`
}

func (m *MobileHandler) Signup(c *gin.Context) {
	var s signup
	if err := c.ShouldBindJSON(&s); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// fetch patient from database
	storedPatient, err := m.Repo.GetPatientByHN(s.Hn)
	if err != nil {
		if errors.Unwrap(err) == gorm.ErrRecordNotFound { // no rows found
			c.Status(http.StatusNotFound)
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// checking
	if storedPatient.Verified { // already verified
		c.JSON(http.StatusConflict, gin.H{"error": "the account has been verified"})
		return
	}
	if s.FirstName != storedPatient.FirstName || s.LastName != storedPatient.LastName {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid first name or last name"})
		return
	}
	if storedPatient.MiddleName != nil && s.MiddleName != *storedPatient.MiddleName {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid middle name"})
		return
	}
	if storedPatient.MiddleName == nil && s.MiddleName != "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "does not require middle name"})
		return
	}
	// update patient info and mark patient as verified
	err = m.Repo.UpdatePatient(
		model.Patient{
			ID:         storedPatient.ID,
			Hn:         storedPatient.Hn,
			FirstName:  storedPatient.FirstName,
			MiddleName: storedPatient.MiddleName,
			LastName:   storedPatient.LastName,
			Email:      &s.Email,
			Phone:      &s.Phone,
			Verified:   true,
		})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusOK)
}

func (m *MobileHandler) Logout(c *gin.Context) {
	tokenString := c.GetHeader("Authorization")
	if tokenString == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No authorization header provided"})
		return
	}
	// parse token
	claims := &auth.PatientClaims{PatientId: -1, DeviceId: -1}
	jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.AppConfig.JWT_KEY), nil
	})
	if claims.DeviceId == -1 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token, no deviceId"})
		return
	}

	err := m.Repo.DeleteDevice(claims.DeviceId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusOK)
}
