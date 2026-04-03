package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func doPost(router http.Handler, path string, body interface{}) *httptest.ResponseRecorder {
	data, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, path, bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	return w
}

// ---- /users ----

func TestCreateUser_Valid(t *testing.T) {
	r := SetupRouter()
	w := doPost(r, "/users", User{Name: "Alice", Email: "alice@example.com", Age: 25})

	assert.Equal(t, http.StatusCreated, w.Code)
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, "user created", resp["message"])
}

func TestCreateUser_MissingName(t *testing.T) {
	r := SetupRouter()
	w := doPost(r, "/users", map[string]interface{}{
		"email": "alice@example.com",
		"age":   25,
	})

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCreateUser_NameTooShort(t *testing.T) {
	r := SetupRouter()
	w := doPost(r, "/users", User{Name: "A", Email: "a@example.com", Age: 20})

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCreateUser_InvalidEmail(t *testing.T) {
	r := SetupRouter()
	w := doPost(r, "/users", User{Name: "Alice", Email: "not-an-email", Age: 25})

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCreateUser_AgeTooYoung(t *testing.T) {
	r := SetupRouter()
	w := doPost(r, "/users", User{Name: "Alice", Email: "alice@example.com", Age: 15})

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// ---- /products ----

func TestCreateProduct_Valid(t *testing.T) {
	r := SetupRouter()
	w := doPost(r, "/products", Product{
		Title:    "Laptop",
		Price:    999.99,
		Quantity: 10,
		Category: "electronics",
	})

	assert.Equal(t, http.StatusCreated, w.Code)
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, "product created", resp["message"])
}

func TestCreateProduct_InvalidCategory(t *testing.T) {
	r := SetupRouter()
	w := doPost(r, "/products", Product{
		Title:    "Laptop",
		Price:    999.99,
		Quantity: 10,
		Category: "weapons",
	})

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCreateProduct_ZeroPrice(t *testing.T) {
	r := SetupRouter()
	w := doPost(r, "/products", map[string]interface{}{
		"title":    "Free Thing",
		"price":    0,
		"quantity": 5,
		"category": "other",
	})

	assert.Equal(t, http.StatusBadRequest, w.Code)
}
