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

func viewItemsHandler(db *sql.DB, ctx *Context) func(http.ResponseWriter, *http.Request) {
	return func(res http.ResponseWriter, req *http.Request) {

		ownedQuery := req.URL.Query().Get("owned")
		var querySuff string
		var view string
		if ownedQuery == "true" {
			view = "owned"
			querySuff = " AND isowned='t'"
		} else if ownedQuery == "false" {
			view = "notowned"
			querySuff = " AND isowned='f'"
		} else {
			view = "all"
		}

		items, err := db.Query("SELECT id, name, key, isowned FROM items WHERE userid=$1 ORDER BY id DESC"+querySuff, ctx.UserID)
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

		tmpl.ExecuteTemplate(res, "itemList",
			struct {
				Context  *Context
				ItemList []Item
				View     string
			}{
				Context:  ctx,
				ItemList: itemList,
				View:     view,
			})
	}
}

func createItemHandler(db *sql.DB, ctx *Context) func(http.ResponseWriter, *http.Request) {
	return func(res http.ResponseWriter, req *http.Request) {

		switch req.Method {

		case "POST":
			req.ParseForm()
			itemName := req.FormValue("name")
			itemKey := req.FormValue("key")
			_, err := db.Query("INSERT INTO items (name, key, userid, isowned) VALUES ($1, $2, $3, FALSE)", itemName, itemKey, ctx.UserID)
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

func deleteItemHandler(db *sql.DB, ctx *Context) func(http.ResponseWriter, *http.Request) {
	return func(res http.ResponseWriter, req *http.Request) {
		itemID := getLastParam(req.URL.Path)
		_, err := db.Query("DELETE FROM items WHERE id=$1 AND userid=$2", itemID, ctx.UserID)
		if err != nil {
			res.Write([]byte("There was an error deleting item number " + itemID + ": " + err.Error()))
		} else {
			http.Redirect(res, req, "/item/list", http.StatusFound)
		}
	}
}

func editItemHandler(db *sql.DB, ctx *Context) func(http.ResponseWriter, *http.Request) {
	return func(res http.ResponseWriter, req *http.Request) {

		itemID := getLastParam(req.URL.Path)

		switch req.Method {

		case "POST":
			req.ParseForm()
			itemName := req.FormValue("name")
			itemKey := req.FormValue("key")
			_, err := db.Query("UPDATE items SET name=$1, key=$2 WHERE id=$3 AND userid=$4", itemName, itemKey, itemID, ctx.UserID)
			if err != nil {
				res.Write([]byte("There was an error updating that item: " + err.Error()))
			} else {
				http.Redirect(res, req, "/item/list", http.StatusFound)
			}
		case "GET":
			var item Item
			row := db.QueryRow("SELECT name, key FROM items WHERE id=$1 AND userid=$2", itemID, ctx.UserID)
			err := row.Scan(&item.Name, &item.Key)
			if err != nil {
				res.Write([]byte("There was an error finding item number " + itemID + ": " + err.Error()))
			} else {
				tmpl.ExecuteTemplate(res, "editItem", struct {
					Context *Context
					Item    Item
				}{
					Context: ctx,
					Item:    item,
				})
			}

		}
	}
}

func toggleItemHandler(db *sql.DB, ctx *Context) func(http.ResponseWriter, *http.Request) {
	return func(res http.ResponseWriter, req *http.Request) {

		itemID := getLastParam(req.URL.Path)

		_, err := db.Query("UPDATE items SET isOwned = NOT isOwned WHERE userID=$1 AND id=$2", ctx.UserID, itemID)

		if err != nil {
			res.Write([]byte("Error toggling item ownership: " + err.Error()))
		} else {
			http.Redirect(res, req, "/item/list", http.StatusTemporaryRedirect)
		}

	}
}
