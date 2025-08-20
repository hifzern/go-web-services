package api

import (
	"log"
	"net/http"
	"os"
	"strconv"
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
		log.Fatal("Failed to connect to database: ", err)
	}

	// migrate once (hapus duplikat)
	if err := DB.AutoMigrate(&Book{}); err != nil {
		log.Fatal("Failed to migrate schema: ", err)
	}
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// helper validasi id
func parseID(c *gin.Context) (uint, bool) {
	idStr := c.Param("id")
	n, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil || n == 0 {
		ResponseJSON(c, http.StatusBadRequest, "Invalid id", nil)
		return 0, false
	}
	return uint(n), true
}

func CreateBook(c *gin.Context) {
	var payload Book

	if err := c.ShouldBindJSON(&payload); err != nil {
		ResponseJSON(c, http.StatusBadRequest, "Invalid input", nil)
		return
	}

	if payload.Title == "" || payload.Author == "" || payload.Year <= 0 {
		ResponseJSON(c, http.StatusBadRequest, "Title, Author, and Year are required", nil)
		return
	}

	if err := DB.Create(&payload).Error; err != nil {
		ResponseJSON(c, http.StatusInternalServerError, "Failed to create book", nil)
		return
	}
	ResponseJSON(c, http.StatusCreated, "Book created successfully", payload)
}

func GetBooks(c *gin.Context) {
	var books []Book
	if err := DB.Find(&books).Error; err != nil {
		ResponseJSON(c, http.StatusInternalServerError, "Failed to fetch books", nil)
		return
	}
	ResponseJSON(c, http.StatusOK, "Books retrieved successfully", books)
}

func GetBook(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}

	var book Book
	if err := DB.First(&book, id).Error; err != nil {
		ResponseJSON(c, http.StatusNotFound, "Book not found", nil)
		return
	}
	ResponseJSON(c, http.StatusOK, "Book retrieved successfully", book)
}

func UpdateBook(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}

	var existing Book
	if err := DB.First(&existing, id).Error; err != nil {
		ResponseJSON(c, http.StatusNotFound, "Book not found", nil)
		return
	}

	var payload Book
	if err := c.ShouldBindJSON(&payload); err != nil {
		ResponseJSON(c, http.StatusBadRequest, "Invalid input", nil)
		return
	}

	if payload.Title == "" || payload.Author == "" || payload.Year <= 0 {
		ResponseJSON(c, http.StatusBadRequest, "Title, Author, and Year are required", nil)
		return
	}

	existing.Title = payload.Title
	existing.Author = payload.Author
	existing.Year = payload.Year

	if err := DB.Save(&existing).Error; err != nil {
		ResponseJSON(c, http.StatusInternalServerError, "Failed to update book", nil)
		return
	}
	ResponseJSON(c, http.StatusOK, "Book updated successfully", existing)
}

func DeleteBook(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}

	res := DB.Delete(&Book{}, id)
	if res.Error != nil {
		ResponseJSON(c, http.StatusInternalServerError, "Failed to delete book", nil)
		return
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

	// demo sederhana
	if loginRequest.Username != "admin" || loginRequest.Password != "Password" {
		ResponseJSON(c, http.StatusUnauthorized, "Invalid credentials", nil)
		return
	}

	secret, err := getJWTSecret()
	if err != nil || len(secret) == 0 {
		ResponseJSON(c, http.StatusInternalServerError, "Server misconfigured", nil)
		return
	}

	now := time.Now()
	expirationTime := now.Add(15 * time.Minute)

	// tambah claim dasar biar standar
	claims := jwt.MapClaims{
		"sub": loginRequest.Username,
		"iat": now.Unix(),
		"exp": expirationTime.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(secret)
	if err != nil {
		ResponseJSON(c, http.StatusInternalServerError, "Couldn't generate token", nil)
		return
	}

	ResponseJSON(c, http.StatusOK, "Token generated successfully", gin.H{"token": tokenString})
}
