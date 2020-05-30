package app

import (
	"errors" // we would need this package
	"fmt"
	"sync"
	_ "time" // we would need this package
	"sync/atomic" // we could need this 
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
	ID       uint64  `json:"id"`
	UserID   int     `json:"user_id"`
}

// Library : struct global
type Library struct {
	Book
	Movie
	BookDB
	MovieDB
	UserDB
	BookCopies		// this was map[string]*int before..but it didnt work	
	MovieCopies	map[uint64]int
	User
	ops			uint64
}


// Book : Book struct
type Book struct {
	Item
	Author	 string
	Category string
	Total	 int
}

// Movie : Movie struct
type Movie struct {
	genre	[]string   //["Drama", "Romance"]
	Item
	total	 int
}

// User : User struct
type User struct {
	userID	int
}


// Item : here you tell us what Item is
type Item struct {
	ID       uint64
	Title	 string
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
	bookDBMap     map[uint64]*Book
	mutex         sync.RWMutex
}

// MovieDB : MovieDB
type MovieDB struct {
	movieDBMap    map[uint64]*Movie
	mutex         sync.RWMutex
}

// UserDB : UserDB
type UserDB struct {
	userDBMap     map[int]*[]Item
	mutex         sync.RWMutex
}

// BookCopies : BookCopies
type BookCopies struct{
	BookCopiesMap map[uint64]int
	mutex         sync.RWMutex
}

// NewLibrary : Constructor of Library struct
func NewLibrary() *Library {


	var bookDBMap = make(map[uint64]*Book)
	var movieDBMap = make(map[uint64]*Movie)
	var userDBMap = make(map[int]*[]Item)

	var BookCopiesMap = make(map[uint64]int)
	var MovieCopies= make(map[uint64]int)
	
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
		Book: Book{Item: Item{},},
		Movie: Movie{Item: Item{},},
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
func (bc *BookCopies) incrementBookCopies(ID uint64, total int) {

	bc.mutex.Lock()
	bc.BookCopiesMap[ID] += total
	bc.mutex.Unlock()

}

// decrementBookCopies :  decrementBookCopies
func (bc *BookCopies) decrementBookCopies(ID uint64) {

	bc.mutex.Lock()
	bc.BookCopiesMap[ID] --
	bc.mutex.Unlock()

}

// NewBook db :
func (bdb *BookDB) addBookDB(ID uint64,title string, author string, category string, total  int) {

	bdb.mutex.Lock()
	bdb.bookDBMap[ID] = &Book{Item: Item{ID: ID,Title: title}, Author: author,Category: category,Total: total}
	bdb.mutex.Unlock()

}

// getBookDB : getBookDB
func (bdb *BookDB) getBookDB(ID uint64) (Book, error) {

	if _, ok := bdb.bookDBMap[ID]; ok {	
		bdb.mutex.RLock()
		
		bookinfo := bdb.bookDBMap[ID] 
		bdb.mutex.Unlock()
		return *bookinfo,nil
	}
	return Book{}, errors.New("Item does not exist")

}

// NewMovie db :
func (mdb *MovieDB) addMovieDB(ID uint64, title string, genre []string, total int) {

	mdb.mutex.Lock()
	mdb.movieDBMap[ID] = &Movie{Item: Item{ID: ID,Title: title}, genre: genre, total: total}
	mdb.mutex.Unlock()

}

// addUserDB : addUserDB
func (udb *UserDB) addUserDB(it *Item, userid int) {

	udb.mutex.Lock()
	*udb.userDBMap[userid] = append(*udb.userDBMap[userid], *it)
	udb.mutex.Unlock()

}

// removeUserDB : removeUserDB
func (udb *UserDB) removeUserDB(it *Item, userid int) {


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


// This function receives an string and generates a Unique ID
func generateAutoIncrement(title string) uint64 {

	var ops uint64

	atomic.AddUint64(&ops, 1)

	return ops

}

// AddBook : This could be a method implementing an interface -> additem
func (l *Library) AddBook(title string, author string, category string, total  int) ResponseAdd {


	atomic.AddUint64(&l.ops, 1)

	ID := l.ops

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


	ID := generateAutoIncrement(title)

	l.MovieDB.addMovieDB(ID,title, genre, total)

	l.MovieCopies[ID] += total 


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
		l.UserDB.userDBMap = map[int]*[]Item{
			userid: {},
		}
	}

	if b,ok := l.BookDB.bookDBMap[ID]; ok {

		item := Item{}

		b.ID = item.ID
		b.Title = item.Title

		l.UserDB.addUserDB(&item,userid)

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

		m.ID = item.ID
		m.Title = item.Title

		l.UserDB.addUserDB(&item,userid)

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
func (l *Library) ReturnBook(ID uint64, userid int) ResponseRent {

	if b,ok := l.BookDB.bookDBMap[ID]; ok {

		item := Item{}

		b.ID = item.ID
		b.Title = item.Title
		

		l.UserDB.removeUserDB(&item,userid)

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
func (l *Library) BookInfo(bookid uint64) ResponseInfo {


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