package app

import (
	"crypto/md5"
	"encoding/hex"
	"errors" // we would need this package
	"fmt"
	"sync"
	_ "time" // we would need this package
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

// IncomingRent : here you tell us what IncomingRent is
type IncomingRent struct {
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
	BookCopies		// this was map[string]*int before..but it didnt work	
	MovieCopies	map[string]int
	User
}


// Book : Book struct
type Book struct {
	ID       string
	Title	 string
	Author	 string
	Category string
	Total	 int
}

// Movie : Movie struct
type Movie struct {
	title	string
	genre	[]string   //["Drama", "Romance"]
	ID 		string
	total	 int
}

// User : User struct
type User struct {
	userID	int
}

// Item : here you tell us what Item is
type Item interface {
	ItemType() string
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

// ResponseInfo : ResponseInfo
type ResponseInfo struct {
	Success bool
	Message string
	Book
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
	userDBMap     map[int]*[]Item
	mutex         sync.RWMutex
}

// BookCopies : BookCopies
type BookCopies struct{
	BookCopiesMap map[string]int
	mutex         sync.RWMutex
}

// NewLibrary : Constructor of Library struct
func NewLibrary() *Library {


	var bookDBMap = make(map[string]*Book)
	var movieDBMap = make(map[string]*Movie)
	var userDBMap = make(map[int]*[]Item)

	var BookCopiesMap = make(map[string]int)
	var MovieCopies= make(map[string]int)
	
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
		BookCopies: BookCopies{
			BookCopiesMap: BookCopiesMap,
		},
		MovieCopies: MovieCopies,
	}
}

// ItemType : Returns item type
func (b Book) ItemType() string {
    return "book"
}

// ItemType : Returns item type
func (m Movie) ItemType() string {
    return "movie"
}

// incrementBookCopies :  incrementBookCopies
func (bc *BookCopies) incrementBookCopies(ID string, total int) {

	bc.mutex.Lock()
	bc.BookCopiesMap[ID] += total
	bc.mutex.Unlock()

}

// decrementBookCopies :  decrementBookCopies
func (bc *BookCopies) decrementBookCopies(ID string) {

	bc.mutex.Lock()
	bc.BookCopiesMap[ID] --
	bc.mutex.Unlock()

}

// NewBook db :
func (bdb *BookDB) addBookDB(ID string,title string, author string, category string, total  int) {

	bdb.mutex.Lock()
	bdb.bookDBMap[ID] = &Book{ID: ID,Title: title, Author: author,Category: category,Total: total}
	bdb.mutex.Unlock()

}

// getBookDB : getBookDB
func (bdb *BookDB) getBookDB(ID string) (Book, error) {

	if _, ok := bdb.bookDBMap[ID]; ok {	
		bdb.mutex.RLock()
		
		bookinfo := bdb.bookDBMap[ID] 
		bdb.mutex.Unlock()
		return *bookinfo,nil
	}
	return Book{}, errors.New("Item does not exist")

}

// NewMovie db :
func (mdb *MovieDB) addMovieDB(ID string, title string, genre []string, total int) {

	mdb.mutex.Lock()
	mdb.movieDBMap[ID] = &Movie{ID: ID,title: title, genre: genre, total: total}
	mdb.mutex.Unlock()

}

// addUserDB : addUserDB
func (udb *UserDB) addUserDB(it Item, userid int) {

	udb.mutex.Lock()
	*udb.userDBMap[userid] = append(*udb.userDBMap[userid], it)
	udb.mutex.Unlock()

}

// addUserDB : addUserDB
func (udb *UserDB) removeUserDB(it Item, userid int) {


	udb.mutex.Lock()

	arr := *udb.userDBMap[userid]

	i := 0 // output index
	for _, x := range *udb.userDBMap[userid] {
		if x == it {
			// copy and increment index
			//arr[i] = x
			//i++
		} else {
			// son diferentes, copia el elemento en el array
			arr[i] = x
			i++
		}
	}


	arr = arr[:i]

	*udb.userDBMap[userid] = arr

	udb.mutex.Unlock()

}


// This function receives an string and generates a Unique ID
func generateHash(title string) string {

	s := title
	bs := md5.New()
	bs.Write([]byte(s))
	hash1 := hex.EncodeToString(bs.Sum(nil)[:3])

	return hash1

}

// AddBook : This could be a method implementing an interface -> additem
func (l *Library) AddBook(title string, author string, category string, total  int) ResponseAdd {

	ID := generateHash(title)

	l.BookDB.addBookDB(ID,title,author,category,total)

	l.BookCopies.incrementBookCopies(ID,total)

	response := ResponseAdd{
		ID:      ID,
		Success: true,
		Message: "",
	}

	return response

}


// AddMovie : AddMovie
func (l *Library) AddMovie(title string, genre []string, total  int) ResponseAdd {


	ID := generateHash(title)

	l.MovieDB.addMovieDB(ID,title, genre, total)

	l.MovieCopies[ID] += total 


	return ResponseAdd{
		ID:      ID,
		Success: true,
		Message: "",
	}

}


// RentBook : RentBook
func (l *Library) RentBook(ID string, userid int) ResponseRent {

	// if the userid does not exists, first initialize !!! 
	// ask Oleg here about the initialization of map

	if _, ok := l.UserDB.userDBMap[userid]; ok {
		fmt.Printf("userdb is present in map\n")
	} else {
		l.UserDB.userDBMap = map[int]*[]Item{
			userid: {},
		}
	}

	if b,ok := l.BookDB.bookDBMap[ID]; ok {

		l.UserDB.addUserDB(*b,userid)

		if l.BookCopies.BookCopiesMap[ID] > 0 {
			l.decrementBookCopies(ID)
		} else {
			return ResponseRent{
				Success: false,
				Message: "Error",
			}
		}



	} else {
		return ResponseRent{
			Success: false,
			Message: "Error",
		}
	}


	return ResponseRent{
		Success: true,
		Message: "",
	}

}

// RentMovie : RentMovie
func (l *Library) RentMovie(ID string, userid int) ResponseRent {


	if _, ok := l.UserDB.userDBMap[userid]; ok {
		fmt.Printf("userdb is present in map\n")
	} else {
		l.UserDB.userDBMap = map[int]*[]Item{
			userid: {},
		}
	}

	if m, ok := l.MovieDB.movieDBMap[ID]; ok {

		l.UserDB.addUserDB(*m,userid)
		l.MovieCopies[ID] --

	} else {
		return ResponseRent{
			Success: false,
			Message: "Error",
		}
	}

	return ResponseRent{
		Success: true,
		Message: "",
	}

}

// ReturnBook : ReturnBook
func (l *Library) ReturnBook(ID string, userid int) ResponseRent {

	if b,ok := l.BookDB.bookDBMap[ID]; ok {


		l.UserDB.removeUserDB(*b,userid)

		l.incrementBookCopies(ID,1)


	} else {
		return ResponseRent{
			Success: false,
			Message: "Error",
		}
	}
	
	return ResponseRent{
		Success: true,
		Message: "",
	}

}

// BookInfo : BookInfo
func (l *Library) BookInfo(bookid string) ResponseInfo {


	bookinfo,err := l.BookDB.getBookDB(bookid)
	if err != nil {
		return ResponseInfo{
			Success: false,
			Message: "Error",
			Book: Book{},
		}
	}


	return ResponseInfo{
		Success: true,
		Message: "",
		Book: bookinfo,
	}

}