package profile

import (
	dbFunc "business-connect/database/dbHelpFunc"
	helperFunc "business-connect/paystack"
	"fmt"

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

func GetFriendsOpen(ctx *fiber.Ctx) error {
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
	friends, hasMore, postErr := dbFunc.DBHelper.GetUsersToConnectOpen(limit, offset)
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
		"hasMore": hasMore,
	})
}

type ConnectRequest struct {
	UserID uint `json:"user_id"` // the user you want to connect to
}

func ConnectFriend(ctx *fiber.Ctx) error {
	// 1️⃣ Get logged in user
	userIDInterface := ctx.Locals("user-id")
	userID, ok := userIDInterface.(uint)
	if !ok {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "not logged in"})
	}

	user, uuidErr := helperFunc.PaystackHelper.FindByUuidFromLocal(userID)
	if uuidErr != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "failed to get user from request"})
	}

	// 2️⃣ Parse request
	var req ConnectRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	if req.UserID == 0 {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "user_id is required"})
	}

	// 3️⃣ Prevent self-follow
	if user.ID == req.UserID {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "cannot connect to yourself"})
	}

	// 4️⃣ Follow
	err := dbFunc.DBHelper.FollowUser(user.ID, req.UserID)
	if err != nil {
		// Could return 409 if already following
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	// 5️⃣ Success
	return ctx.JSON(fiber.Map{
		"message": fmt.Sprintf("You are now following user %d", req.UserID),
	})
}

func IncrementPostViewHandler(ctx *fiber.Ctx) error {
	// Get stored user id from request context
	// userId := ctx.Locals("user-id")
	// if userId == nil {
	// 	return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
	// 		"error": "failed to get user",
	// 	})
	// }

	// // Validate user (optional but recommended)
	// _, uuidErr := helperFunc.PaystackHelper.FindByUuidFromLocal(userId)
	// if uuidErr != nil {
	// 	return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
	// 		"error": "failed to get user from request",
	// 	})
	// }

	// Get post ID from params
	postID, err := ctx.ParamsInt("postId")
	if err != nil || postID <= 0 {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid post id",
		})
	}

	// Increment view count
	if err := dbFunc.DBHelper.IncrementPostView(uint(postID)); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to increment post view",
		})
	}

	return ctx.JSON(fiber.Map{
		"success": true,
		"post_id": postID,
		"type":    "view",
	})
}

func IncrementPostClickHandler(ctx *fiber.Ctx) error {
	// // Get stored user id from request context
	// userId := ctx.Locals("user-id")
	// if userId == nil {
	// 	return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
	// 		"error": "failed to get user",
	// 	})
	// }

	// // Validate user
	// _, uuidErr := helperFunc.PaystackHelper.FindByUuidFromLocal(userId)
	// if uuidErr != nil {
	// 	return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
	// 		"error": "failed to get user from request",
	// 	})
	// }

	// Get post ID from params
	postID, err := ctx.ParamsInt("postId")
	if err != nil || postID <= 0 {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid post id",
		})
	}

	// Increment click count
	if err := dbFunc.DBHelper.IncrementPostClick(uint(postID)); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to increment post click",
		})
	}

	return ctx.JSON(fiber.Map{
		"success": true,
		"post_id": postID,
		"type":    "click",
	})
}
