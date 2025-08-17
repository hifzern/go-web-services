package tests

import (
	"07-Go-Book-Api/api"

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
