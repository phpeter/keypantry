package main

import (
	"database/sql"
	"log"
	"net/http"
	"time"
)

func loginHandler(db *sql.DB) func(http.ResponseWriter, *http.Request) {

	return func(res http.ResponseWriter, req *http.Request) {

		c, err := req.Cookie("session")
		if err == nil && c != nil {
			sessionKey := c.Value
			var userID int
			session := db.QueryRow("SELECT userid FROM usersession WHERE sessionkey=$1", sessionKey)
			err = session.Scan(&userID)
			if err == nil {
				http.Redirect(res, req, "/item/list", http.StatusTemporaryRedirect)
				return
			}
		}

		switch req.Method {

		case "POST":
			req.ParseForm()
			username := req.FormValue("username")
			password := req.FormValue("password")
			passwordHash := pwHash(password, username)

			var userPasswordHash string
			var userID int
			pw := db.QueryRow("SELECT id, passwordHash FROM users WHERE username=$1", username)
			err := pw.Scan(&userID, &userPasswordHash)

			log.Print(username + " : " + passwordHash)

			if err == sql.ErrNoRows {
				res.Write([]byte("Error! Wrong username."))
			} else if err != nil {
				res.Write([]byte("Error! " + err.Error()))
			} else if passwordHash != userPasswordHash {
				res.Write([]byte("Error! Wrong password."))
			} else {
				// generate session key
				key := randKey(16)
				// set session cookie
				c := &http.Cookie{Name: "session", Value: key}
				http.SetCookie(res, c)
				// insert session row in to table
				timestamp := time.Now()
				_, err := db.Query("INSERT INTO usersession (SessionKey, UserID, LoginTime, LastSeenTime) VALUES ($1, $2, $3, $3)", key, userID, timestamp)
				if err != nil {
					res.Write([]byte("Error! " + err.Error()))
				} else {
					http.Redirect(res, req, "/item/list", http.StatusTemporaryRedirect)
				}
			}

		case "GET":
			tmpl.ExecuteTemplate(res, "login", nil)

		}
	}
}

func logoutHandler(db *sql.DB) func(http.ResponseWriter, *http.Request) {
	return func(res http.ResponseWriter, req *http.Request) {

		// check if cookie is actually set
		c, err := req.Cookie("session")
		if err != nil {
			log.Print("Error getting session cookie: " + err.Error())
			loginRedirect(res, req)
			return
		}
		// if cookie is set, get sesion key
		sessionKey := c.Value

		// delete  session key from db
		_, err = db.Query("DELETE FROM usersession WHERE key=$1", sessionKey)
		if err != nil {
			log.Print("Error deleting session from DB: " + err.Error())
		}

		c = &http.Cookie{Name: "session", Value: "", Expires: time.Now()}
		http.SetCookie(res, c)

		loginRedirect(res, req)

	}
}
