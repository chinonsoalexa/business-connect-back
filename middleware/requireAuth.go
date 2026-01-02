package middleware

import (
	"strconv"
	"time"

	myjwt "business-connect/middleware/myjwt"

	"github.com/gofiber/fiber/v2"
)

// WEB AUTHENTICATION CHECKS

func NullifyCookie(ctx *fiber.Ctx) {
	// Set the auth token cookie to expire
	ctx.Cookie(&fiber.Cookie{
		Name:     "__BusinessConnect-Auth-Token",
		Value:    "",
		Expires:  time.Now().Add(-1000 * time.Hour),
		Path:     "/",
		Domain:   ".businessconnectt.com", // Ensure domain is set
		HTTPOnly: true,
		SameSite: "None",
		Secure:   true, // Set to true if using "SameSite: None"
	})

	// Set the refresh token cookie to expire
	ctx.Cookie(&fiber.Cookie{
		Name:     "__BusinessConnect-Refresh-Token",
		Value:    "",
		Expires:  time.Now().Add(-1000 * time.Hour),
		Path:     "/",
		Domain:   ".businessconnectt.com", // Ensure domain is set
		HTTPOnly: true,
		SameSite: "None",
		Secure:   true, // Set to true if using "SameSite: None"
	})

	// Set the CSRF token cookie to expire
	ctx.Cookie(&fiber.Cookie{
		Name:     "__X-Csrf-Token",
		Value:    "",
		Expires:  time.Now().Add(-1000 * time.Hour),
		Path:     "/",
		Domain:   ".businessconnectt.com", // Ensure domain is set
		HTTPOnly: true,
		SameSite: "None",
		Secure:   true, // Set to true if using "SameSite: None"
	})
	// Send cookie in the response header to be saved in the browser for auth
}

func SetAuthAndRefreshCookies(ctx *fiber.Ctx, authToken string, refreshToken string, csrfSecret string) {

	// Set the cookie to be sent
	ctx.Cookie(&fiber.Cookie{
		Name:     "__BusinessConnect-Auth-Token",
		Value:    authToken,
		Expires:  time.Now().Add(14 * 24 * time.Hour),
		Path:     "/",
		Domain:   ".businessconnectt.com",
		HTTPOnly: true,
		SameSite: "None",
		Secure:   true,
	})

	// Set the cookie to be sent
	ctx.Cookie(&fiber.Cookie{
		Name:     "__BusinessConnect-Refresh-Token",
		Value:    refreshToken,
		Expires:  time.Now().Add(14 * 24 * time.Hour),
		Path:     "/",
		Domain:   ".businessconnectt.com",
		HTTPOnly: true,
		SameSite: "None",
		Secure:   true,
	})

	// Set the cookie to be sent
	ctx.Cookie(&fiber.Cookie{
		Name:     "__X-Csrf-Token",
		Value:    csrfSecret,
		Expires:  time.Now().Add(14 * 24 * time.Hour),
		Path:     "/",
		Domain:   ".businessconnectt.com",
		HTTPOnly: true,
		SameSite: "None",
		Secure:   true,
	})
	// send cookie in the response header to be saved in the browser for auth
}

func WebRequireAuth(ctx *fiber.Ctx) error {
	// get the authentication cookie of req
	AuthCookie := ctx.Cookies("__BusinessConnect-Auth-Token")
	// fmt.Println("this is the auth token: ", AuthCookie)
	if AuthCookie == "" {
		// send an error when cookie was not found
		NullifyCookie(ctx)
		// return ctx.Redirect("https://businessconnectt.com/page/signin-new.html")
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "No Authentication cookie found"})
	}

	// get the refresh cookie of req
	RefreshCookie := ctx.Cookies("__BusinessConnect-Refresh-Token")
	// fmt.Println("this is the refresh token", RefreshCookie)
	if RefreshCookie == "" {
		NullifyCookie(ctx)
		// return ctx.Redirect("https://businessconnectt.com/page/signin-new.html")
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "No Refresh cookie found"})
	}

	// get the csrf secrete of req
	CsrfCookie := ctx.Cookies("__X-Csrf-Token")
	// fmt.Println("this is the csrf token", CsrfCookie)
	if CsrfCookie == "" {
		NullifyCookie(ctx)
		// return ctx.Redirect("https://businessconnectt.com/page/signin-new.html")
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "No Csrf secrets found"})
	}

	authTokenString, refreshTokenString, csrfSecret, jwtErr := myjwt.CheckAndRefreshTokens(AuthCookie, RefreshCookie, CsrfCookie)

	if jwtErr != nil && jwtErr.Error() == "Unauthorized" {
		NullifyCookie(ctx)
		// return ctx.Redirect("https://businessconnectt.com/page/signin-new.html")
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized attempt! JWT's not valid!"})
	}

	user_id, authError := myjwt.GrabUUID(authTokenString)
	// fmt.Println("this is the user id: ", user_id)
	if authError != nil {
		if authError.Error() == "error parsing auth token" || authError.Error() == "error fetching claims" {
			// Handle the error and return an appropriate response
			NullifyCookie(ctx)
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "an error occurred while parsing claims and auth token",
			})
		}
		NullifyCookie(ctx)
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized",
		})
	}

	// converting string to uint
	user_idd, uintErr := strconv.ParseUint(user_id, 10, 64)
	if uintErr != nil {
		// Handle the error if the conversion fails
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "an error occurred",
		})
	}
	// fmt.Println("this is the user id 2: ", user_id)

	ctx.Locals("user-id", uint(user_idd))

	// fmt.Println("user  id from the database: ", user_idd)
	// if we've made it this far, everything is valid!
	// And tokens have been refreshed if need-be
	SetAuthAndRefreshCookies(ctx, authTokenString, refreshTokenString, csrfSecret)

	// continue
	return ctx.Next()
}
