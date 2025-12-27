package authentication

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"

	dbFunc "business-connect/database/dbHelpFunc"
	emailValid "business-connect/email"
	reqAuth "business-connect/middleware"
	myjwt "business-connect/middleware/myjwt"
	Data "business-connect/models"
	helperFunc "business-connect/paystack"
)

func SignUp(ctx *fiber.Ctx) error {
	var req Data.SignUpRequest

	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// 1️⃣ Validate input FIRST
	if err := validateSignup(&req); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// 2️⃣ Check existing user
	existingUser, err := dbFunc.DBHelper.FindByEmail(req.Email)
	if err == nil {
		if !existingUser.EmailVerified {
			// let's send token to user to verify the user with email ID
			emailErr := EmailVerification(existingUser.FullName, existingUser.Email)

			if emailErr != nil {
				fmt.Println("this is the email error:", emailErr)
				return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
					"error": "email verification failed",
				})
			}
			return ctx.Status(http.StatusConflict).JSON(fiber.Map{
				"error": "User exists but email is not verified",
			})
		}

		return ctx.Status(http.StatusConflict).JSON(fiber.Map{
			"error": "User already exists, please login",
		})
	}

	// 3️⃣ Disposable email check
	err = emailValid.LoadDisposableList("fakeEmails.json")
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Email validation service error",
		})
	}

	result, err := emailValid.ValidateEmail(req.Email)
	if err != nil || result.RiskScore >= 75 {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid or disposable email address",
		})
	}

	// 4️⃣ Build DB model (explicit mapping)
	newUser := Data.User{
		FullName:     req.FullName,
		BusinessName: req.BusinessName,
		Email:        result.Normalized,
		PhoneNumber:  req.PhoneNumber,
		State:        req.State,
		City:         req.City,
		Longitude:    req.Longitude,
		Latitude:     req.Latitude,
		UserType:     "USER",
	}

	// 5️⃣ Create user
	createdUser, err := dbFunc.DBHelper.CreateNewUser(newUser)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create user",
		})
	}

	// 6️⃣ Send verification email
	if err := EmailVerification(createdUser.FullName, createdUser.Email); err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to send verification email",
		})
	}

	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"success": "Verification email sent to " + createdUser.Email,
	})
}

func validateSignup(req *Data.SignUpRequest) error {
	// Required fields
	if strings.TrimSpace(req.FullName) == "" {
		return errors.New("full name is required")
	}

	if strings.TrimSpace(req.BusinessName) == "" {
		return errors.New("business name is required")
	}

	if strings.TrimSpace(req.Email) == "" {
		return errors.New("email is required")
	}

	if strings.TrimSpace(req.Password) == "" {
		return errors.New("password is required")
	}

	// Email format (basic)
	if !regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`).MatchString(req.Email) {
		return errors.New("invalid email format")
	}

	// Password rules
	if len(req.Password) < 8 {
		return errors.New("password must be at least 8 characters long")
	}

	if !regexp.MustCompile(`[A-Z]`).MatchString(req.Password) {
		return errors.New("password must contain at least one uppercase letter")
	}

	if !regexp.MustCompile(`[0-9]`).MatchString(req.Password) {
		return errors.New("password must contain at least one number")
	}

	// Optional but recommended
	if req.Longitude < -180 || req.Longitude > 180 {
		return errors.New("invalid longitude")
	}

	if req.Latitude < -90 || req.Latitude > 90 {
		return errors.New("invalid latitude")
	}

	return nil
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

	UserBodyReturn.EmailVerified = true

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

	role := "USER"

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
	emailErr = EmailVerification(existingUser.FullName, OTPBody.Email)

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
