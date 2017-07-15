package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
)

func apiHandler(db *sql.DB) func(http.ResponseWriter, *http.Request) {
	return func(res http.ResponseWriter, req *http.Request) {
		// check credentials
		// if credentials are correct, toggle item
		//db.Query("UPDATE items WHERE user = user.id AND key = keyPressed SET own = NOT own")

		log.Print("Connected to api handler")

		for name, value := range res.Header() {
			res.Write([]byte(fmt.Sprintf("%v: %v\n", name, value)))
		}

		res.WriteHeader(200)
	}
}
