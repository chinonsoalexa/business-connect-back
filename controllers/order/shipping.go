package order

import (
	"net/http"

	dbFunc "business-connect/database/dbHelpFunc"
	Data "business-connect/models"

	"github.com/gofiber/fiber/v2"
)

func SetShippingPricePerKm(ctx *fiber.Ctx) error {

	var ShippingFee Data.ShippingFees

	// Check if there is an error binding the request
	if bindErr := ctx.BodyParser(&ShippingFee); bindErr != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to read body",
		})
	}

	// adding new order history
	shippingErr := dbFunc.DBHelper.UpsertShippingFee(ShippingFee.ShippingFeePerKm,
		ShippingFee.ShippingFeeGreater, ShippingFee.ShippingFeeLess, ShippingFee.StoreLatitude, ShippingFee.StoreLongitude,
		ShippingFee.StoreState, ShippingFee.StoreCity, ShippingFee.StateISO, ShippingFee.CalculateUsingKg)

	// checking if there was an error comparing the orders
	if shippingErr != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to update shipping fee",
		})
	}

	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"success": "successfully updated shipping fee",
	})
}

func GetShippingPricePerKm(ctx *fiber.Ctx) error {

	// adding new order history
	shippingFee, shippingErr := dbFunc.DBHelper.GetShippingFee()

	// checking if there was an error comparing the orders
	if shippingErr != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to update shipping fee",
		})
	}

	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"success": shippingFee,
	})
}
