package utils

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

type ICloudinaryUtils interface {
	UploadImage(ctx context.Context, file io.Reader, filename string) (string, error)
	DeleteImage(ctx context.Context, publicID string) error
	ExtractPublicIDFromURL(imageURL string) string
}

type cloudinaryUtils struct {
	cld *cloudinary.Cloudinary
}

func NewCloudinaryUtils() (ICloudinaryUtils, error) {
	cloudName := os.Getenv("CLOUDINARY_CLOUD_NAME")
	apiKey := os.Getenv("CLOUDINARY_API_KEY")
	apiSecret := os.Getenv("CLOUDINARY_API_SECRET")

	if cloudName == "" || apiKey == "" || apiSecret == "" {
		return nil, fmt.Errorf("cloudinary credentials are not set")
	}

	cld, err := cloudinary.NewFromParams(cloudName, apiKey, apiSecret)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize cloudinary: %w", err)
	}

	return &cloudinaryUtils{
		cld: cld,
	}, nil
}

func boolPtr(b bool) *bool {
	return &b
}

func (ch *cloudinaryUtils) UploadImage(ctx context.Context, file io.Reader, filename string) (string, error) {
	uploadResult, err := ch.cld.Upload.Upload(ctx, file, uploader.UploadParams{
		Folder:         "products",
		UniqueFilename: boolPtr(true),
		Overwrite:      boolPtr(false),
		ResourceType:   "image",
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload image: %w", err)
	}

	return uploadResult.SecureURL, nil
}

func (ch *cloudinaryUtils) DeleteImage(ctx context.Context, publicID string) error {
	if publicID == "" {
		return nil
	}

	_, err := ch.cld.Upload.Destroy(ctx, uploader.DestroyParams{
		PublicID: publicID,
	})
	if err != nil {
		return fmt.Errorf("failed to delete image: %w", err)
	}

	return nil
}

func (ch *cloudinaryUtils) ExtractPublicIDFromURL(imageURL string) string {
	if imageURL == "" {
		return ""
	}

	if !strings.Contains(imageURL, "cloudinary.com") {
		return ""
	}

	if strings.Contains(imageURL, "placeholder") {
		return ""
	}

	parts := strings.Split(imageURL, "/upload/")
	if len(parts) < 2 {
		return ""
	}

	afterUpload := parts[1]

	afterUpload = strings.TrimPrefix(afterUpload, "v")
	slashIndex := strings.Index(afterUpload, "/")
	if slashIndex > 0 {
		afterUpload = afterUpload[slashIndex+1:]
	}

	lastDotIndex := strings.LastIndex(afterUpload, ".")
	if lastDotIndex > 0 {
		afterUpload = afterUpload[:lastDotIndex]
	}

	return afterUpload
}
