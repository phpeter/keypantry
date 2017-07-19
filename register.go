package main

import (
	"database/sql"
	"log"
	"net/http"
	"time"
)

func registerHandler(db *sql.DB) func(http.ResponseWriter, *http.Request) {
	return func(res http.ResponseWriter, req *http.Request) {

		var pageError string

		switch req.Method {
		case "POST":

			req.ParseForm()
			username := req.FormValue("username")
			password := req.FormValue("password")
			passwordConf := req.FormValue("passwordConf")

			userLookup := db.QueryRow("SELECT FROM users WHERE username=$1", username).Scan()

			if password != passwordConf {
				pageError = "Error, password and password confirmation do not match"
			} else if userLookup == nil {
				pageError = "Username already taken!"
			} else if userLookup != sql.ErrNoRows {
				log.Print("Error looking up user for register: " + userLookup.Error())
				pageError = "An unknown error occurred, please try again."
			} else {

				var userID int
				err := db.QueryRow("INSERT INTO users (username, passwordhash) VALUES ($1, $2) RETURNING id", username, pwHash(password, username)).Scan(&userID)
				if err != nil {
					log.Print("Error creating new user: " + err.Error())
					pageError = "An unknown error occurred, please try again."
				} else {
					// generate session key
					key := randKey(16)
					// set session cookie
					c := &http.Cookie{Name: "session", Value: key}
					http.SetCookie(res, c)
					// insert session row in to table
					timestamp := time.Now()

					rows, err := db.Query("INSERT INTO usersession (SessionKey, UserID, LoginTime, LastSeenTime) VALUES ($1, $2, $3, $3)", key, userID, timestamp)
					if err != nil {
						log.Print("Error logging user in: " + err.Error())
						http.Redirect(res, req, "/login", http.StatusTemporaryRedirect)
						return
					}
					defer rows.Close()

					http.Redirect(res, req, "/item/list", http.StatusTemporaryRedirect)
					return

				}

			}

			fallthrough

		case "GET":
			tmpl.ExecuteTemplate(res, "register", struct{ Error string }{Error: pageError})

		}

	}
}
