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
	Word string `json:"word"`
	Translation string `json:"translation"`
}

type SearchResponse struct {
    Success bool `json:"success"`
    Message string `json:"message,omitempty"`
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

func search(writer http.ResponseWriter, r *http.Request) {
	var items []SearchItem
	var word, translation string
	hasQuery := r.URL.Query().Has("query")

	if (hasQuery) {
		query := r.URL.Query().Get("query")

		const file string = "dict.db"
		db, err := sql.Open("sqlite3", file)

		if err != nil {
			WriteResponse(writer, SearchResponse{ Success: false, Message: "Database opening error"})
	  		return
		}
		
		defer db.Close()		

		var sql = fmt.Sprintf("SELECT ru, ge  FROM words WHERE ru LIKE '%s%%' LIMIT 10", query)
		fmt.Printf("Query: %s\n", sql)

		rows, err := db.Query(sql)

		if err != nil {
			WriteResponse(writer, SearchResponse{ Success: false, Message: "Database query error"})
	  		return
		}


		for rows.Next() {
			err = rows.Scan(&word, &translation)
			item := SearchItem { word, translation }
			items = append(items, item)
		}
	}

	WriteResponse(writer, SearchResponse{ Success: true, Items: items })
	
}

func main() {
	http.HandleFunc("/", getRoot)
	http.HandleFunc("/search", search)

	err := http.ListenAndServe(":8080", nil)

	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server closed\n")
	} else if err != nil {
		fmt.Printf("error starting server: %s\n", err)
		os.Exit(1)
	}
}
