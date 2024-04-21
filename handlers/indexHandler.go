package handlers

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"

	"github.com/kataras/go-sessions/v3"
)

type IndexHandler struct {
	DB   *sql.DB
	Sess *sessions.Sessions
}

type IndexData struct {
	Author        User
	Post          Post
	Authenticated bool
}

type Data struct {
	IndexData     []IndexData
	Authenticated bool
}

func (ih *IndexHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		ih.ServerGet(w, r)
	} else {
		http.Error(w, "Unsupported Method", http.StatusMethodNotAllowed)
	}
}

func (ih *IndexHandler) ServerGet(w http.ResponseWriter, r *http.Request) {
	posts, err := ih.getPosts()
	if err != nil {
		fmt.Println(err)
		fmt.Fprintln(w, "internal server error")
		return
	}
	sess := ih.Sess.Start(w, r)
	auth, err := sess.GetBoolean("is_authenticated")
	if err != nil {
		auth = false
	}
	data := Data{IndexData: posts, Authenticated: auth}
	tpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		fmt.Println(err)
		fmt.Fprintln(w, "There was an error")
		return
	}
	err = tpl.Execute(w, data)
	if err != nil {
		fmt.Println(err)
		fmt.Fprintln(w, "There was an error")
		return
	}
}

func (ih *IndexHandler) getPosts() ([]IndexData, error) {
	rows, err := ih.DB.Query("SELECT * from posts")
	if err != nil {
		return nil, err
	}
	var posts []IndexData
	for rows.Next() {
		var post Post
		err = rows.Scan(&post.ID, &post.Title, &post.Content, &post.Date, &post.User_ID)
		if err != nil {
			return nil, err
		}
		var author User
		author, err = ih.getUser(post.User_ID)
		if err != nil {
			return nil, err
		}
		pst := IndexData{Author: author, Post: post}
		posts = append(posts, pst)
	}
	return posts, nil
}

func (ih *IndexHandler) getUser(id int) (User, error) {
	query := "SELECT * from users WHERE id = $1"
	rows, err := ih.DB.Query(query, id)
	if err != nil {
		return User{}, err
	}
	var user User
	for rows.Next() {
		err = rows.Scan(&user.ID, &user.Name, &user.Email, &user.Password)
		if err != nil {
			return User{}, err
		}
	}
	return user, nil
}
