package mobile

import (
	"database/sql"

	"github.com/PhasitWo/duchenne-server/repository"
	_ "github.com/go-sql-driver/mysql"
)

type MobileHandler struct {
	Repo   *repository.Repo
	DBConn *sql.DB
}

func Init(db *sql.DB) *MobileHandler {
	return &MobileHandler{Repo: repository.New(db), DBConn: db}
}
