package main

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"
)

var tmpl, tmplErr = template.ParseGlob("templates/*.html")

func main() {

	port := os.Getenv("PORT")
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	// catch DB connection error
	if err != nil {
		log.Fatal(err)
	}

	if tmplErr != nil {
		panic("Error parsing template: " + tmplErr.Error())
	}

	http.HandleFunc("/toggleitem/", apiHandler(db))

	http.HandleFunc("/login", loginHandler(db))

	http.HandleFunc("/item/create", createItemHandler(db))

	http.HandleFunc("/item/list", viewItemsHandler(db))

	log.Print("Running on port " + port)
	http.ListenAndServe(":"+port, nil)
}
