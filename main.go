package main

import (
	"encoding/json"
	"fmt"
	app "github.com/santinoncs/LibraryApi/app"
	"net/http"
	"net/url"
	"strconv"
)

var library *app.Library

func main() {

	library = app.NewLibrary()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		handler(w, r)
	})
	err := http.ListenAndServe(":8080", nil)

	fmt.Println(err)

}

func handler(w http.ResponseWriter, r *http.Request) {

	var incomingaddbook app.IncomingAddBook
	var incomingaddmovie app.IncomingAddMovie
	var incomingrent app.IncomingRent
	var responseAdd app.ResponseAdd
	var responseRent app.ResponseRent
	var responseInfoBook app.ResponseInfoBook

	if r.URL.Path == "/book/add" {

		err := json.NewDecoder(r.Body).Decode(&incomingaddbook)
		if err != nil {
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		responseAdd = library.AddBook(incomingaddbook.Title, incomingaddbook.Author, incomingaddbook.Category, incomingaddbook.Total)
		responseJSON, _ := json.Marshal(responseAdd)
		fmt.Fprintf(w, "Response: %s\n", responseJSON)

	}

	if r.URL.Path == "/movie/add" {

		err := json.NewDecoder(r.Body).Decode(&incomingaddmovie)
		if err != nil {
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		responseAdd = library.AddMovie(incomingaddmovie.Title, incomingaddmovie.Genre, incomingaddmovie.Total)
		responseJSON, _ := json.Marshal(responseAdd)
		fmt.Fprintf(w, "Response: %s\n", responseJSON)

	}

	if r.URL.Path == "/book/rent" {

		err := json.NewDecoder(r.Body).Decode(&incomingrent)
		if err != nil {
			fmt.Println(err)
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		responseRent = library.Rent(incomingrent.ID, incomingrent.UserID)
		responseJSON, _ := json.Marshal(responseRent)
		fmt.Fprintf(w, "Response: %s\n", responseJSON)

	}

	if r.URL.Path == "/movie/rent" {

		err := json.NewDecoder(r.Body).Decode(&incomingrent)
		if err != nil {
			fmt.Println(err)
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		responseRent = library.Rent(incomingrent.ID, incomingrent.UserID)
		responseJSON, _ := json.Marshal(responseRent)
		fmt.Fprintf(w, "Response: %s\n", responseJSON)

	}

	if r.URL.Path == "/book/return" {

		err := json.NewDecoder(r.Body).Decode(&incomingrent)
		if err != nil {
			fmt.Println(err)
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		responseRent = library.Return(incomingrent.ID, incomingrent.UserID)
		responseJSON, _ := json.Marshal(responseRent)
		fmt.Fprintf(w, "Response: %s\n", responseJSON)

	}

	m, _ := url.ParseQuery(r.URL.RawQuery)

	if r.URL.Path == "/book" {
		bookid := m["id"]

		n, _ := strconv.ParseUint(bookid[0], 10, 64)

		fmt.Println(n)

		responseInfoBook = library.BookInfo(n)
		responseJSON, _ := json.Marshal(responseInfoBook)
		fmt.Fprintf(w, "Response: %s\n", responseJSON)

	}

}
