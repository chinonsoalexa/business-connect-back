package profile

import (
	dbFunc "business-connect/database/dbHelpFunc"
	helperFunc "business-connect/paystack"

	"github.com/gofiber/fiber/v2"
)

func GetFriends(ctx *fiber.Ctx) error {
	// Get stored user id from request context
	userId := ctx.Locals("user-id")
	if userId == nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "failed to get user",
		})
	}

	user, uuidErr := helperFunc.PaystackHelper.FindByUuidFromLocal(userId)
	if uuidErr != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "failed to get user from request",
		})
	}

	// Default pagination values
	page := ctx.QueryInt("page", 1)
	limit := ctx.QueryInt("limit", 10)

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 50 {
		limit = 20
	}

	offset := (page - 1) * limit

	// Fetch posts using limit+1 for hasMore
	friends, hasMore, postErr := dbFunc.DBHelper.GetUsersToConnect(user.ID, limit, offset)
	if postErr != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch posts",
		})
	}

	// Return JSON
	return ctx.JSON(fiber.Map{
		"page":    page,
		"limit":   limit,
		"friends": friends,
		"user":    user,
		"hasMore": hasMore,
	})
}

type ConnectRequest struct {
	UserID uint `json:"user_id"` // the user you want to connect to
}

func ConnectFriend(ctx *fiber.Ctx) error {
	// Get current logged in user-id
	userId := ctx.Locals("user-id")
	if userId == nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "not logged in"})
	}

	user, uuidErr := helperFunc.PaystackHelper.FindByUuidFromLocal(userId)
	if uuidErr != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "failed to get user from request",
		})
	}

	var req ConnectRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	if req.UserID == 0 {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "user_id is required"})
	}

	err := dbFunc.DBHelper.ConnectToUser(user.ID, req.UserID)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.JSON(fiber.Map{"message": "connection request sent"})
}
