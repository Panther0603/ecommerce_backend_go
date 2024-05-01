package middleware

import (
	"Ecommerce-Backend/tokens"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func Authentication() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Exclude authentication for specific routes
		if strings.HasPrefix(c.Request.URL.Path, "/user/signup") || strings.HasPrefix(c.Request.URL.Path, "/login") {
			c.Next()
			return
		}

		ClientToken := c.Request.Header.Get("token")

		if ClientToken == "" {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "No authorization token found"})
			c.Abort()
			return
		}

		claims, err := tokens.ValidateToken(ClientToken)

		if err != "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			c.Abort()
			return
		}

		c.Set("email", claims.Email)
		c.Set("uid", claims.Uid)
		c.Next()
	}
}
