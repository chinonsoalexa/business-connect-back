package profile

import (
	"net/http"

	dbFunc "business-connect/database/dbHelpFunc"
	helperFunc "business-connect/paystack"

	"github.com/gofiber/fiber/v2"
)

type Password struct {
	OldPassword string
	NewPassword string
}

func UpdatePassword(ctx *fiber.Ctx) error {
	// get stored user id from request time line
	userId := ctx.Locals("user-id")

	if userId == nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "failed to get user",
		})
	}

	user, uuidErr := helperFunc.PaystackHelper.FindByUuidFromLocal(userId)
	if uuidErr != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "failed to get user from request",
		})
	}

	var password Password

	// Check if there is an error binding the request
	if bindErr := ctx.BodyParser(&password); bindErr != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to read body",
		})
	}

	// comparing existing password with user old password in form of hash
	hashErr := dbFunc.DBHelper.ComparePasswordHash(user.Password, password.OldPassword)

	// checking if there was an error comparing the hashes
	if hashErr != nil {
		// Check if the error is due to a password mismatch
		if hashErr.Error() == "password does not match" {
			// Handle password mismatch error here
			return ctx.Status(http.StatusNotFound).JSON(fiber.Map{
				"error": "Invalid password",
			})
		} else if hashErr.Error() == "password too short to be a bcrypt password" {
			// Handle password mismatch error here
			return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid password",
			})
		}
		// Handle other bcrypt-related errors
	}

	NewPassword, err := dbFunc.DBHelper.CreatePasswordHash(password.NewPassword)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "error creating new password",
		})
	}

	// if all is successful let's update the users password to the new one
	user.Password = string(NewPassword)

	// update user's profile in the database
	dbAddErr := dbFunc.DBHelper.UpdateUser(user)
	if dbAddErr != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "error updating user's password"})
	}

	return ctx.Status(http.StatusOK).JSON(fiber.Map{"success": "successfully updated user's password"})

}
