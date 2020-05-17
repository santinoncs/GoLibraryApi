package main

import (
	"fmt"
	app "github.com/santinoncs/LibraryApi/app"
	"encoding/json"
	"net/http"
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
	var incomingrentbook app.IncomingRentBook
	var responseAdd	app.ResponseAdd
	var responseRent app.ResponseRent

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

		err := json.NewDecoder(r.Body).Decode(&incomingrentbook)
		if err != nil {
			fmt.Println(err)
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.WriteHeader(http.StatusBadRequest)
			return
		}


		responseRent = library.RentBook(incomingrentbook.ID,incomingrentbook.UserID)
		responseJSON, _ := json.Marshal(responseRent)
		fmt.Fprintf(w, "Response: %s\n", responseJSON)

	}


}
