package adapters

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/Capmus-Team/supost-cli/internal/domain"
	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
)

const (
	defaultS3PhotoRegion = "us-east-1"
	defaultS3PhotoPrefix = "v2/posts"
	defaultPhotoExt      = ".jpg"
)

// S3PostPhotoUploader stores post photos in S3 under v2/posts/{post_id}/{uuid}.{ext}.
type S3PostPhotoUploader struct {
	client *s3.Client
	bucket string
	prefix string
}

// NewS3PostPhotoUploader builds an uploader using AWS default credentials or an optional profile.
func NewS3PostPhotoUploader(
	ctx context.Context,
	region string,
	bucket string,
	prefix string,
	profile string,
) (*S3PostPhotoUploader, error) {
	region = strings.TrimSpace(region)
	if region == "" {
		region = defaultS3PhotoRegion
	}
	bucket = strings.TrimSpace(bucket)
	if bucket == "" {
		return nil, fmt.Errorf("s3 photo bucket is required")
	}

	prefix = strings.Trim(strings.TrimSpace(prefix), "/")
	if prefix == "" {
		prefix = defaultS3PhotoPrefix
	}

	loadOptions := []func(*awsconfig.LoadOptions) error{
		awsconfig.WithRegion(region),
	}
	if trimmedProfile := strings.TrimSpace(profile); trimmedProfile != "" {
		loadOptions = append(loadOptions, awsconfig.WithSharedConfigProfile(trimmedProfile))
	}

	cfg, err := awsconfig.LoadDefaultConfig(ctx, loadOptions...)
	if err != nil {
		return nil, fmt.Errorf("loading aws config: %w", err)
	}

	return &S3PostPhotoUploader{
		client: s3.NewFromConfig(cfg),
		bucket: bucket,
		prefix: prefix,
	}, nil
}

// UploadPostPhoto uploads a single photo and returns the S3 metadata for public.photo.
func (u *S3PostPhotoUploader) UploadPostPhoto(
	ctx context.Context,
	postID int64,
	photo domain.PostCreatePhotoUpload,
) (domain.PostCreateSavedPhoto, error) {
	if postID <= 0 {
		return domain.PostCreateSavedPhoto{}, fmt.Errorf("post id must be positive")
	}
	if len(photo.Content) == 0 {
		return domain.PostCreateSavedPhoto{}, fmt.Errorf("photo content is empty")
	}

	ext := photoFileExtension(photo.FileName, photo.ContentType)
	key := fmt.Sprintf("%s/%d/%s%s", u.prefix, postID, uuid.NewString(), ext)

	contentType := strings.TrimSpace(photo.ContentType)
	if contentType == "" {
		contentType = http.DetectContentType(photo.Content)
	}

	_, err := u.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(u.bucket),
		Key:         aws.String(key),
		Body:        bytes.NewReader(photo.Content),
		ContentType: aws.String(contentType),
	})
	if err != nil {
		return domain.PostCreateSavedPhoto{}, fmt.Errorf("put object %q: %w", key, err)
	}

	return domain.PostCreateSavedPhoto{
		PostID:      postID,
		S3Key:       key,
		TickerS3Key: "",
		Position:    photo.Position,
	}, nil
}

func photoFileExtension(fileName string, contentType string) string {
	ext := strings.ToLower(strings.TrimSpace(filepath.Ext(strings.TrimSpace(fileName))))
	if validPhotoExtension(ext) {
		return ext
	}

	switch strings.ToLower(strings.TrimSpace(contentType)) {
	case "image/jpeg", "image/jpg":
		return ".jpg"
	case "image/png":
		return ".png"
	case "image/webp":
		return ".webp"
	case "image/gif":
		return ".gif"
	case "image/heic":
		return ".heic"
	case "image/heif":
		return ".heif"
	case "image/avif":
		return ".avif"
	default:
		return defaultPhotoExt
	}
}

func validPhotoExtension(ext string) bool {
	if len(ext) < 2 || len(ext) > 10 || ext[0] != '.' {
		return false
	}
	for _, r := range ext[1:] {
		if (r < 'a' || r > 'z') && (r < '0' || r > '9') {
			return false
		}
	}
	return true
}
