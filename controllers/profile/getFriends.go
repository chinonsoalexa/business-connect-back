package profile

import (
	dbFunc "business-connect/database/dbHelpFunc"
	helperFunc "business-connect/paystack"
	"fmt"
	"net/url"

	"github.com/gofiber/fiber/v2"
	Data "business-connect/models"
	conn "business-connect/database"
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

	user, err := helperFunc.PaystackHelper.FindByUuidFromLocal(userId)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "failed to get user",
		})
	}

	var req ConnectRequest
	if err := ctx.BodyParser(&req); err != nil || req.UserID == 0 {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request",
		})
	}

	// 🔒 CHECK CONNECT LIMIT
	allowed, remaining, err := dbFunc.DBHelper.CanAndConsumeConnect(user.ID)
	if err != nil {
		return ctx.Status(500).JSON(fiber.Map{
			"error": "limit check failed",
		})
	}

	if !allowed {
		return ctx.Status(403).JSON(fiber.Map{
			"status": "limit_reached",
			"message": "Daily connect limit reached. Refer a friend to unlock 5 more connects.",
		})
	}

	// Proceed with follow
	err = dbFunc.DBHelper.FollowUser(user.ID, req.UserID)
	if err != nil {
		if err == dbFunc.ErrAlreadyFollowing {
			return ctx.JSON(fiber.Map{
				"status": "already_following",
			})
		}
		return ctx.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Receiver
	userRec, err := helperFunc.PaystackHelper.FindByUuidFromLocal(req.UserID)
	if err != nil {
		return ctx.Status(400).JSON(fiber.Map{
			"error": "receiver not found",
		})
	}

	err = dbFunc.DBHelper.CreateNotification(
		userRec.ID,          // recipient
		user.ID,             // actor
		"connect",           // type
		fmt.Sprintf("%s has connected with you!", user.FullName),
		user.ProfilePhotoURL, // avatar
	)

	device := req.Device
	if device != "mobile" && device != "tablet" && device != "desktop" {
		device = "desktop"
	}

	appLink, webLink := GenerateWhatsAppLinks(
		userRec.PhoneNumber,
		user.FullName,
		userRec.FullName,
		userRec.BusinessName,
		device,
	)

	return ctx.JSON(fiber.Map{
		"status":          "ok",
		"remainingSlots":  remaining,
		"appLink":         appLink,
		"webLink":         webLink,
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

// GetNotifications Handler with Load More / Pagination
func GetNotifications(ctx *fiber.Ctx) error {
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

	var notifications []Data.Notification
	var total int64

	// Count total notifications
	if err := conn.DB.Model(&Data.Notification{}).
		Where("user_id = ?", user.ID).
		Count(&total).Error; err != nil {
		return ctx.Status(500).JSON(fiber.Map{"error": "failed to count notifications"})
	}

	// Fetch notifications with limit + offset, newest first
	if err := conn.DB.Where("user_id = ?", user.ID).
		Order("created_at desc").
		Limit(limit).
		Offset(offset).
		Find(&notifications).Error; err != nil {
		return ctx.Status(500).JSON(fiber.Map{"error": "failed to fetch notifications"})
	}

	hasMore := int64(page*limit) < total

	// Default message if no notifications
	if len(notifications) == 0 {
		return ctx.JSON(fiber.Map{
			"page":          page,
			"limit":         limit,
			"hasMore":       false,
			"total":         0,
			"notifications": []map[string]string{
				{
					"message":   "No notifications yet! Connect with friends, post a product, or check out items to purchase.",
					"avatarURL": "/assets/images/default-notif.png", // default placeholder
					"type":      "info",
				},
			},
		})
	}

	return ctx.JSON(fiber.Map{
		"page":          page,
		"limit":         limit,
		"hasMore":       hasMore,
		"total":         total,
		"notifications": notifications,
		"user":          user,
	})
}
