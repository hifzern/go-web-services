package main

import (
	"07-Go-Book-Api/api"

	"github.com/gin-gonic/gin"
)

func main() {
	api.InitDB()
	r := gin.Default()

	//public routes
	r.POST("/token", api.GenerateJWT)

	// protected routes
	protected := r.Group("/", api.JWTAuthMiddleware())
	{
		protected.POST("/book", api.CreateBook)
		protected.GET("/books", api.GetBooks)
		protected.GET("/book/:id", api.GetBook)
		protected.PUT("/book/:id", api.UpdateBook)
		protected.DELETE("/book/:id", api.DeleteBook)
	}
	

	r.Run(":8080")
}
