package cmd

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadPostCreatePhotos_ReadsPathsAndPositions(t *testing.T) {
	dir := t.TempDir()
	photoOne := filepath.Join(dir, "one.jpg")
	photoTwo := filepath.Join(dir, "two.png")

	if err := os.WriteFile(photoOne, []byte("photo-1"), 0o600); err != nil {
		t.Fatalf("writing photo one: %v", err)
	}
	if err := os.WriteFile(photoTwo, []byte("photo-2"), 0o600); err != nil {
		t.Fatalf("writing photo two: %v", err)
	}

	photos, err := loadPostCreatePhotos([]string{photoOne, photoTwo})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(photos) != 2 {
		t.Fatalf("expected 2 photos, got %d", len(photos))
	}
	if photos[0].Position != 0 || photos[1].Position != 1 {
		t.Fatalf("unexpected positions: %+v", photos)
	}
}

func TestLoadPostCreatePhotos_RejectsMoreThanFour(t *testing.T) {
	_, err := loadPostCreatePhotos([]string{"1.jpg", "2.jpg", "3.jpg", "4.jpg", "5.jpg"})
	if err == nil {
		t.Fatalf("expected validation error")
	}
}
