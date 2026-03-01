package cmd

import (
	"bytes"
	"strings"
	"testing"

	"github.com/spf13/viper"
)

func TestSearchCommand_RunE_KeywordQueryRendersMatchingPost(t *testing.T) {
	viper.Set("database_url", "")
	viper.Set("format", "json")
	t.Cleanup(func() {
		viper.Set("database_url", "")
		viper.Set("format", "json")
	})

	var out bytes.Buffer
	searchCmd.SetOut(&out)
	searchCmd.SetErr(&out)

	if err := searchCmd.RunE(searchCmd, []string{"buddy"}); err != nil {
		t.Fatalf("unexpected error running search command: %v", err)
	}

	rendered := stripANSI(out.String())
	for _, needle := range []string{"search: buddy", "Looking for a buddy to go to the movies"} {
		if !strings.Contains(rendered, needle) {
			t.Fatalf("expected rendered output to contain %q; output was %q", needle, rendered)
		}
	}
}

func stripANSI(input string) string {
	var b strings.Builder
	b.Grow(len(input))

	inEscape := false
	for i := 0; i < len(input); i++ {
		ch := input[i]
		if inEscape {
			if (ch >= 'A' && ch <= 'Z') || (ch >= 'a' && ch <= 'z') {
				inEscape = false
			}
			continue
		}
		if ch == 0x1b {
			inEscape = true
			continue
		}
		b.WriteByte(ch)
	}
	return b.String()
}
