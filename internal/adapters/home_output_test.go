package adapters

import "testing"

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
