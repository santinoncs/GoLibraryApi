Library

Let's create an API for a library with books and movies. Let's build the following API methods:

1. POST /book/add - add a new book to the library.
Request:
{
  "title": "Book Title",
  "author": "Author Name",
  "category": "fiction",
  "total": 1               // number of physical copies of this book that will be available in the library
}

Response:
{
  "success": true,          // true on success, false on error
  "id": 1,                  // book id on success, 0 on error
  "message": ""             // empty string on success, error message if the request is invalid
}


2. POST /movie/add - add a new movie to the library
Request:
{
  "title": "Movie Title",
  "genre": ["Drama", "Romance"],  // array of genres, may be several, one, or none (empty array)
  "total": 2                      // number of physical copies of this movie (assume those are DVDs)
}

Response:
{
  "success": true,                // true on success, false on error
  "id": 1,                        // movie id on success, 0 on error
  "message": ""                   // empty string on success, error message if the request is invalid
}

3. POST /book/rent - rent a book by the user
If the book with the given id does not exist or is not available right now (all copies are already rented out), the corresponding error message should be returned.
Each user may rent no more than 2 items at the same time (either books or movies), so if the user already has 2 items currently in rent, he/she cannot rent more.

Request:
{
  "id": 1,                        // book id that the user wants to rent
  "user_id": 12345                // user id who rents the book
}

Response:
{
  "success": true,                // true on success, false on error
  "message": ""                   // empty string on success, error message if error occurred
}

4. POST /movie/rent - rent a movie by the user
If the movie with the given id does not exist or is not available right now (all copies are already rented out), the corresponding error message should be returned.
Each user may rent no more than 2 items at the same time (either books or movies), so if the user already has 2 items currently in rent, he/she cannot rent more.
Request:
{
  "id": 1,                        // movie id that the user wants to rent
  "user_id": 12345                // user id who rents the movie
}

Response:
{
  "success": true,                // true on success, false on error
  "message": ""                   // empty string on success, error message if error occurred
}

5. POST /book/return - return the book by the user
Request:
{
  "id": 1,                        // book id that the user wants to return
  "user_id": 12345                // user id who returns the book
}

Response:
{
  "success": true,                // true on success, false on error
  "message": ""                   // empty string on success, error message if error occurred
}

6. POST /movie/return - return the movie by the user
Request:
{
  "id": 1,                        // book id that the user wants to return
  "user_id": 12345                // user id who returns the book
}

Response:
{
  "success": true,                // true on success, false on error
  "message": ""                   // empty string on success, error message if error occurred
}

7. GET /book/?id=<book_id> - return book information
Response format:
{
  "success": true,                // true on success, false on error
  "message": ""                   // empty string on success, error message if the request is invalid
  "book": {                       // book json on success, empty json on error
    "title": "Book Title",
    "author": "Author Name",
    "category": "fiction",
    "total": 2 // total number of physical copies of this book in the library
    "available": 1 // the number of physical copies of the book that are currently not rented out
  }
}

8. GET /movie/?id=<movie_id> - return movie information
Response format:
{
  "success": true,                  // true on success, false on error
  "message": ""                     // empty string on success, error message if the request is invalid
  "movie": {                        // movie json on success, empty json on error
    "title": "Movie Title",
    "genre": ["Drama", "Romance"],  // array of genres, may be several, one, or none (empty array)
    "total": 2                      // number of physical copies of this movie (assume those are DVDs)
    "available": 1                  // the number of physical copies of the book that are currently not rented out
  }
}






curl \
--request POST \
--header "Content-Type: application/json" \
--data '
{
   "title": "any",
   "author":  "joe doe",
   "category": "scifi",
   "total": 1
}
' http://localhost:8080/book/add