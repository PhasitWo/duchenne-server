package common

import (
	"github.com/PhasitWo/duchenne-server/notification"
	"github.com/PhasitWo/duchenne-server/repository"
	"gorm.io/gorm"
)

type CommonHandler struct {
	Repo        repository.IRepo
	DBConn      *gorm.DB
	NotiService notification.INotificationService
}

func Init(db *gorm.DB) *CommonHandler {
	return &CommonHandler{
		Repo:        repository.New(db),
		DBConn:      db,
		NotiService: notification.NewService(db),
	}
}
