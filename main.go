package main

import (
	"net/http"
)

func main() {

	http.HandleFunc("/login", loginHandler)

	http.HandleFunc("/toggleitem", apiHandler)

	http.ListenAndServe(":3000", nil)
}
