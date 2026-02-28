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
	ansiBlue    = "\033[1;34m"
	ansiGray    = "\033[0;37m"
	ansiMagenta = "\033[0;35m"
	ansiHeader  = "\033[48;5;153m\033[1;34m"

	homePageWidth      = 118
	homeRecentWidth    = 54
	homeStripGap       = 2
	homeCalloutWidth   = 28
	homePhotoColumns   = 4
	homePhotoColumnGap = 2
)

type styledWord struct {
	text  string
	color string
}

// RenderHomePosts renders the terminal homepage list.
func RenderHomePosts(w io.Writer, posts []domain.Post) error {
	now := time.Now()

	if err := RenderPageHeader(w, PageHeaderOptions{
		Width:      homePageWidth,
		Location:   "Stanford, California",
		RightLabel: "post",
		Now:        now,
	}); err != nil {
		return err
	}
	if err := renderHomePhotoStrip(w, posts, now, homePageWidth); err != nil {
		return err
	}

	if err := renderHomeOverviewAndRecent(w, posts, now, homePageWidth); err != nil {
		return err
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

func renderHomePhotoStrip(w io.Writer, posts []domain.Post, now time.Time, width int) error {
	photos := selectRecentImagePosts(posts, homePhotoColumns)
	if len(photos) == 0 {
		return nil
	}

	calloutWidth, rightWidth := calculateStripWidths(width)
	columnWidth := photoColumnWidth(rightWidth, homePhotoColumns, homePhotoColumnGap)

	imageURLs := make([]string, 0, len(photos))
	titles := make([]string, 0, len(photos))
	timeAgo := make([]string, 0, len(photos))

	for _, post := range photos {
		imageURLs = append(imageURLs, formatTickerImageURL(post, now))
		titles = append(titles, strings.TrimSpace(post.Name))
		timeAgo = append(timeAgo, formatRelativeTime(postTimestamp(post), now))
	}

	rightRows := make([]string, 0, 8)
	rightRows = append(rightRows, renderWrappedColumnRows(imageURLs, columnWidth, "", homePhotoColumns)...)
	rightRows = append(rightRows, renderColumnRow(titles, columnWidth, ansiBlue, homePhotoColumns))
	rightRows = append(rightRows, renderColumnRow(timeAgo, columnWidth, ansiMagenta, homePhotoColumns))

	leftRows := renderHomeCalloutRows(calloutWidth)
	totalRows := len(rightRows)
	if len(leftRows) > totalRows {
		totalRows = len(leftRows)
	}

	for i := 0; i < totalRows; i++ {
		left := strings.Repeat(" ", calloutWidth)
		right := strings.Repeat(" ", rightWidth)
		if i < len(leftRows) {
			left = leftRows[i]
		}
		if i < len(rightRows) {
			right = rightRows[i]
		}
		if _, err := fmt.Fprintln(w, left+strings.Repeat(" ", homeStripGap)+right); err != nil {
			return err
		}
	}
	return nil
}

func renderHomeOverviewAndRecent(w io.Writer, posts []domain.Post, now time.Time, width int) error {
	leftWidth, rightWidth := calculateStripWidths(width)
	leftRows := renderHomeOverviewRows(leftWidth)
	recentWrapWidth := minInt(homeRecentWidth, rightWidth)
	rightRows := renderRecentPostRows(posts, now, recentWrapWidth, rightWidth)

	totalRows := len(rightRows)
	if len(leftRows) > totalRows {
		totalRows = len(leftRows)
	}

	for i := 0; i < totalRows; i++ {
		left := strings.Repeat(" ", leftWidth)
		right := strings.Repeat(" ", rightWidth)
		if i < len(leftRows) {
			left = leftRows[i]
		}
		if i < len(rightRows) {
			right = rightRows[i]
		}
		if _, err := fmt.Fprintln(w, left+strings.Repeat(" ", homeStripGap)+right); err != nil {
			return err
		}
	}
	return nil
}

func renderHomeOverviewRows(width int) []string {
	rows := []string{
		ansiHeader + centerText("overview", width) + ansiReset,
		renderOverviewRow("housing", "2 hours", width),
		renderOverviewRow("for sale", "5 hours", width),
		renderOverviewRow("jobs", "22 hours", width),
		renderOverviewRow("personals", "19 days", width),
		renderOverviewRow("campus job", "3 hours", width),
		renderOverviewRow("community", "3 hours", width),
		renderOverviewRow("services", "3 hours", width),
	}
	return rows
}

func renderOverviewRow(label, age string, width int) string {
	if width <= 0 {
		return ""
	}

	label = strings.TrimSpace(label)
	age = strings.TrimSpace(age)
	labelLen := len([]rune(label))
	ageLen := len([]rune(age))

	minGap := 1
	if labelLen+minGap+ageLen > width {
		available := width - ageLen - minGap
		if available < 1 {
			return fitText(label+" "+age, width)
		}
		label = fitText(label, available)
		labelLen = len([]rune(label))
	}

	gap := width - labelLen - ageLen
	if gap < 1 {
		gap = 1
	}

	return ansiBlue + label + ansiReset + strings.Repeat(" ", gap) + ansiMagenta + age + ansiReset
}

func renderRecentPostRows(posts []domain.Post, now time.Time, wrapWidth int, sectionWidth int) []string {
	rows := make([]string, 0, len(posts)+1)
	rows = append(rows, ansiHeader+renderHomeHeader("recently posted", sectionWidth)+ansiReset)

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

		lines := wrapStyledWords(words, wrapWidth)
		for _, lineWords := range lines {
			rows = append(rows, renderStyledLine(lineWords))
		}
	}

	return rows
}

func calculateStripWidths(totalWidth int) (calloutWidth, rightWidth int) {
	calloutWidth = homeCalloutWidth
	if totalWidth <= 0 {
		return calloutWidth, 0
	}

	minCallout := 18
	maxCallout := totalWidth / 2
	if maxCallout < minCallout {
		maxCallout = minCallout
	}
	if calloutWidth > maxCallout {
		calloutWidth = maxCallout
	}
	if calloutWidth < minCallout {
		calloutWidth = minCallout
	}

	rightWidth = totalWidth - calloutWidth - homeStripGap
	if rightWidth < 1 {
		rightWidth = 1
	}
	return calloutWidth, rightWidth
}

func renderHomeCalloutRows(width int) []string {
	return []string{
		styleCentered("post to classifieds", width, ansiBlue),
		styleCentered("@stanford.edu required", width, ansiGray),
		strings.Repeat(" ", width),
		styleCentered("post a job", width, ansiBlue),
		styleCentered("post housing", width, ansiBlue),
		styleCentered("post a car", width, ansiBlue),
		strings.Repeat(" ", width),
		styleCentered("open for all emails", width, ansiGray),
	}
}

func styleCentered(text string, width int, color string) string {
	cell := centerText(text, width)
	if strings.TrimSpace(text) == "" || color == "" {
		return cell
	}
	return color + cell + ansiReset
}

func centerText(text string, width int) string {
	if width <= 0 {
		return ""
	}
	trimmed := strings.TrimSpace(text)
	runes := []rune(trimmed)
	if len(runes) > width {
		return fitText(trimmed, width)
	}
	padding := width - len(runes)
	left := padding / 2
	right := padding - left
	return strings.Repeat(" ", left) + trimmed + strings.Repeat(" ", right)
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func selectRecentImagePosts(posts []domain.Post, limit int) []domain.Post {
	if limit <= 0 {
		return nil
	}
	selected := make([]domain.Post, 0, limit)
	for _, post := range posts {
		if !post.HasImage {
			continue
		}
		selected = append(selected, post)
		if len(selected) == limit {
			break
		}
	}
	return selected
}

func formatTickerImageURL(post domain.Post, now time.Time) string {
	timestamp := post.TimePosted
	if timestamp <= 0 {
		if postedAt := postTimestamp(post); !postedAt.IsZero() {
			timestamp = postedAt.Unix()
		}
	}
	if timestamp <= 0 {
		timestamp = now.Unix()
	}
	return fmt.Sprintf("https://supost-prod.s3.amazonaws.com/posts/%d/ticker_%da?%d", post.ID, post.ID, timestamp)
}

func photoColumnWidth(totalWidth, columns, gap int) int {
	if columns <= 0 {
		return totalWidth
	}
	usable := totalWidth - ((columns - 1) * gap)
	if usable < columns {
		return 1
	}
	return usable / columns
}

func renderColumnRow(values []string, width int, color string, columns int) string {
	if columns <= 0 {
		columns = 1
	}
	cells := make([]string, 0, columns)
	for i := 0; i < columns; i++ {
		value := ""
		if i < len(values) {
			value = strings.TrimSpace(values[i])
		}
		cell := fitText(value, width)
		if color != "" && value != "" {
			cell = color + cell + ansiReset
		}
		cells = append(cells, cell)
	}
	return strings.Join(cells, "  ")
}

func renderWrappedColumnRows(values []string, width int, color string, columns int) []string {
	if columns <= 0 {
		columns = 1
	}

	columnLines := make([][]string, columns)
	maxLines := 1
	for i := 0; i < columns; i++ {
		value := ""
		if i < len(values) {
			value = strings.TrimSpace(values[i])
		}
		lines := wrapColumnValue(value, width)
		columnLines[i] = lines
		if len(lines) > maxLines {
			maxLines = len(lines)
		}
	}

	rows := make([]string, 0, maxLines)
	for line := 0; line < maxLines; line++ {
		cells := make([]string, 0, columns)
		for col := 0; col < columns; col++ {
			segment := ""
			if line < len(columnLines[col]) {
				segment = columnLines[col][line]
			}
			cell := fitText(segment, width)
			if color != "" && segment != "" {
				cell = color + cell + ansiReset
			}
			cells = append(cells, cell)
		}
		rows = append(rows, strings.Join(cells, "  "))
	}
	return rows
}

func wrapColumnValue(value string, width int) []string {
	if width <= 0 {
		return []string{""}
	}

	runes := []rune(strings.TrimSpace(value))
	if len(runes) == 0 {
		return []string{""}
	}

	lines := make([]string, 0, (len(runes)/width)+1)
	for len(runes) > width {
		lines = append(lines, string(runes[:width]))
		runes = runes[width:]
	}
	lines = append(lines, string(runes))
	return lines
}

func fitText(value string, width int) string {
	if width <= 0 {
		return ""
	}
	runes := []rune(value)
	if len(runes) > width {
		if width == 1 {
			return "…"
		}
		return string(runes[:width-1]) + "…"
	}
	return value + strings.Repeat(" ", width-len(runes))
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
