package authentication

// imported packages to be used
import (
	// "fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"

	dbFunc "business-connect/database/dbHelpFunc"
	reqAuth "business-connect/middleware"
	myjwt "business-connect/middleware/myjwt"
	Data "business-connect/models"
)

// declared variables
var (
	RandomNum string
	Err       error
	user      Data.User
	OldUser   = struct {
		// taking users login credentials info
		Email    string `json:"email"`
		Password string `json:"password"`
	}{}
	MagicLinkOldUser = struct {
		// taking users login credentials info
		Email string `json:"email"`
	}{}
	MagicLinkOldUserVerified = struct {
		// taking users login credentials info
		Email string `json:"email"`
	}{}
)

// login function to authenticate user and create a middleware token for subsequent request authentication
func SignIn(ctx *fiber.Ctx) error {
	var (
		dbErr   error
		hashErr error
	)
	// binding json body and also checking if an error exists in the users login info
	if err := ctx.BodyParser(&OldUser); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to read body (email/password)",
		})
	}

	// checking emails that contains special characters for security reasons
	if dbFunc.DBHelper.CheckSpecialCharacters(OldUser.Email) {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "This email is invalid because it uses illegal characters. Please enter a valid email",
		})
	}

	// check if user exists in the db by email in the data base
	user, dbErr = dbFunc.DBHelper.FindByEmail(OldUser.Email)

	if dbErr != nil {
		if dbErr.Error() == "error retrieving user" {
			// An error occurred while querying the database
			return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to check user existence",
			})
		}
	}

	if dbErr != nil {
		// checking if there are any error encountered
		if dbErr.Error() == "user not found" {
			// User do not exists, send an error message
			return ctx.Status(http.StatusNotFound).JSON(fiber.Map{
				"error": "User do not exist, please sign up",
			})
		}
	}

	// let's check if the user's email is verified
	// if !user.EmailActivated {
	// 	// User needs to verify email address, send an error message
	// 	return ctx.Status(http.StatusConflict).JSON(fiber.Map{
	// 		"error": "User already exist, please verify your email ID",
	// 	})
	// }

	// comparing existing password with user login password in form of hash
	hashErr = dbFunc.DBHelper.ComparePasswordHash(user.Password, OldUser.Password)

	// checking if there was an error comparing the hashes
	if hashErr != nil {
		// Check if the error is due to a password mismatch
		if hashErr.Error() == "password does not match" {
			// Handle password mismatch error here
			return ctx.Status(http.StatusNotFound).JSON(fiber.Map{
				"error": "Invalid email or password",
			})
		} else if hashErr.Error() == "password too short to be a bcrypt password" {
			// Handle password mismatch error here
			return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid email or password",
			})
		}
		// Handle other bcrypt-related errors
	}

	// let's check if the user's email is suspended
	if user.Suspended {
		// User needs to verify email address, send an error message
		return ctx.Status(http.StatusConflict).JSON(fiber.Map{
			"error": "Your account has been suspended. Please contact support for more details.",
		})
	}

	role := "USER"

	// now generate cookies for this user
	authTokenString, refreshTokenString, csrfSecret, errJwt := myjwt.CreateNewTokens(ctx, strconv.FormatUint(uint64(user.ID), 10), role)
	if errJwt != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Error creating cookies"})
	}

	// fmt.Println("this is the Refresh Token String: ", refreshTokenString)

	// set the cookies to these newly created jwt's and also set's the csrf token
	reqAuth.SetAuthAndRefreshCookies(ctx, authTokenString, refreshTokenString, csrfSecret)

	// Returning a success message (user logged in successfully) and redirecting user to restricted page

	// send response after successful login
	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"success":   OldUser.Email + " Successfully logged in",
		"dashboard": "admin/dashboard-02.html",
	})
}

// login function to authenticate user and create a middleware token for subsequent request authentication using a magic link sent to the email address
func MagicLinkSignIn(ctx *fiber.Ctx) error {
	var (
		dbErr error
	)
	// binding json body and also checking if an error exists in the users login info
	if err := ctx.BodyParser(&MagicLinkOldUser); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to read body (email/password)",
		})
	}

	// checking emails that contains special characters for security reasons
	if dbFunc.DBHelper.CheckSpecialCharacters(MagicLinkOldUser.Email) {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "This email is invalid because it uses illegal characters. Please enter a valid email",
		})
	}

	// check if user exists in the db by email in the data base
	user, dbErr = dbFunc.DBHelper.FindByEmail(MagicLinkOldUser.Email)

	if dbErr != nil {
		if dbErr.Error() == "error retrieving user" {
			// An error occurred while querying the database
			return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to check user existence",
			})
		}
	}

	if dbErr != nil {
		// checking if there are any error encountered
		if dbErr.Error() == "user not found" {
			// User do not exists, send an error message
			return ctx.Status(http.StatusNotFound).JSON(fiber.Map{
				"error": "User do not exist, please sign up",
			})
		}
	}

	// let's check if the user's email is verified
	if !user.EmailVerified {
		// User needs to verify email address, send an error message
		return ctx.Status(http.StatusConflict).JSON(fiber.Map{
			"error": "User already exist, please verify your email ID",
		})
	}

	// let's check if the user's email is suspended
	if user.Suspended {
		// User needs to verify email address, send an error message
		return ctx.Status(http.StatusConflict).JSON(fiber.Map{
			"error": "Your account has been suspended. Please contact support for more details.",
		})
	}

	// let's send the magic link to the user's email
	magicError := MagicLinkEmailVerification(user.FullName, user.Email)
	if magicError != nil {
		return ctx.Status(http.StatusConflict).JSON(fiber.Map{
			"error": "there was an error sending mail",
		})
	}

	// Returning a success message (user logged in successfully) and redirecting user to restricted page

	// send response after successful login
	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"success": "Successfully sent a magic login link to " + MagicLinkOldUser.Email,
	})
}

// this is to verify the magic login link if it's correct
func VerifySignInMagicLink(ctx *fiber.Ctx) error {
	var (
		SentOTPData struct {
			Email   string
			SentOTP string
		}
		bindErr    error
		hashOTPErr error
		dbEmailErr error
		OTPBody    Data.OTP
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
	otpExpiryDuration := 5 * time.Minute

	// Check if OTP has expired
	createdTime := time.Unix(OTPBody.CreatedAT, 0)
	expiryTime := createdTime.Add(otpExpiryDuration)

	if time.Now().After(expiryTime) {
		return ctx.Status(http.StatusRequestTimeout).JSON(fiber.Map{"error": "Magic Link Expired"})
	}

	// after all successful checks let's delete the magic otp login link from our db
	delOtpErr := dbFunc.DBHelper.DeleteExistingOTPByID(OTPBody.CustomID)
	if delOtpErr != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "an error occurred",
		})
	}

	// fmt.Println("this is the user email being verified for MagicLinkOldUser.Email: ", MagicLinkOldUser.Email)
	// fmt.Println("this is the user email being verified for SentOTPData.Email: ", SentOTPData.Email)

	// check if user exists in the db by email in the data base
	user, dbEmailErr = dbFunc.DBHelper.FindByEmail(SentOTPData.Email)

	// fmt.Println("this is the user id being verified: ", user)

	if dbEmailErr != nil {
		if dbEmailErr.Error() == "error retrieving user" {
			// An error occurred while querying the database
			return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to check user existence",
			})
		}
	}

	role := "USER"

	// fmt.Println("this is the user id id verified: ", user.ID)

	// now generate cookies for this user
	authTokenString, refreshTokenString, csrfSecret, errJwt := myjwt.CreateNewTokens(ctx, strconv.FormatUint(uint64(user.ID), 10), role)
	if errJwt != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Error creating cookies"})
	}

	// set the cookies to these newly created jwt's and also set's the csrf token
	reqAuth.SetAuthAndRefreshCookies(ctx, authTokenString, refreshTokenString, csrfSecret)

	// Returning a success message (user successfully logged in using the magic link provided)
	return ctx.Status(http.StatusCreated).JSON(fiber.Map{
		"success": MagicLinkOldUser.Email + " Magic Login successfully",
	})
}
