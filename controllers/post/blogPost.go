package post

import (
	"fmt"
	"mime/multipart"

	dbFunc "business-connect/database/dbHelpFunc"
	Data "business-connect/models"
	upload "business-connect/upload"

	"github.com/gofiber/fiber/v2"
)

var (
	Ok                  bool
	BlogBody            Data.Blog
	blogImageUploads    []Data.BlogImage
	blogReceivedFiles   *multipart.Form
	blogImageUploadsErr error
	blogFileParseError  error
)

func BlogPost(ctx *fiber.Ctx) error {

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
	blogReceivedFiles, blogFileParseError = ctx.MultipartForm()
	if blogFileParseError != nil {
		fmt.Println("Error parsing multipart form:", blogFileParseError)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to parse images",
		})
	}

	// Parse product details
	var product Data.Blog

	product.UserID = user.ID
	product.Title = ctx.FormValue("blogTitle")
	product.Description1 = ctx.FormValue("blogDescription1")
	product.Description2 = ctx.FormValue("blogDescription2")
	product.BlogCategory = ctx.FormValue("blogCategory")

	// Access the uploaded files by name
	files := blogReceivedFiles.File["imageArray"]
	if len(files) == 0 {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "No files uploaded",
		})
	}

	if len(files) != 3 {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Image upload must be three",
		})
	}

	// Upload the images to Backblaze B2
	blogImageUploads, blogImageUploadsErr = upload.UploadBlogFiles(files)
	if blogImageUploadsErr != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to upload image to Backblaze B2",
		})
	}

	// for _, procductImage := range blogImageUploads {
	product.Image1 = blogImageUploads[0].URL
	product.Image2 = blogImageUploads[1].URL
	product.Image3 = blogImageUploads[2].URL
	// }

	// Save the product details to the database after successful image upload
	savedProduct, err := dbFunc.DBHelper.AddBlog(product, user)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "error occurred while adding product to database",
		})
	}

	// Add uploaded images to the database
	for _, eachImage := range blogImageUploads {
		err := dbFunc.DBHelper.AddBlogImage(eachImage, savedProduct.ID)
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
