package post

import (
	"fmt"
	"mime/multipart"
	"strconv"

	dbFunc "business-connect/database/dbHelpFunc"
	Data "business-connect/models"
	upload "business-connect/upload"

	"github.com/gofiber/fiber/v2"
)

var (
	OK              bool
	PostBody        Data.Productt
	imageUploads    []Data.ProductImage
	recievedFiles   *multipart.Form
	imageUploadsErr error
	fileParseError  error
)

func PostProduct(ctx *fiber.Ctx) error {

	// Get stored user id from request timeline
	userId := ctx.Locals("user-id")

	if userId == nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "failed to get user",
		})
	}

	user, uuidErr := dbFunc.DBHelper.FindByUuidFromLocal(userId)

	if uuidErr != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": uuidErr.Error(),
		})
	}

	// Parse the form data
	recievedFiles, fileParseError = ctx.MultipartForm()
	if fileParseError != nil {
		fmt.Println("Error parsing multipart form:", fileParseError)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to parse images",
		})
	}

	// Parse product details
	var product Data.Productt

	product.UserID = user.ID
	product.Title = ctx.FormValue("productTitle")
	product.Description = ctx.FormValue("productDescription")
	product.InitialCost, _ = ConvertStringToFloat64(ctx.FormValue("initialCost"))
	product.SellingPrice, _ = ConvertStringToFloat64(ctx.FormValue("sellingPrice"))
	product.Currency = ctx.FormValue("currency")
	product.ProductStock, _ = ConvertStringToInt64(ctx.FormValue("productStock"))
	product.StockRemaining, _ = ConvertStringToInt64(ctx.FormValue("productStock"))
	product.NetWeight, _ = ConvertStringToInt64(ctx.FormValue("netWeight"))
	product.ProductRank, _ = ConvertStringToInt(ctx.FormValue("productRank"))
	product.DiscountType = ctx.FormValue("discountType")
	product.Category = ctx.FormValue("category")
	product.Tags = ctx.FormValue("tags")
	product.PublishStatus = ctx.FormValue("publishStatus")
	var featuredProduct = ctx.FormValue("featuredStatus")
	if featuredProduct == "featured" {
		product.Featured = true
	} else {
		product.Featured = false
	}

	// Access the uploaded files by name
	files := recievedFiles.File["imageArray"]
	if len(files) == 0 {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "No files uploaded",
		})
	}

	if len(files) < 2 {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "At least two images are required",
		})
	}

	// Upload the images to Backblaze B2
	imageUploads, imageUploadsErr = upload.UploadFiles(files)
	if imageUploadsErr != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to upload image to Backblaze B2",
		})
	}

	// for _, procductImage := range imageUploads {
	product.Image1 = imageUploads[0].URL
	product.Image2 = imageUploads[1].URL
	// }

	if product.InitialCost != product.SellingPrice {
		product.OnSale = true
	}

	// Save the product details to the database after successful image upload
	savedProduct, err := dbFunc.DBHelper.AddProduct(product, user)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "error occurred while adding product to database",
		})
	}

	// Add uploaded images to the database
	for _, eachImage := range imageUploads {
		err := dbFunc.DBHelper.AddProductImage(eachImage, savedProduct.ID)
		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "error occurred while adding image to database",
			})
		}
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "successfully added product",
	})
}

// ConvertStringToInt64 converts a string to int64 and returns the value and an error if any.
func ConvertStringToInt64(s string) (int64, error) {
	value, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to convert string to int64: %v", err)
	}
	return value, nil
}

// ConvertStringToInt converts a string to int and returns the value and an error if any.
func ConvertStringToInt(s string) (int, error) {
	value, err := strconv.Atoi(s)
	if err != nil {
		return 0, fmt.Errorf("failed to convert string to int: %v", err)
	}
	return value, nil
}

// ConvertStringToFloat64 converts a string to float64 and returns the value and an error if any.
func ConvertStringToFloat64(s string) (float64, error) {
	value, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to convert string to float64: %v", err)
	}
	return value, nil
}
