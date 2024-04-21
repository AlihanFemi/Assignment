package handlers

import (
	"database/sql"
	"net/http"

	"github.com/kataras/go-sessions/v3"
)

type UserHandler struct {
	DB   *sql.DB
	Sess *sessions.Sessions // need to import it
}

type User struct {
	ID       int
	Name     string
	Email    string
	Password string
}

func (uh *UserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

}
