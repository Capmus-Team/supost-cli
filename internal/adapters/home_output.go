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
)

// RenderHomePosts renders the terminal homepage list.
func RenderHomePosts(w io.Writer, posts []domain.Post) error {
	if _, err := fmt.Fprintf(w, "%s%s%s\n", ansiHeader, renderHomeHeader("recently posted", 72), ansiReset); err != nil {
		return err
	}

	for _, post := range posts {
		title := formatPostTitle(post)
		photo := ""
		if post.HasImage {
			photo = " 📷"
		}
		timeAgo := formatRelativeTime(postTimestamp(post), time.Now())

		if _, err := fmt.Fprintf(
			w,
			"%s%s%s %s%s%s%s %s%s%s\n",
			ansiBlue, title, ansiReset,
			ansiGray, post.Email, ansiReset,
			photo,
			ansiMagenta, timeAgo, ansiReset,
		); err != nil {
			return err
		}
	}

	return nil
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
