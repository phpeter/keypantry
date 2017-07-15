package main

import (
	"database/sql"
	"html/template"
	"net/http"
)

func creatItemHandler(db *sql.DB) func(http.ResponseWriter, *http.Request) {

	return func(res http.ResponseWriter, req *http.Request) {
		tmpl, err := template.ParseGlob("templates/*.html")

		if err != nil {
			panic("error parsing templates " + err.Error())
		}

		switch req.Method {

		case "POST":
		case "GET":
			tmpl.ExecuteTemplate(res, "createItem", nil)

		}
	}
}
