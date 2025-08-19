package main

import (
	"07-Go-Book-Api/api"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()
	api.InitDB()
	r := gin.Default()

	//public routes
	r.POST("/token", api.GenerateJWT)

	// auth routes
	auth := r.Group("/", api.JWTAuthMiddleware())
	{
		auth.POST("/book", api.CreateBook)
		auth.GET("/books", api.GetBooks)
		auth.GET("/book/:id", api.GetBook)
		auth.PUT("/book/:id", api.UpdateBook)
		auth.DELETE("/book/:id", api.DeleteBook)
	}

	r.Run(":8080")
}
