package api
import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"os"

	var jwtSecret = []byte(os.Getenv("SECRET_TOKEN"))

	func JWTAuthMiddleware() gin.HandlerFunc {
		return func(c *gin.Context) {
			tokenString := c.GetHeader("Authorization")
			if tokenString == "" {
				ResponseJSON(c, http.StatusUnauthorized, "Auth token required", nil)
				c.Abort()
				return
			}

			//parse and validate the token
			_, err := jwt.Parse(tokenString)
		}
	}
)