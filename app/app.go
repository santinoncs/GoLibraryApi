package app

import (
	"crypto/md5"
	"encoding/hex"
	_ "errors" // we would need this package
	"fmt"
	_ "fmt"   // we would need this package
	_ "strconv" // we would need this package
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
	BookCopies	map[string]int
	MovieCopies	map[string]int
	User
}


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

// NewLibrary : Constructor of Library struct
func NewLibrary() *Library {

	var bookDBMap = make(map[string]*Book)
	var movieDBMap = make(map[string]*Movie)
	var userDBMap = make(map[int]*[]Item)

	var BookCopies = make(map[string]int)
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
		BookCopies: BookCopies,
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

// NewBook db :
func (bdb *BookDB) addBookDB(ID string,title string, author string, category string, total  int) {

	bdb.mutex.Lock()
	bdb.bookDBMap[ID] = &Book{ID: ID,title: title, author: author,category: category,total: total}
	bdb.mutex.Unlock()

}

// NewMovie db :
func (mdb *MovieDB) addMovieDB(ID string, title string, genre []string, total int) {

	fmt.Println("antes de add movie db",  mdb.movieDBMap[ID])


	mdb.mutex.Lock()
	mdb.movieDBMap[ID] = &Movie{ID: ID,title: title, genre: genre, total: total}
	mdb.mutex.Unlock()

	fmt.Println("despues de add movie db",  mdb.movieDBMap[ID])


}


// addUserDB : addUserDB
func (udb *UserDB) addUserDB(it Item, userid int) {

	fmt.Println("antes de setear", *udb.userDBMap[userid])

	udb.mutex.Lock()
	*udb.userDBMap[userid] = append(*udb.userDBMap[userid], it)
	udb.mutex.Unlock()

	fmt.Println("despues de seetear",*udb.userDBMap[userid])

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

	// if does not exists, you can add it

	ID := generateHash(title)

	l.BookDB.addBookDB(ID,title,author,category,total)

	// Add book copies to BookCopies map

	fmt.Printf("Total book copies before adding %s is : %d\n", ID, l.BookCopies[ID])


	l.BookCopies[ID] += total 

	fmt.Printf("Total book copies after adding %s is : %d\n", ID, l.BookCopies[ID])


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

	l.MovieDB.addMovieDB(ID,title, genre, total)

	fmt.Printf("Total movie copies before adding %s is : %d\n", ID, l.MovieCopies[ID])


	l.MovieCopies[ID] += total 

	fmt.Printf("Total movie copies after adding %s is : %d\n", ID, l.MovieCopies[ID])


	return ResponseAdd{
		ID:      ID,
		Success: true,
		Message: "",
	}

}


// RentBook : RentBook
func (l *Library) RentBook(ID string, userid int) ResponseRent {

	// if the userid does not exists, first initialize

	if _, ok := l.UserDB.userDBMap[userid]; ok {
		fmt.Printf("userdb is present in map\n")
	} else {
		fmt.Printf("userdb is NOT present in map.\n")
		l.UserDB.userDBMap = map[int]*[]Item{
			userid: {},
		}
		fmt.Printf("userdb is JUST present in map.\n")
	}

	if b,ok := l.BookDB.bookDBMap[ID]; ok {
		fmt.Println("este es la movie a rentar", b)

		l.UserDB.addUserDB(*b,userid)

		fmt.Printf("Total book copies before removing %s is : %d\n", ID, l.BookCopies[ID])

		if l.BookCopies[ID] > 0 {
			l.BookCopies[ID] --
		} else {
			fmt.Printf("There is no more copies of this ID: %s", ID)
			return ResponseRent{
				Success: false,
				Message: "Error",
			}
		}

		fmt.Printf("Total book copies after removing %s is : %d\n", ID, l.BookCopies[ID])


	} else {
		fmt.Println("This id does not exists")
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
		fmt.Printf("userdb is NOT present in map.\n")
		l.UserDB.userDBMap = map[int]*[]Item{
			userid: {},
		}
	}

	if m, ok := l.MovieDB.movieDBMap[ID]; ok {
		fmt.Println("este es la movie a rentar", m)

		l.UserDB.addUserDB(*m,userid)
		
		fmt.Printf("Total movie copies before removing %s is : %d\n", ID, l.MovieCopies[ID])


		l.MovieCopies[ID] --

		fmt.Printf("Total movie copies after adding %s is : %d\n", ID, l.MovieCopies[ID])

	} else {
		fmt.Println("This id does not exists")
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


