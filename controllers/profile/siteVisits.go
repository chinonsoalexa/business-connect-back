package profile

import (
	dbFunc "business-connect/database/dbHelpFunc"
	"fmt"
	"net/http"

	Data "business-connect/models"

	"github.com/gofiber/fiber/v2"
)

func AddSiteVisit(ctx *fiber.Ctx) error {

	siteVisitErr := dbFunc.DBHelper.UpdateSiteVisits()
	if siteVisitErr != nil {
		if siteVisitErr.Error() == "failed to update site visit" {
			return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "site analytics error"})
		}
	}

	return ctx.Status(http.StatusOK).JSON(fiber.Map{"": ""})
}

func GetBusinessConnectUserByFingerprint(ctx *fiber.Ctx) error {

	fingerprintHash := ctx.Params("fingerprint")

	fingerprintHashErr := dbFunc.DBHelper.CreateBusinessConnectDeviceFingerprint(fingerprintHash)
	if fingerprintHashErr != nil {
		return ctx.Status(http.StatusNotFound).JSON(fiber.Map{
			"error": fingerprintHashErr,
		})
	}

	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"success": "successfully created user hash",
	})
}

func AddClickHistory(ctx *fiber.Ctx) error {

	var NewClick Data.BusinessConnectUserActivity

	// Check if there is an error binding the request
	if bindErr := ctx.BodyParser(&NewClick); bindErr != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to read body",
		})
	}

	// adding user click history
	clickErr := dbFunc.DBHelper.LogUserClickData(NewClick.FingerprintHash, NewClick.ProductID, NewClick.ActivityType, NewClick.Category, NewClick.TitleOrSearchQuery)

	// checking if there was an error comparing the orders
	if clickErr != nil {
		fmt.Println("click error: ", clickErr)
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to add new click history",
		})
	}

	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"success": "success",
	})
}
