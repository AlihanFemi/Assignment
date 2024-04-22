package handlers

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"
	"time"

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

type Comment struct {
	ID      int
	Content string
	Date    string
	Post_ID int
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
		path := strings.Split(r.URL.Path, "/")
		var err error
		PostIndex, err = strconv.Atoi(path[2])
		if err != nil {
			fmt.Println(err)
			fmt.Fprintln(w, err)
		}
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
	r.ParseForm()
	content := r.Form.Get("comment")
	data := time.Now().Format("2006-01-02")
	sess := ph.Sess.Start(w, r)
	user := sess.Get("User").(User)
	post_id := PostIndex
	user_id := user.ID
	err := ph.insertComment(content, data, post_id, user_id)
	if err != nil {
		fmt.Println(err)
		fmt.Fprintln(w, "Internal server error!")
		return
	}
	ph.ServerGet(w, r)
}

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

func (ph *PostHandler) insertComment(content string, date string, post_id int, user_id int) error {
	const query = "insert into comments (content, date, post_id, user_id) VALUES ($1, $2, $3, $4)"
	_, err := ph.DB.Exec(query, content, date, post_id, user_id)
	if err != nil {
		return err
	}
	return nil
}
