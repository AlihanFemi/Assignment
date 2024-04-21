package handlers

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"

	"github.com/kataras/go-sessions/v3"
)

type PostHandler struct {
	DB   *sql.DB
	Sess *sessions.Sessions
}

type Post struct {
	ID      int
	Title   string
	Content string
	Date    string
	User_ID int
}

type Comments struct {
	Author  User
	Comment Comment
}

type PostData struct {
	Post          Post
	Author        User
	Comments      []Comments
	Authenticated bool
}

var PostIndex int = 0

func (ph *PostHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		path := strings.Split(r.URL.Path, "/")
		var err error
		PostIndex, err = strconv.Atoi(path[2])
		if err != nil {
			fmt.Println(err)
			fmt.Fprintln(w, err)
		}
		ph.ServerGet(w, r)
	} else if r.Method == http.MethodPost {
		ph.ServerPost(w, r)
	} else {
		http.Error(w, "Unsupported Method", http.StatusMethodNotAllowed)
	}
}

func (ph *PostHandler) ServerGet(w http.ResponseWriter, r *http.Request) {
	post, err := ph.getPost()
	if err != nil {
		fmt.Println(err)
		fmt.Fprintln(w, "internal server error")
		return
	}
	comments, err := ph.getComments()
	if err != nil {
		fmt.Println(err)
		fmt.Fprintln(w, "internal server error")
		return
	}
	author, err := ph.getUser(post.User_ID)
	if err != nil {
		fmt.Println(err)
		fmt.Fprint(w, "Internal server error")
	}
	sess := ph.Sess.Start(w, r)
	auth, err := sess.GetBoolean("is_authenticated")
	if err != nil {
		auth = false
	}
	data := PostData{Post: post, Comments: comments, Author: author, Authenticated: auth}
	tpl, err := template.ParseFiles("templates/post.html")
	if err != nil {
		fmt.Println(err)
		fmt.Fprintln(w, "There was an error! TPL is NIL")
		return
	}
	err = tpl.Execute(w, data)
	if err != nil {
		fmt.Println(err)
		fmt.Fprintln(w, "There was an error")
		return
	}
}

func (ph *PostHandler) ServerPost(w http.ResponseWriter, r *http.Request) {

}

// func (ph *PostHandler) insertPost(title string, content string, date string, user User) error {
// 	// insert into posts VALUES(content, date, user_ID) ($1, $2, $2)
// 	const query = "INSERT INTO posts (title, content, date, user_ID) VALUES ($1, $2, $3, $4)"
// 	id := user.ID
// 	_, err := ph.DB.Exec(query, title, content, date, id)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

func (ph *PostHandler) getPost() (Post, error) {
	query := "SELECT * from posts WHERE id = $1"
	rows, err := ph.DB.Query(query, PostIndex)
	if err != nil {
		return Post{}, err
	}
	var post Post
	for rows.Next() {
		err = rows.Scan(&post.ID, &post.Title, &post.Content, &post.Date, &post.User_ID)
		if err != nil {
			return Post{}, err
		}
	}

	return post, nil
}

func (ph *PostHandler) getComments() ([]Comments, error) {
	query := "SELECT * from comments WHERE post_id = $1"
	rows, err := ph.DB.Query(query, PostIndex)
	if err != nil {
		return nil, err
	}
	var comments []Comments
	for rows.Next() {
		var comment Comment
		err = rows.Scan(&comment.ID, &comment.Content, &comment.Date, &comment.Post_ID, &comment.User_ID)
		if err != nil {
			return nil, err
		}
		var user User
		user, err = ph.getUser(comment.User_ID)
		if err != nil {
			return nil, err
		}
		coms := Comments{Author: user, Comment: comment}
		comments = append(comments, coms)
	}
	return comments, nil
}

func (ph *PostHandler) getUser(id int) (User, error) {
	query := "SELECT * from users WHERE id = $1"
	rows, err := ph.DB.Query(query, id)
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
