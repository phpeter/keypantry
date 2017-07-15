package main

import (
	"database/sql"
	"log"
	"net/http"
	"net/http/httputil"
)

func apiHandler(db *sql.DB) func(res http.ResponseWriter, req *http.Request) {
	return func(res http.ResponseWriter, req *http.Request) {
		// check credentials
		// if credentials are correct, toggle item
		// db.Query("UPDATE items WHERE user = user.id AND key = keyPressed SET own = NOT own")

		x, _ := httputil.DumpRequest(req, true)
		log.Print(string(x))

		res.WriteHeader(200)
	}
}
