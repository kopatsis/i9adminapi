package platform

import (
	"i9-adminapi/email"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func New(database *mongo.Database) *gin.Engine {
	router := gin.Default()

	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "blop",
		})
	})

	router.POST("/verified", AuthRequired(), email.VerifiedEmail(database))

	return router
}
