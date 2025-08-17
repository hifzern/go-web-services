package tests

import (
	"07-Go-Book-Api/api"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
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

	if status := w.Code; status != http.StatusCreated {
		t.Errorf("Expected status %d you got %d", http.StatusCreated, status)
	}
	var response api.JsonResponse
	json.NewDecoder(w.Body).Decode(&response)

	if response.Data == nil {
		t.Errorf("Expected book data, got nil")
	}
}

func TestGetBooks(t *testing.T) {
	setupTestDB()
	addBook()
	router := gin.Default()
	router.GET("/books", api.GetBooks)

	req, _ := http.NewRequest("GET", "/books", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if status := w.Code; status != http.StatusOK {
		t.Errorf("Expected sttaus %d you got %d", http.StatusOK, status)
	}
	var response api.JsonResponse
	json.NewDecoder(w.Body).Decode(&response)

	if len(response.Data.([]interface{})) == 0 {
		t.Errorf("Expected non-empty books list")
	}
}

func TestGetBook(t *testing.T) {
	setupTestDB()
	book := addBook()
	router := gin.Default()
	router.POST("/book/:id", api.GetBook)

	req, _ := http.NewRequest("GET", "/book/"+strconv.Itoa(int(book.ID)), nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if status := w.Code; status != http.StatusOK {
		t.Errorf("Expected status %d got %d", http.StatusOK, status)
	}

	var response api.JsonResponse
	json.NewDecoder(w.Body).Decode(&response)

	if response.Data == nil || response.Data.(map[string]interface{})["id"] != float64(book.ID) {
		t.Errorf("Expected book ID %d, got nil or wrong ID", book.ID)
	}
}

func TestUpdateBook(t *testing.T) {
	setupTestDB()
	book := addBook()
	router := gin.Default()
	router.PUT("/book/:id", api.UpdateBook)

	updateBook := api.Book{
		Title: "Lorem Ipsum", Author: "Dolor Amet", Year: 2024,
	}
	jsonValue, _ := json.Marshal(updateBook)

	req, _ := http.NewRequest("PUT", "/book/"+strconv.Itoa(int(book.ID)), bytes.NewBuffer(jsonValue))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if status := w.Code; status != http.StatusOK {
		t.Errorf("Expected status %d got %d", http.StatusOK, status)
	}

	var response api.JsonResponse
	json.NewDecoder(w.Body).Decode(&response)

	if response.Data == nil || response.Data.(map[string]interface{})["title"] != "Lorem Ipsum" {
		t.Errorf("Expected updated book title 'Lorem Ipsum', got %v", response.Data)
	}
}

func TestDeleteBook(t *testing.T) {
	setupTestDB()
	book := addBook()
	router := gin.Default()
	router.DELETE("/book/:id", api.DeleteBook)

	req, _ := http.NewRequest("DELETE", "/book/"+strconv.Itoa(int(book.ID)), nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if status := w.Code; status != http.StatusOK {
		t.Errorf("Expected status %d got %d", http.StatusOK, status)
	}

	var response api.JsonResponse
	json.NewDecoder(w.Body).Decode(&response)

	if response.Message != "Book deleted successfully" {
		t.Errorf("Expected delete message 'Book deleted successfully', got %v", response.Message)
	}

	//verify that the book was deleted
	var deletedBook api.Book
	result := api.DB.First(&deletedBook, book.ID)
	if result.Error == nil {
		t.Errorf("Expected book to be deleted but it still exists")
	}

}
