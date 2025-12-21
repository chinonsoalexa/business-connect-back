package retrieveTransaction

import (
	// "net/http"
	// helperFunc "business-connect/paystack"

	"github.com/gofiber/fiber/v2"
)

func RetrieveLatestTransactionByUserID(ctx *fiber.Ctx) error {
	// get stored user id from request time line
	// userId := ctx.Locals("user-id")

	// if userId == nil {
	// 	return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
	// 		"error": "failed to get user",
	// 	})
	// }

	// user, uuidErr := helperFunc.PaystackHelper.FindByUuidFromLocal(userId)
	// if uuidErr != nil {
	// 	return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
	// 		"error": "failed to get user from request",
	// 	})
	// }

	// transactionDetail, latestTransactionService, walletBalance, TransErr := helperFunc.PaystackHelper.GetLatestTransactionHistory(user.ID)

	// if TransErr != nil {
	// 	return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "error getting transaction details"})
	// }

	// // var amount float64

	// // if transactionDetail.TransactionType == "paystack" {
	// // 	amount = ConvertToRealAmount(walletBalance)
	// // }

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": "transactionDetail",
		"service": "latestTransactionService",
		"balance": "walletBalance",
	})
}

func ConvertToRealAmount(originalAmount int) float64 {
	percentage := 1.5
	additionalCharge := 20.0

	// Calculate 1.5% of the original number
	transactionCharge := (percentage / 100.00) * float64(originalAmount)

	// Add additional charge (20 in this case)
	updatedTransactionCharge := transactionCharge + additionalCharge

	return updatedTransactionCharge
}
