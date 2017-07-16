package main

import (
	"database/sql"
	"log"
	"net/http"
	"strings"
)

func getLastParam(path string) string {
	params := strings.Split(path, "/")
	return params[len(params)-1]
}

func apiHandler(db *sql.DB) func(res http.ResponseWriter, req *http.Request) {
	return func(res http.ResponseWriter, req *http.Request) {
		// x, _ := httputil.DumpRequest(req, true)
		// log.Print(string(x))

		username := req.Header.Get("Username")
		password := req.Header.Get("Password")

		row := db.QueryRow("SELECT passwordHash FROM users WHERE username=$1", username)
		var userPw string
		err := row.Scan(&userPw)

		if err == sql.ErrNoRows || pwHash(password) != userPw {
			res.Write([]byte("Bad credentials!"))
		} else {

			row := db.QueryRow("SELECT id FROM users WHERE username=$1", username)
			var userID int
			err := row.Scan(&userID)

			if err != nil {
				log.Print("Error, no user found")
			}

			keyPressed := getLastParam(req.URL.Path)

			log.Print("user id is " + string(userID))

			_, err = db.Query("UPDATE items SET isOwned = NOT isOwned WHERE userID=$1 AND key=$2", userID, keyPressed)
			if err != nil {
				log.Print(err)
			}
			res.Write([]byte("Good!"))
		}

	}
}
