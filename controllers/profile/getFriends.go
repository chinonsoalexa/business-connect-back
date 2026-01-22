package profile

import (
	dbFunc "business-connect/database/dbHelpFunc"
	helperFunc "business-connect/paystack"
	"fmt"

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
	friends, hasMore, postErr := dbFunc.DBHelper.GetRecommendedUsersWithFallback(user, limit, offset)
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
	friends, hasMore, postErr := dbFunc.DBHelper.GetOpenRecommendedUsers(limit, offset)
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
func GenerateWhatsAppLinks(phoneNumber, senderName, receiverName, businessName, device string) (string, string) {
	var message string
	if businessName != "" {
		message = BuildBusinessMessage(senderName, receiverName, businessName)
	} else {
		message = BuildProfileMessage(senderName, receiverName)
	}

	encoded := url.QueryEscape(message)

	// Mobile/Tablet app-first
	return fmt.Sprintf("whatsapp://send?phone=%s&text=%s", phoneNumber, encoded), fmt.Sprintf("https://wa.me/%s?text=%s", phoneNumber, encoded)
	// return encoded
}

func ConnectFriend(ctx *fiber.Ctx) error {
	userId := ctx.Locals("user-id")
	if userId == nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "not logged in",
		})
	}

	user, uuidErr := helperFunc.PaystackHelper.FindByUuidFromLocal(userId)
	if uuidErr != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "failed to get user from request",
		})
	}

	var req ConnectRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	if req.UserID == 0 {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "user_id is required",
		})
	}

	// Follow/connect
	err := dbFunc.DBHelper.FollowUser(user.ID, req.UserID)
	if err != nil {
		switch err {
		case dbFunc.ErrAlreadyFollowing:
			// User is already connected
			return ctx.JSON(fiber.Map{
				"status": "already_following",
			})
		case dbFunc.ErrSelfFollow:
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status": "self_follow",
				"error":  err.Error(),
			})
		default:
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status": "error",
				"error":  err.Error(),
			})
		}
	}

	// Lookup receiver
	userRec, uuidErr := helperFunc.PaystackHelper.FindByUuidFromLocal(req.UserID)
	if uuidErr != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "failed to get user from request",
		})
	}

	// Determine device type
	device := req.Device
	if device != "mobile" && device != "tablet" && device != "desktop" {
		device = "desktop"
	}

	// Generate WhatsApp links
	appLink, webLink := GenerateWhatsAppLinks(
		userRec.PhoneNumber,
		user.FullName,
		userRec.FullName,
		userRec.BusinessName,
		device,
	)

	return ctx.JSON(fiber.Map{
		"status":  "ok",
		"appLink": appLink,
		"webLink": webLink,
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

func SearchUsers(ctx *fiber.Ctx) error {
    // Get query and pagination
    query := ctx.Query("q", "")
    page := ctx.QueryInt("page", 1)
    limit := ctx.QueryInt("limit", 10)

    if page < 1 {
        page = 1
    }
    if limit < 1 || limit > 50 {
        limit = 20
    }
    offset := (page - 1) * limit

    if query == "" {
        return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "Query parameter 'q' is required",
        })
    }

    // Call database helper
    users, hasMore, err := dbFunc.DBHelper.SearchUsers(query, limit, offset)
    if err != nil {
        return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": "Failed to search users",
        })
    }

    return ctx.JSON(fiber.Map{
        "page":    page,
        "limit":   limit,
        "query":   query,
        "friends": users,
        "hasMore": hasMore,
    })
}
