package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"
)

var port = os.Getenv("PORT")

func main() {

	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/login", loginHandler(db))

	http.HandleFunc("/toggleitem/	", apiHandler(db))

	http.ListenAndServe(":"+port, nil)
}
