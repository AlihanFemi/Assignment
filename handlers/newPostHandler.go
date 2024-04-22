package handlers

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/kataras/go-sessions/v3"
)

type NewPostHandler struct {
	DB   *sql.DB
	Sess *sessions.Sessions
}

func (nph *NewPostHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		http.ServeFile(w, r, "templates/post-in.html")
	} else if r.Method == http.MethodPost {
		nph.newPost(w, r)
	} else {
		http.Error(w, "Unsupported Method", http.StatusMethodNotAllowed)
	}
}

func (nph *NewPostHandler) newPost(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	title := r.Form.Get("title")
	content := r.Form.Get("content")
	data := time.Now().Format("2006-01-02")
	sess := nph.Sess.Start(w, r)
	user := sess.Get("User").(User)
	err := nph.insertPost(title, content, data, user)
	if err != nil {
		fmt.Println(err)
		fmt.Fprintln(w, "Internal server error!")
		return
	}
	http.Redirect(w, r, "index", http.StatusFound)
}

func (nph *NewPostHandler) insertPost(title string, content string, date string, user User) error {
	// insert into posts VALUES(content, date, user_ID) ($1, $2, $2)
	const query = "INSERT INTO posts (title, content, date, user_ID) VALUES ($1, $2, $3, $4)"
	id := user.ID
	_, err := nph.DB.Exec(query, title, content, date, id)
	if err != nil {
		return err
	}
	return nil
}
