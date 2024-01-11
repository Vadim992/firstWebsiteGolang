package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

type Article struct {
	Id        uint16
	Title     string
	Anons     string
	Full_text string
}

type DataTemplate struct {
	Articles   []Article
	ActivePage string
}

// var articles = []Article{}

func index(w http.ResponseWriter, r *http.Request) {

	t, err := template.ParseFiles("templates/index.html", "templates/header.html", "templates/footer.html")

	if err != nil {
		fmt.Fprint(w, err.Error())
	}

	db, err := sql.Open("mysql", "firstAppUser:mypass@tcp(localhost:3306)/golang")
	if err != nil {
		panic(err)
	}

	defer db.Close()

	if err = db.Ping(); err != nil {
		panic(err)
	}
	q := "SELECT * FROM `articles`"

	querry, err := db.Query(q)
	if err != nil {
		panic(err)
	}
	defer querry.Close()

	articles := []Article{}

	for querry.Next() {
		var article Article
		err = querry.Scan(&article.Id, &article.Title, &article.Anons, &article.Full_text)
		if err != nil {
			panic(err)
		}

		articles = append(articles, article)
	}

	activePage := "/"

	dataTemplate := DataTemplate{
		articles,
		activePage,
	}

	err = t.ExecuteTemplate(w, "index", dataTemplate)

	if err != nil {
		fmt.Fprint(w, err.Error())
	}
}

func create(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/create.html", "templates/header.html", "templates/footer.html")

	if err != nil {
		fmt.Fprint(w, err.Error())
	}

	activePage := "/create/"

	dataTemplate := DataTemplate{
		ActivePage: activePage,
	}

	err = t.ExecuteTemplate(w, "create", dataTemplate)

	if err != nil {
		fmt.Fprint(w, err.Error())
	}

}

func save_article(w http.ResponseWriter, r *http.Request) {
	title := r.FormValue("title")
	anons := r.FormValue("anons")
	full_text := r.FormValue("full_text")

	columns := []string{title, anons, full_text}
	var flag bool

	for i := range columns {
		columns[i] = strings.TrimSpace(columns[i])
		if columns[i] == "" {
			flag = true
			break
		}
	}

	if flag {
		fmt.Fprint(w, "Please fill in all fields of the form")
	} else {
		db, err := sql.Open("mysql", "firstAppUser:mypass@tcp(localhost:3306)/golang")
		if err != nil {
			panic(err)
		}
		defer db.Close()

		if err = db.Ping(); err != nil {
			panic(err)
		}

		q := fmt.Sprintf("INSERT INTO `articles` (`title`, `anons`, `full_text`) VALUES('%s',' %s', '%s')", title, anons, full_text)

		insert, err := db.Query(q)

		if err != nil {
			panic(err)
		}

		defer insert.Close()

		http.Redirect(w, r, "/", http.StatusSeeOther)
	}

}

func show_post(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	t, err := template.ParseFiles("templates/show_post.html", "templates/header.html", "templates/footer.html")

	if err != nil {
		fmt.Fprint(w, err.Error())
	}

	db, err := sql.Open("mysql", "firstAppUser:mypass@tcp(localhost:3306)/golang")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	if err = db.Ping(); err != nil {
		panic(err)
	}

	q := fmt.Sprintf("SELECT * FROM `articles` WHERE `id` = '%s'", vars["id"])

	querry, err := db.Query(q)
	if err != nil {
		panic(err)
	}
	defer querry.Close()

	article := Article{}

	for querry.Next() {
		err = querry.Scan(&article.Id, &article.Title, &article.Anons, &article.Full_text)
		if err != nil {
			panic(err)
		}

	}

	articles := []Article{article}

	dataTemplate := DataTemplate{
		Articles: articles,
	}

	t.ExecuteTemplate(w, "show_post", dataTemplate)
}

func handleFunc() {
	r := mux.NewRouter()
	r.HandleFunc("/", index).Methods("GET")
	r.HandleFunc("/create/", create).Methods("GET")
	r.HandleFunc("/save_article/", save_article).Methods("POST")
	r.HandleFunc("/post/{id:[0-9]+}", show_post).Methods("GET")

	// http.Handle("/", r)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))
	log.Fatal(http.ListenAndServe(":8080", r))

}

func main() {
	handleFunc()
}
