package main

import (
	"database/sql"
	"net/http"
)

func apiHandler(db *sql.DB) func(res http.ResponseWriter, req *http.Request) {
	return func(res http.ResponseWriter, req *http.Request) {
		// check credentials
		// if credentials are correct, toggle item
		// db.Query("UPDATE items WHERE user = user.id AND key = keyPressed SET own = NOT own")

		username := req.Header.Get("Username")
		password := req.Header.Get("Password")

		// x, _ := httputil.DumpRequest(req, true)
		// log.Print(string(x))

		res.Write([]byte(username + " - " + password))

		res.WriteHeader(200)
	}
}
