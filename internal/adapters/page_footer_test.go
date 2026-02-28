package adapters

import (
	"bytes"
	"strings"
	"testing"
)

func TestRenderPageFooter_ContainsExpectedContent(t *testing.T) {
	var out bytes.Buffer
	if err := RenderPageFooter(&out, PageFooterOptions{Width: 118}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	plain := stripANSI(out.String())
	for _, needle := range []string{
		"post a job",
		"post housing",
		"post a car",
		"about",
		"contact",
		"privacy",
		"terms",
		"help",
		"a Greg Wientjes production",
		"SUpost is not sponsored by, endorsed by, or affiliated with Stanford University.",
		"SUpost © 2009",
	} {
		if !strings.Contains(plain, needle) {
			t.Fatalf("missing %q in footer", needle)
		}
	}
}
