package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

//Book Model
type Book struct {
	ID     int    `json:"id,omitempty"`
	Title  string `json:"title,omitempty"`
	Author string `json:"author,omitempty"`
	Year   uint8  `json:"year,omitempty"`
}

//BookList Model
type BookList struct {
	Books []Book
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/books/", getBooks)
	router.HandleFunc("/books/{int:id}", getBook).Methods("GET")
	router.HandleFunc("/books/{int:id}", createBook).Methods("POST")
	router.HandleFunc("/books/{int:id}", updateBook).Methods("PUT")
	router.HandleFunc("/books/{int:id}", deleteBook).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8000", router))
}
func getBooks(w http.ResponseWriter, r *http.Request) {
	fmt.Print("books")
}
func getBook(w http.ResponseWriter, r *http.Request) {
	fmt.Print("book")
}
func createBook(w http.ResponseWriter, r *http.Request) {
	fmt.Print("create")
}
func updateBook(w http.ResponseWriter, r *http.Request) {
	fmt.Print("update")
}
func deleteBook(w http.ResponseWriter, r *http.Request) {
	fmt.Print("delete")
}
