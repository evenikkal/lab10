// Package main — точка входа Go Gin API сервера.
package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"lab10/task_m1_gin_api/internal/handler"
	"lab10/task_m1_gin_api/internal/repository"
)

const port = ":8080"

func main() {
	repo := repository.NewBookRepository()
	h := handler.NewBookHandler(repo)

	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	books := r.Group("/books")
	{
		books.GET("", h.GetBooks)
		books.GET("/:id", h.GetBook)
		books.POST("", h.CreateBook)
	}

	fmt.Printf("Books API starting on %s\n", port)
	fmt.Println("GET  /books      — list all books")
	fmt.Println("GET  /books/:id  — get book by id")
	fmt.Println("POST /books      — create book")
	log.Fatal(r.Run(port))
}
