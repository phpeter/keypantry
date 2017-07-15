package main

import "net/http"

func apiHandler(res http.ResponseWriter, req *http.Request) {
	res.WriteHeader(200)
}
