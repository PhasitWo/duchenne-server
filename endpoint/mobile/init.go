package mobile

import (
	"database/sql"
	"github.com/PhasitWo/duchenne-server/repository"
	_ "github.com/go-sql-driver/mysql"
)

type mobileHandler struct {
	repo *repository.Repo
	dbConn *sql.DB
}

func Init(db *sql.DB) *mobileHandler {
	return &mobileHandler{repo: repository.New(db), dbConn : db}
}
