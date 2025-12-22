package upload

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	Data "business-connect/models"

	"github.com/joho/godotenv"
	"github.com/kurin/blazer/b2"
)

// Function to map file extensions to MIME types
func getContentType(filename string) string {
	ext := strings.ToLower(filepath.Ext(filename))
	switch ext {
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".png":
		return "image/png"
	case ".gif":
		return "image/gif"
	case ".bmp":
		return "image/bmp"
	case ".webp":
		return "image/webp"
	default:
		return "application/octet-stream" // Default to binary stream
	}
}

func UploadFiles(fileHeader []*multipart.FileHeader) ([]Data.PostImage, error) {
	var (
		b2Client   *b2.Client
		bucket     *b2.Bucket
		results    []Data.PostImage
		EnvErr     error
		B2LinkErr  error
		bucketName string
		// accountID      string
		applicationKey string
		keyID          string
	)

	// Load environment variables from the .env file
	EnvErr = godotenv.Load(".env")
	if EnvErr != nil {
		return nil, errors.New("error loading .env file")
	}

	// Get B2 credentials from environment variables
	bucketName = os.Getenv("B2_BUCKET_NAME")
	applicationKey = os.Getenv("B2_APPLICATION_KEY")
	keyID = os.Getenv("B2_KEY_ID")
	// accountID = os.Getenv("B2_ACCOUNT_ID")

	// Create a new B2 client
	b2Client, B2LinkErr = b2.NewClient(context.Background(), keyID, applicationKey)
	if B2LinkErr != nil {
		return nil, errors.New("error occurred while setting up B2 client")
	}

	// Get the bucket instance
	bucket, B2LinkErr = b2Client.Bucket(context.Background(), bucketName)
	if B2LinkErr != nil {
		return nil, errors.New("error occurred while getting B2 bucket")
	}

	// Define a temporary directory where you'll save the uploaded files
	tempDir := "./temp"

	// Create the temporary directory if it doesn't exist
	if err := os.MkdirAll(tempDir, os.ModePerm); err != nil {
		return nil, err
	}

	folderPath := "business-connect-store/"

	// Iterate through the provided files and upload each one
	for _, fileHeader := range fileHeader {
		// Generate a unique filename for the temporary file
		tempFilename := filepath.Join(tempDir, fileHeader.Filename)

		// Open the temporary file for writing
		tempFile, err := os.Create(tempFilename)
		if err != nil {
			return nil, err
		}
		defer tempFile.Close()

		// Open the uploaded file for reading
		uploadedFile, err := fileHeader.Open()
		if err != nil {
			return nil, err
		}
		defer uploadedFile.Close()

		// Copy the contents of the uploaded file to the temporary file
		_, err = io.Copy(tempFile, uploadedFile)
		if err != nil {
			return nil, err
		}

		// Open the temporary file for reading again
		tempFileForUpload, err := os.Open(tempFilename)
		if err != nil {
			return nil, err
		}
		defer tempFileForUpload.Close()

		// Generate a unique file name for the B2 object
		fileName := fmt.Sprintf("%s%d_%s", folderPath, time.Now().UnixNano(), fileHeader.Filename)

		// Create a new writer with attributes
		obj := bucket.Object(fileName)
		contentType := getContentType(fileHeader.Filename)
		attrs := &b2.Attrs{ContentType: contentType}
		writer := obj.NewWriter(context.Background()).WithAttrs(attrs)

		// Upload the file to Backblaze B2
		if _, err := io.Copy(writer, tempFileForUpload); err != nil {
			return nil, errors.New("error occurred while uploading file to B2")
		}
		if err := writer.Close(); err != nil {
			return nil, errors.New("error occurred while closing writer after uploading file to B2")
		}

		// Build the download URL
		// downloadURL := fmt.Sprintf("https://shopsphereafrica.com/image/%s", fileName)

		var imageResp = Data.PostImage{
			URL:              fileName,
			OriginalFilename: fileHeader.Filename,
		}

		// Append the upload result to the results slice
		results = append(results, imageResp)

		// Delete the temporary file corresponding to the uploaded file
		if err := os.Remove(tempFilename); err != nil {
			return nil, errors.New("error occurred while removing temporary file")
		}
	}

	// Return the results of all files uploaded
	return results, nil
}

func UploadEmailFiles(htmlContent string) (string, error) {
	var (
		b2Client       *b2.Client
		bucket         *b2.Bucket
		EnvErr         error
		B2LinkErr      error
		bucketName     string
		applicationKey string
		keyID          string
	)

	// Load environment variables from the .env file
	EnvErr = godotenv.Load(".env")
	if EnvErr != nil {
		return "", errors.New("error loading .env file")
	}
	// log.Println("Loading environment variables")

	// Get B2 credentials from environment variables
	bucketName = os.Getenv("B2_BUCKET_NAME")
	applicationKey = os.Getenv("B2_APPLICATION_KEY")
	keyID = os.Getenv("B2_KEY_ID")

	// log.Printf("Bucket Name: %s\n", bucketName)
	// log.Printf("Application Key: %s\n", applicationKey)
	// log.Printf("Key ID: %s\n", keyID)

	// Create a new B2 client
	log.Println("Creating B2 client")
	b2Client, B2LinkErr = b2.NewClient(context.Background(), keyID, applicationKey)
	if B2LinkErr != nil {
		return "", errors.New("error occurred while setting up B2 client")
	}

	// Get the bucket instance
	log.Println("Getting B2 bucket")
	bucket, B2LinkErr = b2Client.Bucket(context.Background(), bucketName)
	if B2LinkErr != nil {
		return "", errors.New("error occurred while getting B2 bucket")
	}

	// Create a temporary directory for file storage
	tempDir := "./temp"
	if err := os.MkdirAll(tempDir, os.ModePerm); err != nil {
		return "", err
	}
	defer os.RemoveAll(tempDir) // Clean up temporary directory
	// log.Printf("Creating temporary directory: %s\n", tempDir)

	// Define the folder path in the bucket where the files will be uploaded
	folderPath := "business-connect-email/"
	// log.Printf("Using folder path for upload: %s\n", folderPath)

	// Use regular expression to find all base64 images in the HTML content
	base64ImgRegex := regexp.MustCompile(`(?i)<img\s+[^>]*src="data:image/[^;]+;base64,([^"]+)"[^>]*>`)
	matches := base64ImgRegex.FindAllStringSubmatch(htmlContent, -1)

	if matches == nil {
		log.Println("No base64 images found")
		return htmlContent, nil
	}

	// Iterate over all base64 image matches
	for _, match := range matches {
		base64Data := match[1]
		// log.Printf("Found base64 image with data length: %d\n", len(base64Data))

		// Decode the base64 image
		imageData, err := base64.StdEncoding.DecodeString(base64Data)
		if err != nil {
			return "", errors.New("error decoding base64 image")
		}

		// Generate a temporary filename
		tempFilename := fmt.Sprintf("%s%d_image.png", tempDir, time.Now().UnixNano())
		err = os.WriteFile(tempFilename, imageData, 0644)
		if err != nil {
			return "", errors.New("error writing image to temporary file")
		}
		// log.Printf("Writing image to temporary file: %s\n", tempFilename)

		// Generate a unique file name for the B2 object
		fileName := fmt.Sprintf("%s%d_image.png", folderPath, time.Now().UnixNano())

		// Upload the file to Backblaze B2
		tempFileForUpload, err := os.Open(tempFilename)
		if err != nil {
			return "", errors.New("error opening temporary file for upload")
		}
		defer tempFileForUpload.Close()
		// log.Printf("Uploading file to B2 with name: %s\n", fileName)

		// Create the object writer with content type
		obj := bucket.Object(fileName)
		contentType := getContentType(tempFilename)
		attrs := &b2.Attrs{ContentType: contentType}
		writer := obj.NewWriter(context.Background()).WithAttrs(attrs)

		_, err = io.Copy(writer, tempFileForUpload)
		if err != nil {
			return "", errors.New("error uploading file to B2")
		}
		if err := writer.Close(); err != nil {
			return "", errors.New("error closing writer after upload")
		}

		// Build the public URL for the uploaded image
		imageURL := fmt.Sprintf("https://shopsphereafrica.com/image/%s", fileName)
		// log.Printf("Image URL: %s\n", imageURL)

		// Replace the base64 data with the new URL in the HTML content
		htmlContent = strings.Replace(htmlContent, match[0], fmt.Sprintf(`<img src="%s">`, imageURL), 1)
	}

	// Return the modified HTML content with the base64 images replaced by URLs
	return htmlContent, nil
}

func UploadBlogFiles(fileHeader []*multipart.FileHeader) ([]Data.BlogImage, error) {
	var (
		b2Client   *b2.Client
		bucket     *b2.Bucket
		results    []Data.BlogImage
		EnvErr     error
		B2LinkErr  error
		bucketName string
		// accountID      string
		applicationKey string
		keyID          string
	)

	// Load environment variables from the .env file
	EnvErr = godotenv.Load(".env")
	if EnvErr != nil {
		return nil, errors.New("error loading .env file")
	}

	// Get B2 credentials from environment variables
	bucketName = os.Getenv("B2_BUCKET_NAME")
	applicationKey = os.Getenv("B2_APPLICATION_KEY")
	keyID = os.Getenv("B2_KEY_ID")
	// accountID = os.Getenv("B2_ACCOUNT_ID")

	// Create a new B2 client
	b2Client, B2LinkErr = b2.NewClient(context.Background(), keyID, applicationKey)
	if B2LinkErr != nil {
		return nil, errors.New("error occurred while setting up B2 client")
	}

	// Get the bucket instance
	bucket, B2LinkErr = b2Client.Bucket(context.Background(), bucketName)
	if B2LinkErr != nil {
		return nil, errors.New("error occurred while getting B2 bucket")
	}

	// Define a temporary directory where you'll save the uploaded files
	tempDir := "./temp"

	// Create the temporary directory if it doesn't exist
	if err := os.MkdirAll(tempDir, os.ModePerm); err != nil {
		return nil, err
	}

	// Specify the folder inside the bucket where the files will be uploaded
	folderPath := "business-connect-blog/"

	// Iterate through the provided files and upload each one
	for _, fileHeader := range fileHeader {
		// Generate a unique filename for the temporary file
		tempFilename := filepath.Join(tempDir, fileHeader.Filename)

		// Open the temporary file for writing
		tempFile, err := os.Create(tempFilename)
		if err != nil {
			return nil, err
		}
		defer tempFile.Close()

		// Open the uploaded file for reading
		uploadedFile, err := fileHeader.Open()
		if err != nil {
			return nil, err
		}
		defer uploadedFile.Close()

		// Copy the contents of the uploaded file to the temporary file
		_, err = io.Copy(tempFile, uploadedFile)
		if err != nil {
			return nil, err
		}

		// Open the temporary file for reading again
		tempFileForUpload, err := os.Open(tempFilename)
		if err != nil {
			return nil, err
		}
		defer tempFileForUpload.Close()

		// Generate a unique file name for the B2 object and include the folder path
		fileName := fmt.Sprintf("%s%d_%s", folderPath, time.Now().UnixNano(), fileHeader.Filename)

		// Create a new writer with attributes
		obj := bucket.Object(fileName)
		contentType := getContentType(fileHeader.Filename)
		attrs := &b2.Attrs{ContentType: contentType}
		writer := obj.NewWriter(context.Background()).WithAttrs(attrs)

		// Upload the file to Backblaze B2
		if _, err := io.Copy(writer, tempFileForUpload); err != nil {
			return nil, errors.New("error occurred while uploading file to B2")
		}
		if err := writer.Close(); err != nil {
			return nil, errors.New("error occurred while closing writer after uploading file to B2")
		}

		// Build the download URL
		// downloadURL := fmt.Sprintf("https://shopsphereafrica.com/image/%s", fileName)

		var imageResp = Data.BlogImage{
			URL:              fileName,
			OriginalFilename: fileHeader.Filename,
		}

		// Append the upload result to the results slice
		results = append(results, imageResp)

		// Delete the temporary file corresponding to the uploaded file
		if err := os.Remove(tempFilename); err != nil {
			return nil, errors.New("error occurred while removing temporary file")
		}
	}

	// Return the results of all files uploaded
	return results, nil
}
