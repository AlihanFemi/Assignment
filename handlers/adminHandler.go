package handlers

import (
	"database/sql"

	"github.com/kataras/go-sessions/v3"
)

type AdminHandler struct {
	DB   *sql.DB
	Sess *sessions.Sessions
}
