package main

import (
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"fmt"
	"net/http"
)

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

func auth(handler func(db *sql.DB) func(res http.ResponseWriter, req *http.Request), db *sql.DB) func(res http.ResponseWriter, req *http.Request) {
	return func(res http.ResponseWriter, req *http.Request) {

		c, err := req.Cookie("session")
		if err != nil {
			loginRedirect(res, req)
		} else {
			sessionKey := c.Value
			session := db.QueryRow("SELECT FROM usersession WHERE sessionkey=$1", sessionKey)
			err = session.Scan()
			if err == sql.ErrNoRows {
				loginRedirect(res, req)
			} else if err != nil {
				res.Write([]byte("Error! " + err.Error()))
			} else {
				handler(db)(res, req)
			}
		}
	}
}
