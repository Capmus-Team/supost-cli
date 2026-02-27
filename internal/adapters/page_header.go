package adapters

import (
	"fmt"
	"io"
	"strings"
	"time"
)

const (
	ansiReset   = "\033[0m"
	ansiTopBar  = "\033[48;5;24m\033[1;37m"
	ansiMetaBar = "\033[48;5;252m\033[1;34m"

	defaultPageWidth      = 118
	defaultPageLocation   = "Stanford, California"
	defaultPageRightLabel = "post"
)

// PageHeaderOptions configures the reusable SUPost page header.
type PageHeaderOptions struct {
	Width      int
	Location   string
	RightLabel string
	Now        time.Time
}

// RenderPageHeader renders the shared SUPost page header for home/search pages.
func RenderPageHeader(w io.Writer, opts PageHeaderOptions) error {
	width := opts.Width
	if width <= 0 {
		width = defaultPageWidth
	}

	now := opts.Now
	if now.IsZero() {
		now = time.Now()
	}

	location := strings.TrimSpace(opts.Location)
	if location == "" {
		location = defaultPageLocation
	}

	rightLabel := strings.TrimSpace(opts.RightLabel)
	if rightLabel == "" {
		rightLabel = defaultPageRightLabel
	}

	top := renderThreePartLine(
		" SUPost  [__________] [Search]",
		location,
		" "+rightLabel+" ",
		width,
	)
	meta := renderSplitLine(
		" SUPost » "+location,
		formatUpdatedTimestamp(now),
		width,
	)

	if _, err := fmt.Fprintf(w, "%s%s%s\n", ansiTopBar, top, ansiReset); err != nil {
		return err
	}
	if _, err := fmt.Fprintf(w, "%s%s%s\n", ansiMetaBar, meta, ansiReset); err != nil {
		return err
	}
	return nil
}

func formatUpdatedTimestamp(now time.Time) string {
	return now.Format("Mon, Jan 2, 2006 03:04 PM") + " - Updated"
}

func renderThreePartLine(left, center, right string, width int) string {
	if width <= 0 {
		return strings.TrimSpace(left + " " + center + " " + right)
	}

	leftLen := len([]rune(left))
	centerLen := len([]rune(center))
	rightLen := len([]rune(right))
	if leftLen+centerLen+rightLen > width {
		return renderSplitLine(left+" "+center, right, width)
	}

	remaining := width - leftLen - rightLen
	if centerLen > remaining {
		return renderSplitLine(left+" "+center, right, width)
	}

	spacing := remaining - centerLen
	leftPad := spacing / 2
	rightPad := spacing - leftPad
	return left + strings.Repeat(" ", leftPad) + center + strings.Repeat(" ", rightPad) + right
}

func renderSplitLine(left, right string, width int) string {
	if width <= 0 {
		return strings.TrimSpace(left + " " + right)
	}

	rightRunes := []rune(right)
	if len(rightRunes) >= width {
		return string(rightRunes[:width])
	}

	leftRunes := []rune(left)
	availableLeft := width - len(rightRunes)
	if len(leftRunes) > availableLeft {
		leftRunes = leftRunes[:availableLeft]
	}

	return string(leftRunes) + strings.Repeat(" ", availableLeft-len(leftRunes)) + right
}
