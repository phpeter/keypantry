package main

import (
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"fmt"
	"net/http"
)

type Context struct {
	UserID   int
	LoggedIn bool
}

func randKey(chars int) string {
	b := make([]byte, chars)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}
	return fmt.Sprintf("%X", b)
}

func pwHash(password string, salt string) string {
	// need to add some salt
	return fmt.Sprintf("%x", sha256.Sum256([]byte(password+salt)))
}

func loginRedirect(res http.ResponseWriter, req *http.Request) {
	http.Redirect(res, req, "/login", http.StatusUnauthorized)
}

func auth(handler func(*sql.DB, *Context) func(res http.ResponseWriter, req *http.Request), db *sql.DB) func(res http.ResponseWriter, req *http.Request) {
	return func(res http.ResponseWriter, req *http.Request) {

		ctx := &Context{}

		c, err := req.Cookie("session")
		if err != nil {
			loginRedirect(res, req)
		} else {
			sessionKey := c.Value
			var userID int
			session := db.QueryRow("SELECT userid FROM usersession WHERE sessionkey=$1", sessionKey)
			err = session.Scan(&userID)
			if err == sql.ErrNoRows {
				loginRedirect(res, req)
			} else if err != nil {
				res.Write([]byte("Error! " + err.Error()))
			} else {
				ctx.LoggedIn = true
				ctx.UserID = userID
				handler(db, ctx)(res, req)
			}
		}
	}
}
