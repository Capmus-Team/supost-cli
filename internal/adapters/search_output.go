package adapters

import (
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/Capmus-Team/supost-cli/internal/domain"
)

const (
	searchPageWidth  = 118
	ansiSearchHeader = "\033[48;5;194m\033[1;30m"
)

// RenderSearchResults renders the search result page with date-grouped sections.
func RenderSearchResults(w io.Writer, result domain.SearchResultPage) error {
	now := time.Now()
	if err := RenderPageHeader(w, PageHeaderOptions{
		Width:      searchPageWidth,
		Location:   "Stanford, California",
		RightLabel: "post",
		Now:        now,
		Breadcrumb: &BreadcrumbOptions{
			CategoryID:    result.CategoryID,
			SubcategoryID: result.SubcategoryID,
		},
	}); err != nil {
		return err
	}

	title := searchResultTitle(result.Query)
	if _, err := fmt.Fprintln(w, ansiSearchHeader+renderHomeHeader(title, searchPageWidth)+ansiReset); err != nil {
		return err
	}
	if _, err := fmt.Fprintln(w); err != nil {
		return err
	}

	if err := renderSearchGroupedPosts(w, result.Posts, searchPageWidth); err != nil {
		return err
	}

	paginationLabels := make([]string, 0, 2)
	if result.Page > 1 {
		paginationLabels = append(paginationLabels, fmt.Sprintf("previous %d posts", result.PerPage))
	}
	if result.HasMore {
		paginationLabels = append(paginationLabels, fmt.Sprintf("next %d posts", result.PerPage))
	}
	if len(paginationLabels) > 0 {
		if _, err := fmt.Fprintln(w); err != nil {
			return err
		}
		for _, label := range paginationLabels {
			if _, err := fmt.Fprintln(w, styleCentered(label, searchPageWidth, ansiBlue)); err != nil {
				return err
			}
		}
	}

	if _, err := fmt.Fprintln(w); err != nil {
		return err
	}
	return RenderPageFooter(w, PageFooterOptions{Width: searchPageWidth})
}

func renderSearchGroupedPosts(w io.Writer, posts []domain.Post, width int) error {
	if len(posts) == 0 {
		_, err := fmt.Fprintln(w, styleCentered("No posts found for this page.", width, ansiGray))
		return err
	}

	lastHeader := ""
	for _, post := range posts {
		header := formatSearchDateHeader(post)
		if header != lastHeader {
			if _, err := fmt.Fprintln(w, ansiSearchHeader+fitText(" "+header, width)+ansiReset); err != nil {
				return err
			}
			lastHeader = header
		}

		line := renderSearchPostLine(post)
		lines := wrapStyledWords(line, width)
		for _, words := range lines {
			if _, err := fmt.Fprintln(w, renderStyledLine(words)); err != nil {
				return err
			}
		}
	}
	return nil
}

func renderSearchPostLine(post domain.Post) []styledWord {
	words := make([]styledWord, 0, 20)
	words = append(words, splitStyledWords(formatPostTitle(post), ansiBlue)...)
	words = append(words, splitStyledWords(formatDisplayEmail(post.Email), ansiGray)...)
	if post.HasImage {
		words = append(words, styledWord{text: "📷"})
	}
	return words
}

func formatSearchDateHeader(post domain.Post) string {
	ts := postTimestamp(post)
	if ts.IsZero() {
		return "Unknown date"
	}
	return ts.Format("Mon, Jan 2")
}

func searchResultTitle(query string) string {
	query = strings.TrimSpace(query)
	if query == "" {
		return "all posts"
	}
	return fmt.Sprintf("search: %s", query)
}
