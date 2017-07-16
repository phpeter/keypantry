package main

import (
	"database/sql"
	"log"
	"net/http"
)

type Item struct {
	ID      int
	Name    string
	Key     rune
	IsOwned bool
}

func viewItemsHandler(db *sql.DB) func(http.ResponseWriter, *http.Request) {
	return func(res http.ResponseWriter, req *http.Request) {

		items, err := db.Query("SELECT id, name, key, isowned FROM items")
		var itemList []Item

		if err != nil {
			log.Print("No items, or error: " + err.Error())
		} else {

			var tempItem Item

			for items.Next() {
				items.Scan(&tempItem.ID, &tempItem.Name, &tempItem.Key, &tempItem.IsOwned)
				itemList = append(itemList, tempItem)
			}

		}

		log.Print(itemList)

		tmpl.ExecuteTemplate(res, "itemList", itemList)
	}
}

func createItemHandler(db *sql.DB) func(http.ResponseWriter, *http.Request) {

	return func(res http.ResponseWriter, req *http.Request) {

		switch req.Method {

		case "POST":
		case "GET":
			tmpl.ExecuteTemplate(res, "createItem", nil)

		}
	}
}
