package main

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"
)

var port = os.Getenv("PORT")

var tmpl, tmplErr = template.ParseGlob("templates/*.html")

func main() {

	if tmplErr != nil {
		panic("Error parsing template: " + tmplErr.Error())
	}

	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/toggleitem/", apiHandler(db))

	http.HandleFunc("/login", loginHandler(db))

	http.HandleFunc("/item/create", createItemHandler(db))

	http.HandleFunc("/item/list", viewItemsHandler(db))

	http.ListenAndServe(":"+port, nil)
}
