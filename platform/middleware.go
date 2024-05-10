package platform

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {

		apiKey := os.Getenv("SPECIAL_KEY")
		if apiKey == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Inexistent API key"})
			return
		}

		requestKey := c.GetHeader("X-API-KEY")
		if requestKey != apiKey {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid API key"})
			return
		}
		c.Next()
	}
}
