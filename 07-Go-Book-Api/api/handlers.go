package api

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	dsn := os.Getenv("DB_URL")
	if dsn == "" {
		log.Fatal("DB_URL is empty")
	}

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database : ", err)
	}

	//migrate the schema
	if err := DB.AutoMigrate(&Book{}); err != nil {
		log.Fatal("Failed to migrate schema : ", err)
	}

	if err := DB.AutoMigrate(&Book{}); err != nil {
		log.Fatal("Failed to migrate schema : ", err)
	}
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func CreateBook(c *gin.Context) {
	var payload Book

	//bind the request body
	if err := c.ShouldBindJSON(&payload); err != nil {
		ResponseJSON(c, http.StatusBadRequest, "Invalid Input", nil)
		return
	}

	if payload.Title == "" || payload.Author == "" || payload.Year == 0 {
		ResponseJSON(c, http.StatusBadRequest, "Title, Author, and Year required", nil)
		return
	}

	if err := DB.Create(&payload).Error; err != nil {
		ResponseJSON(c, http.StatusInternalServerError, "Failed to create book", nil)
		return
	}
	ResponseJSON(c, http.StatusCreated, "Book created sucessfully", payload)
}

// getting list of books
func GetBooks(c *gin.Context) {
	var books []Book
	if err := DB.Find(&books).Error; err != nil {
		ResponseJSON(c, http.StatusInternalServerError, "Failed to detch books", nil)
		return
	}
	ResponseJSON(c, http.StatusOK, "Books retrieved successfully", books)
}

// get a single book
func GetBook(c *gin.Context) {
	var book Book
	if err := DB.First(&book, c.Param("id")).Error; err != nil {
		ResponseJSON(c, http.StatusNotFound, "Book Not Found", nil)
		return
	}
	ResponseJSON(c, http.StatusOK, "Book retrieved succesfully", book)
}

// Update a book
func UpdateBook(c *gin.Context) {
	var existing Book
	if err := DB.First(&existing, c.Param("id")).Error; err != nil {
		ResponseJSON(c, http.StatusNotFound, "Book not found", nil)
		return
	}

	var payload Book
	//bind the request body
	if err := c.ShouldBindJSON(&payload); err != nil {
		ResponseJSON(c, http.StatusBadRequest, "Invalid input", nil)
		return
	}

	existing.Title = payload.Title
	existing.Author = payload.Author
	existing.Year = payload.Year

	if err := DB.Save(&existing).Error; err != nil {
		ResponseJSON(c, http.StatusInternalServerError, "Failed to update book", nil)
	}
	ResponseJSON(c, http.StatusOK, "Book updated successfully", existing)
}

// delete a book
func DeleteBook(c *gin.Context) {
	res := DB.Delete(&Book{}, c.Param("id"))
	if res.Error != nil {
		ResponseJSON(c, http.StatusInternalServerError, "Failed to delete book ", nil)
	}

	if res.RowsAffected == 0 {
		ResponseJSON(c, http.StatusNotFound, "Book not found", nil)
		return
	}
	ResponseJSON(c, http.StatusOK, "Book deleted successfully", nil)
}

func GenerateJWT(c *gin.Context) {
	var loginRequest LoginRequest
	if err := c.ShouldBindJSON(&loginRequest); err != nil {
		ResponseJSON(c, http.StatusBadRequest, "Invalid request payload", nil)
		return
	}

	if loginRequest.Username != "admin" || loginRequest.Password != "Password" {
		ResponseJSON(c, http.StatusUnauthorized, "Invalid credentials", nil)
		return
	}
	secret, err := getJWTSecret()
	if err != nil {
		ResponseJSON(c, http.StatusInternalServerError, "Server misconfiguired", nil)
		return
	}

	expirationTime := time.Now().Add(15 * time.Minute)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp": expirationTime.Unix(),
	})

	//sign the token
	tokenString, err := token.SignedString(secret)
	if err != nil {
		ResponseJSON(c, http.StatusInternalServerError, "Couldn't generate token", nil)
		return
	}

	ResponseJSON(c, http.StatusOK, "Token generated successfully", gin.H{"token": tokenString})
}
