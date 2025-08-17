package tests

import (
	"07-Go-Book-Api/api"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB() {
	var err error
	api.DB, err = gorm.Open(sqlite.Open(":memory"), &gorm.Config{})
	if err != nil {
		panic("failed to connect test database")
	}
	api.DB.AutoMigrate(&api.Book{})
}

func addBook() api.Book {
	book := api.Book{Title: "Menolak Ngoding", Author: "Jakwan Bagung", Year: 2025}
	api.DB.Create(&book)
	return book
}

func TestCreateBook(t *testing.T) {
	setupTestDB()
	router := gin.Default()
	router.POST("/book", api.CreateBook)
	book := api.Book{
		Title: "Test Demo", Author: "Kakso Bontol", Year: 2024,
	}
	jsonValue, _ := json.Marshal(book)
	req, _ := http.NewRequest("POST", "/book", bytes.NewBuffer(jsonValue))

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

}
