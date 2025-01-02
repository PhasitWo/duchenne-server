package web

import (
	"database/sql"

	"github.com/PhasitWo/duchenne-server/repository"
	_ "github.com/go-sql-driver/mysql"
)

type WebHandler struct {
	Repo   *repository.Repo
	DBConn *sql.DB
}

func Init(db *sql.DB) *WebHandler {
	return &WebHandler{Repo: repository.New(db), DBConn: db}
}
