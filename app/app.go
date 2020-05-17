package app

import (
	"crypto/md5"
	"encoding/hex"
	_ "errors" // we would need this package
	"fmt"
	_ "fmt"   // we would need this package
	"strconv" // we would need this package
	"sync"
	"time" // we would need this package
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
	Title    string   `json:"title"`
	Genre    []string   `json:"genre"`
	Total    int      `json:"total"`
}

// IncomingRentBook : here you tell us what IncomingRentBook is
type IncomingRentBook struct {
	ID       string `json:"id"`
	UserID   int `json:"user_id"`
}

// Library : struct global
type Library struct {
	Book
	Movie
	BookDB
	MovieDB
	UserDB
	BookCopies	map[string]int
	MovieCopies	map[string]int
	User
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

// User : User struct
type User struct {
	userID	int
}

// ResponseAdd : ResponseAdd ack to add movie/book
type ResponseAdd struct {
	ID      string
	Success bool
	Message string
}

// ResponseRent : ResponseRent 
type ResponseRent struct {
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

// UserDB : UserDB
type UserDB struct {
	userDBMap     map[int]*[]string
	mutex         sync.RWMutex
}

// NewLibrary : Constructor of Library struct
func NewLibrary() *Library {

	

	var bookDBMap = make(map[string]*Book)
	var movieDBMap = make(map[string]*Movie)
	var userDBMap = make(map[int]*[]string)

	

	return &Library{
		BookDB: BookDB{
			bookDBMap: bookDBMap,
		},
		MovieDB: MovieDB{
			movieDBMap: movieDBMap,
		},
		UserDB: UserDB{
			userDBMap: userDBMap,
		},
		Book: Book{},
		Movie: Movie{},
		User: User{},
	}
}

// NewBook db :
func (bdb *BookDB) addBookDB(ID string) {

	bdb.mutex.Lock()
	bdb.bookDBMap[ID] = &Book{ID: ID}
	bdb.mutex.Unlock()

}

// NewMovie db :
func (mdb *MovieDB) addMovieDB(ID string) {

	mdb.mutex.Lock()
	mdb.movieDBMap[ID] = &Movie{ID: ID}
	mdb.mutex.Unlock()

}

// addUserDB : addUserDB
func (udb *UserDB) addUserDB(ID string, userid int) {

	//var s []string

	fmt.Println("antes de setear", *udb.userDBMap[userid])

	udb.mutex.Lock()
	*udb.userDBMap[userid] = append(*udb.userDBMap[userid], ID)
	udb.mutex.Unlock()

	fmt.Println("despues de seetear",*udb.userDBMap[userid])

}


// This function receives an string and generates a Unique ID
func generateHash(question string) string {

	now := time.Now().UnixNano()
	t := strconv.FormatInt(now, 10)
	s := question + t
	bs := md5.New()
	bs.Write([]byte(s))
	hash1 := hex.EncodeToString(bs.Sum(nil)[:3])

	return hash1

}

// AddBook : This could be a method implementing an interface -> additem
func (l *Library) AddBook(title string, author string, category string, total  int) ResponseAdd {

	// if does not exists, you can add it

	ID := generateHash(title)

	l.BookDB.addBookDB(ID)

	response := ResponseAdd{
		ID:      ID,
		Success: true,
		Message: "",
	}

	return response

}


// AddMovie : AddMovie
func (l *Library) AddMovie(title string, genre []string, total  int) ResponseAdd {


	// if does not exists, you can add it

	ID := generateHash(title)

	l.MovieDB.addMovieDB(ID)

	response := ResponseAdd{
		ID:      ID,
		Success: true,
		Message: "",
	}

	return response

}

// RentBook : RentBook
func (l *Library) RentBook(ID string, userid int) ResponseRent {

	response := ResponseRent{
		Success: true,
		Message: "",
	}

	// if the userid does not exists, first initialize

	if _, ok := l.UserDB.userDBMap[userid]; ok {
		fmt.Printf("userdb is present in map")
	} else {
		fmt.Printf("userdb is NOT present in map.")
		l.UserDB.userDBMap = map[int]*[]string{
			userid: {},
		}
	}

	// add book to user id array 

	l.UserDB.addUserDB(ID,userid)


	return response

}


