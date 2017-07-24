package main

import (
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

func viewItemsHandler(ctx *appContext, res http.ResponseWriter, req *http.Request) (int, error) {

	db := ctx.db
	tmpl := ctx.tmpl

	owned := req.URL.Query().Get("owned")
	var querySuffix string
	var view string
	if owned == "true" {
		view = "owned"
		querySuffix = " AND isowned='t'"
	} else if owned == "false" {
		view = "notowned"
		querySuffix = " AND isowned='f'"
	} else {
		view = "all"
	}

	items, err := db.Query("SELECT id, name, key, isowned FROM items WHERE userid=$1"+querySuffix+" ORDER BY id DESC", ctx.user.id)
	defer items.Close()

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
			User     user
			ItemList []Item
			View     string
		}{
			User:     ctx.user,
			ItemList: itemList,
			View:     view,
		})

	return http.StatusOK, nil
}

func createItemHandler(ctx *appContext, res http.ResponseWriter, req *http.Request) (int, error) {

	db := ctx.db
	tmpl := ctx.tmpl

	switch req.Method {

	case "POST":
		req.ParseForm()
		itemName := req.FormValue("name")
		itemKey := req.FormValue("key")

		rows, err := db.Query("INSERT INTO items (name, key, userid, isowned) VALUES ($1, $2, $3, FALSE)", itemName, itemKey, ctx.user.id)
		defer rows.Close()

		if err != nil {
			res.Write([]byte("There was an error creating that item: " + err.Error()))
			return http.StatusInternalServerError, err
		}
		http.Redirect(res, req, req.Referer(), http.StatusFound)
		return http.StatusOK, nil

	default:
		tmpl.ExecuteTemplate(res, "createItem", nil)
		return http.StatusOK, nil
	}
}

func deleteItemHandler(ctx *appContext, res http.ResponseWriter, req *http.Request) (int, error) {
	db := ctx.db
	itemID := getLastParam(req.URL.Path)
	rows, err := db.Query("DELETE FROM items WHERE id=$1 AND userid=$2", itemID, ctx.user.id)
	defer rows.Close()
	if err != nil {
		return http.StatusInternalServerError, err
	}

	http.Redirect(res, req, req.Referer(), http.StatusFound)
	return http.StatusOK, nil

}

func editItemHandler(ctx *appContext, res http.ResponseWriter, req *http.Request) (int, error) {

	db := ctx.db
	tmpl := ctx.tmpl

	itemID := getLastParam(req.URL.Path)

	switch req.Method {

	case "POST":
		req.ParseForm()
		itemName := req.FormValue("name")
		itemKey := req.FormValue("key")
		referer := req.FormValue("referer")
		rows, err := db.Query("UPDATE items SET name=$1, key=$2 WHERE id=$3 AND userid=$4", itemName, itemKey, itemID, ctx.user.id)
		defer rows.Close()
		if err != nil {
			return http.StatusInternalServerError, err
		}

		http.Redirect(res, req, referer, http.StatusFound)
		return http.StatusOK, nil

	default:
		var item Item
		row := db.QueryRow("SELECT name, key FROM items WHERE id=$1 AND userid=$2", itemID, ctx.user.id)
		err := row.Scan(&item.Name, &item.Key)
		if err != nil {
			return http.StatusInternalServerError, err
		}

		tmpl.ExecuteTemplate(res, "editItem", struct {
			User    user
			Item    Item
			Referer string
		}{
			User:    ctx.user,
			Item:    item,
			Referer: req.Referer(),
		})

		return http.StatusOK, nil
	}

}

func toggleItemHandler(ctx *appContext, res http.ResponseWriter, req *http.Request) (int, error) {

	db := ctx.db

	itemID := getLastParam(req.URL.Path)

	rows, err := db.Query("UPDATE items SET isOwned = NOT isOwned WHERE userID=$1 AND id=$2", ctx.user.id, itemID)
	defer rows.Close()

	if err != nil {
		return http.StatusInternalServerError, err
	}

	http.Redirect(res, req, req.Referer(), http.StatusTemporaryRedirect)
	return http.StatusTemporaryRedirect, nil
}
