package authentication

import (
	"net/http"

	"github.com/gofiber/fiber/v2"

	cookieNul "business-connect/middleware"
)

func Logout(ctx *fiber.Ctx) error {

	cookieNul.NullifyCookie(ctx)

	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"success": "Logged out successfully",
	})
}
