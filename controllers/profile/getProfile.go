package profile

import (
	dbFunc "business-connect/database/dbHelpFunc"
	helperFunc "business-connect/paystack"

	"github.com/gofiber/fiber/v2"
)

func GetMyProfile(ctx *fiber.Ctx) error {
	// 1️⃣ Get logged-in user
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

	// 2️⃣ Pagination for posts and groups
	page := ctx.QueryInt("page", 1)
	limit := ctx.QueryInt("limit", 10)
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 50 {
		limit = 20
	}
	offset := (page - 1) * limit

	// 3️⃣ Fetch profile data
	profile, err := dbFunc.DBHelper.GetUserProfile(user.UniqueName, limit, offset)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	// 4️⃣ Return JSON
	return ctx.JSON(fiber.Map{
		"user":     profile.User,
		"timeline": profile.Posts,
		"friends":  profile.Connections,
		"groups":   profile.Groups,
		"about":    profile.About,
	})
}

func GetOthersProfile(ctx *fiber.Ctx) error {
	// 1️⃣ Get logged-in user
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

	// 2️⃣ Pagination for posts and groups
	page := ctx.QueryInt("page", 1)
	limit := ctx.QueryInt("limit", 10)
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 50 {
		limit = 20
	}
	offset := (page - 1) * limit

	// 3️⃣ Fetch profile data
	profile, err := dbFunc.DBHelper.GetUserProfile(user.UniqueName, limit, offset)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	// 4️⃣ Return JSON
	return ctx.JSON(fiber.Map{
		"user":     profile.User,
		"timeline": profile.Posts,
		"friends":  profile.Connections,
		"groups":   profile.Groups,
		"about":    profile.About,
	})
}

func GetProfileOpen(ctx *fiber.Ctx) error {
	// 1️⃣ Get target user ID from query param
	uniqueName := ctx.Params("name")

	// 2️⃣ Pagination for posts and groups
	page := ctx.QueryInt("page", 1)
	limit := ctx.QueryInt("limit", 10)
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 50 {
		limit = 20
	}
	offset := (page - 1) * limit

	// 3️⃣ Fetch profile data
	profile, err := dbFunc.DBHelper.GetUserProfile(uniqueName, limit, offset)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	// 4️⃣ Return JSON (no sensitive info)
	return ctx.JSON(fiber.Map{
		"user":     profile.User,
		"timeline": profile.Posts,
		"friends":  profile.Connections,
		"groups":   profile.Groups,
		"about":    profile.About,
	})
}
