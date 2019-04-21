package main

import (
	"database/sql"
	"log"
	"net/http"

	"./controllers"
	"./driver"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/subosito/gotenv"
)

var db *sql.DB

func init() {
	err := gotenv.Load()
	if err != nil {
		log.Fatalf("Can not get env variables: %s\n", err)
	}
	db = driver.ConnectDB()
	con := controllers.Controller{DB: db}
	controllers.DefaultController = con
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/books", controllers.GetBooks).Methods("GET")
	router.HandleFunc("/books/{id:[0-9]+}", controllers.GetBook).Methods("GET")
	router.HandleFunc("/books", controllers.CreateBook).Methods("POST")
	router.HandleFunc("/books/{id:[0-9]+}", controllers.UpdateBook).Methods("PUT")
	router.HandleFunc("/books/{id:[0-9]+}", controllers.DeleteBook).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8000", router))
}
