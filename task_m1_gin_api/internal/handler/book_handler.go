// Package handler содержит HTTP-обработчики Gin.
package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"lab10/task_m1_gin_api/internal/model"
	"lab10/task_m1_gin_api/internal/repository"
)

// BookHandler обрабатывает HTTP-запросы для работы с книгами.
type BookHandler struct {
	repo *repository.BookRepository
}

// NewBookHandler создаёт новый BookHandler.
func NewBookHandler(repo *repository.BookRepository) *BookHandler {
	return &BookHandler{repo: repo}
}

// GetBooks godoc
// @Summary Получить все книги
// @GET /books
func (h *BookHandler) GetBooks(c *gin.Context) {
	books := h.repo.FindAll()
	c.JSON(http.StatusOK, gin.H{
		"count": len(books),
		"data":  books,
	})
}

// GetBook godoc
// @Summary Получить книгу по ID
// @GET /books/:id
func (h *BookHandler) GetBook(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	book, err := h.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "book not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusOK, book)
}

// CreateBook godoc
// @Summary Создать новую книгу
// @POST /books
func (h *BookHandler) CreateBook(c *gin.Context) {
	var book model.Book
	if err := c.ShouldBindJSON(&book); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	created := h.repo.Create(book)
	c.JSON(http.StatusCreated, created)
}
