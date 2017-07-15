package main

import (
	"net/http"
	"os"
)

var port = os.Getenv("PORT")

func main() {

	http.HandleFunc("/login", loginHandler)

	http.HandleFunc("/toggleitem", apiHandler)

	http.ListenAndServe(":"+port, nil)
}
