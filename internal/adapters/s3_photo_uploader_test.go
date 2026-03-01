package adapters

import "testing"

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
