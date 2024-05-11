package email

import (
	"context"
	"os"

	firebase "firebase.google.com/go"
	"github.com/gin-gonic/gin"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"google.golang.org/api/iterator"
)

func SendVerifiedEmail(firebase *firebase.App) gin.HandlerFunc {
	return func(c *gin.Context) {

		authClient, err := firebase.Auth(context.TODO())
		if err != nil {
			c.JSON(400, gin.H{
				"Error": "Issue sending followup email",
				"Exact": err.Error(),
			})
			return
		}

		userIterator := authClient.Users(context.TODO(), "")
		var verifiedEmails []string
		for {
			user, err := userIterator.Next()
			if err != nil {

				if err == iterator.Done {
					break
				} else {
					c.JSON(400, gin.H{
						"Error": "Issue sending followup email",
						"Exact": err.Error(),
					})
					return
				}
			}

			if user.EmailVerified {
				verifiedEmails = append(verifiedEmails, user.Email)
			}
		}

		failedEmails := []string{}
		successEmails := []string{}

		apiKey := os.Getenv("SENDGRID_KEY")
		fromEmail := "main@i9fit.co"
		fromName := "i9 Team"
		subject := "i9 Yeah!"
		htmlContent := "<p>Example email</p>"

		client := sendgrid.NewSendClient(apiKey)
		for _, email := range verifiedEmails {
			message := mail.NewSingleEmail(
				mail.NewEmail(fromName, fromEmail),
				subject,
				mail.NewEmail("", email),
				"",
				htmlContent,
			)
			_, err := client.Send(message)
			if err != nil {
				failedEmails = append(failedEmails, email)
			} else {
				successEmails = append(successEmails, email)
			}
		}

		c.JSON(200, gin.H{
			"Success": successEmails,
			"Fail":    failedEmails,
		})

	}
}
