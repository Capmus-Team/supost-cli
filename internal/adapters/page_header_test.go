package adapters

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func TestFormatUpdatedTimestamp(t *testing.T) {
	now := time.Date(2026, time.February, 27, 14, 25, 0, 0, time.FixedZone("PST", -8*60*60))
	got := formatUpdatedTimestamp(now)
	want := "Fri, Feb 27, 2026 02:25 PM - Updated"
	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestRenderPageHeader_UsesRightAlignedTimestamp(t *testing.T) {
	now := time.Date(2026, time.February, 27, 14, 25, 0, 0, time.FixedZone("PST", -8*60*60))
	var out bytes.Buffer

	err := RenderPageHeader(&out, PageHeaderOptions{
		Width:      90,
		Location:   "Stanford, California",
		RightLabel: "post",
		Now:        now,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	lines := strings.Split(strings.TrimSpace(out.String()), "\n")
	if len(lines) != 2 {
		t.Fatalf("expected 2 header lines, got %d", len(lines))
	}

	metaPlain := stripANSI(lines[1])
	if got := len([]rune(metaPlain)); got != 90 {
		t.Fatalf("line width mismatch: got %d want %d", got, 90)
	}
	wantSuffix := "Fri, Feb 27, 2026 02:25 PM - Updated"
	if !strings.HasSuffix(metaPlain, wantSuffix) {
		t.Fatalf("line %q does not end with %q", metaPlain, wantSuffix)
	}
}

func TestBuildAdaptiveBreadcrumb_UsesIDsAndTruncatesTitle(t *testing.T) {
	got := buildAdaptiveBreadcrumb("Stanford, California", BreadcrumbOptions{
		CategoryID:    5,
		SubcategoryID: 9,
		PostID:        130031605,
		PostTitle:     "New Artisan Handmade Ultra Comfortable Women's Dance Boots, US 8.5-9",
	})

	want := "Stanford, California » for sale » clothing & accessories » New Artisan Handmade Ult..."
	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestRenderPageHeader_RendersAdaptiveBreadcrumbInMeta(t *testing.T) {
	now := time.Date(2026, time.February, 27, 20, 8, 0, 0, time.FixedZone("PST", -8*60*60))
	var out bytes.Buffer

	err := RenderPageHeader(&out, PageHeaderOptions{
		Width:      118,
		Location:   "Stanford, California",
		RightLabel: "post",
		Now:        now,
		Breadcrumb: &BreadcrumbOptions{
			CategoryID:    3,
			SubcategoryID: 59,
			PostID:        130031900,
			PostTitle:     "Large furnished room in quiet home close to campus",
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	lines := strings.Split(strings.TrimSpace(out.String()), "\n")
	if len(lines) != 2 {
		t.Fatalf("expected 2 header lines, got %d", len(lines))
	}
	metaPlain := stripANSI(lines[1])

	if !strings.Contains(metaPlain, "SUPost » Stanford, California » housing » rooms & shares »") {
		t.Fatalf("missing category/subcategory breadcrumb in %q", metaPlain)
	}
	if !strings.Contains(metaPlain, "...") {
		t.Fatalf("expected truncated title ellipsis in %q", metaPlain)
	}
	if !strings.Contains(metaPlain, "Fri, Feb 27, 2026 08:08 PM - Updated") {
		t.Fatalf("missing timestamp in %q", metaPlain)
	}
}
