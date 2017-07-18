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

	var databaseURL = os.Getenv("DATABASE_URL")
	var db, err = sql.Open("postgres", databaseURL)
	defer db.Close()

	port := os.Getenv("PORT")
	// catch DB connection error
	if err != nil {
		log.Fatal(err)
	}

	if tmplErr != nil {
		panic("Error parsing template: " + tmplErr.Error())
	}

	http.HandleFunc("/login", loginHandler(db))
	http.HandleFunc("/logout", logoutHandler(db))

	http.HandleFunc("/toggleitem/", apiHandler(db))

	http.HandleFunc("/item/list", auth(viewItemsHandler, db))
	http.HandleFunc("/item/create", auth(createItemHandler, db))
	http.HandleFunc("/item/delete/", auth(deleteItemHandler, db))
	http.HandleFunc("/item/edit/", auth(editItemHandler, db))
	http.HandleFunc("/item/toggle/", auth(toggleItemHandler, db))

	// redirect from root URL to login screen or item list
	http.HandleFunc("/", func(res http.ResponseWriter, req *http.Request) {
		http.Redirect(res, req, "/item/list", http.StatusPermanentRedirect)
	})

	log.Print("Running on port " + port)
	http.ListenAndServe(":"+port, nil)
}
