package driver

import (
	"database/sql"
	"fmt"
	"log"
	"os"
)

//ConnectDB allows to to connect to a database
func ConnectDB() *sql.DB {
	password := os.Getenv("PASSWORD")
	connStr := fmt.Sprintf("user=postgres password=%s dbname=library sslmode=disable", password)

	database, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Could not connect to database: %s\n", err)
	}
	return database
}
