package handlers

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/kataras/go-sessions/v3"
	"golang.org/x/crypto/bcrypt"
)

type RegisterHandler struct {
	Sess *sessions.Sessions
	DB   *sql.DB
}

func (rh *RegisterHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		rh.ServerGet(w, r)
	} else if r.Method == http.MethodPost {
		rh.Register(w, r)
	} else {
		http.Error(w, "Unsupported Method", http.StatusMethodNotAllowed)
	}
}

func (rh *RegisterHandler) ServerGet(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "templates/register.html")
}

func (rh *RegisterHandler) Register(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	fmt.Println(r.Form)
	username := r.Form.Get("username")
	email := r.Form.Get("email")
	password := r.Form.Get("password")
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Can't Register!", http.StatusBadRequest)
		return
	}
	err = rh.insertIntoDB(username, email, string(hashedPassword))
	if err != nil {
		http.Error(w, "Already exists!", http.StatusBadRequest)
		fmt.Println(err)
		return
	}

	http.Redirect(w, r, "login", http.StatusFound)
}

func (rh *RegisterHandler) insertIntoDB(name string, email string, password string) error {
	const query = "INSERT INTO users (name, email, password) VALUES ($1, $2, $3)"
	_, err := rh.DB.Exec(query, name, email, password)
	if err != nil {
		return err
	}
	return nil
}
