package adapters

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/Capmus-Team/supost-cli/internal/domain"
)

func TestFormatPostPhotoURL_UsesPostPattern(t *testing.T) {
	now := time.Date(2026, time.February, 27, 14, 25, 0, 0, time.UTC)
	post := domain.Post{
		ID:         130031961,
		TimePosted: 1772238444,
	}
	got := formatPostPhotoURL(post, 0, now)
	want := "https://supost-prod.s3.amazonaws.com/posts/130031961/post_130031961a?1772238444"
	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestPostPhotoQuadrantURLs_MapsThirdPhotoToBottomLeft(t *testing.T) {
	now := time.Date(2026, time.February, 27, 14, 25, 0, 0, time.UTC)
	post := domain.Post{
		ID:         130031961,
		TimePosted: 1772238444,
		Photo1File: "a.jpg",
		Photo2File: "b.jpg",
		Photo3File: "c.jpg",
	}

	urls := postPhotoQuadrantURLs(post, now)
	if urls[2] == "" {
		t.Fatalf("expected third slot (bottom-left) to be populated")
	}
	if !strings.Contains(urls[2], "post_130031961c") {
		t.Fatalf("expected c-suffix url in third slot, got %q", urls[2])
	}
}

func TestRenderPostPage_ContainsHeaderBodyAndPosterBox(t *testing.T) {
	var out bytes.Buffer
	post := domain.Post{
		ID:           130031961,
		Name:         "Shared House - $700",
		Email:        "owner@stanford.edu",
		CategoryID:   3,
		Body:         "Room w/ Private Bathroom for Rent in Quiet Home | Menlo Park",
		TimePosted:   1772238444,
		HasPrice:     true,
		Price:        700,
		Photo1File:   "a.jpg",
		Photo2File:   "b.jpg",
		Status:       domain.PostStatusActive,
		TimePostedAt: time.Date(2026, time.February, 27, 16, 26, 0, 0, time.UTC),
	}

	if err := RenderPostPage(&out, post); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	plain := stripANSI(out.String())
	for _, needle := range []string{
		"SUPost",
		"Message Poster",
		"Reply to: Use the form at the right to send messages to this user.",
		"Date:",
		"Price: $700",
		"post_130031961a",
		"post_130031961b",
		"Be sure to follow official Student Housing Sublicensing policies:",
		"http://sublicense.Stanford.edu",
		"please do not message this poster about other commercial services",
	} {
		if !strings.Contains(plain, needle) {
			t.Fatalf("missing %q in rendered post page", needle)
		}
	}
}

func TestRenderPostPage_NonHousingSkipsHousingNotice(t *testing.T) {
	var out bytes.Buffer
	post := domain.Post{
		ID:         130031961,
		CategoryID: 5,
		Name:       "Bike for Sale",
		Body:       "Clean bike, ready to ride.",
		TimePosted: 1772238444,
	}

	if err := RenderPostPage(&out, post); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	plain := stripANSI(out.String())
	if strings.Contains(plain, "Student Housing Sublicensing policies") {
		t.Fatalf("did not expect housing notice for non-housing category")
	}
}
