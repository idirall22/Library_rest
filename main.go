package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/subosito/gotenv"
)

//Book Model
type Book struct {
	ID     int    `json:"id,omitempty"`
	Title  string `json:"title,omitempty"`
	Author string `json:"author,omitempty"`
	Year   uint   `json:"year,omitempty"`
}

var db *sql.DB

func init() {
	err := gotenv.Load()
	if err != nil {
		log.Fatalf("Can not get env variables: %s\n", err)
	}
	connectDB()
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/books", getBooks).Methods("GET")
	router.HandleFunc("/books/{id:[0-9]+}", getBook).Methods("GET")
	router.HandleFunc("/books", createBook).Methods("POST")
	router.HandleFunc("/books/{id:[0-9]+}", updateBook).Methods("PUT")
	router.HandleFunc("/books/{id:[0-9]+}", deleteBook).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8000", router))
}
func getBooks(w http.ResponseWriter, r *http.Request) {
	var books []Book
	query := "SELECT * FROM books"
	rows, err := db.QueryContext(r.Context(), query)
	if err != nil {
		log.Fatal(err)
	}
	for rows.Next() {
		var book Book
		err := rows.Scan(&book.ID, &book.Title, &book.Author, &book.Year)
		if err != nil {
			log.Fatal(err)
		}
		books = append(books, book)
	}
	json.NewEncoder(w).Encode(&books)
}
func getBook(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	bookID := params["id"]

	query := fmt.Sprintf("SELECT * FROM books WHERE id=%s", bookID)
	var book Book
	err := db.QueryRowContext(r.Context(), query).Scan(
		&book.ID, &book.Title, &book.Author, &book.Year)
	if err != nil {
		switch {
		case err == sql.ErrNoRows:
			fmt.Fprintf(w, "There is not a book with id : %s", bookID)
			return
		default:
			log.Printf("There is an error with QueryRowContext : %s\n", err)
			return
		}
	}
	json.NewEncoder(w).Encode(&book)
}
func createBook(w http.ResponseWriter, r *http.Request) {
	var book Book
	if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
		fmt.Fprintf(w, "Can not create the book : %s", err)
		return
	}
	defer r.Body.Close()
	lastID := 0
	err := db.QueryRowContext(r.Context(),
		"INSERT INTO books (title, author, year) VALUES($1, $2, $3) RETURNING id",
		&book.Title, &book.Author, &book.Year).Scan(&lastID)
	if err != nil {
		log.Fatalf("error with QueryRowContext : %s", err)
		return
	}
	book.ID = int(lastID)
	json.NewEncoder(w).Encode(&book)
}
func updateBook(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	bookID, errID := strconv.Atoi(params["id"])
	if errID != nil {
		fmt.Fprintf(w, "the id must be an integer")
		return
	}
	var book Book
	if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
		fmt.Fprintf(w, "Can not Update the book : %s", err)
		return
	}
	defer r.Body.Close()
	//Check if book exist
	var bookExist bool
	row := db.QueryRowContext(r.Context(),
		"SELECT EXISTS (SELECT * FROM books WHERE id=$1)", bookID)
	if errE := row.Scan(&bookExist); errE != nil || bookExist == false {
		switch {
		case errE == sql.ErrNoRows || !bookExist:
			fmt.Fprintf(w, "Book Does not Existe")
			return
		default:
			log.Fatalf("error with QueryRowContext : %s", errE)
			return
		}
	}
	fmt.Println(bookExist)
	_, err := db.ExecContext(r.Context(),
		"UPDATE books SET title=$1, author=$2, year=$3",
		&book.Title, &book.Author, &book.Year)
	if err != nil {
		log.Fatalf("error with ExecContext : %s", err)
		return
	}
	json.NewEncoder(w).Encode(&book)
}
func deleteBook(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	bookID, errID := strconv.Atoi(params["id"])
	if errID != nil {
		fmt.Fprintf(w, "the id must be an integer")
		return
	}
	//check if book exist
	var bookExist bool
	row := db.QueryRowContext(r.Context(),
		"SELECT EXISTS (SELECT * FROM books WHERE id=$1)", bookID)
	if errE := row.Scan(&bookExist); errE != nil || bookExist == false {
		switch {
		case errE == sql.ErrNoRows || !bookExist:
			fmt.Fprintf(w, "Book Does not Existe")
			return
		default:
			log.Fatalf("error with QueryRowContext Check if book exist : %s", errE)
			return
		}
	}
	//Delete book
	_, err := db.ExecContext(r.Context(), "DELETE FROM books WHERE id=$1", bookID)
	if err != nil {
		log.Fatalf("error with ExecContext delete book : %s", err)
		return
	}
	fmt.Fprintf(w, "Book deleted")
}

func connectDB() {
	password := os.Getenv("PASSWORD")
	connStr := fmt.Sprintf("user=postgres password=%s dbname=library sslmode=disable", password)

	database, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Could not connect to database: %s\n", err)
	}
	db = database
}
