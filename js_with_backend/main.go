package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"database/sql" 
//	"strings"
//	"strconv"
	"encoding/json"
	 _ "github.com/mattn/go-sqlite3"	
)

type SearchItem struct {
	Ru string `json:"ru"`
	Ge string `json:"ge"`
}

type SearchResponse struct {
    Success bool `json:"success"`
    Message string `json:"message,omitempty"`
    Lang string `json:"lang"`
	Items []SearchItem `json:"items,omitempty"`   
}

func WriteResponse(writer http.ResponseWriter, response SearchResponse) {
	bytes, err := json.Marshal(response)

   	if err != nil {
    	panic(err)
   	}

	if (response.Success) {
		writer.WriteHeader(http.StatusOK)
	} else {
		writer.WriteHeader(http.StatusInternalServerError)	
  	}

  	writer.Header().Set("Access-Control-Allow-Origin", "http://32cats.com")
	writer.Header().Set("Access-Control-Allow-Methods", "GET")
	writer.Header().Set("Content-Type", "application/json")

  	io.WriteString(writer, string(bytes))
}

func getRoot(writer http.ResponseWriter, r *http.Request) {
	io.WriteString(writer, "=^_^=")
}

func getRows(rows *sql.Rows) []SearchItem {
	var items []SearchItem
	var ru, ge string

	for rows.Next() {
		err := rows.Scan(&ru, &ge)
		if (err == nil) {
			item := SearchItem { ru, ge }
			items = append(items, item)
		}
	}
	return items
}

func search(writer http.ResponseWriter, r *http.Request) {
	var items []SearchItem

	fmt.Printf("Search %v\n", r.URL)

	hasQuery := r.URL.Query().Has("query")

	if (hasQuery) {
		query := r.URL.Query().Get("query")
		
		fmt.Printf("query: %s\n", query)

		const file string = "dict.db"
		db, err := sql.Open("sqlite3", file)

		if err != nil {
			WriteResponse(writer, SearchResponse{ Success: false, Message: "Database opening error"})
			fmt.Printf("Database opening error\n")
	  		return
		}
		
		defer db.Close()		

		var sql = fmt.Sprintf("SELECT ru, ge  FROM words WHERE ru LIKE '%s%%' LIMIT 10", query)
		fmt.Printf("sql: %s\n", sql)
		rows, err := db.Query(sql)

		if err != nil {
			WriteResponse(writer, SearchResponse{ Success: false, Message: "Database query error"})
			fmt.Printf("Database query error\n")
	  		return
		}

		items = getRows(rows)
		if (len(items) > 0) {
			fmt.Printf("%d items found\n", len(items))
			WriteResponse(writer, SearchResponse{ Success: true, Lang: "ru", Items: items })
			return
		}

		sql = fmt.Sprintf("SELECT ru, ge  FROM words WHERE ge LIKE '%%%s%%' LIMIT 10", query)
		fmt.Printf("sql: %s\n", sql)
		rows, err = db.Query(sql)

		if err != nil {
			WriteResponse(writer, SearchResponse{ Success: false, Message: "Database query error"})
			fmt.Printf("Database query error\n")
	  		return
		}

		items = getRows(rows)
		fmt.Printf("%d items found\n", len(items))
		WriteResponse(writer, SearchResponse{ Success: true, Lang: "ge", Items: items })
	}

	
	
}

func main() {
    fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)

	http.HandleFunc("/search", search)

	err := http.ListenAndServe(":8080", nil)

	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server closed\n")
	} else if err != nil {
		fmt.Printf("error starting server: %s\n", err)
		os.Exit(1)
	}	
}
