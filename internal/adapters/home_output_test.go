package adapters

import (
	"strings"
	"testing"
	"time"
)

func TestWrapStyledWords_RespectsWidth(t *testing.T) {
	words := []styledWord{
		{text: "one", color: ansiBlue},
		{text: "two", color: ansiBlue},
		{text: "three", color: ansiGray},
		{text: "four", color: ansiMagenta},
		{text: "five"},
		{text: "six"},
		{text: "seven"},
	}
	lines := wrapStyledWords(words, 10)
	if len(lines) == 0 {
		t.Fatalf("expected wrapped lines")
	}

	for _, lineWords := range lines {
		rendered := renderStyledLine(lineWords)
		plain := stripANSI(rendered)
		if got := len([]rune(plain)); got > 10 {
			t.Fatalf("line %q exceeds width: %d", plain, got)
		}
	}
}

func TestWrapStyledWords_SplitsLongWord(t *testing.T) {
	lines := wrapStyledWords([]styledWord{{text: "supercalifragilisticexpialidocious", color: ansiBlue}}, 8)
	if len(lines) < 2 {
		t.Fatalf("expected long word to be split, got %d lines", len(lines))
	}
	for _, lineWords := range lines {
		rendered := renderStyledLine(lineWords)
		plain := stripANSI(rendered)
		if got := len([]rune(plain)); got > 8 {
			t.Fatalf("line %q exceeds width: %d", plain, got)
		}
	}
}

func TestFormatDisplayEmail_StanfordDomainCollapses(t *testing.T) {
	got := formatDisplayEmail("wientjes@alumni.stanford.edu")
	if got != "@stanford.edu" {
		t.Fatalf("got %q, want %q", got, "@stanford.edu")
	}
}

func TestFormatDisplayEmail_NonStanfordUnchanged(t *testing.T) {
	got := formatDisplayEmail("person@example.com")
	if got != "" {
		t.Fatalf("got %q, want empty string", got)
	}
}

func TestFormatHomeUpdatedTimestamp(t *testing.T) {
	now := time.Date(2026, time.February, 27, 14, 25, 0, 0, time.FixedZone("PST", -8*60*60))
	got := formatHomeUpdatedTimestamp(now)
	want := "Fri, Feb 27, 2026 02:25 PM - Updated"
	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestRenderHomeMetaBar_RightAlignsTimestamp(t *testing.T) {
	now := time.Date(2026, time.February, 27, 14, 25, 0, 0, time.FixedZone("PST", -8*60*60))
	width := 90
	line := renderHomeMetaBar(now, width)
	wantSuffix := "Fri, Feb 27, 2026 02:25 PM - Updated"

	if got := len([]rune(line)); got != width {
		t.Fatalf("line width mismatch: got %d want %d", got, width)
	}
	if !strings.HasSuffix(line, wantSuffix) {
		t.Fatalf("line %q does not end with %q", line, wantSuffix)
	}
}

func stripANSI(s string) string {
	// Minimal scrubber for tests.
	out := make([]rune, 0, len(s))
	inEscape := false
	for _, r := range s {
		if r == 0x1b {
			inEscape = true
			continue
		}
		if inEscape {
			if r == 'm' {
				inEscape = false
			}
			continue
		}
		out = append(out, r)
	}
	return string(out)
}
