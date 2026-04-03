package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	"lab10/task_m1_gin_api/internal/handler"
	"lab10/task_m1_gin_api/internal/model"
	"lab10/task_m1_gin_api/internal/repository"
)

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	repo := repository.NewBookRepository()
	h := handler.NewBookHandler(repo)

	r := gin.New()
	r.GET("/books", h.GetBooks)
	r.GET("/books/:id", h.GetBook)
	r.POST("/books", h.CreateBook)
	return r
}

func TestGetBooks_ReturnsOK(t *testing.T) {
	r := setupRouter()
	req := httptest.NewRequest(http.MethodGet, "/books", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}
}

func TestGetBooks_ReturnsBooks(t *testing.T) {
	r := setupRouter()
	req := httptest.NewRequest(http.MethodGet, "/books", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	var resp map[string]any
	json.NewDecoder(rec.Body).Decode(&resp)

	count := resp["count"].(float64)
	if count < 1 {
		t.Errorf("expected at least 1 book, got %v", count)
	}
}

func TestGetBook_ValidID(t *testing.T) {
	r := setupRouter()
	req := httptest.NewRequest(http.MethodGet, "/books/1", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}

	var book model.Book
	json.NewDecoder(rec.Body).Decode(&book)
	if book.ID != 1 {
		t.Errorf("expected book id 1, got %d", book.ID)
	}
}

func TestGetBook_NotFound(t *testing.T) {
	r := setupRouter()
	req := httptest.NewRequest(http.MethodGet, "/books/999", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", rec.Code)
	}
}

func TestGetBook_InvalidID(t *testing.T) {
	r := setupRouter()
	req := httptest.NewRequest(http.MethodGet, "/books/abc", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rec.Code)
	}
}

func TestCreateBook(t *testing.T) {
	r := setupRouter()
	body := `{"title":"Test Book","author":"Test Author","year":2024}`
	req := httptest.NewRequest(http.MethodPost, "/books", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Errorf("expected 201, got %d", rec.Code)
	}

	var book model.Book
	json.NewDecoder(rec.Body).Decode(&book)
	if book.Title != "Test Book" {
		t.Errorf("expected title 'Test Book', got '%s'", book.Title)
	}
	if book.ID == 0 {
		t.Error("expected non-zero ID")
	}
}
