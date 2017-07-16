package main

import (
	"database/sql"
	"log"
	"net/http"
)

type Item struct {
	ID      int
	Name    string
	Key     byte
	IsOwned bool
	KeyChar string
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
				tempItem.KeyChar = string(tempItem.Key)
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
			req.ParseForm()
			itemName := req.FormValue("name")
			itemKey := req.FormValue("key")
			_, err := db.Query("INSERT INTO items (name, key) VALUES ($1, $2)", itemName, itemKey)
			if err != nil {
				res.Write([]byte("There was an error creating that item: " + err.Error()))
			} else {
				http.Redirect(res, req, "/item/list", http.StatusFound)
			}
		case "GET":
			tmpl.ExecuteTemplate(res, "createItem", nil)

		}
	}
}

func deleteItemHandler(db *sql.DB) func(http.ResponseWriter, *http.Request) {
	return func(res http.ResponseWriter, req *http.Request) {
		itemID := getLastParam(req.URL.Path)
		_, err := db.Query("DELETE FROM items WHERE id=$1", itemID)
		if err != nil {
			res.Write([]byte("There was an error deleting item number " + itemID + ": " + err.Error()))
		} else {
			http.Redirect(res, req, "/item/list", http.StatusFound)
		}
	}
}

func editItemHandler(db *sql.DB) func(http.ResponseWriter, *http.Request) {
	return func(res http.ResponseWriter, req *http.Request) {

		itemID := getLastParam(req.URL.Path)

		switch req.Method {

		case "POST":
			req.ParseForm()
			itemName := req.FormValue("name")
			itemKey := req.FormValue("key")
			_, err := db.Query("UPDATE items SET name=$1, key=$2 WHERE id=$3", itemName, itemKey, itemID)
			if err != nil {
				res.Write([]byte("There was an error updating that item: " + err.Error()))
			} else {
				http.Redirect(res, req, "/item/list", http.StatusFound)
			}
		case "GET":
			var item Item
			row := db.QueryRow("SELECT name, key FROM items WHERE id=$1", itemID)
			err := row.Scan(&item.Name, &item.Key)
			if err != nil {
				res.Write([]byte("There was an error finding item number " + itemID + ": " + err.Error()))
			} else {
				tmpl.ExecuteTemplate(res, "editItem", item)
			}

		}
	}
}
