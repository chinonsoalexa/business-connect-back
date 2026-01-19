package profile

import (
	dbFunc "business-connect/database/dbHelpFunc"
	helperFunc "business-connect/paystack"

	Data "business-connect/models"

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

func UpdateMyProfile(ctx *fiber.Ctx) error {
	// 1️⃣ Get logged-in user ID from middleware
	userId := ctx.Locals("user-id")
	if userId == nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized",
		})
	}

	// 2️⃣ Resolve user from UUID
	user, uuidErr := helperFunc.PaystackHelper.FindByUuidFromLocal(userId)
	if uuidErr != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "failed to resolve user",
		})
	}

	// 3️⃣ Parse request body
	var req Data.UpdateProfileRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	// 4️⃣ Update user profile
	err := dbFunc.DBHelper.UpdateUserProfile(
		user.ID,
		req.FullName,
		req.BusinessName,
		req.BioDescription,
		req.PhoneNumber,
		req.Address,
		req.State,
		req.Language,
		req.PreferredLanguage,
		req.PreferredCurrency,
		req.ProfilePhotoURL,
		req.CoverPhotoURL,
	)

	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// 5️⃣ Fetch updated profile (fresh data)
	profile, err := dbFunc.DBHelper.GetUserProfile(user.UniqueName, 10, 0)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "profile updated but failed to reload",
		})
	}

	// 6️⃣ Return updated profile
	return ctx.JSON(fiber.Map{
		"message": "profile updated successfully",
		"user":    profile.User,
	})
}

func GetOthersProfile(ctx *fiber.Ctx) error {
	// 1️⃣ Get logged-in user
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
