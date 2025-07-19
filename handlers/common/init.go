package common

import (
	"cloud.google.com/go/storage"
	cloudstorage "github.com/PhasitWo/duchenne-server/cloud-storage"
	"github.com/PhasitWo/duchenne-server/notification"
	"github.com/PhasitWo/duchenne-server/repository"
	"gorm.io/gorm"
)

type CommonHandler struct {
	Repo                repository.IRepo
	DBConn              *gorm.DB
	NotiService         notification.INotificationService
	CloudStorageService cloudstorage.ICloudStorageService
}

func Init(db *gorm.DB, gcsClient *storage.Client) *CommonHandler {
	return &CommonHandler{
		Repo:                repository.New(db),
		DBConn:              db,
		NotiService:         notification.NewService(db),
		CloudStorageService: cloudstorage.NewService(gcsClient),
	}
}
