package home

import (
	"net/http"

	dbFunc "business-connect/database/dbHelpFunc"

	"github.com/gofiber/fiber/v2"
)

func GetBusinessConnectAnalytics(ctx *fiber.Ctx) error {

	// Call the database helper function to retrieve the order
	siteVisits, err2 := dbFunc.DBHelper.GetLast12DaysSiteVisits()
	if err2 != nil {
		// Handle potential errors based on the get site visits analytics function implementation
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to retrieve site visits",
		})
	}

	// Call the database helper function to retrieve the order
	userAnalytics, err3 := dbFunc.DBHelper.GetAnalyticsData()
	if err3 != nil {
		// Handle potential errors based on the get site visits analytics function implementation
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to retrieve user sales analytics",
		})
	}

	// Return successful response with order details
	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"site_visits":    siteVisits,
		"user_analytics": userAnalytics,
	})
}
