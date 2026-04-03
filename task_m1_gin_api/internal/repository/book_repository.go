// Package repository содержит слой хранения данных.
package repository

import (
	"errors"
	"sync"

	"lab10/task_m1_gin_api/internal/model"
)

// ErrNotFound возвращается когда книга не найдена.
var ErrNotFound = errors.New("book not found")

// BookRepository хранит книги в памяти.
type BookRepository struct {
	mu      sync.RWMutex
	books   map[int]model.Book
	counter int
}

// NewBookRepository создаёт репозиторий с тестовыми данными.
func NewBookRepository() *BookRepository {
	r := &BookRepository{books: make(map[int]model.Book)}
	r.counter = 1
	r.books[1] = model.Book{ID: 1, Title: "The Go Programming Language", Author: "Donovan & Kernighan", Year: 2015}
	r.books[2] = model.Book{ID: 2, Title: "Clean Code", Author: "Robert Martin", Year: 2008}
	r.counter = 3
	return r
}

// FindAll возвращает все книги.
func (r *BookRepository) FindAll() []model.Book {
	r.mu.RLock()
	defer r.mu.RUnlock()
	books := make([]model.Book, 0, len(r.books))
	for _, b := range r.books {
		books = append(books, b)
	}
	return books
}

// FindByID возвращает книгу по ID или ErrNotFound.
func (r *BookRepository) FindByID(id int) (model.Book, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	b, ok := r.books[id]
	if !ok {
		return model.Book{}, ErrNotFound
	}
	return b, nil
}

// Create добавляет новую книгу и возвращает её с присвоенным ID.
func (r *BookRepository) Create(b model.Book) model.Book {
	r.mu.Lock()
	defer r.mu.Unlock()
	b.ID = r.counter
	r.counter++
	r.books[b.ID] = b
	return b
}
