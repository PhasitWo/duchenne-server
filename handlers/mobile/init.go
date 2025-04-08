package mobile

import (
	// "database/sql"

	"github.com/PhasitWo/duchenne-server/repository"
	_ "github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

type MobileHandler struct {
	Repo   repository.IRepo
	DBConn repository.IGorm
}

func Init(db *gorm.DB) *MobileHandler {
	return &MobileHandler{Repo: repository.New(db), DBConn: db}
}
