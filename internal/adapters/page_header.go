package adapters

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"
)

const (
	ansiReset   = "\033[0m"
	ansiTopBar  = "\033[48;5;24m\033[1;37m"
	ansiMetaBar = "\033[48;5;252m\033[1;34m"

	defaultPageWidth      = 118
	defaultPageLocation   = "Stanford, California"
	defaultPageRightLabel = "post"
	maxBreadcrumbTitleLen = 24
)

type categorySeedRow struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	ShortName string `json:"short_name"`
}

type subcategorySeedRow struct {
	ID         int64  `json:"id"`
	CategoryID int64  `json:"category_id"`
	Name       string `json:"name"`
}

// BreadcrumbOptions configures adaptive breadcrumb behavior for page headers.
type BreadcrumbOptions struct {
	CategoryID    int64
	SubcategoryID int64
	PostID        int64
	PostTitle     string
}

// PageHeaderOptions configures the reusable SUPost page header.
type PageHeaderOptions struct {
	Width      int
	Location   string
	RightLabel string
	Now        time.Time
	Breadcrumb *BreadcrumbOptions
}

var (
	breadcrumbTaxonomyOnce  sync.Once
	categoryNameByID        map[int64]string
	subcategoryNameByID     map[int64]string
	subcategoryCategoryByID map[int64]int64
)

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

	metaRight := formatUpdatedTimestamp(now)
	metaLeft := " SUPost » " + location
	if opts.Breadcrumb != nil {
		adaptive := buildAdaptiveBreadcrumb(location, *opts.Breadcrumb)
		if width > 0 {
			allowedLeft := width - len([]rune(metaRight))
			if allowedLeft > 1 {
				allowedLeft--
			}
			prefix := " SUPost » "

			if opts.Breadcrumb.PostID > 0 {
				titleLimit := maxBreadcrumbTitleLen
				for titleLimit > 0 && len([]rune(prefix+adaptive)) > allowedLeft {
					titleLimit--
					adaptive = buildAdaptiveBreadcrumbWithTitleLimit(location, *opts.Breadcrumb, titleLimit)
				}
			}
			if allowedLeft > 0 && len([]rune(prefix+adaptive)) > allowedLeft {
				adaptive = truncateToWidth(adaptive, allowedLeft-len([]rune(prefix)))
			}
		}
		if adaptive != "" {
			metaLeft = " SUPost » " + adaptive
		}
	}

	top := renderThreePartLine(
		" SUPost  [__________] [Search]",
		location,
		" "+rightLabel+" ",
		width,
	)
	meta := renderSplitLine(
		metaLeft,
		metaRight,
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

func buildAdaptiveBreadcrumb(location string, breadcrumb BreadcrumbOptions) string {
	return buildAdaptiveBreadcrumbWithTitleLimit(location, breadcrumb, maxBreadcrumbTitleLen)
}

func buildAdaptiveBreadcrumbWithTitleLimit(location string, breadcrumb BreadcrumbOptions, titleLimit int) string {
	parts := []string{strings.TrimSpace(location)}

	categoryID := breadcrumb.CategoryID
	if categoryID <= 0 && breadcrumb.SubcategoryID > 0 {
		categoryID = lookupSubcategoryCategoryID(breadcrumb.SubcategoryID)
	}

	if name := lookupCategoryName(categoryID); name != "" {
		parts = append(parts, name)
	} else if categoryID > 0 {
		parts = append(parts, fmt.Sprintf("category %d", categoryID))
	}

	if name := lookupSubcategoryName(breadcrumb.SubcategoryID); name != "" {
		parts = append(parts, name)
	} else if breadcrumb.SubcategoryID > 0 {
		parts = append(parts, fmt.Sprintf("subcategory %d", breadcrumb.SubcategoryID))
	}

	if breadcrumb.PostID > 0 {
		title := strings.TrimSpace(breadcrumb.PostTitle)
		if title == "" {
			title = fmt.Sprintf("post %d", breadcrumb.PostID)
		}
		title = truncateWithEllipsis(title, titleLimit)
		if title != "" {
			parts = append(parts, title)
		}
	}

	out := make([]string, 0, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part != "" {
			out = append(out, part)
		}
	}
	return strings.Join(out, " » ")
}

func truncateWithEllipsis(value string, maxChars int) string {
	if maxChars <= 0 {
		return ""
	}
	runes := []rune(strings.TrimSpace(value))
	if len(runes) <= maxChars {
		return string(runes)
	}
	return string(runes[:maxChars]) + "..."
}

func truncateToWidth(value string, maxWidth int) string {
	if maxWidth <= 0 {
		return ""
	}
	runes := []rune(strings.TrimSpace(value))
	if len(runes) <= maxWidth {
		return string(runes)
	}
	if maxWidth <= 3 {
		return string(runes[:maxWidth])
	}
	return string(runes[:maxWidth-3]) + "..."
}

func lookupCategoryName(categoryID int64) string {
	loadBreadcrumbTaxonomy()
	if categoryID <= 0 || categoryNameByID == nil {
		return ""
	}
	return categoryNameByID[categoryID]
}

func lookupSubcategoryName(subcategoryID int64) string {
	loadBreadcrumbTaxonomy()
	if subcategoryID <= 0 || subcategoryNameByID == nil {
		return ""
	}
	return subcategoryNameByID[subcategoryID]
}

func lookupSubcategoryCategoryID(subcategoryID int64) int64 {
	loadBreadcrumbTaxonomy()
	if subcategoryID <= 0 || subcategoryCategoryByID == nil {
		return 0
	}
	return subcategoryCategoryByID[subcategoryID]
}

func loadBreadcrumbTaxonomy() {
	breadcrumbTaxonomyOnce.Do(func() {
		categoryNameByID = make(map[int64]string, 32)
		subcategoryNameByID = make(map[int64]string, 128)
		subcategoryCategoryByID = make(map[int64]int64, 128)

		categories, err := readCategorySeedRows()
		if err == nil {
			for _, row := range categories {
				name := strings.TrimSpace(row.ShortName)
				if name == "" {
					name = strings.TrimSpace(row.Name)
				}
				if row.ID > 0 && name != "" {
					categoryNameByID[row.ID] = name
				}
			}
		}

		subcategories, err := readSubcategorySeedRows()
		if err == nil {
			for _, row := range subcategories {
				name := strings.TrimSpace(row.Name)
				if row.ID > 0 && name != "" {
					subcategoryNameByID[row.ID] = name
				}
				if row.ID > 0 && row.CategoryID > 0 {
					subcategoryCategoryByID[row.ID] = row.CategoryID
				}
			}
		}
	})
}

func readCategorySeedRows() ([]categorySeedRow, error) {
	path, err := adapterSeedFilePath("category_rows.json")
	if err != nil {
		return nil, err
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var rows []categorySeedRow
	if err := json.Unmarshal(data, &rows); err != nil {
		return nil, err
	}
	return rows, nil
}

func readSubcategorySeedRows() ([]subcategorySeedRow, error) {
	path, err := adapterSeedFilePath("subcategory_rows.json")
	if err != nil {
		return nil, err
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var rows []subcategorySeedRow
	if err := json.Unmarshal(data, &rows); err != nil {
		return nil, err
	}
	return rows, nil
}

func adapterSeedFilePath(filename string) (string, error) {
	_, srcFile, _, ok := runtime.Caller(0)
	if !ok {
		return "", fmt.Errorf("resolving adapter seed path")
	}
	repoRoot := filepath.Clean(filepath.Join(filepath.Dir(srcFile), "..", ".."))
	return filepath.Join(repoRoot, "testdata", "seed", filename), nil
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
