package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

//Book Model
type Book struct {
	ID     int    `json:"id,omitempty"`
	Title  string `json:"title,omitempty"`
	Author string `json:"author,omitempty"`
	Year   uint   `json:"year,omitempty"`
}

var books = []Book{
	{ID: 1, Title: "Start python", Author: "Idir", Year: 2016},
	{ID: 2, Title: "Golang From Scratch", Author: "Idir", Year: 2017},
	{ID: 3, Title: "TypeScript VS JavaScript", Author: "Idir", Year: 2017},
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
	json.NewEncoder(w).Encode(&books)
}
func getBook(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	bookID, err := strconv.Atoi(params["id"])
	if err != nil {
		log.Fatal(err)
	}
	var book Book
	for _, b := range books {
		if b.ID == bookID {
			book = b
		}
	}
	json.NewEncoder(w).Encode(&book)
}
func createBook(w http.ResponseWriter, r *http.Request) {
	var book Book
	defer r.Body.Close()
	err := json.NewDecoder(r.Body).Decode(&book)
	if err != nil {
		fmt.Fprint(w, "Could not create the book")
		return
	}
	books = append(books, book)
	json.NewEncoder(w).Encode(&book)
}
func updateBook(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	defer r.Body.Close()
	bookID, err := strconv.Atoi(params["id"])
	if err != nil {
		fmt.Fprint(w, "book id must be an integer")
		return
	}
	var book Book
	errD := json.NewDecoder(r.Body).Decode(&book)
	book.ID = bookID
	if errD != nil {
		fmt.Fprint(w, "Could not update the book")
		return
	}
	for i, item := range books {
		if item.ID == bookID {
			books[i] = book
		}
	}
	json.NewEncoder(w).Encode(&book)
}
func deleteBook(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	defer r.Body.Close()
	bookID, err := strconv.Atoi(params["id"])
	if err != nil {
		fmt.Fprint(w, "book id must be an integer")
		return
	}
	for i, item := range books {
		if item.ID == bookID {
			books = append(books[0:i], books[i+1:]...)
		}
	}
	json.NewEncoder(w).Encode(&books)
}
