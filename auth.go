package main

import (
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"fmt"
	"log"
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
	return fmt.Sprintf("%x", sha256.Sum256([]byte(password+salt)))
}

func isAuthorized(req *http.Request, user *user, db *sql.DB) bool {
	c, err := req.Cookie("session")
	if err != nil {
		return false
	}

	sessionKey := c.Value
	session := db.QueryRow("SELECT userid FROM usersession WHERE sessionkey=$1", sessionKey)
	err = session.Scan(&user.id)
	if err != nil {
		if err != sql.ErrNoRows {
			log.Print("Error retrieving user session: " + err.Error())
		}
		return false
	}

	return true
}
