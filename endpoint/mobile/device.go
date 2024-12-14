package mobile

import (
	"database/sql"
	"errors"
	"net/http"
	"time"

	"github.com/PhasitWo/duchenne-server/auth"
	"github.com/PhasitWo/duchenne-server/config"
	"github.com/PhasitWo/duchenne-server/model"
	"github.com/PhasitWo/duchenne-server/repository"
	"github.com/gin-gonic/gin"
)

func (m *MobileHandler) GetAllDevice(c *gin.Context) {
	i, exists := c.Get("patientId")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "no 'patientId' from auth middleware"})
		return
	}
	id := i.(int)
	dv, err := m.Repo.GetAllDevice(id, repository.PATIENTID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, dv)
}

type device struct {
	DeviceName string `json:"deviceName" binding:"required"`
	ExpoToken  string `json:"expoToken" binding:"required"`
}

func (m *MobileHandler) CreateDevice(c *gin.Context) {
	i, exists := c.Get("patientId")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "no 'patientId' from auth middleware"})
		return
	}
	id := i.(int)
	// binding input
	var dv device
	if err := c.ShouldBindJSON(&dv); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// save this device for notification stuff
	newDevice := model.Device{
		Id:         -1,
		LoginAt:    int(time.Now().Unix()),
		DeviceName: dv.DeviceName,
		ExpoToken:  dv.ExpoToken,
		PatientId:  id,
	}
	devices, err := m.Repo.GetAllDevice(id, repository.PATIENTID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	tx, err := m.DBConn.Begin()
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
	token, err := auth.GeneratePatientToken(id, deviceId)
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
