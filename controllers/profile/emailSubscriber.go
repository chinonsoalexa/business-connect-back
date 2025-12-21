package profile

import (
	SendEmail "business-connect/controllers/authentication/emails"
	dbFunc "business-connect/database/dbHelpFunc"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type Email struct {
	Email string
}

func AddEmailSubscription(ctx *fiber.Ctx) error {

	var email Email

	// Check if there is an error binding the request
	if bindErr := ctx.BodyParser(&email); bindErr != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to read body",
		})
	}

	emailErr := dbFunc.DBHelper.AddEmailSubscriber(email.Email)
	if emailErr != nil {
		if emailErr.Error() == "user with email already exists" {
			return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"error": "user with email already exists",
			})
		}
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "error adding new email subscriber",
		})
	}

	// send a confirmation email
	emailErr = SendEmail.TodacWelcomeEmail(email.Email)
	if emailErr != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to send order email",
		})
	}

	return ctx.Status(http.StatusOK).JSON(fiber.Map{"success": "successfully added " + email.Email + " to our mail list"})
}
