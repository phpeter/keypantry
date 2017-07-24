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

func apiHandler(ctx *appContext, res http.ResponseWriter, req *http.Request) (int, error) {

	db := ctx.db

	username := req.Header.Get("Username")
	password := req.Header.Get("Password")

	row := db.QueryRow("SELECT passwordHash FROM users WHERE username=$1", username)
	var userPw string
	err := row.Scan(&userPw)

	if err == sql.ErrNoRows || pwHash(password, username) != userPw {
		res.Write([]byte("Bad credentials!"))
		return http.StatusUnauthorized, err
	}

	row = db.QueryRow("SELECT id FROM users WHERE username=$1", username)

	var userID int
	err = row.Scan(&userID)

	if err != nil {
		log.Print("Error, no user found")
		return http.StatusInternalServerError, err
	}

	keyPressed := getLastParam(req.URL.Path)

	log.Print("user id is " + string(userID))

	rows, err := db.Query("UPDATE items SET isOwned = NOT isOwned WHERE userID=$1 AND key=$2", userID, keyPressed)
	defer rows.Close()
	if err != nil {
		log.Print(err)
		return http.StatusInternalServerError, err
	}

	res.Write([]byte("Good!"))
	return http.StatusOK, nil
}
