package post

import (
	"net/http"
	// "strconv"

	// SendEmail "business-connect/controllers/authentication/emails"
	// dbFunc "business-connect/database/dbHelpFunc"
	// Data "business-connect/models"

	"github.com/gofiber/fiber/v2"
)

func AddBusinessConnectProductComment(ctx *fiber.Ctx) error {

	// Parse product details
	// var ProductComment Data.CustomerReview

	// Check if there is an error binding the request
	// if bindErr := ctx.BodyParser(&ProductComment); bindErr != nil {
	// 	return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
	// 		"error": "Failed to read body",
	// 	})
	// }

	// // Save the product details to the database after successful image upload
	// savedProduct, err := dbFunc.DBHelper.SaveCustomerReview(ProductComment.ProducttID, ProductComment.Email, ProductComment.Name, ProductComment.Review, ProductComment.Rating)
	// if err != nil {
	// 	return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
	// 		"error": "error occurred while adding comment to database",
	// 	})
	// }

	// if ProductComment.AddEmail {
	// 	emailErr := dbFunc.DBHelper.AddEmailSubscriber(ProductComment.Email)
	// 	if emailErr != nil {
	// 		if emailErr.Error() == "user with email already exists" {
	// 			return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
	// 				"error": "user with email already exists",
	// 			})
	// 		}
	// 		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
	// 			"error": "error adding new email subscriber",
	// 		})
	// 	}

	// 	// send a confirmation email
	// 	emailErr = SendEmail.TodacWelcomeEmail(ProductComment.Email)
	// 	if emailErr != nil {
	// 		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
	// 			"error": "failed to send subscriber email",
	// 		})
	// 	}
	// }

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		// "success": savedProduct,
	})
}

func GetBusinessConnectProductCommentsByLimit(ctx *fiber.Ctx) error {
	// this is to get the page number to route to
	// pageNumber := 1
	// productNumber := 1
	// limitNumber := ctx.Params("idLimit")
	// ProductId := ctx.Params("proId")
	// eachPage := 4

	// if limitNumber != "" {
	// 	pageNumber, _ = strconv.Atoi(limitNumber)
	// }
	// if ProductId != "" {
	// 	productNumber, _ = strconv.Atoi(ProductId)
	// }

	// offset := (pageNumber - 1) * eachPage

	// productCommentRecords, reviewCount, productRecordsErr := dbFunc.DBHelper.GetCustomerReviewsByProduct(uint(productNumber), eachPage, offset)
	// if productRecordsErr != nil {
	// 	return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
	// 		"error": "failed to get product history",
	// 	})
	// }

	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		// "success": productCommentRecords,
		// "count":   reviewCount,
	})
}
