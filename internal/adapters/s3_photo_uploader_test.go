package adapters

import (
	"context"
	"image"
	"strings"
	"testing"
)

func TestPhotoFileExtension_FromFileName(t *testing.T) {
	got := photoFileExtension("bike.PNG", "image/jpeg")
	if got != ".png" {
		t.Fatalf("expected .png, got %q", got)
	}
}

func TestPhotoFileExtension_FromContentTypeFallback(t *testing.T) {
	got := photoFileExtension("bike", "image/webp")
	if got != ".webp" {
		t.Fatalf("expected .webp, got %q", got)
	}
}

func TestPhotoFileExtension_DefaultFallback(t *testing.T) {
	got := photoFileExtension("bike", "application/octet-stream")
	if got != ".jpg" {
		t.Fatalf("expected .jpg, got %q", got)
	}
}

func TestValidPhotoExtension(t *testing.T) {
	if !validPhotoExtension(".jpg") {
		t.Fatalf("expected .jpg to be valid")
	}
	if validPhotoExtension(".JPG") {
		t.Fatalf("expected uppercase extension to be invalid")
	}
	if validPhotoExtension(".jp*g") {
		t.Fatalf("expected special chars to be invalid")
	}
}

func TestNewS3PostPhotoUploader_RequiresBucket(t *testing.T) {
	_, err := NewS3PostPhotoUploader(context.Background(), "us-east-1", "", "v2/posts", "")
	if err == nil {
		t.Fatalf("expected bucket validation error")
	}
	if !strings.Contains(err.Error(), "bucket") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestResizeToMaxWidth_Downscales(t *testing.T) {
	src := image.NewRGBA(image.Rect(0, 0, 1000, 500))
	out := resizeToMaxWidth(src, 340)
	if out.Bounds().Dx() != 340 {
		t.Fatalf("expected width 340, got %d", out.Bounds().Dx())
	}
	if out.Bounds().Dy() != 170 {
		t.Fatalf("expected proportional height 170, got %d", out.Bounds().Dy())
	}
}

func TestResizeToMaxWidth_DoesNotUpscale(t *testing.T) {
	src := image.NewRGBA(image.Rect(0, 0, 200, 100))
	out := resizeToMaxWidth(src, 340)
	if out.Bounds().Dx() != 200 {
		t.Fatalf("expected unchanged width 200, got %d", out.Bounds().Dx())
	}
	if out.Bounds().Dy() != 100 {
		t.Fatalf("expected unchanged height 100, got %d", out.Bounds().Dy())
	}
}

func TestImageOutputFormat_PrefersDecodedFormat(t *testing.T) {
	contentType, ext := imageOutputFormat("png", "photo.jpg", "image/jpeg")
	if contentType != "image/png" || ext != ".png" {
		t.Fatalf("unexpected format mapping: %s %s", contentType, ext)
	}
}

func TestImageOutputFormat_FallbacksToJpeg(t *testing.T) {
	contentType, ext := imageOutputFormat("", "photo.unknown", "application/octet-stream")
	if contentType != "image/jpeg" || ext != ".jpg" {
		t.Fatalf("unexpected fallback mapping: %s %s", contentType, ext)
	}
}
