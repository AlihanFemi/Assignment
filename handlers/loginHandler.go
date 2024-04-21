package handlers

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/kataras/go-sessions/v3"
	"golang.org/x/crypto/bcrypt"
)

type LoginHandler struct {
	Sess *sessions.Sessions
	DB   *sql.DB
}

func (lh *LoginHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		lh.ServerGet(w, r)
	} else if r.Method == http.MethodPost {
		lh.Login(w, r)
	} else {
		http.Error(w, "Unsupported Method", http.StatusMethodNotAllowed)
	}
}

func (lh *LoginHandler) ServerGet(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "templates/login.html")
}

func (lh *LoginHandler) Login(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	username := r.Form.Get("username")
	password := r.Form.Get("password")
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Can't Register!", http.StatusBadRequest)
		return
	}
	user, err := lh.authenticateFromDB(username, string(hashedPassword))
	if err != nil {
		http.Error(w, "Can't authenticate", http.StatusUnauthorized)
		fmt.Println(err)
		return
	}

	sess := lh.Sess.Start(w, r)
	sess.Set("User", user)
	sess.Set("is_authenticated", true)
	http.Redirect(w, r, "index", http.StatusFound)
}

func (lh *LoginHandler) authenticateFromDB(name string, password string) (User, error) {
	const query = "Select * from users WHERE name = $1"
	rows, err := lh.DB.Query(query, name)
	if err != nil {
		return User{}, err
	}
	var user User
	for rows.Next() {
		err = rows.Scan(&user.ID, &user.Name, &user.Email, &user.Password)
		if err != nil {
			return User{}, nil
		}
		if user.Password == password {
			break
		}
	}
	return user, nil
}
