package authentication

import (
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"

	dbFunc "business-connect/database/dbHelpFunc"
	reqAuth "business-connect/middleware"
	myjwt "business-connect/middleware/myjwt"
	Data "business-connect/models"
	helperFunc "business-connect/paystack"
)

func SignUp(ctx *fiber.Ctx) error {
	var (
		bindErr      error
		dbErr        error
		dbAddErr     error
		existingUser Data.User
		emailErr     error
		NewUser      Data.User
		createdUser  Data.User
	)

	// Check if there is an error binding the request
	if bindErr = ctx.BodyParser(&NewUser); bindErr != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to read body",
		})
	}

	NewUser.UserType = "USER"

	// Check if the user already exists by email in the data base
	existingUser, dbErr = dbFunc.DBHelper.FindByEmail(NewUser.Email)

	if dbErr != nil {
		if dbErr.Error() == "error retrieving user" {
			// An error occurred while querying the database
			return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to check user existence",
			})
		}
	}

	if dbErr == nil {
		// let's check if the email is valid
		if dbFunc.DBHelper.CheckSpecialCharacters(NewUser.Email) {
			return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
				"error": "This email is invalid because it uses illegal characters. Please enter a valid email",
			})
		}
		// let's check if the user's email is verified
		if !existingUser.EmailActivated {
			// User needs to verify email address, send an error message
			return ctx.Status(http.StatusConflict).JSON(fiber.Map{
				"error": "User already exist, please verify your email ID",
			})
		}
		// User already exists, send an error message
		return ctx.Status(http.StatusConflict).JSON(fiber.Map{
			"error": "User already exist, please login",
		})
	}

	usersIP := ctx.IP()

	if usersIP == "102.89.33.226" {
		NewUser.Suspended = true
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Sorry your account has been suspended",
		})
	}

	// Add new user to the database
	createdUser, dbAddErr = dbFunc.DBHelper.CreateNewUser(NewUser)
	if dbAddErr != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Error creating new user"})
	}

	// let's send token to user to verify the user with email ID
	emailErr = EmailVerification(createdUser.FirstName+" "+createdUser.LastName, NewUser.Email)

	if emailErr != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": emailErr.Error(),
		})
	}

	// let's return a success message to the client side
	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"success": "Email verification link sent to " + NewUser.Email,
	})
}

func EmailAuthentication(ctx *fiber.Ctx) error {
	// Separate variables for different error checks
	var (
		otp            struct{ Email, SentOTP string }
		bindErr        error
		hashOTPErr     error
		OTPBody        Data.OTP
		UserBodyReturn Data.User
		EmailErr       error
	)

	// Check if there is an error binding the request
	if bindErr = ctx.BodyParser(&otp); bindErr != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to read request body",
		})
	}

	// let's check if the otp exists and check it's validity
	OTPBody, hashOTPErr = dbFunc.DBHelper.GetAndCheckOTPByEmail(otp.Email, otp.SentOTP)
	// Check if OTPBody is the zero-value of Data.OTP so that we can check the max tries
	if OTPBody != (Data.OTP{}) {
		// let's check the max tries of an otp
		if OTPBody.MaxTry >= 5 {
			return ctx.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "email limit check exceeded"})
		}
	}

	if hashOTPErr != nil {
		// let's get this otp body to check if the max try is greater or equals to 5
		otpMaxTryBody, err := dbFunc.DBHelper.GetOTPByEmail(otp.Email)
		if err != nil {
			return ctx.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "error getting otp by email for limit check"})
		}
		if otpMaxTryBody.MaxTry >= 5 && hashOTPErr.Error() == "incorrect otp value" {
			return ctx.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "email limit check exceeded"})
		}
	}

	// log.Println("where i verify the otp:", OTPBody)
	if hashOTPErr != nil {
		// Handle different error cases
		switch {
		case hashOTPErr.Error() == "otp not found":
			return ctx.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "OTP not found"})
		case hashOTPErr.Error() == "incorrect otp value":
			updateMaxErr := dbFunc.DBHelper.UpdateMaxTry(otp.Email)
			if updateMaxErr != nil {
				return ctx.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "an error occurred"})
			}
			return ctx.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "Wrong OTP"})
		default:
			return ctx.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "an error occurred"})
		}
	}

	// OTP expiry duration is 5 minutes
	otpExpiryDuration := 24 * time.Hour

	// Check if OTP has expired
	createdTime := time.Unix(OTPBody.CreatedAT, 0)
	expiryTime := createdTime.Add(otpExpiryDuration)

	if time.Now().After(expiryTime) {
		return ctx.Status(http.StatusRequestTimeout).JSON(fiber.Map{"error": "Verification Code Expired"})
	}

	// Activate user's email
	UserBodyReturn, EmailErr = dbFunc.DBHelper.FindByEmail(OTPBody.Email)
	if EmailErr != nil {
		return ctx.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "Failed to find user by email"})
	}

	UserBodyReturn.EmailActivated = true

	// Update the user
	updateErr := dbFunc.DBHelper.UpdateUser(UserBodyReturn)
	if updateErr != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update user"})
	}

	// Delete the OTP after successful validation
	delOtpErr := dbFunc.DBHelper.DeleteExistingOTPByID(OTPBody.CustomID)
	if delOtpErr != nil {
		log.Println("error occurred here:", delOtpErr)
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "An error occurred"})
	}

	// let's create a paystack visual account for this user

	// let's get user by email to authenticate the user
	userEmail, userErr := dbFunc.DBHelper.FindByEmail(otp.Email)
	if userErr != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Error getting user"})
	}

	// fmt.Println("this is the users email: ", otp.Email)
	// fmt.Println("this is the users id: ", userEmail.ID)

	role := "user"

	// now generate cookies for this user
	authTokenString, refreshTokenString, csrfSecret, errJwt := myjwt.CreateNewTokens(ctx, strconv.FormatUint(uint64(userEmail.ID), 10), role)
	if errJwt != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Error creating cookies"})
	}

	// set the cookies to these newly created jwt's and also set's the csrf token
	reqAuth.SetAuthAndRefreshCookies(ctx, authTokenString, refreshTokenString, csrfSecret)

	// Return a success message
	return ctx.Status(http.StatusCreated).JSON(fiber.Map{"success": "Email verification successful"})
}

func ResendEmailVerification(ctx *fiber.Ctx) error {
	// Separate variables for different error checks
	var (
		otp struct {
			Email string
		}
		bindErr  error
		emailErr error
		OtpErr   error
		OTPBody  Data.OTP
	)

	// Check if there is an error binding the request
	if bindErr = ctx.BodyParser(&otp); bindErr != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to read request body",
		})
	}

	// let's check if user is in our unverified email list
	OTPBody, OtpErr = dbFunc.DBHelper.GetOTPByEmail(otp.Email)
	if OtpErr != nil {
		return ctx.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"error": "failed to get previous email OTP",
		})
	}

	// Check if the user already exists by email in the data base so we can take the users name
	existingUser, dbErr := dbFunc.DBHelper.FindByEmail(OTPBody.Email)

	if dbErr != nil {
		if dbErr.Error() == "user not found" {
			// An error occurred while querying the database
			return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"error": "user not found in the db",
			})
		} else if dbErr.Error() == "error retrieving user" {
			// An error occurred while querying the database
			return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to check user existence",
			})
		}
	}

	// let's send token to user to verify the user with email ID
	emailErr = EmailVerification(existingUser.FirstName+" "+existingUser.LastName, OTPBody.Email)

	if emailErr != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "email verification failed",
		})
	}

	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"success": "Email re-sent successfully to " + otp.Email,
	})
}

const (
	referralCodeLength = 10
)

func GenerateReferralCode() string {
	source := rand.NewSource(time.Now().UnixNano())
	rand := rand.New(source)

	code := make([]byte, referralCodeLength)
	for i := 0; i < referralCodeLength; i++ {
		code[i] = byte(rand.Intn(10) + '0') // '0' to '9'
	}

	return string(code)
}

func CreateNewOtp(ctx *fiber.Ctx) error {
	var (
		bindErr error
		NewUser Data.User
	)

	// get stored user id from request time line
	userId := ctx.Locals("user-id")

	if userId == nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "User ID is < nil >",
		})
	}

	_, uuidErr := helperFunc.PaystackHelper.FindByUuidFromLocal(userId)
	if uuidErr != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "failed to get user from request",
		})
	}

	// Check if there is an error binding the request
	if bindErr = ctx.BodyParser(&NewUser); bindErr != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to read body",
		})
	}

	// let's return a success message to the client side
	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"success": "Transaction Code Updated Successfuly for " + NewUser.Email,
	})
}
