package models

//Book model
type Book struct {
	ID     int    `json:"id,omitempty"`
	Title  string `json:"title,omitempty"`
	Author string `json:"author,omitempty"`
	Year   uint   `json:"year,omitempty"`
}
