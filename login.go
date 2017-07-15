package main

import (
	"html/template"
	"net/http"
)

func login(res http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	username := req.FormValue("username")
	password := req.FormValue("password")
	res.Write([]byte(username + " - " + password))
}

func loginHandler(res http.ResponseWriter, req *http.Request) {
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
