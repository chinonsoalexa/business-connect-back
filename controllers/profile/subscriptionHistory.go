package profile

// import (
// 	// "fmt"
// 	dbFunc "business-connect/database/dbHelpFunc"
// 	helperFunc "business-connect/paystack"
// 	"fmt"
// 	"math"
// 	"net/http"
// 	"strconv"

// 	"github.com/gofiber/fiber/v2"
// )

// type SubscriptionPaginationData struct {
// 	NextPage     int
// 	PreviousPage int
// 	CurrentPage  int
// 	TotalPages   int
// 	TwoAfter     int
// 	TwoBefore    int
// 	ThreeAfter   int
// 	Offset       int
// 	AllRecords   int64
// }

// func GetSubscriptionHistoryByLimit(ctx *fiber.Ctx) error {
// 	// get stored user id from request time line
// 	userId := ctx.Locals("user-id")

// 	if userId == nil {
// 		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
// 			"error": "failed to get user",
// 		})
// 	}

// 	user, uuidErr := helperFunc.PaystackHelper.FindByUuidFromLocal(userId)
// 	if uuidErr != nil {
// 		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
// 			"error": "failed to get user from request",
// 		})
// 	}

// 	// this is to get the page number to route to
// 	pageNumber := 1
// 	limitNumber := ctx.Params("idLimit")
// 	eachPage := 6

// 	if limitNumber != "" {
// 		pageNumber, _ = strconv.Atoi(limitNumber)
// 	}

// 	offset := (pageNumber - 1) * eachPage

// 	subscriptionHistory, totalRecords, subscriptionHistoryErr := dbFunc.DBHelper.GetSubscriptionHistoryByLimit(uint64(user.ID), eachPage, offset)
// 	if subscriptionHistoryErr != nil {
// 		if subscriptionHistoryErr.Error() == "transaction record not found" {
// 			return ctx.Status(http.StatusOK).JSON(fiber.Map{
// 				"error": "transaction record do not exist",
// 			})
// 		}
// 		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
// 			"error": "failed to get transaction history",
// 		})
// 	}

// 	totalPages := int(math.Ceil(float64(totalRecords) / float64(eachPage)))

// 	return ctx.Status(http.StatusOK).JSON(fiber.Map{
// 		"success": subscriptionHistory,
// 		"pagination": SubscriptionPaginationData{
// 			NextPage:     pageNumber + 1,
// 			PreviousPage: pageNumber - 1,
// 			CurrentPage:  pageNumber,
// 			TotalPages:   totalPages,
// 			TwoAfter:     pageNumber + 2,
// 			TwoBefore:    pageNumber - 2,
// 			ThreeAfter:   pageNumber + 4,
// 			AllRecords:   totalRecords,
// 		},
// 	})
// }

// func UpdateSubscriptionStatus(ctx *fiber.Ctx) error {
// 	// get stored user id from request time line
// 	userId := ctx.Locals("user-id")

// 	if userId == nil {
// 		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
// 			"error": "failed to get user",
// 		})
// 	}

// 	_, uuidErr := helperFunc.PaystackHelper.FindByUuidFromLocal(userId)
// 	if uuidErr != nil {
// 		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
// 			"error": "failed to get user from request",
// 		})
// 	}

// 	// this is to get the subscription id to update
// 	subscriptionID := ctx.Params("subscriptionID")
// 	var subscriptionIDNumb int

// 	if subscriptionID != "" {
// 		subscriptionIDNumb, _ = strconv.Atoi(subscriptionID)
// 	}

// 	subscriptionHistoryErr := dbFunc.DBHelper.UpdateSubscriptionStatus(uint64(subscriptionIDNumb))
// 	if subscriptionHistoryErr != nil {
// 		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
// 			"error": "failed to get transaction history",
// 		})
// 	}

// 	return ctx.Status(http.StatusOK).JSON(fiber.Map{
// 		"success": "successfully canceled subscription",
// 	})
// }

// func RechargeSubscriptionNow(ctx *fiber.Ctx) error {
// 	// get stored user id from request time line
// 	userId := ctx.Locals("user-id")

// 	if userId == nil {
// 		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
// 			"error": "failed to get user",
// 		})
// 	}

// 	user, uuidErr := helperFunc.PaystackHelper.FindByUuidFromLocal(userId)
// 	if uuidErr != nil {
// 		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
// 			"error": "failed to get user from request",
// 		})
// 	}

// 	// this is to get the subscription id to update
// 	subscriptionID := ctx.Params("subscriptionID")
// 	var subscriptionIDNumb int

// 	if subscriptionID != "" {
// 		subscriptionIDNumb, _ = strconv.Atoi(subscriptionID)
// 	}

// 	_, subscriptionHistoryErr := dbFunc.DBHelper.GetAndRechargeWithTransactionID(uint64(subscriptionIDNumb), user)
// 	fmt.Println("this is the subscription error: ", subscriptionHistoryErr)
// 	if subscriptionHistoryErr != nil {
// 		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
// 			"error": subscriptionHistoryErr.Error(),
// 		})
// 	}

// 	return ctx.Status(http.StatusOK).JSON(fiber.Map{
// 		"success": "successfully renewed subscription",
// 		"id":      subscriptionID,
// 	})
// }
