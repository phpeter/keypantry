package main

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"
)

type user struct {
	id           int
	username     string
	passwordHash string
}

type appContext struct {
	db   *sql.DB
	tmpl *template.Template
	user user
}

type appHandler struct {
	*appContext
	H    func(*appContext, http.ResponseWriter, *http.Request) (int, error)
	auth bool
}

func (ah appHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {

	if ah.auth && !isAuthorized(req, &ah.user, ah.db) {
		http.Redirect(res, req, "/login", http.StatusTemporaryRedirect)
		return
	}

	status, err := ah.H(ah.appContext, res, req)

	if err != nil {
		log.Printf("HTTP %d: %q", status, err)
		switch status {
		case http.StatusNotFound:
			http.NotFound(res, req)
		case http.StatusInternalServerError:
			http.Error(res, http.StatusText(status), status)
		default:
			http.Error(res, http.StatusText(status), status)
		}
	}
}

func main() {

	var databaseURL = os.Getenv("DATABASE_URL")
	var db, err = sql.Open("postgres", databaseURL)
	defer db.Close()

	var tmpl, tmplErr = template.ParseGlob("templates/*.html")
	if tmplErr != nil {
		panic("Error parsing template: " + tmplErr.Error())
	}

	var context = &appContext{db: db, tmpl: tmpl}

	port := os.Getenv("PORT")
	// catch DB connection error
	if err != nil {
		log.Fatal(err)
	}

	server := http.NewServeMux()

	server.Handle("/login", appHandler{context, loginHandler, false})

	server.Handle("/logout", appHandler{context, logoutHandler, false})
	server.Handle("/register", appHandler{context, registerHandler, false})

	server.Handle("/toggleitem/", appHandler{context, apiHandler, false})

	server.Handle("/item/list", appHandler{context, viewItemsHandler, true})
	server.Handle("/item/create", appHandler{context, createItemHandler, true})
	server.Handle("/item/delete/", appHandler{context, deleteItemHandler, true})
	server.Handle("/item/edit/", appHandler{context, editItemHandler, true})
	server.Handle("/item/toggle/", appHandler{context, toggleItemHandler, true})

	// redirect from root URL to login screen or item list
	server.Handle("/", appHandler{nil, func(ctx *appContext, res http.ResponseWriter, req *http.Request) (int, error) {
		http.Redirect(res, req, "/item/list", http.StatusPermanentRedirect)
		return http.StatusTemporaryRedirect, nil
	}, false})

	log.Print("Running on port " + port)
	http.ListenAndServe(":"+port, server)
}
