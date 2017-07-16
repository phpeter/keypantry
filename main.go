package main

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"
)

var environment = os.Getenv("ENVIRONMENT")

var tmpl, tmplErr = template.ParseGlob("templates/*.html")

func main() {

	// declare port, db, and err variables so we can use them later on
	var port string
	var db *sql.DB
	var err error

	// set up port and DB connection based on environment
	if environment == "production" {
		port = os.Getenv("PORT")
		db, err = sql.Open("postgres", os.Getenv("DATABASE_URL"))
	} else {
		port = "8080"
		db, err = sql.Open("postgres", "postgresql://localhost:5432/peter?sslmode=disable")
	}
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
