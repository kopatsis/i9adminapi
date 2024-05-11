package platform

import (
	"i9-adminapi/email"
	"net/http"

	firebase "firebase.google.com/go"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func New(database *mongo.Database, firebase *firebase.App) *gin.Engine {
	router := gin.Default()

	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "blop",
		})
	})

	router.POST("/verified", AuthRequired(), email.SendVerifiedEmail(firebase))

	return router
}
