package email

import (
	"context"
	"errors"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Email struct {
	ID      primitive.ObjectID `bson:"_id,omitempty"`
	Address string             `bson:"address"`
}

type EmailRoute struct {
	Email string `json:"status" binding:"required"`
}

func VerifiedEmail(database *mongo.Database) gin.HandlerFunc {
	return func(c *gin.Context) {

		var emailBody EmailRoute

		if err := c.ShouldBindJSON(&emailBody); err != nil {
			c.JSON(400, gin.H{
				"Error": "Issue with body binding",
				"Exact": err.Error(),
			})
			return
		}

		email := Email{Address: emailBody.Email}
		err := addUniqueEmail(database, email)
		if err != nil {
			c.JSON(400, gin.H{
				"Error": "Issue with insertion",
				"Exact": err.Error(),
			})
			return
		}

		if err := sendEmail(email.Address); err != nil {
			c.JSON(400, gin.H{
				"Error": "Issue sending followup email",
				"Exact": err.Error(),
			})
			return
		}

	}
}

func addUniqueEmail(db *mongo.Database, email Email) error {
	collection := db.Collection("email")

	var result Email
	err := collection.FindOne(context.TODO(), bson.M{"address": email.Address}).Decode(&result)

	if err == mongo.ErrNoDocuments {
		_, err := collection.InsertOne(context.TODO(), email)
		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	}

	return nil
}

func sendEmail(email string) error {

	from := mail.NewEmail("i9 Fitness", "main@i9fit.co")
	subject := "i9: You are Valid and Verified"
	to := mail.NewEmail("", email)
	plainTextContent := "Your email has been verified by i9! We won't pester you with emails too much..."
	message := mail.NewSingleEmail(from, subject, to, plainTextContent, "")

	key := os.Getenv("SENDGRID_KEY")
	if key == "" {
		return errors.New("no sendgrid api key")
	}
	client := sendgrid.NewSendClient(key)

	_, err := client.Send(message)
	if err != nil {
		return err
	}

	return nil
}
