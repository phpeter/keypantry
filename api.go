package main

import (
	"database/sql"
	"net/http"
)

func apiHandler(*sql.DB) func(http.ResponseWriter, *http.Request) {
	return func(res http.ResponseWriter, req *http.Request) {
		res.WriteHeader(200)
	}
}
