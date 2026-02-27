package adapters

import (
	"fmt"
	"io"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/Capmus-Team/supost-cli/internal/domain"
)

const (
	ansiReset   = "\033[0m"
	ansiBlue    = "\033[1;34m"
	ansiGray    = "\033[0;37m"
	ansiMagenta = "\033[0;35m"
	ansiHeader  = "\033[48;5;153m\033[1;34m"
	ansiTopBar  = "\033[48;5;24m\033[1;37m"
	ansiMetaBar = "\033[48;5;252m\033[1;34m"

	homeRowWidth = 118
)

type styledWord struct {
	text  string
	color string
}

// RenderHomePosts renders the terminal homepage list.
func RenderHomePosts(w io.Writer, posts []domain.Post) error {
	now := time.Now()

	if _, err := fmt.Fprintf(w, "%s%s%s\n", ansiTopBar, renderHomeTopBar(homeRowWidth), ansiReset); err != nil {
		return err
	}
	if _, err := fmt.Fprintf(w, "%s%s%s\n", ansiMetaBar, renderHomeMetaBar(now, homeRowWidth), ansiReset); err != nil {
		return err
	}
	if _, err := fmt.Fprintf(w, "%s%s%s\n", ansiHeader, renderHomeHeader("recently posted", homeRowWidth), ansiReset); err != nil {
		return err
	}

	for _, post := range posts {
		title := formatPostTitle(post)
		email := formatDisplayEmail(post.Email)
		timeAgo := formatRelativeTime(postTimestamp(post), now)

		words := make([]styledWord, 0, 16)
		words = append(words, splitStyledWords(title, ansiBlue)...)
		words = append(words, splitStyledWords(email, ansiGray)...)
		if post.HasImage {
			words = append(words, styledWord{text: "📷"})
		}
		words = append(words, splitStyledWords(timeAgo, ansiMagenta)...)

		lines := wrapStyledWords(words, homeRowWidth)
		for _, lineWords := range lines {
			if _, err := fmt.Fprintln(w, renderStyledLine(lineWords)); err != nil {
				return err
			}
		}
	}

	return nil
}

func renderHomeTopBar(width int) string {
	left := " SUPost  [__________] [Search]"
	center := "Stanford, California"
	right := "post "
	return renderThreePartLine(left, center, right, width)
}

func renderHomeMetaBar(now time.Time, width int) string {
	left := " SUPost » Stanford, California"
	right := formatHomeUpdatedTimestamp(now)
	return renderSplitLine(left, right, width)
}

func formatHomeUpdatedTimestamp(now time.Time) string {
	return now.Format("Mon, Jan 2, 2006 03:04 PM") + " - Updated"
}

func renderHomeHeader(text string, width int) string {
	if width < len(text)+2 {
		return " " + text + " "
	}

	padding := width - len(text)
	left := padding / 2
	right := padding - left
	return strings.Repeat(" ", left) + text + strings.Repeat(" ", right)
}

func formatPostTitle(post domain.Post) string {
	title := strings.TrimSpace(post.Name)
	if title == "" {
		title = "(untitled post)"
	}

	price := formatPrice(post.Price, post.HasPrice)
	if price == "" {
		return title
	}
	return title + " - " + price
}

func formatPrice(price float64, hasPrice bool) string {
	if !hasPrice {
		return ""
	}
	if price <= 0 {
		return "Free"
	}

	dollars := int64(math.Round(price))
	if math.Abs(price-float64(dollars)) < 0.001 {
		return "$" + formatIntWithCommas(dollars)
	}
	return fmt.Sprintf("$%.2f", price)
}

func formatIntWithCommas(v int64) string {
	s := strconv.FormatInt(v, 10)
	if len(s) <= 3 {
		return s
	}

	var out []byte
	prefix := len(s) % 3
	if prefix == 0 {
		prefix = 3
	}
	out = append(out, s[:prefix]...)
	for i := prefix; i < len(s); i += 3 {
		out = append(out, ',')
		out = append(out, s[i:i+3]...)
	}
	return string(out)
}

func postTimestamp(post domain.Post) time.Time {
	if !post.TimePostedAt.IsZero() {
		return post.TimePostedAt
	}
	if post.TimePosted > 0 {
		return time.Unix(post.TimePosted, 0)
	}
	if !post.CreatedAt.IsZero() {
		return post.CreatedAt
	}
	return time.Time{}
}

func formatRelativeTime(from, now time.Time) string {
	if from.IsZero() {
		return "about 0 minutes"
	}

	if from.After(now) {
		from = now
	}
	d := now.Sub(from)

	if d < time.Minute {
		return "about 1 minute"
	}
	if d < time.Hour {
		minutes := int(d.Minutes())
		if minutes <= 1 {
			return "about 1 minute"
		}
		return fmt.Sprintf("about %d minutes", minutes)
	}
	if d < 24*time.Hour {
		hours := int(d.Hours())
		if hours <= 1 {
			return "about 1 hour"
		}
		return fmt.Sprintf("about %d hours", hours)
	}

	days := int(d.Hours() / 24)
	if days <= 1 {
		return "about 1 day"
	}
	return fmt.Sprintf("about %d days", days)
}

func formatDisplayEmail(email string) string {
	normalized := strings.TrimSpace(email)
	if strings.Contains(strings.ToLower(normalized), "stanford.edu") {
		return "@stanford.edu"
	}
	return ""
}

func splitStyledWords(text, color string) []styledWord {
	fields := strings.Fields(strings.TrimSpace(text))
	words := make([]styledWord, 0, len(fields))
	for _, word := range fields {
		words = append(words, styledWord{text: word, color: color})
	}
	return words
}

func wrapStyledWords(words []styledWord, width int) [][]styledWord {
	if width <= 0 {
		return [][]styledWord{words}
	}
	if len(words) == 0 {
		return [][]styledWord{{}}
	}

	lines := make([][]styledWord, 0, len(words))
	current := make([]styledWord, 0, 8)
	currentWidth := 0

	for _, word := range words {
		wordLen := len([]rune(word.text))
		if wordLen > width {
			if len(current) > 0 {
				lines = append(lines, current)
				current = make([]styledWord, 0, 8)
				currentWidth = 0
			}

			runes := []rune(word.text)
			for len(runes) > width {
				chunk := string(runes[:width])
				lines = append(lines, []styledWord{{text: chunk, color: word.color}})
				runes = runes[width:]
			}
			if len(runes) > 0 {
				current = append(current, styledWord{text: string(runes), color: word.color})
				currentWidth = len(runes)
			}
			continue
		}

		needed := wordLen
		if len(current) > 0 {
			needed++
		}
		if currentWidth+needed <= width {
			current = append(current, word)
			currentWidth += needed
			continue
		}

		lines = append(lines, current)
		current = []styledWord{word}
		currentWidth = wordLen
	}

	if len(current) > 0 {
		lines = append(lines, current)
	}
	return lines
}

func renderStyledLine(words []styledWord) string {
	if len(words) == 0 {
		return ""
	}
	var b strings.Builder
	for i, word := range words {
		if i > 0 {
			b.WriteByte(' ')
		}
		if word.color != "" {
			b.WriteString(word.color)
			b.WriteString(word.text)
			b.WriteString(ansiReset)
		} else {
			b.WriteString(word.text)
		}
	}
	return b.String()
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
