package blog

import (
	"errors"
	"math"

	// "fmt"
	// "math"
	"net/http"
	"strconv"

	dbFunc "business-connect/database/dbHelpFunc"
	Data "business-connect/models"

	SendEmail "business-connect/controllers/authentication/emails"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type PaginationData struct {
	NextPage     int
	PreviousPage int
	CurrentPage  int
	TotalPages   int
	TwoAfter     int
	TwoBefore    int
	ThreeAfter   int
	Offset       int
	AllRecords   int64
}

func GetBlogPost(ctx *fiber.Ctx) error {
	// Extract order ID from path or query parameter (adjust based on your implementation)
	blogID, err := strconv.Atoi(ctx.Params("blogID"))
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid order ID",
		})
	}

	// Call the database helper function to retrieve the order
	orderHistory, err := dbFunc.DBHelper.GetBlogPostById(uint(blogID))
	if err != nil {
		// Handle potential errors based on your GetOrder function implementation
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ctx.Status(http.StatusNotFound).JSON(fiber.Map{
				"error": "order not found",
			})
		}
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to retrieve order",
		})
	}

	// Return successful response with order details
	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"success": *orderHistory,
	})
}

func GetBlogPosts(ctx *fiber.Ctx) error {
	// this is to get the page number to route to
	pageNumber := 1
	limitNumber := ctx.Params("idLimit")
	eachPage := 9

	if limitNumber != "" {
		pageNumber, _ = strconv.Atoi(limitNumber)
	}

	offset := (pageNumber - 1) * eachPage

	productRecords, totalRecords, productRecordsErr := dbFunc.DBHelper.GetBusinessConnectBlogByLimit(eachPage, offset)
	if productRecordsErr != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to get product history",
		})
	}

	totalPages := int(math.Ceil(float64(totalRecords) / float64(eachPage)))

	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"success": productRecords,
		"pagination": PaginationData{
			NextPage:     pageNumber + 1,
			PreviousPage: pageNumber - 1,
			CurrentPage:  pageNumber,
			TotalPages:   totalPages,
			TwoAfter:     pageNumber + 2,
			TwoBefore:    pageNumber - 2,
			ThreeAfter:   pageNumber + 4,
			AllRecords:   totalRecords,
		},
	})
}

func UpdateBusinessConnectBlog(ctx *fiber.Ctx) error {

	// // Get stored user id from request timeline
	// userId := ctx.Locals("user-id")

	// if userId == nil {
	// 	return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
	// 		"error": "failed to get user",
	// 	})
	// }

	// _, uuidErr := dbFunc.DBHelper.FindByUuidFromLocal(userId)

	// if uuidErr != nil {
	// 	return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
	// 		"error": uuidErr.Error(),
	// 	})
	// }

	type BlogToUpdate struct {
		BlogID           int    `json:"blog_id"`
		BlogTitle        string `json:"blog_title"`
		BlogDescription1 string `json:"blog_description1"`
		BlogDescription2 string `json:"blog_description2"`
		BlogCategory     string `json:"blog_category"`
	}

	var BlogUpdate BlogToUpdate
	var UpdatedBlog Data.Blog

	// Check if there is an error binding the request
	if bindErr := ctx.BodyParser(&BlogUpdate); bindErr != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to read body",
		})
	}

	UpdatedBlog.Title = BlogUpdate.BlogTitle
	UpdatedBlog.Description1 = BlogUpdate.BlogDescription1
	UpdatedBlog.Description2 = BlogUpdate.BlogDescription2
	UpdatedBlog.BlogCategory = BlogUpdate.BlogCategory
	// UpdatedBlog.ProductStock = BlogUpdate.ProductStock

	// Call the database helper function to retrieve the order
	err := dbFunc.DBHelper.UpdateBusinessConnectBlog(UpdatedBlog, uint(BlogUpdate.BlogID))
	if err != nil {
		// Handle potential errors based on your GetOrder function implementation
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ctx.Status(http.StatusNotFound).JSON(fiber.Map{
				"error": "error retrieving blog",
			})
		}
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to retrieve blog",
		})
	}

	// Return successful response with order details
	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"success": "successfully updated blog",
	})
}

func DeleteBusinessConnectBlog(ctx *fiber.Ctx) error {

	// // Get stored user id from request timeline
	// userId := ctx.Locals("user-id")

	// if userId == nil {
	// 	return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
	// 		"error": "failed to get user",
	// 	})
	// }

	// _, uuidErr := dbFunc.DBHelper.FindByUuidFromLocal(userId)

	// if uuidErr != nil {
	// 	return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
	// 		"error": uuidErr.Error(),
	// 	})
	// }

	// Extract order ID from path or query parameter (adjust based on your implementation)
	blogID, err := strconv.Atoi(ctx.Params("blogID"))
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid order ID",
		})
	}

	// Call the database helper function to retrieve the order
	err2 := dbFunc.DBHelper.DeleteBusinessConnectBlog(uint(blogID))
	if err2 != nil {
		// Handle potential errors based on your GetOrder function implementation
		if errors.Is(err2, gorm.ErrRecordNotFound) {
			return ctx.Status(http.StatusNotFound).JSON(fiber.Map{
				"error": "blog not found",
			})
		}
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to retrieve blog",
		})
	}

	// Return successful response with order details
	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"success": "successfully deleted blog",
	})
}

func AddBusinessConnectBlogComment(ctx *fiber.Ctx) error {

	// Parse product details
	var BlogComment Data.CustomerBlogReview

	// Check if there is an error binding the request
	if bindErr := ctx.BodyParser(&BlogComment); bindErr != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to read body",
		})
	}

	// Save the product details to the database after successful image upload
	savedProduct, err := dbFunc.DBHelper.SaveCustomerBlogReview(BlogComment.BlogID, BlogComment.Email, BlogComment.Name, BlogComment.Review, BlogComment.Rating)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "error occurred while adding comment to database",
		})
	}

	if BlogComment.AddEmail {
		emailErr := dbFunc.DBHelper.AddEmailSubscriber(BlogComment.Email)
		if emailErr != nil {
			if emailErr.Error() == "user with email already exists" {
				return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
					"error": "user with email already exists",
				})
			}
			return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"error": "error adding new email subscriber",
			})
		}

		// send a confirmation email
		emailErr = SendEmail.TodacWelcomeEmail(BlogComment.Email)
		if emailErr != nil {
			return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"error": "failed to send subscriber email",
			})
		}
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": savedProduct,
	})
}

func GetBusinessConnectBlogCommentsByLimit(ctx *fiber.Ctx) error {
	// this is to get the page number to route to
	pageNumber := 1
	productNumber := 1
	limitNumber := ctx.Params("idLimit")
	ProductId := ctx.Params("proId")
	eachPage := 4

	if limitNumber != "" {
		pageNumber, _ = strconv.Atoi(limitNumber)
	}
	if ProductId != "" {
		productNumber, _ = strconv.Atoi(ProductId)
	}

	offset := (pageNumber - 1) * eachPage

	productCommentRecords, reviewCount, productRecordsErr := dbFunc.DBHelper.GetCustomerBlogReviewsByBlogPost(uint(productNumber), eachPage, offset)
	if productRecordsErr != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to get product history",
		})
	}

	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"success": productCommentRecords,
		"count":   reviewCount,
	})
}
