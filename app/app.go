package app

import (
	"crypto/md5"
	"encoding/hex"
	_ "errors" // we would need this package
	_ "fmt"  // we would need this package
	"sync"
	_ "time" // we would need this package
	_ "strconv" // we would need this package
)

// IncomingAddBook : here you tell us what IncomingAddBook is
type IncomingAddBook struct {
	Title    string `json:"title"`
	Author   string `json:"author"`
	Category string `json:"category"`
	Total    int    `json:"total"`
}

// IncomingAddMovie : here you tell us what IncomingAddMovie is
type IncomingAddMovie struct {
	Total    int    `json:"total"`
	Author   string `json:"author"`
	Category string `json:"category"`
	Title    string `json:"title"`
}

// Library : struct global
type Library struct {
	Book
	Movie
	BookDB
	MovieDB
	BookCopies	map[string]int
	MovieCopies	map[string]int
}

type item interface{}

// Book : Book struct
type Book struct {
	ID       string
	title	 string
	author	 string
	category string
	total	 int
}

// Movie : Movie struct
type Movie struct {
	title	string
	genre	[]string   //["Drama", "Romance"]
	ID 		string
}


// Response : Response ack to add movie/book
type Response struct {
	ID      string
	Success bool
	Message string
}

// BookDB : BookDB
type BookDB struct {
	bookDBMap map[string]*Book
	mutex         sync.RWMutex
}

// MovieDB : MovieDB
type MovieDB struct {
	movieDBMap    map[string]*Movie
	mutex         sync.RWMutex
}


// NewLibrary : Constructor of Library struct
func NewLibrary() *Library {

	var bookDBMap = make(map[string]*Book)
	var movieDBMap = make(map[string]*Movie)

	return &Library{
		BookDB: BookDB{
			bookDBMap: bookDBMap,
		},
		MovieDB: MovieDB{
			movieDBMap: movieDBMap,
		},
		Book: Book{},
		Movie: Movie{},
	}
}

// NewBook db :
func (bdb *BookDB) addBookDB(ID string) {

	bdb.mutex.Lock()
	bdb.bookDBMap[ID] = &Book{ID: ID}
	bdb.mutex.Unlock()

}

// This function receives an string and generates a Unique ID
func generateHash(title string,author string) string {

	s := title + author
	bs := md5.New()
	bs.Write([]byte(s))
	hash1 := hex.EncodeToString(bs.Sum(nil)[:3])

	return hash1

}


// AddBook : AddBook
func (l *Library) AddBook(title string, author string, category string, total  int) Response {


	// if does not exists, you can add it

	ID := generateHash(title,author)

	response := Response{
		ID:      ID,
		Success: true,
		Message: "",
	}

	return response

}




