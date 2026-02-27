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
