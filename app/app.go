package app

import (
	"errors" // we would need this package
	"fmt"
	"sync"
	"sync/atomic" // we could need this
	_ "time"      // we would need this package
)

// IncomingAddBook : here you tell us what IncomingAddBook is
type IncomingAddBook struct {
	Title    string `json:"title"`
	Author   string `json:"author"`
	Category string `json:"category"`
	Total    uint64 `json:"total"`
}

// IncomingAddMovie : here you tell us what IncomingAddMovie is
type IncomingAddMovie struct {
	Title string   `json:"title"`
	Genre []string `json:"genre"`
	Total uint64   `json:"total"`
}

// IncomingRent : here you tell us what IncomingRent is
type IncomingRent struct {
	ID     uint64 `json:"id"`
	UserID int    `json:"user_id"`
}

// Library : struct global
type Library struct {
	Book
	Movie
	BookDB
	MovieDB
	UserDB
	User
	autoinc uint64
}

// Book : Book struct
type Book struct {
	Item
	Author   string
	Category string
	copies   uint64
}

// Movie : Movie struct
type Movie struct {
	genre []string //["Drama", "Romance"]
	Item
	copies uint64
}

// User : User struct
type User struct {
	userID int
}

// itemBehave : here you tell us what Item is
type itemBehave interface {
	RentItem()
	ReturnItem()
}

// Item : here you tell us what Item is
type Item struct {
	ID    uint64
	Title string
	itemBehave
}

// ResponseAdd : ResponseAdd ack to add movie/book
type ResponseAdd struct {
	ID      uint64
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
	bookDBMap map[uint64]*Book
	mutex     sync.RWMutex
}

// MovieDB : MovieDB
type MovieDB struct {
	movieDBMap map[uint64]*Movie
	mutex      sync.RWMutex
}

// UserDB : UserDB
type UserDB struct {
	userDBMap map[int]*[]Item
	mutex     sync.RWMutex
}

// NewLibrary : Constructor of Library struct
func NewLibrary() *Library {

	var bookDBMap = make(map[uint64]*Book)
	var movieDBMap = make(map[uint64]*Movie)
	var userDBMap = make(map[int]*[]Item)

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
		Book:  Book{Item: Item{}},
		Movie: Movie{Item: Item{}},
		User:  User{},
	}
}

// NewBook db :
func (bdb *BookDB) addBookDB(ID uint64, title string, author string, category string, total uint64) {

	bdb.mutex.Lock()
	bdb.bookDBMap[ID] = &Book{Item: Item{ID: ID, Title: title}, Author: author, Category: category}
	bdb.mutex.Unlock()

	atomic.AddUint64(&bdb.bookDBMap[ID].copies, total)

	fmt.Printf("this is the book added: %v\n", bdb.bookDBMap[ID])

}

// getBookDB : getBookDB
func (bdb *BookDB) getBookDB(ID uint64) (Book, error) {

	if _, ok := bdb.bookDBMap[ID]; ok {
		bdb.mutex.Lock()

		bookinfo := bdb.bookDBMap[ID]
		bdb.mutex.Unlock()
		return *bookinfo, nil
	}
	return Book{}, errors.New("Item does not exist")

}

// NewMovie db :
func (mdb *MovieDB) addMovieDB(ID uint64, title string, genre []string, total uint64) {

	mdb.mutex.Lock()
	mdb.movieDBMap[ID] = &Movie{Item: Item{ID: ID, Title: title}, genre: genre}
	mdb.mutex.Unlock()

	atomic.AddUint64(&mdb.movieDBMap[ID].copies, total)

}

// addUserDB : addUserDB
func (udb *UserDB) addUserDB(it *Item, userid int) {

	fmt.Printf("Adding the item %v:\n", *it)

	udb.mutex.Lock()
	*udb.userDBMap[userid] = append(*udb.userDBMap[userid], *it)
	udb.mutex.Unlock()

	fmt.Printf("this is the item %v for userid: %d\n", *udb.userDBMap[userid], userid)

}

// removeUserDB : removeUserDB
func (udb *UserDB) removeUserDB(it *Item, userid int) {

	fmt.Printf("Removing in removeuserdb the item %v:\n", *it)

	udb.mutex.Lock()

	arr := *udb.userDBMap[userid]

	i := 0 // output index
	for _, x := range *udb.userDBMap[userid] {
		if x == *it {
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

// AddBook : This could be a method implementing an interface -> additem
func (l *Library) AddBook(title string, author string, category string, total uint64) ResponseAdd {

	atomic.AddUint64(&l.autoinc, 1)

	ID := l.autoinc

	l.BookDB.addBookDB(ID, title, author, category, total)

	response := ResponseAdd{
		ID:      ID,
		Success: true,
		Message: "",
	}

	return response

}

// AddMovie : AddMovie
func (l *Library) AddMovie(title string, genre []string, total uint64) ResponseAdd {

	atomic.AddUint64(&l.autoinc, 1)

	ID := l.autoinc

	l.MovieDB.addMovieDB(ID, title, genre, total)

	return ResponseAdd{
		ID:      ID,
		Success: true,
		Message: "",
	}

}

// RentBook : RentBook
func (l *Library) RentBook(ID uint64, userid int) ResponseRent {

	// if the userid does not exists, first initialize !!!
	// ask Oleg here about the initialization of map

	if _, ok := l.UserDB.userDBMap[userid]; ok {
		fmt.Printf("userdb is present in map\n")
	} else {
		fmt.Printf("Initialize userdb\n")
		l.UserDB.userDBMap = map[int]*[]Item{
			userid: {},
		}
	}

	if b, ok := l.BookDB.bookDBMap[ID]; ok {

		fmt.Printf("Exists, you can rent it. With ID: %d and Title: %s\n", b.ID, b.Title)

		fmt.Printf("the actual slice lenght is: %d\n", len(*l.UserDB.userDBMap[userid]))

		if len(*l.UserDB.userDBMap[userid]) >= 2 {
			return ResponseRent{
				Success: false,
				Message: "You cannot rent more than 2 copies",
			}
		}

		item := Item{}

		item.ID = b.ID
		item.Title = b.Title

		l.UserDB.addUserDB(&item, userid)

		if l.BookDB.bookDBMap[ID].copies > 0 {
			fmt.Printf("Exists, and there are more than 0 copies\n")
			atomic.AddUint64(&l.BookDB.bookDBMap[ID].copies, ^uint64(0))

		} else {
			fmt.Printf("0 copies\n")
			return ResponseRent{
				Success: false,
				Message: "Not enough copies",
			}
		}

	} else {
		return ResponseRent{
			Success: false,
			Message: "Does no exists in the book db",
		}
	}

	fmt.Printf("this is the book rented: %v\n", l.BookDB.bookDBMap[ID])

	return ResponseRent{
		Success: true,
		Message: "",
	}

}

// RentMovie : RentMovie
func (l *Library) RentMovie(ID uint64, userid int) ResponseRent {

	if _, ok := l.UserDB.userDBMap[userid]; ok {
		fmt.Printf("userdb is present in map\n")
	} else {
		l.UserDB.userDBMap = map[int]*[]Item{
			userid: {},
		}
	}

	if m, ok := l.MovieDB.movieDBMap[ID]; ok {

		item := Item{}

		item.ID = m.ID
		item.Title = m.Title

		l.UserDB.addUserDB(&item, userid)

		// Decrement

		atomic.AddUint64(&l.MovieDB.movieDBMap[ID].copies, ^uint64(0))

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
func (l *Library) ReturnBook(ID uint64, userid int) ResponseRent {

	if b, ok := l.BookDB.bookDBMap[ID]; ok {

		item := Item{}

		item.ID = b.ID
		item.Title = b.Title

		l.UserDB.removeUserDB(&item, userid)

		atomic.AddUint64(&l.BookDB.bookDBMap[ID].copies, 1)

		fmt.Printf("this is the book for return: %v\n", l.BookDB.bookDBMap[ID])

	} else {
		return ResponseRent{
			Success: false,
			Message: "Does not exists in bookDB",
		}
	}

	return ResponseRent{
		Success: true,
		Message: "",
	}

}

// BookInfo : BookInfo
func (l *Library) BookInfo(bookid uint64) ResponseInfo {

	bookinfo, err := l.BookDB.getBookDB(bookid)
	if err != nil {
		return ResponseInfo{
			Success: false,
			Message: "Error",
			Book:    Book{},
		}
	}

	return ResponseInfo{
		Success: true,
		Message: "",
		Book:    bookinfo,
	}

}
