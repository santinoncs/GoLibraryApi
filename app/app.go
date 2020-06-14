package app

import (
	_ "errors" // we would need this package
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
	BookDB []Book
	ItemDB
	MovieDB []Movie
	UserDB
	autoinc uint64
}

// Book : Book struct
type Book struct {
	Item
	Author   string
	Category string
}

// RentableItem : here you tell us what RentableItem is
type RentableItem interface {
	getItemID() uint64
	getItemTitle() string
	getItemCopies() uint64
}

// Movie : Movie struct
type Movie struct {
	genre []string //["Drama", "Romance"]
	Item
}

// Item : here you tell us what Item is
type Item struct {
	ID     uint64
	Title  string
	Copies uint64
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

// ResponseInfoBook : ResponseInfoBook
type ResponseInfoBook struct {
	Success bool
	Message string
	Book
}

// ResponseInfoMovie : ResponseInfoMovie
type ResponseInfoMovie struct {
	Success bool
	Message string
	Movie
}

// BookDB : BookDB
type BookDB struct {
	bookDBMap map[uint64]*Book
	mutex     sync.RWMutex
}

// ItemDB : ItemDB
type ItemDB struct {
	itemDBMap map[uint64]RentableItem
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

func (i Item) getItemID() uint64 {
	return i.ID
}

func (i Item) getItemTitle() string {
	return i.Title
}

func (i Item) getItemCopies() uint64 {
	return i.Copies
}

// NewLibrary : Constructor of Library struct
func NewLibrary() *Library {

	var itemDBMap = make(map[uint64]RentableItem)
	var userDBMap = make(map[int]*[]Item)

	var BookDB = []Book{}
	var MovieDB = []Movie{}

	return &Library{
		BookDB: BookDB,
		ItemDB: ItemDB{
			itemDBMap: itemDBMap,
		},
		MovieDB: MovieDB,
		UserDB: UserDB{
			userDBMap: userDBMap,
		},
		Book:  Book{Item: Item{}},
		Movie: Movie{Item: Item{}},
	}
}

// addItemDB db :
func (idb *ItemDB) addItemDB(ID uint64, title string, copies uint64) {

	if _, ok := idb.itemDBMap[ID]; ok {

		idb.mutex.Lock()
		copiesCurr := idb.itemDBMap[ID].getItemCopies()
		copies = copiesCurr + copies
		idb.itemDBMap[ID] = Item{ID: ID, Title: title, Copies: copies}
		idb.mutex.Unlock()

	} else {
		idb.mutex.Lock()
		idb.itemDBMap[ID] = Item{ID: ID, Title: title, Copies: copies}
		idb.mutex.Unlock()
	}

	fmt.Printf("Method addItemDB. Item Added: %v\n", idb.itemDBMap[ID])

}

// addItemDB db :
func (idb *ItemDB) removeItemDB(ID uint64, title string, copies uint64) {

	if _, ok := idb.itemDBMap[ID]; ok {
		idb.mutex.Lock()

		copiesCurr := idb.itemDBMap[ID].getItemCopies()

		copies = copiesCurr - copies
		idb.itemDBMap[ID] = Item{ID: ID, Title: title, Copies: copies}
		idb.mutex.Unlock()
	} else {

		idb.mutex.Lock()
		idb.itemDBMap[ID] = Item{ID: ID, Title: title, Copies: copies}
		idb.mutex.Unlock()

	}

	fmt.Printf("Method removeItemDB. Item removed: %v\n", idb.itemDBMap[ID])

}

// Initialize userdbmap
func (udb *UserDB) init(userid int) {

	if _, ok := udb.userDBMap[userid]; ok {
		fmt.Printf("userdb is present in map\n")
	} else {
		fmt.Printf("Initialize userdb\n")
		udb.userDBMap = map[int]*[]Item{
			userid: {},
		}
	}
}

// FindUserDB : find item in the userdbMap slice
func (udb *UserDB) FindUserDB(it *Item, userid int) bool {

	arr := *udb.userDBMap[userid]

	for _, item := range arr {
		if item == *it {
			return true
		}
	}

	return false

}

// addUserDB : addUserDB
func (udb *UserDB) addUserDB(it *Item, userid int) {

	it.Copies = 0

	udb.mutex.Lock()
	*udb.userDBMap[userid] = append(*udb.userDBMap[userid], *it)
	udb.mutex.Unlock()

	fmt.Printf("Method addUserDB. This is the item %v for userid: %d after adding\n", *udb.userDBMap[userid], userid)

}

// removeUserDB : removeUserDB
func (udb *UserDB) removeUserDB(it *Item, userid int) {

	udb.mutex.Lock()

	arr := *udb.userDBMap[userid]

	i := 0 // output index
	for _, x := range *udb.userDBMap[userid] {
		if x == *it {
			fmt.Println("Match with an element")
		} else {
			arr[i] = x
			i++
		}
	}

	arr = arr[:i]

	*udb.userDBMap[userid] = arr

	udb.mutex.Unlock()

	fmt.Printf("Method removeUserDB. This is the item %v for userid: %d after removal\n", *udb.userDBMap[userid], userid)

}

// AddBook :
func (l *Library) AddBook(title string, author string, category string, total uint64) ResponseAdd {

	for _, book := range l.BookDB {
		if title == book.Title {
			return ResponseAdd{
				ID:      book.ID,
				Success: false,
				Message: "Already exists",
			}
		}
	}

	atomic.AddUint64(&l.autoinc, 1)

	ID := l.autoinc

	l.ItemDB.addItemDB(ID, title, total)

	l.Book.Author = author
	l.Book.Category = category
	l.Book.Item.ID = ID
	l.Book.Item.Title = title
	l.Book.Item.Copies = total

	l.BookDB = append(l.BookDB, l.Book)

	response := ResponseAdd{
		ID:      ID,
		Success: true,
		Message: "",
	}

	return response

}

// AddMovie : AddMovie
func (l *Library) AddMovie(title string, genre []string, total uint64) ResponseAdd {

	for _, movie := range l.MovieDB {
		if title == movie.Title {
			return ResponseAdd{
				ID:      movie.ID,
				Success: false,
				Message: "Already exists",
			}
		}
	}

	atomic.AddUint64(&l.autoinc, 1)

	ID := l.autoinc

	l.ItemDB.addItemDB(ID, title, total)

	l.Movie.genre = genre
	l.Movie.Item.ID = ID
	l.Movie.Item.Title = title

	l.MovieDB = append(l.MovieDB, l.Movie)

	response := ResponseAdd{
		ID:      ID,
		Success: true,
		Message: "",
	}

	return response

}

// Rent : Rent
func (l *Library) Rent(ID uint64, userid int) ResponseRent {


	l.UserDB.init(userid)

	if b, ok := l.ItemDB.itemDBMap[ID]; ok {
		fmt.Println("You can rent it")

		if len(*l.UserDB.userDBMap[userid]) >= 2 {
			return ResponseRent{
				Success: false,
				Message: "You cannot rent more than 2 copies",
			}
		}

		item := Item{}

		item.ID = b.getItemID()
		item.Title = b.getItemTitle()
		item.Copies = b.getItemCopies()

		l.UserDB.addUserDB(&item, userid)

		if (l.ItemDB.itemDBMap[ID].getItemCopies()) > 0 {
			fmt.Printf("Exists, and there are more than 0 copies\n")

			l.ItemDB.removeItemDB(ID, item.Title, 1)

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
			Message: "Does no exists in the item db",
		}
	}

	return ResponseRent{
		Success: true,
		Message: "",
	}

}

// Return : Return
func (l *Library) Return(ID uint64, userid int) ResponseRent {

	if b, ok := l.ItemDB.itemDBMap[ID]; ok {

		item := Item{}

		item.ID = b.getItemID()
		item.Title = b.getItemTitle()
		item.Copies = b.getItemCopies()

		l.ItemDB.addItemDB(ID, item.Title, 1)

		found := l.UserDB.FindUserDB(&item, userid)
		if !found {
			fmt.Println("Value not found in slice")
			return ResponseRent{
				Success: false,
				Message: "You cannot return as it is not rented",
			}
		}

		l.UserDB.removeUserDB(&item, userid)

	} else {
		return ResponseRent{
			Success: false,
			Message: "Does not exists in itemDB",
		}
	}

	return ResponseRent{
		Success: true,
		Message: "",
	}

}

// BookInfo : BookInfo
func (l *Library) BookInfo(bookid uint64) ResponseInfoBook {

	for _, n := range l.BookDB {
		if n.ID == bookid {
			return ResponseInfoBook{
				Success: true,
				Message: "",
				Book:    Book{Item: Item{ID: bookid, Title: l.ItemDB.itemDBMap[bookid].getItemTitle(), Copies: l.ItemDB.itemDBMap[bookid].getItemCopies()}, Author: n.Author, Category: n.Category},
			}
		}
	}

	return ResponseInfoBook{
		Success: false,
		Message: "Error",
		Book:    Book{},
	}

}

// MovieInfo : MovieInfo
func (l *Library) MovieInfo(movieid uint64) ResponseInfoMovie {

	for _, n := range l.MovieDB {
		if n.ID == movieid {
			return ResponseInfoMovie{
				Success: true,
				Message: "",
				Movie:   Movie{Item: Item{ID: movieid, Title: l.ItemDB.itemDBMap[movieid].getItemTitle(), Copies: l.ItemDB.itemDBMap[movieid].getItemCopies()}, genre: n.genre},
			}
		}
	}

	return ResponseInfoMovie{
		Success: false,
		Message: "Error",
		Movie:   Movie{},
	}

}
