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

		for name, value := range res.Header() {
			log.Printf("%v: %v\n", name, value)
		}

		res.WriteHeader(200)
	}
}
