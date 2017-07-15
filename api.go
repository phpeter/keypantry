package main

import (
	"database/sql"
	"log"
	"net/http"
)

func apiHandler(db *sql.DB) func(res http.ResponseWriter, req *http.Request) {
	return func(res http.ResponseWriter, req *http.Request) {
		// check credentials
		// if credentials are correct, toggle item
		// db.Query("UPDATE items WHERE user = user.id AND key = keyPressed SET own = NOT own")

		username := req.Header.Get("Username")
		password := req.Header.Get("Password")

		row := db.QueryRow("SELECT passwordHash FROM users WHERE username=$1", username)
		var userPw string
		err := row.Scan(&userPw)

		if err != nil {
			log.Fatal("Error getting user password for API call " + err.Error())
		}

		if pwHash(password) == userPw {
			res.Write([]byte("Correct pw"))
		} else {
			res.Write([]byte("Incorrect pw!!!!!!!"))
		}

		// x, _ := httputil.DumpRequest(req, true)
		// log.Print(string(x))

		log.Print(username + " - " + password)

		// res.WriteHeader(200)
	}
}
