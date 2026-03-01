package adapters

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
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
	maxPostPhotoWidth    = 340
	maxTickerPhotoWidth  = 220
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

	decoded, formatName, err := image.Decode(bytes.NewReader(photo.Content))
	if err != nil {
		return domain.PostCreateSavedPhoto{}, fmt.Errorf("decoding image: %w", err)
	}

	contentType, ext := imageOutputFormat(formatName, photo.FileName, photo.ContentType)
	postImage := resizeToMaxWidth(decoded, maxPostPhotoWidth)
	tickerImage := resizeToMaxWidth(decoded, maxTickerPhotoWidth)

	postBytes, err := encodeImageBytes(postImage, contentType)
	if err != nil {
		return domain.PostCreateSavedPhoto{}, fmt.Errorf("encoding post image: %w", err)
	}
	tickerBytes, err := encodeImageBytes(tickerImage, contentType)
	if err != nil {
		return domain.PostCreateSavedPhoto{}, fmt.Errorf("encoding ticker image: %w", err)
	}

	objectID := uuid.NewString()
	key := fmt.Sprintf("%s/%d/%s%s", u.prefix, postID, objectID, ext)
	tickerKey := fmt.Sprintf("%s/%d/ticker_%s%s", u.prefix, postID, objectID, ext)

	if err := u.putObject(ctx, key, postBytes, contentType); err != nil {
		return domain.PostCreateSavedPhoto{}, fmt.Errorf("put object %q: %w", key, err)
	}
	if err := u.putObject(ctx, tickerKey, tickerBytes, contentType); err != nil {
		return domain.PostCreateSavedPhoto{}, fmt.Errorf("put object %q: %w", tickerKey, err)
	}

	return domain.PostCreateSavedPhoto{
		PostID:      postID,
		S3Key:       key,
		TickerS3Key: tickerKey,
		Position:    photo.Position,
	}, nil
}

func (u *S3PostPhotoUploader) putObject(ctx context.Context, key string, content []byte, contentType string) error {
	_, err := u.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(u.bucket),
		Key:         aws.String(key),
		Body:        bytes.NewReader(content),
		ContentType: aws.String(contentType),
	})
	return err
}

func resizeToMaxWidth(src image.Image, maxWidth int) image.Image {
	if maxWidth <= 0 {
		return src
	}
	bounds := src.Bounds()
	srcW := bounds.Dx()
	srcH := bounds.Dy()
	if srcW <= 0 || srcH <= 0 || srcW <= maxWidth {
		return src
	}

	dstW := maxWidth
	dstH := srcH * dstW / srcW
	if dstH <= 0 {
		dstH = 1
	}

	dst := image.NewRGBA(image.Rect(0, 0, dstW, dstH))
	for y := 0; y < dstH; y++ {
		srcY := bounds.Min.Y + (y * srcH / dstH)
		for x := 0; x < dstW; x++ {
			srcX := bounds.Min.X + (x * srcW / dstW)
			dst.Set(x, y, src.At(srcX, srcY))
		}
	}
	return dst
}

func imageOutputFormat(decodedFormat string, fileName string, suppliedContentType string) (contentType string, ext string) {
	switch strings.ToLower(strings.TrimSpace(decodedFormat)) {
	case "png":
		return "image/png", ".png"
	case "gif":
		return "image/gif", ".gif"
	case "jpeg":
		return "image/jpeg", ".jpg"
	}

	ext = photoFileExtension(fileName, suppliedContentType)
	switch ext {
	case ".png":
		return "image/png", ".png"
	case ".gif":
		return "image/gif", ".gif"
	default:
		return "image/jpeg", ".jpg"
	}
}

func encodeImageBytes(img image.Image, contentType string) ([]byte, error) {
	var buf bytes.Buffer
	switch strings.ToLower(strings.TrimSpace(contentType)) {
	case "image/png":
		if err := png.Encode(&buf, img); err != nil {
			return nil, err
		}
	case "image/gif":
		if err := gif.Encode(&buf, img, nil); err != nil {
			return nil, err
		}
	default:
		if err := jpeg.Encode(&buf, img, &jpeg.Options{Quality: 85}); err != nil {
			return nil, err
		}
	}
	return buf.Bytes(), nil
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
