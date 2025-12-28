package profile

import (
	Data "business-connect/models"
	dbFunc "business-connect/database/dbHelpFunc"
	helperFunc "business-connect/paystack"

	"github.com/gofiber/fiber/v2"
)

func GetGroups(ctx *fiber.Ctx) error {
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
	groups, hasMore, postErr := dbFunc.DBHelper.GetAvailableGroups(limit, offset)
	if postErr != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch posts",
		})
	}

	// Return JSON
	return ctx.JSON(fiber.Map{
		"page":    page,
		"limit":   limit,
		"groups": groups,
		"user":    user,
		"hasMore": hasMore,
	})
}

type JoinGroupRequest struct {
	GroupPostID uint `json:"group_post_id"`
}

func JoinGroupHandler(ctx *fiber.Ctx) error {
	// Get current user from context
	userCtx := ctx.Locals("user-id")
	if userCtx == nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "user not logged in",
		})
	}
	user, ok := userCtx.(Data.User)
	if !ok {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to parse user",
		})
	}

	// Parse request body
	req := new(JoinGroupRequest)
	if err := ctx.BodyParser(req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	if req.GroupPostID == 0 {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "group_post_id is required",
		})
	}

	// Call DB helper
	participant, created, err := dbFunc.DBHelper.JoinGroup(user, req.GroupPostID)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to join group",
		})
	}

	if !created {
		return ctx.Status(fiber.StatusConflict).JSON(fiber.Map{
			"message":     "user already joined this group",
			"participant": participant,
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":     "successfully joined group",
		"participant": participant,
	})
}
