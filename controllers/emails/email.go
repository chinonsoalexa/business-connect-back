package emails

import (
	email "business-connect/controllers/authentication"
	dbFunc "business-connect/database/dbHelpFunc"
	Data "business-connect/models"
	upload "business-connect/upload"
	"log"
	"time"

	// "regexp"
	"strings"
	// "sync"

	"github.com/gofiber/fiber/v2"
)

func SendEmails(ctx *fiber.Ctx) error {
	// log.Println("Starting SendEmails function")

	// Get stored user id from request timeline
	userId := ctx.Locals("user-id")
	// log.Printf("User ID: %v\n", userId)

	if userId == nil {
		log.Println("User ID is nil")
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "failed to get user",
		})
	}

	user, uuidErr := dbFunc.DBHelper.FindByUuidFromLocal(userId)
	// log.Printf("Retrieved user: %v\n", user)
	if uuidErr != nil {
		log.Printf("Error retrieving user: %v\n", uuidErr)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": uuidErr.Error(),
		})
	}

	var EmailToSend Data.Email

	// Check if there is an error binding the request
	if bindErr := ctx.BodyParser(&EmailToSend); bindErr != nil {
		log.Printf("Error binding request body: %v\n", bindErr)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to read body",
		})
	}
	// log.Printf("Email to send: %+v\n", EmailToSend)

	// Call the database helper function to retrieve the order
	emailContent, err := upload.UploadEmailFiles(EmailToSend.Content)
	if err != nil {
		log.Printf("Error uploading email files: %v\n", err)
		// error uploading email images
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to upload email images",
		})
	}
	// log.Println("Email content after upload: ", emailContent)

	// Modify the <img> tags by adding the style attribute
	updatedHTML := AddStyleToImgTags(emailContent)
	// if err != nil {
	// 	// error updating email images styles
	// 	return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
	// 		"error": "failed to update image style",
	// 	})
	// }

	// Define a goroutine to send the emails asynchronously
	go func() {
		if EmailToSend.SendTo == "me" {
			// Send to only the current user
			emailErr := email.SendEmailToSubscribers(EmailToSend.Subject, updatedHTML, user.Email)
			if emailErr != nil {
				log.Printf("Failed to send email to user %s: %v\n", user.Email, emailErr)
			}
		} else {
			// Send to all subscribers
			emailSubscribers, emailErr := dbFunc.DBHelper.GetBusinessConnectEmailSubscribers()
			if emailErr != nil {
				log.Printf("Failed to retrieve email subscribers: %v\n", emailErr)
				return
			}

			for i, emaill := range emailSubscribers {
				// Wait 5 minutes between sends, except before the first one
				if i > 0 {
					time.Sleep(5 * time.Minute)
				}

				// Send email (in the same goroutine sequentially)
				emailSubErr := email.SendEmailToSubscribers(EmailToSend.Subject, updatedHTML, emaill.Email)
				if emailSubErr != nil {
					log.Printf("Failed to send email to %s: %v\n", emaill.Email, emailSubErr)
				}
			}
		}
	}()

	// Save sent emails
	var sentEmail Data.Email
	sentEmail.Subject = EmailToSend.Subject
	sentEmail.Content = updatedHTML
	sentEmail.SendTo = EmailToSend.SendTo
	saveSentEmailErr := dbFunc.DBHelper.SaveBusinessConnectSentEmail(sentEmail)
	if saveSentEmailErr != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to save email copy",
		})
	}

	// Immediately return a success response without waiting for the emails to be sent
	// log.Println("Returning success response")
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": "email sending in progress",
	})
}

func AddStyleToImgTags(htmlContent string) string {
	// Define the style you want to add to each <img> tag
	style := `style="max-width: 100%; height: auto; border-radius: 8px;"`

	// Split the content by <img
	parts := strings.Split(htmlContent, "<img")

	// Rebuild HTML content with added style to <img> tags
	var updatedHTML strings.Builder
	updatedHTML.WriteString(parts[0]) // Add the part before the first <img

	// Iterate over the parts to add style to each <img> tag
	for i := 1; i < len(parts); i++ {
		// Add style attribute to the start of each <img> tag
		updatedHTML.WriteString("<img ")
		updatedHTML.WriteString(style)
		updatedHTML.WriteString(" ")

		// Find the end of the <img> tag (until '>')
		tagEnd := strings.Index(parts[i], ">")
		if tagEnd != -1 {
			updatedHTML.WriteString(parts[i][:tagEnd])
			updatedHTML.WriteString(">")
			updatedHTML.WriteString(parts[i][tagEnd+1:])
		} else {
			updatedHTML.WriteString(parts[i])
		}
	}

	return updatedHTML.String()
}
