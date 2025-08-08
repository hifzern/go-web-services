package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// album represent data about a record album
type album struct {
	ID     string `json:"id"`
	Title  string `json:"title"`
	Artist string `json:"artist"`
	Price  string `json:"price"`
}

// albums slice to seed record album data
var albums = []album{
	{ID: "1", Title: "Jakwan Bagung", Artist: "Jawgger", Price: "69000"},
	{ID: "2", Title: "Kakso Bontol", Artist: "Chigga", Price: "78000"},
	{ID: "3", Title: "Pasi Nadang", Artist: "Batackz", Price: "88000"},
	{ID: "4", Title: "Dawet Jembud Kecabut", Artist: "Sundgga", Price: "98000"},
}

// getalbums respond with the list of all albums as JSON
func getAlbums(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, albums)
}

func main() {
	router := gin.Default()
	router.GET("/albums", getAlbums)

	router.Run("localhost:8888")
}
