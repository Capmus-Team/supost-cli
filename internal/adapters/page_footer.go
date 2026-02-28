package adapters

import (
	"fmt"
	"io"
	"strings"
)

const (
	ansiFooterNav = "\033[0;37m"
	ansiFooter    = "\033[0;37m"
)

// PageFooterOptions configures the reusable SUPost footer.
type PageFooterOptions struct {
	Width int
}

// RenderPageFooter renders the shared footer used across page views.
func RenderPageFooter(w io.Writer, opts PageFooterOptions) error {
	width := opts.Width
	if width <= 0 {
		width = defaultPageWidth
	}

	lines := []string{
		styleCentered(strings.Join([]string{
			"post a job",
			"post housing",
			"post a car",
			"about",
			"contact",
			"privacy",
			"terms",
			"help",
		}, "  "), width, ansiFooterNav),
		"",
		styleCentered("a Greg Wientjes production", width, ansiFooter),
		styleCentered("SUpost is not sponsored by, endorsed by, or affiliated with Stanford University.", width, ansiFooter),
		styleCentered("SUpost © 2009", width, ansiFooter),
	}

	for _, line := range lines {
		if _, err := fmt.Fprintln(w, line); err != nil {
			return err
		}
	}
	return nil
}
