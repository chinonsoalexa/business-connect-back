package profile

import (
	dbFunc "business-connect/database/dbHelpFunc"
	helperFunc "business-connect/paystack"
	"fmt"

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
	groups, hasMore, postErr := dbFunc.DBHelper.GetRecommendedGroupsWithFallback(user, limit, offset)
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

func GetGroupsOpen(ctx *fiber.Ctx) error {
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
	groups, hasMore, postErr := dbFunc.DBHelper.GetOpenRecommendedGroups(limit, offset)
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
		"hasMore": hasMore,
	})
}

type JoinGroupRequest struct {
	GroupPostID uint `json:"group_post_id"`
}

func JoinGroupHandler(ctx *fiber.Ctx) error {
	// 1️⃣ Get current logged in user
	userIDInterface := ctx.Locals("user-id")
	userID, ok := userIDInterface.(uint)
	if !ok {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "not logged in"})
	}

	user, uuidErr := helperFunc.PaystackHelper.FindByUuidFromLocal(userID)
	if uuidErr != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "failed to get user from request",
		})
	}

	// 2️⃣ Parse request body
	var req JoinGroupRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	if req.GroupPostID == 0 {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "group_post_id is required"})
	}

	// 3️⃣ Call DB helper to join group
	participant, created, membersCount, err := dbFunc.DBHelper.JoinGroup(user, req.GroupPostID)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	if !created {
		return ctx.Status(fiber.StatusConflict).JSON(fiber.Map{
			"message":       "user already joined this group",
			"participant":   participant,
			"members_count": membersCount,
		})
	}

	// 4️⃣ Create notification for group owner
	group, err := dbFunc.DBHelper.GetGroupByID(req.GroupPostID)
	if err != nil {
		// Non-fatal, just log and continue
		fmt.Println("Failed to fetch group for notification:", err)
	} else {
		err = dbFunc.DBHelper.CreateNotification(
			group.UserID,                     // recipient = group owner
			user.ID,                           // actor = user who joined
			"group_join",                      // type
			fmt.Sprintf("%s has joined your group %s!", user.FullName, group.Title),
			user.ProfilePhotoURL,              // actor avatar
		)
		if err != nil {
			fmt.Println("Failed to create group join notification:", err)
		}
	}

	// 5️⃣ Return success response
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":       "successfully joined group",
		"participant":   participant,
		"members_count": membersCount,
	})
}

