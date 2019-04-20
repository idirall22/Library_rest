package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

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
	lastID := 0
	err := db.QueryRowContext(r.Context(),
		"INSERT INTO books (title, author, year) VALUES($1, $2, $3) RETURNING id",
		&book.Title, &book.Author, &book.Year).Scan(&lastID)
	if err != nil {
		log.Fatalf("error with ExecContext : %s", err)
		return
	}
	book.ID = int(lastID)
	json.NewEncoder(w).Encode(&book)
}
func updateBook(w http.ResponseWriter, r *http.Request) {

}
func deleteBook(w http.ResponseWriter, r *http.Request) {

}

func connectDB() {
	password := os.Getenv("PASSWORD")
	// password := "password"
	connStr := fmt.Sprintf("user=postgres password=%s dbname=library sslmode=disable", password)

	database, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Could not connect to database: %s\n", err)
	}
	db = database
}
