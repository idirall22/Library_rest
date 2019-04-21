package bookRepository

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"../../models"
)

//GETBookDB sql
func GETBookDB(ctx context.Context, db *sql.DB, books *[]models.Book) bool {
	query := "SELECT * FROM books"
	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		log.Fatal(err)
	}
	for rows.Next() {
		var book models.Book
		err := rows.Scan(&book.ID, &book.Title, &book.Author, &book.Year)
		switch {
		case err == sql.ErrNoRows:
			return false
		}
		if err != nil {
			log.Fatal(err)
		}
		*books = append(*books, book)
	}
	return true
}

//GetSingleBookDB model
func GetSingleBookDB(ctx context.Context, db *sql.DB,
	book *models.Book, bookID int) bool {
	query := fmt.Sprintf("SELECT * FROM books WHERE id=%v", bookID)
	err := db.QueryRowContext(ctx, query).Scan(
		&book.ID, &book.Title, &book.Author, &book.Year)
	if err != nil {
		switch {
		case err == sql.ErrNoRows:
			return false
		default:
			log.Fatalf("There is an error with QueryRowContext in GetSingleBookDB : %s\n", err)
			return false
		}
	}
	return true
}

func checkIfBookExist(ctx context.Context, db *sql.DB, bookID int) bool {
	bookExist := false
	row := db.QueryRowContext(ctx,
		"SELECT EXISTS (SELECT * FROM books WHERE id=$1)", bookID)
	if errE := row.Scan(&bookExist); errE != nil || bookExist == false {
		switch {
		case errE == sql.ErrNoRows || !bookExist:
			return false
		default:
			log.Fatalf("error with QueryRowContext in checkIfBookExist: %s", errE)
			return false
		}
	}
	return true
}

//CreateBookDB model
func CreateBookDB(ctx context.Context, db *sql.DB, book *models.Book, lastID *int) bool {
	err := db.QueryRowContext(ctx,
		"INSERT INTO books (title, author, year) VALUES($1, $2, $3) RETURNING id",
		&book.Title, &book.Author, &book.Year).Scan(&lastID)
	if err != nil {
		log.Fatalf("error with QueryRowContext : %s", err)
		return false
	}
	return true
}

//UpdateBookDB model
func UpdateBookDB(ctx context.Context, db *sql.DB, book *models.Book, bookID int) bool {
	bookExist := checkIfBookExist(ctx, db, bookID)
	if !bookExist {
		return false
	}
	_, err := db.ExecContext(ctx,
		"UPDATE books SET title=$1, author=$2, year=$3",
		&book.Title, &book.Author, &book.Year)
	if err != nil {
		log.Fatalf("error with ExecContext UpdateBookDB : %s", err)
		return false
	}
	return true
}

//DeleteBookDB model
func DeleteBookDB(ctx context.Context, db *sql.DB, bookID int) bool {
	bookExist := checkIfBookExist(ctx, db, bookID)
	if !bookExist {
		return false
	}
	_, err := db.ExecContext(ctx,
		"DELETE FROM books WHERE id=$1", bookID)
	if err != nil {
		log.Fatalf("error with ExecContext delete book : %s", err)
		return false
	}
	return true
}
