package mobile

import (
	"database/sql"
	"errors"
	"net/http"
	"time"

	"github.com/PhasitWo/duchenne-server/auth"
	"github.com/PhasitWo/duchenne-server/config"
	"github.com/PhasitWo/duchenne-server/repository"
	"github.com/golang-jwt/jwt/v4"

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

func (m *mobileHandler) Login(c *gin.Context) {
	var input login
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// fetch patient from database
	storedPatient, err := m.repo.GetPatient(input.Hn)
	if err != nil {
		if errors.Unwrap(err) == sql.ErrNoRows { // no rows found
			c.Status(http.StatusNotFound)
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// checking
	if !storedPatient.Verified {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unverified account"})
		return
	}
	if input.FirstName != storedPatient.FirstName || input.LastName != storedPatient.LastName {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credential"})
		return
	}
	// save this device for notification stuff
	devices, err := m.repo.GetAllDevice(storedPatient.Id, repository.PATIENTID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	newDevice := model.Device{
		Id:         -1,
		LoginAt:    int(time.Now().Unix()),
		DeviceName: input.DeviceName,
		ExpoToken:  input.ExpoToken,
		PatientId:  storedPatient.Id,
	}
	tx, err := m.dbConn.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer func() {
		if err := tx.Rollback(); err != nil && !errors.Is(err, sql.ErrTxDone) {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "tx can't rollback"})
		}
	}()
	repoWithTx := repository.New(tx)
	if len(devices) >= config.AppConfig.MAX_DEVICE {
		// remove the oldest login device
		toRemoveDeviceId := devices[0].Id
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
	token, err := auth.GeneratePatientToken(storedPatient.Id, deviceId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// commit tx
	if err := tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Tx can't commit"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": token})
}

type signup struct {
	Hn         string `json:"hn" binding:"required"`
	FirstName  string `json:"firstName" binding:"required"`
	MiddleName string `json:"middleName" binding:"required"`
	LastName   string `json:"lastName" binding:"required"`
	Phone      string `json:"phone" binding:"required"`
	Email      string `json:"email" binding:"required"`
}

func (m *mobileHandler) Signup(c *gin.Context) {
	var s signup
	if err := c.ShouldBindJSON(&s); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// fetch patient from database
	storedPatient, err := m.repo.GetPatient(s.Hn)
	if err != nil {
		if errors.Unwrap(err) == sql.ErrNoRows { // no rows found
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
	// update patient info and mark patient as verified
	err = m.repo.UpdatePatient(
		model.Patient{
			Id:         storedPatient.Id,
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

func (m *mobileHandler) Logout(c *gin.Context) {
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

	err := m.repo.DeleteDevice(claims.DeviceId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusOK)
}
