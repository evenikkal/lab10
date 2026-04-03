// Package model содержит доменные модели приложения.
package model

// Book представляет книгу в библиотеке.
type Book struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Author string `json:"author"`
	Year   int    `json:"year"`
}
