package api

import (
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	DB        *gorm.DB
	jwtSecret = []byte(os.Getenv("JWT_SECRET"))
)

type Book struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Title     string    `json:"title"`
	Author    string    `json:"author"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func ResponseJSON(c *gin.Context, code int, msg string, data any) {
	c.JSON(code, gin.H{"message": msg, "data": data})
}

func InitDB() {
	_ = godotenv.Load() // jangan fatal; di production .env bisa tidak ada

	dsn := os.Getenv("DB_URL") // pastikan benar
	if dsn == "" {
		log.Fatal("DB_URL is empty")
	}

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}

	if err := DB.AutoMigrate(&Book{}); err != nil {
		log.Fatal("Failed to migrate schema: ", err)
	}
}

func CreateBook(c *gin.Context) {
	var in Book
	if err := c.ShouldBindJSON(&in); err != nil {
		ResponseJSON(c, http.StatusBadRequest, "Invalid input", nil)
		return
	}
	if err := DB.Create(&in).Error; err != nil {
		ResponseJSON(c, http.StatusInternalServerError, "DB error", nil)
		return
	}
	ResponseJSON(c, http.StatusCreated, "Book created successfully", in)
}

func GetBooks(c *gin.Context) {
	var books []Book
	if err := DB.Find(&books).Error; err != nil {
		ResponseJSON(c, http.StatusInternalServerError, "DB error", nil)
		return
	}
	ResponseJSON(c, http.StatusOK, "Books retrieved successfully", books)
}

func GetBook(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		ResponseJSON(c, http.StatusBadRequest, "Invalid ID", nil)
		return
	}
	var book Book
	if err := DB.First(&book, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ResponseJSON(c, http.StatusNotFound, "Book not found", nil)
		} else {
			ResponseJSON(c, http.StatusInternalServerError, "DB error", nil)
		}
		return
	}
	ResponseJSON(c, http.StatusOK, "Book retrieved successfully", book)
}

func UpdateBook(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		ResponseJSON(c, http.StatusBadRequest, "Invalid ID", nil)
		return
	}

	var book Book
	if err := DB.First(&book, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ResponseJSON(c, http.StatusNotFound, "Book not found", nil)
		} else {
			ResponseJSON(c, http.StatusInternalServerError, "DB error", nil)
		}
		return
	}

	// bind ke struct input agar tidak menimpa field tak terkirim
	var in struct {
		Title  *string `json:"title"`
		Author *string `json:"author"`
	}
	if err := c.ShouldBindJSON(&in); err != nil {
		ResponseJSON(c, http.StatusBadRequest, "Invalid input", nil)
		return
	}

	tx := DB.Model(&book)
	if in.Title != nil {
		tx = tx.Update("title", *in.Title)
	}
	if in.Author != nil {
		tx = tx.Update("author", *in.Author)
	}
	if tx.Error != nil {
		ResponseJSON(c, http.StatusInternalServerError, "DB error", nil)
		return
	}

	ResponseJSON(c, http.StatusOK, "Book updated successfully", book)
}

func DeleteBook(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		ResponseJSON(c, http.StatusBadRequest, "Invalid ID", nil)
		return
	}
	res := DB.Delete(&Book{}, id)
	if res.Error != nil {
		ResponseJSON(c, http.StatusInternalServerError, "DB error", nil)
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

	// TODO: ganti dengan verifikasi ke DB
	if loginRequest.Username != os.Getenv("ADMIN_USER") || loginRequest.Password != os.Getenv("ADMIN_PASS") {
		ResponseJSON(c, http.StatusUnauthorized, "Invalid credentials", nil)
		return
	}

	if len(jwtSecret) == 0 {
		ResponseJSON(c, http.StatusInternalServerError, "JWT secret not configured", nil)
		return
	}

	now := time.Now()
	exp := now.Add(15 * time.Minute)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": loginRequest.Username,
		"iat": now.Unix(),
		"exp": exp.Unix(),
	})

	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		ResponseJSON(c, http.StatusInternalServerError, "Couldn't generate token", nil)
		return
	}

	ResponseJSON(c, http.StatusOK, "Token generated successfully", gin.H{"token": tokenString})
}
