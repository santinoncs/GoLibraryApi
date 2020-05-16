package main

import (
	"fmt"
	app "github.com/santinoncs/LibraryApi/app"
	//"sync"
	"encoding/json"
	//"log"
	"net/http"
)

var library *app.Library

func main() {

//	application = app.NewApp()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		handler(w, r)
	})
	err := http.ListenAndServe(":8080", nil)

	fmt.Println(err)

}

func handler(w http.ResponseWriter, r *http.Request) {

	var content app.IncomingAddBook
	var responseAdd	app.Response

	if r.URL.Path == "/book/add" {

		err := json.NewDecoder(r.Body).Decode(&content)
		if err != nil {
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		responseAdd = library.AddBook(content.Title, content.Author, content.Category, content.Total)
		responseJSON, _ := json.Marshal(responseAdd)
		fmt.Fprintf(w, "Response: %s\n", responseJSON)

	}


}
