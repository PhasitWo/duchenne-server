package web

import (
	// "database/sql"

	"github.com/PhasitWo/duchenne-server/services/notification"
	"github.com/PhasitWo/duchenne-server/repository"
	_ "github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

type WebHandler struct {
	Repo   repository.IRepo
	DBConn *gorm.DB
	NotiService notification.INotificationService
}

func Init(db *gorm.DB) *WebHandler {
	return &WebHandler{
		Repo: repository.New(db),
		DBConn: db,
		NotiService: notification.NewService(db),
	}
}
