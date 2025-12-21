package authentication

import (
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"

	dbFunc "business-connect/database/dbHelpFunc"
	Data "business-connect/models"
)

func SendEmailPasswordChange(ctx *fiber.Ctx) error {
	// Separate variables for different error checks
	var (
		otp struct {
			Email string
		}
		bindErr        error
		UserBodyReturn Data.User
		EmailErr       error
		emailUserErr   error
	)

	// Check if there is an error binding the request
	if bindErr = ctx.BodyParser(&otp); bindErr != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to read request body",
		})
	}

	// let's get user by email ID
	UserBodyReturn, EmailErr = dbFunc.DBHelper.FindByEmail(otp.Email)
	if EmailErr != nil {
		return ctx.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"error": "failed to find user by email",
		})
	}

	// if !UserBodyReturn.EmailActivated {
	// 	// User needs to verify email address, send an error message
	// 	return ctx.Status(http.StatusConflict).JSON(fiber.Map{
	// 		"error": "User already exist, please verify your email ID",
	// 	})
	// }

	// let's send a token with users first and last name to verify the user with email ID
	emailUserErr = ForgotPasswordEmailVerification(UserBodyReturn.FirstName+" "+UserBodyReturn.LastName, UserBodyReturn.Email)

	if emailUserErr != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "email verification failed",
		})
	}

	// send response after successful login
	return ctx.Status(http.StatusCreated).JSON(fiber.Map{
		"success": "Email verification sent successfully to " + UserBodyReturn.Email,
	})
}

func VerifyForgotPassword(ctx *fiber.Ctx) error {
	var (
		SentOTPData struct {
			Email    string
			SentOTP  string
			Password string
		}
		bindErr error
		hashErr error
		// updateMaxErr   error
		hashOTPErr     error
		OTPBody        Data.OTP
		hashedPassword string
		existingUser   Data.User
	)

	// Check if there is an error binding the request
	if bindErr = ctx.BodyParser(&SentOTPData); bindErr != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to read body",
		})
	}

	if dbFunc.DBHelper.CheckSpecialCharacters(SentOTPData.Email) {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "This email is invalid because it uses illegal characters. Please enter a valid email.",
		})
	}

	// let's check if the otp exists and check it's validity
	OTPBody, hashOTPErr = dbFunc.DBHelper.GetAndCheckOTPByEmail(SentOTPData.Email, SentOTPData.SentOTP)
	// Check if OTPBody is the zero-value of Data.OTP so that we can check the max tries
	if OTPBody != (Data.OTP{}) {
		// let's check the max tries of an otp
		if OTPBody.MaxTry >= 5 {
			return ctx.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "email limit check exceeded"})
		}
	}

	if hashOTPErr != nil {
		// let's get this otp body to check if the max try is greater or equals to 5
		if OTPBody.MaxTry >= 5 {
			return ctx.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "email limit check exceeded"})
		}
	}

	if hashOTPErr != nil {
		// Handle different error cases
		switch {
		case hashOTPErr.Error() == "otp not found":
			return ctx.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "OTP not found"})
		case hashOTPErr.Error() == "incorrect otp value":
			updateMaxErr := dbFunc.DBHelper.UpdateMaxTry(SentOTPData.Email)
			if updateMaxErr != nil {
				return ctx.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "an error occurred"})
			}
			return ctx.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "Wrong OTP"})
		default:
			return ctx.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "an error occurred"})
		}
	}

	// OTP expiry duration is 60 minutes
	otpExpiryDuration := 60 * time.Minute

	// Check if OTP has expired
	createdTime := time.Unix(OTPBody.CreatedAT, 0)
	expiryTime := createdTime.Add(otpExpiryDuration)

	if time.Now().After(expiryTime) {
		return ctx.Status(http.StatusRequestTimeout).JSON(fiber.Map{"error": "Reset Link Expired"})
	}

	existingUser, userErr := dbFunc.DBHelper.FindByEmail(SentOTPData.Email)
	if userErr != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "an error occurred",
		})
	}

	// Create password hash
	hashedPassword, hashErr = dbFunc.DBHelper.CreatePasswordHash(SentOTPData.Password)
	if hashErr != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to create password hash",
		})
	}

	existingUser.Password = hashedPassword

	// Save the updated user to the database
	updateErr := dbFunc.DBHelper.UpdateUser(existingUser)

	if updateErr != nil {
		if updateErr.Error() == "user to update not found" {
			// The record with the specified email was not found
			return ctx.Status(http.StatusCreated).JSON(fiber.Map{
				"error": "user not found",
			})
		} else {
			// Some other error occurred
			return ctx.Status(http.StatusCreated).JSON(fiber.Map{
				"error": "error updating user",
			})
		}
	}

	// after all successful checks let's delete the otp reset link from our db
	delOtpErr := dbFunc.DBHelper.DeleteExistingOTPByID(OTPBody.CustomID)
	if delOtpErr != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "an error occurred",
		})
	}

	// Returning a success message (user's password updated successfully)
	return ctx.Status(http.StatusCreated).JSON(fiber.Map{
		"success": "Password updated successfully",
	})
}
