package main

import (
	"crypto/sha256"
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
)

func pwHash(password string) string {
	// need to add some salt
	return fmt.Sprintf("%x", sha256.Sum256([]byte(password)))
}

func login(res http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	username := req.FormValue("username")
	password := req.FormValue("password")

	passwordHash := pwHash(password)

	res.Write([]byte(username + " - " + passwordHash))
}

func loginHandler(db *sql.DB) func(http.ResponseWriter, *http.Request) {

	return func(res http.ResponseWriter, req *http.Request) {
		tmpl, err := template.ParseGlob("templates/*.html")

		if err != nil {
			panic("error parsing templates " + err.Error())
		}

		switch req.Method {

		case "POST":
			login(res, req)
		case "GET":
			tmpl.ExecuteTemplate(res, "login", nil)

		}
	}
}
