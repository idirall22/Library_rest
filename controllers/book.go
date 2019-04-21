package controllers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"../models"
	"../repository/book"
	"github.com/gorilla/mux"
)

//Controller to get default database
type Controller struct {
	DB *sql.DB
}

//DefaultController var
var DefaultController Controller

//GetBooks handler
func GetBooks(w http.ResponseWriter, r *http.Request) {
	var books []models.Book
	resp := bookRepository.GETBookDB(r.Context(), DefaultController.DB, &books)
	if resp == false {
		fmt.Fprintf(w, "There are no books")
		return
	}
	json.NewEncoder(w).Encode(&books)

}

//GetBook model
func GetBook(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	bookID, _ := strconv.Atoi(params["id"])
	var book models.Book
	bookRepository.GetSingleBookDB(r.Context(), DefaultController.DB, &book, bookID)
	json.NewEncoder(w).Encode(&book)
}

//CreateBook model
func CreateBook(w http.ResponseWriter, r *http.Request) {
	var book models.Book
	if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
		fmt.Fprintf(w, "Can not create the book : %s", err)
		return
	}
	defer r.Body.Close()
	lastID := 0
	bookRepository.CreateBookDB(r.Context(), DefaultController.DB, &book, &lastID)
	book.ID = int(lastID)
	json.NewEncoder(w).Encode(&book)
}

//UpdateBook model
func UpdateBook(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	bookID, errID := strconv.Atoi(params["id"])
	if errID != nil {
		fmt.Fprintf(w, "the id must be an integer")
		return
	}
	var book models.Book
	if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
		fmt.Fprintf(w, "Can not Update the book : %s", err)
		return
	}
	defer r.Body.Close()
	res := bookRepository.UpdateBookDB(r.Context(), DefaultController.DB, &book, bookID)
	if res == false {
		fmt.Fprintf(w, "Can not Update the book")
	}
	json.NewEncoder(w).Encode(&book)
}

//DeleteBook model
func DeleteBook(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	bookID, errID := strconv.Atoi(params["id"])
	if errID != nil {
		fmt.Fprintf(w, "the id must be an integer")
		return
	}
	res := bookRepository.DeleteBookDB(r.Context(), DefaultController.DB, bookID)
	if !res {
		fmt.Fprintf(w, "Can not Delete the book")
	}
	fmt.Fprintf(w, "Book deleted")
}
