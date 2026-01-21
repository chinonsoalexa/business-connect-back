package profile

import (
	dbFunc "business-connect/database/dbHelpFunc"
	helperFunc "business-connect/paystack"
	// "fmt"
	"net/url"

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

	user.PhoneNumber = "-"

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
	UserID uint   `json:"user_id"` // the user you want to connect to
	Device string `json:"device"`  // "mobile", "tablet", "desktop"
}

// Example WhatsApp link generator
func GenerateWhatsAppLinks(phoneNumber, senderName, receiverName, businessName, device string) string {
	var message string
	if businessName != "" {
		message = BuildBusinessMessage(senderName, receiverName, businessName)
	} else {
		message = BuildProfileMessage(senderName, receiverName)
	}

	encoded := url.QueryEscape(message)

	// if device == "desktop" {
	// 	return fmt.Sprintf("https://wa.me/%s?text=%s", phoneNumber, encoded)
	// }

	// Mobile/Tablet app-first
	// return fmt.Sprintf("whatsapp://send?phone=%s&text=%s", phoneNumber, encoded)
	return encoded
}

// ConnectFriend handles the friend connection + WhatsApp redirect
func ConnectFriend(ctx *fiber.Ctx) error {
	// 1️⃣ Get logged in user
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

	// 2️⃣ Parse request
	var req ConnectRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).SendString("Invalid request body")
	}

	if req.UserID == 0 {
		return ctx.Status(fiber.StatusBadRequest).SendString("user_id is required")
	}

	// 3️⃣ Prevent self-follow
	if user.ID == req.UserID {
		return ctx.Status(fiber.StatusBadRequest).SendString("Cannot connect to yourself")
	}

	// 4️⃣ Follow/connect
	if err := dbFunc.DBHelper.FollowUser(user.ID, req.UserID); err != nil {
		return ctx.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	// 5️⃣ Lookup the receiver
	userRec, uuidErr := helperFunc.PaystackHelper.FindByUuidFromLocal(req.UserID)
	if uuidErr != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "failed to get user from request",
		})
	}

	// 6️⃣ Determine device type, default to desktop if invalid
	device := req.Device
	if device != "mobile" && device != "tablet" && device != "desktop" {
		device = "desktop"
	}

	// 7️⃣ Generate WhatsApp link
	whatsappLink := GenerateWhatsAppLinks(
		userRec.PhoneNumber,
		user.FullName,
		userRec.FullName,
		userRec.BusinessName,
		device,
	)

	// 8️⃣ Redirect the user directly to WhatsApp (app-first or web)
	return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
		"message": whatsappLink,
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
