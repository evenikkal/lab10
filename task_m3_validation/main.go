package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// User represents a user registration request.
type User struct {
	Name  string `json:"name"  binding:"required,min=2,max=50"`
	Email string `json:"email" binding:"required,email"`
	Age   int    `json:"age"   binding:"required,min=18,max=100"`
}

// Product represents a product creation request with nested validation.
type Product struct {
	Title    string  `json:"title"    binding:"required,min=3,max=100"`
	Price    float64 `json:"price"    binding:"required,gt=0"`
	Quantity int     `json:"quantity" binding:"required,min=0"`
	Category string  `json:"category" binding:"required,oneof=electronics food clothing other"`
}

func SetupRouter() *gin.Engine {
	r := gin.Default()
	r.POST("/users", createUserHandler)
	r.POST("/products", createProductHandler)
	return r
}

func createUserHandler(c *gin.Context) {
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "validation failed",
			"details": err.Error(),
		})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"message": "user created",
		"user":    user,
	})
}

func createProductHandler(c *gin.Context) {
	var product Product
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "validation failed",
			"details": err.Error(),
		})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"message": "product created",
		"product": product,
	})
}

func main() {
	r := SetupRouter()
	r.Run(":8081")
}
