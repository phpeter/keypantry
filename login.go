package main

import (
	"database/sql"
	"log"
	"net/http"
	"time"
)

func loginRedirect(res http.ResponseWriter, req *http.Request) {
	http.Redirect(res, req, "/login", http.StatusTemporaryRedirect)
}

func loginHandler(ctx *appContext, res http.ResponseWriter, req *http.Request) (int, error) {

	db := ctx.db
	tmpl := ctx.tmpl

	c, err := req.Cookie("session")
	if err == nil && c != nil {
		sessionKey := c.Value
		session := db.QueryRow("SELECT userid FROM usersession WHERE sessionkey=$1", sessionKey)
		err = session.Scan()
		if err == nil {
			http.Redirect(res, req, "/item/list", http.StatusTemporaryRedirect)
			return http.StatusTemporaryRedirect, nil
		}
	}

	var pageError string

	switch req.Method {

	case "POST":
		// get username and pw from form
		req.ParseForm()
		username := req.FormValue("username")
		password := req.FormValue("password")

		// calculate the password hash
		passwordHash := pwHash(password, username)

		pw := db.QueryRow("SELECT id, passwordHash FROM users WHERE username=$1", username)
		err := pw.Scan(&ctx.user.id, &ctx.user.passwordHash)

		if err == sql.ErrNoRows || passwordHash != ctx.user.passwordHash {
			pageError = "Error! Bad credentials."
		} else if err != nil {
			log.Print("Error looking up user for login: " + err.Error())
			pageError = "An unknown error occurred, please try again."
		} else {
			// generate session key
			key := randKey(16)
			// set session cookie
			c := &http.Cookie{Name: "session", Value: key}
			http.SetCookie(res, c)
			// insert session row in to table
			timestamp := time.Now()

			rows, err := db.Query("INSERT INTO usersession (SessionKey, UserID, LoginTime, LastSeenTime) VALUES ($1, $2, $3, $3)", key, ctx.user.id, timestamp)
			defer rows.Close()

			if err != nil {
				log.Print("Error assigning a user session " + err.Error())
				pageError = "An unknown error occurred, please try again."
				return http.StatusInternalServerError, nil
			}

			http.Redirect(res, req, "/item/list", http.StatusTemporaryRedirect)
			return http.StatusTemporaryRedirect, nil

		}

		fallthrough
	default:
		tmpl.ExecuteTemplate(res, "login", struct{ Error string }{Error: pageError})
		return http.StatusOK, nil
	}

}

func logoutHandler(ctx *appContext, res http.ResponseWriter, req *http.Request) (int, error) {

	db := ctx.db

	// check if cookie is actually set
	c, err := req.Cookie("session")
	if err != nil {
		log.Print("Error getting session cookie: " + err.Error())
		loginRedirect(res, req)
		return http.StatusInternalServerError, err
	}
	// if cookie is set, get sesion key
	sessionKey := c.Value

	// delete  session key from db
	rows, err := db.Query("DELETE FROM usersession WHERE sessionkey=$1", sessionKey)
	defer rows.Close()

	if err != nil {
		log.Print("Error deleting session from DB: " + err.Error())
	}

	c = &http.Cookie{Name: "session", Value: "", Expires: time.Now()}
	http.SetCookie(res, c)

	loginRedirect(res, req)
	return http.StatusTemporaryRedirect, nil
}
