package adapters

import (
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/Capmus-Team/supost-cli/internal/domain"
)

const (
	postPageWidth        = homePageWidth
	postContentGap       = 2
	postPosterWidth      = 26
	postPhotoGridGap     = 2
	postPosterMsgRows    = 6
	postCommercialNotice = "please do not message this poster about other commercial services"
)

// RenderPostPage renders a single post view with shared page chrome.
func RenderPostPage(w io.Writer, post domain.Post) error {
	now := time.Now()
	if err := RenderPageHeader(w, PageHeaderOptions{
		Width:      postPageWidth,
		Location:   "Stanford, California",
		RightLabel: "post",
		Now:        now,
	}); err != nil {
		return err
	}

	if err := renderPostTopBlock(w, post, postPageWidth); err != nil {
		return err
	}

	leftWidth, rightWidth := splitPostContentWidths(postPageWidth)
	leftRows := renderPostMainRows(post, now, leftWidth)
	rightRows := renderPostMessagePosterRows(rightWidth)

	totalRows := len(leftRows)
	if len(rightRows) > totalRows {
		totalRows = len(rightRows)
	}

	for i := 0; i < totalRows; i++ {
		left := strings.Repeat(" ", leftWidth)
		right := strings.Repeat(" ", rightWidth)
		if i < len(leftRows) {
			left = padANSIVisibleWidth(leftRows[i], leftWidth)
		}
		if i < len(rightRows) {
			right = padANSIVisibleWidth(rightRows[i], rightWidth)
		}
		if _, err := fmt.Fprintln(w, left+strings.Repeat(" ", postContentGap)+right); err != nil {
			return err
		}
	}

	if _, err := fmt.Fprintln(w); err != nil {
		return err
	}
	return RenderPageFooter(w, PageFooterOptions{Width: postPageWidth})
}

func renderPostTopBlock(w io.Writer, post domain.Post, width int) error {
	title := strings.TrimSpace(post.Name)
	if title == "" {
		title = "(untitled post)"
	}
	if email := formatDisplayEmail(post.Email); email != "" {
		title += " " + email
	}
	if _, err := fmt.Fprintln(w, ansiHeader+fitText(title, width)+ansiReset); err != nil {
		return err
	}

	replyLine := "Reply to: Use the form at the right to send messages to this user."
	if _, err := fmt.Fprintln(w, fitText(replyLine, width)); err != nil {
		return err
	}

	dateValue := "Date: " + formatPostPageDate(post)
	if post.HasPrice {
		dateValue += "  Price: " + formatPrice(post.Price, true)
	}
	_, err := fmt.Fprintln(w, fitText(dateValue, width))
	return err
}

func splitPostContentWidths(width int) (int, int) {
	rightWidth := postPosterWidth
	if width <= postPosterWidth+postContentGap+20 {
		rightWidth = width / 3
	}
	if rightWidth < 18 {
		rightWidth = 18
	}
	leftWidth := width - rightWidth - postContentGap
	if leftWidth < 1 {
		leftWidth = width
		rightWidth = 0
	}
	return leftWidth, rightWidth
}

func renderPostMainRows(post domain.Post, now time.Time, width int) []string {
	rows := make([]string, 0, 24)
	gridRows := renderPostPhotoGridRows(post, now, width)
	rows = append(rows, gridRows...)
	if len(gridRows) > 0 {
		rows = append(rows, "")
	}

	body := strings.TrimSpace(post.Body)
	if body == "" {
		body = "(no body provided)"
	}
	rows = append(rows, wrapPlainText(body, width)...)
	rows = append(rows, "")
	rows = append(rows, wrapPlainText(postCommercialNotice, width)...)
	return rows
}

func renderPostPhotoGridRows(post domain.Post, now time.Time, width int) []string {
	quadrants := postPhotoQuadrantURLs(post, now)
	hasAny := false
	for _, value := range quadrants {
		if value != "" {
			hasAny = true
			break
		}
	}
	if !hasAny {
		return nil
	}

	colWidth := (width - postPhotoGridGap) / 2
	if colWidth < 1 {
		colWidth = width
	}
	rows := make([]string, 0, 8)
	for row := 0; row < 2; row++ {
		leftValue := ""
		if quadrants[row*2] != "" {
			leftValue = fmt.Sprintf("[%d] %s", (row*2)+1, quadrants[row*2])
		}
		rightValue := ""
		if quadrants[(row*2)+1] != "" {
			rightValue = fmt.Sprintf("[%d] %s", (row*2)+2, quadrants[(row*2)+1])
		}

		leftLines := wrapColumnValue(leftValue, colWidth)
		rightLines := wrapColumnValue(rightValue, colWidth)
		maxLines := len(leftLines)
		if len(rightLines) > maxLines {
			maxLines = len(rightLines)
		}

		for lineIdx := 0; lineIdx < maxLines; lineIdx++ {
			left := ""
			if lineIdx < len(leftLines) {
				left = leftLines[lineIdx]
			}
			right := ""
			if lineIdx < len(rightLines) {
				right = rightLines[lineIdx]
			}
			rows = append(rows, fitText(left, colWidth)+strings.Repeat(" ", postPhotoGridGap)+fitText(right, colWidth))
		}
	}
	return rows
}

func postPhotoQuadrantURLs(post domain.Post, now time.Time) [4]string {
	flags := [4]bool{
		strings.TrimSpace(post.Photo1File) != "" || strings.TrimSpace(post.ImageSource1) != "",
		strings.TrimSpace(post.Photo2File) != "" || strings.TrimSpace(post.ImageSource2) != "",
		strings.TrimSpace(post.Photo3File) != "" || strings.TrimSpace(post.ImageSource3) != "",
		strings.TrimSpace(post.Photo4File) != "" || strings.TrimSpace(post.ImageSource4) != "",
	}
	if post.HasImage && !flags[0] && !flags[1] && !flags[2] && !flags[3] {
		flags[0] = true
	}

	var urls [4]string
	for i := 0; i < len(flags); i++ {
		if flags[i] {
			urls[i] = formatPostPhotoURL(post, i, now)
		}
	}
	return urls
}

func formatPostPhotoURL(post domain.Post, index int, now time.Time) string {
	timestamp := post.TimePosted
	if timestamp <= 0 {
		if ts := postTimestamp(post); !ts.IsZero() {
			timestamp = ts.Unix()
		}
	}
	if timestamp <= 0 {
		timestamp = now.Unix()
	}
	suffix := rune('a' + index)
	return fmt.Sprintf("https://supost-prod.s3.amazonaws.com/posts/%d/post_%d%c?%d", post.ID, post.ID, suffix, timestamp)
}

func formatPostPageDate(post domain.Post) string {
	ts := postTimestamp(post)
	if ts.IsZero() {
		return ""
	}
	return ts.Format("Mon, Jan 2, 2006 03:04 PM")
}

func renderPostMessagePosterRows(width int) []string {
	if width <= 0 {
		return nil
	}

	rows := make([]string, 0, 18)
	rows = append(rows, ansiHeader+renderHomeHeader("Message Poster", width)+ansiReset)
	rows = append(rows, fitText("Message:", width))
	for i := 0; i < postPosterMsgRows; i++ {
		rows = append(rows, fitText("|", width))
	}
	rows = append(rows, fitText("Your Email:", width))
	rows = append(rows, fitText("[you@example.com]", width))
	rows = append(rows, fitText("Send!", width))
	rows = append(rows, "")
	rows = append(rows, wrapPlainText(postCommercialNotice, width)...)
	return rows
}

func wrapPlainText(text string, width int) []string {
	lines := wrapStyledWords(splitStyledWords(text, ""), width)
	out := make([]string, 0, len(lines))
	for _, line := range lines {
		out = append(out, renderStyledLine(line))
	}
	return out
}
