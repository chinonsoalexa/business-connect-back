package profile

import (
	dbFunc "business-connect/database/dbHelpFunc"
	Data "business-connect/models"
	helperFunc "business-connect/paystack"
	"fmt"

	// "fmt"
	"math"
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
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

func GetBusinessConnectProductsByLimit(ctx *fiber.Ctx) error {
	var totalRecords int64
	var productRecords []Data.Product
	var err error

	// Parse query params
	category := ctx.Query("category", "na")
	limit := ctx.QueryInt("limit", 12)
	sortField := ctx.Query("sort", "created_at")
	sortOrder := ctx.Query("order", "asc")
	page := ctx.Params("page", "1")
	pageNumber, _ := strconv.Atoi(page)

	offset := (pageNumber - 1) * limit

	if category != "na" {
		productRecords, totalRecords, err = dbFunc.DBHelper.GetProductsByCategory(category, limit, offset, sortField, sortOrder)
	} else {
		productRecords, totalRecords, err = dbFunc.DBHelper.GetProductsAll(limit, offset, sortField, sortOrder)
	}

	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch products"})
	}

	totalPages := int(math.Ceil(float64(totalRecords) / float64(limit)))

	return ctx.JSON(fiber.Map{
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

func GetBusinessConnectAdminProductsByLimit(ctx *fiber.Ctx) error {
	// this is to get the page number to route to
	pageNumber := 1
	limitNumber := ctx.Params("idLimit")
	eachPage := 12

	if limitNumber != "" {
		pageNumber, _ = strconv.Atoi(limitNumber)
	}

	offset := (pageNumber - 1) * eachPage

	productRecords, totalRecords, productRecordsErr := dbFunc.DBHelper.GetBusinessConnectAdminProductsByLimit(eachPage, offset)
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

func GetBusinessConnectProductByID(ctx *fiber.Ctx) error {

	productID := ctx.Params("id")
	convertedTransactionID, err := strconv.ParseUint(productID, 10, 64)
	if err != nil {
		// Handle error
		fmt.Println(err)
	}

	productDetail, TransErr := dbFunc.DBHelper.GetBusinessConnectProductByID(convertedTransactionID)
	if TransErr != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to get product history",
		})
	}

	relatedProducts, _, CatErr := dbFunc.DBHelper.GetBusinessConnectRecommendedProductsByLimit(convertedTransactionID, productDetail.Category, 12)
	if CatErr != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to get recommended product",
		})
	}

	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"success": productDetail,
		"related": relatedProducts,
	})
}

func GetBusinessConnectAdminProductByID(ctx *fiber.Ctx) error {

	productID := ctx.Params("id")
	convertedTransactionID, err := strconv.ParseUint(productID, 10, 64)
	if err != nil {
		// Handle error
		fmt.Println(err)
	}

	productDetail, TransErr := dbFunc.DBHelper.GetBusinessConnectProductByID(convertedTransactionID)
	if TransErr != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to get product history",
		})
	}

	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"success": productDetail,
	})
}

func GetNextProductID(ctx *fiber.Ctx) error {

	productID := ctx.Params("id")
	convertedTransactionID, err := strconv.ParseUint(productID, 10, 64)
	if err != nil {
		// Handle error
		fmt.Println(err)
	}

	productDetail, TransErr := dbFunc.DBHelper.GetNextProductID(uint64(convertedTransactionID))
	if TransErr != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to get product history",
		})
	}

	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"success": productDetail,
	})
}

func GetPreviousProductID(ctx *fiber.Ctx) error {

	productID := ctx.Params("id")
	convertedTransactionID, err := strconv.ParseUint(productID, 10, 64)
	if err != nil {
		// Handle error
		fmt.Println(err)
	}

	productDetail, TransErr := dbFunc.DBHelper.GetPreviousProductID(uint64(convertedTransactionID))
	if TransErr != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to get product history",
		})
	}

	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"success": productDetail,
	})
}

func SearchProductsByTitleAndCategory(ctx *fiber.Ctx) error {
	query := ctx.Query("q")
	categorySlug := ctx.Query("category") // This will be like "household-essentials"

	productDetail, TransErr := dbFunc.DBHelper.SearchProductsByTitleAndCategory(query, categorySlug)
	if TransErr != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to get product search",
		})
	}

	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"success": productDetail,
	})
}

func SearchAdminProductsByTitle(ctx *fiber.Ctx) error {
	type RequestBody struct {
		SearchTerm string `json:"search_term"`
	}

	var body RequestBody
	if err := ctx.BodyParser(&body); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to read body",
		})
	}

	productDetail, TransErr := dbFunc.DBHelper.SearchAdminOrderByTitle(body.SearchTerm)
	if TransErr != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to get product history",
		})
	}

	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"success": productDetail,
	})
}

func SearchAdminOrderByTitle(ctx *fiber.Ctx) error {
	type RequestBody struct {
		SearchTerm string `json:"search_term"`
	}

	var body RequestBody
	if err := ctx.BodyParser(&body); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to read body",
		})
	}

	productDetail, TransErr := dbFunc.DBHelper.SearchAdminOrderByTitle(body.SearchTerm)
	if TransErr != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to get product history",
		})
	}

	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"success": productDetail,
	})
}

func GetTransactionHistoryByDate(ctx *fiber.Ctx) error {
	type RequestBody struct {
		Date string `json:"date"`
	}

	// get stored user id from request time line
	userId := ctx.Locals("user-id")

	if userId == nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "failed to get user",
		})
	}

	user, uuidErr := helperFunc.PaystackHelper.FindByUuidFromLocal(userId)
	if uuidErr != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "failed to get user from request",
		})
	}

	var body RequestBody
	if err := ctx.BodyParser(&body); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to read body",
		})
	}

	productHistory, productHistoryErr := dbFunc.DBHelper.GetTransactionsByUserAndDateWithLimit(user.ID, body.Date)
	if productHistoryErr != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to get product history",
		})
	}

	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"success": productHistory,
	})
}

func GetBusinessConnectHomePageProducts(ctx *fiber.Ctx) error {

	eachPage := 8

	// Get all products
	allProductRecords, allProductRecordsErr := dbFunc.DBHelper.GetBusinessConnectHomeAllProductsByLimit(eachPage)
	if allProductRecordsErr != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to get all products",
		})
	}

	// Get featured product
	featuredProductRecords, featuredProductRecordsErr := dbFunc.DBHelper.GetBusinessConnectHomeFeaturedProductsByLimit(eachPage)
	if featuredProductRecordsErr != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to get featured products",
		})
	}

	// Get best seller (most sold products)
	bestSellingProductRecords, bestSellingProductRecordsErr := dbFunc.DBHelper.GetBusinessConnectHomeBestSellingProductsByLimit(eachPage)
	if bestSellingProductRecordsErr != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to get best selling products",
		})
	}

	// Get sales (products on promo)
	OnSaleProductRecords, OnSaleProductRecordsErr := dbFunc.DBHelper.GetBusinessConnectHomeOnSaleProductsByLimit(eachPage)
	if OnSaleProductRecordsErr != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to get on sale products",
		})
	}

	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"all":          allProductRecords,
		"featured":     featuredProductRecords,
		"best_selling": bestSellingProductRecords,
		"on_sale":      OnSaleProductRecords,
	})
}
