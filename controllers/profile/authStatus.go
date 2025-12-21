package profile

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
)

func CheckAuthStatus(ctx *fiber.Ctx) error {

	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"success": "user authenticated",
	})
}
