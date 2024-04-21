package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/alihanfemi/server/config"
	"github.com/alihanfemi/server/handlers"
	"github.com/kataras/go-sessions/v3"
	_ "github.com/lib/pq"
)

func main() {
	cfg, err := config.Parse()
	if err != nil {
		panic(err)
	}
	db, err := sql.Open("postgres", cfg.Database.BuildConnectionString())
	if err != nil {
		panic(err)
	}
	fs := http.FileServer(http.Dir("assets"))
	sess := sessions.New(sessions.Config{Expires: time.Hour * 2})
	indexHandler := handlers.IndexHandler{Sess: sess, DB: db}
	postHandler := handlers.PostHandler{Sess: sess, DB: db}
	registerHandler := handlers.RegisterHandler{Sess: sess, DB: db}
	loginHandler := handlers.LoginHandler{Sess: sess, DB: db}

	http.Handle("/assets/", http.StripPrefix("/assets/", fs))
	http.Handle("/", &indexHandler)
	http.Handle("/register", &registerHandler)
	http.Handle("/login", &loginHandler)
	http.Handle("/view/", &postHandler)
	http.Handle("/post-in", &postHandler)
	fmt.Println("Started")
	http.ListenAndServe(":42069", nil)
}
