package upload

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	Data "business-connect/models"

	"github.com/kurin/blazer/b2"
)

var b2Client *b2.Client

func initB2() (*b2.Bucket, error) {
	bucketName := os.Getenv("B2_BUCKET_NAME")
	applicationKey := os.Getenv("B2_APPLICATION_KEY")
	keyID := os.Getenv("B2_KEY_ID")

	if bucketName == "" || applicationKey == "" || keyID == "" {
		return nil, fmt.Errorf("missing B2 environment variables")
	}

	client, err := b2.NewClient(context.Background(), keyID, applicationKey)
	if err != nil {
		return nil, fmt.Errorf("failed to init B2 client: %w", err)
	}

	b2Client = client

	bucket, err := client.Bucket(context.Background(), bucketName)
	if err != nil {
		return nil, fmt.Errorf("failed to get bucket: %w", err)
	}

	return bucket, nil
}

// -----------------------------
// Helpers
// -----------------------------

func sanitizeFilename(name string) string {
	name = filepath.Base(name)
	name = strings.ReplaceAll(name, " ", "_")
	return regexp.MustCompile(`[^a-zA-Z0-9._-]`).ReplaceAllString(name, "")
}

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
		return "application/octet-stream"
	}
}

func uploadToB2(bucket *b2.Bucket, folder string, filename string, r io.Reader, contentType string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	filePath := fmt.Sprintf("%s%d_%s", folder, time.Now().UnixNano(), filename)

	writer := bucket.Object(filePath).NewWriter(ctx).WithAttrs(&b2.Attrs{
		ContentType: contentType,
	})
	defer writer.Close()

	_, err := io.Copy(writer, r)
	if err != nil {
		return "", err
	}

	return filePath, nil
}

// -----------------------------
// CORE FILE UPLOAD (generic)
// -----------------------------

func uploadMultipartFiles[T any](
	files []*multipart.FileHeader,
	folder string,
	mapper func(url, filename string) T,
) ([]T, error) {

	bucket, err := initB2()
	if err != nil {
		return nil, err
	}

	var results []T

	for _, fh := range files {

		file, err := fh.Open()
		if err != nil {
			return nil, err
		}

		filename := sanitizeFilename(fh.Filename)

		path, err := uploadToB2(
			bucket,
			folder,
			filename,
			file,
			getContentType(fh.Filename),
		)

		file.Close()

		if err != nil {
			return nil, err
		}

		results = append(results, mapper(path, fh.Filename))
	}

	return results, nil
}

// -----------------------------
// POSTS
// -----------------------------

func UploadFiles(fileHeader []*multipart.FileHeader) ([]Data.PostImage, error) {
	return uploadMultipartFiles(fileHeader, "business-connect-store/",
		func(url, name string) Data.PostImage {
			return Data.PostImage{
				URL:              url,
				OriginalFilename: name,
			}
		},
	)
}

// -----------------------------
// STATUS
// -----------------------------

func UploadStatusFiles(fileHeader []*multipart.FileHeader) ([]Data.StatusImage, error) {
	return uploadMultipartFiles(fileHeader, "business-connect-status/",
		func(url, name string) Data.StatusImage {
			return Data.StatusImage{
				URL:              url,
				OriginalFilename: name,
			}
		},
	)
}

// -----------------------------
// BLOG
// -----------------------------

func UploadBlogFiles(fileHeader []*multipart.FileHeader) ([]Data.BlogImage, error) {
	return uploadMultipartFiles(fileHeader, "business-connect-blog/",
		func(url, name string) Data.BlogImage {
			return Data.BlogImage{
				URL:              url,
				OriginalFilename: name,
			}
		},
	)
}

// -----------------------------
// PROFILE
// -----------------------------

func UploadProfileFiles(fileHeader []*multipart.FileHeader) ([]Data.ProfileImage, error) {
	return uploadMultipartFiles(fileHeader, "business-connect-profile-images/",
		func(url, name string) Data.ProfileImage {
			return Data.ProfileImage{
				URL:              url,
				OriginalFilename: name,
			}
		},
	)
}

// -----------------------------
// DELETE
// -----------------------------

func DeleteB2Files(filePaths []string) error {
	bucket, err := initB2()
	if err != nil {
		return err
	}

	ctx := context.Background()

	for _, path := range filePaths {
		if path == "" {
			continue
		}

		if err := bucket.Object(path).Delete(ctx); err != nil {
			fmt.Printf("failed to delete %s: %v\n", path, err)
		}
	}

	return nil
}

// -----------------------------
// EMAIL BASE64 IMAGES
// -----------------------------

func UploadEmailFiles(htmlContent string) (string, error) {
	bucket, err := initB2()
	if err != nil {
		return "", err
	}

	regex := regexp.MustCompile(`(?i)<img\s+[^>]*src="data:image/[^;]+;base64,([^"]+)"[^>]*>`)
	matches := regex.FindAllStringSubmatch(htmlContent, -1)

	if matches == nil {
		return htmlContent, nil
	}

	folder := "business-connect-email/"

	for _, match := range matches {

		data, err := base64.StdEncoding.DecodeString(match[1])
		if err != nil {
			return "", err
		}

		fileName := fmt.Sprintf("%d.png", time.Now().UnixNano())

		path, err := uploadToB2(
			bucket,
			folder,
			fileName,
			bytes.NewReader(data),
			"image/png",
		)

		if err != nil {
			return "", err
		}

		imageURL := fmt.Sprintf("https://shopsphereafrica.com/image/%s", path)

		htmlContent = strings.Replace(htmlContent, match[0],
			fmt.Sprintf(`<img src="%s">`, imageURL), 1)
	}

	return htmlContent, nil
}