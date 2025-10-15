package utils

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

type ICloudinaryUtils interface {
	UploadImage(ctx context.Context, file io.Reader, filename string) (string, error)
	DeleteImage(ctx context.Context, publicID string) error
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
		Folder:         "/",
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
	_, err := ch.cld.Upload.Destroy(ctx, uploader.DestroyParams{
		PublicID: publicID,
	})
	if err != nil {
		return fmt.Errorf("failed to delete image: %w", err)
	}

	return nil
}
