package adapters

import (
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/Capmus-Team/supost-cli/internal/domain"
)

const postCreatePageWidth = homePageWidth

var postCreateCategoryMenuOrder = []int64{5, 8, 3, 4, 9, 7, 1, 2}

var postCreateCategoryLabels = map[int64]string{
	1: "campus job",
	2: "job off-campus",
	3: "housing (offering)",
	4: "housing (need)",
	5: "for sale / wanted",
	7: "service offered",
	8: "personals",
	9: "community",
}

// RenderPostCreatePage renders the staged post-create page.
func RenderPostCreatePage(w io.Writer, page domain.PostCreatePage) error {
	if err := RenderPageHeader(w, PageHeaderOptions{
		Width:      postCreatePageWidth,
		Location:   "Stanford, California",
		RightLabel: "post",
		Now:        time.Now(),
		Breadcrumb: postCreateBreadcrumb(page),
	}); err != nil {
		return err
	}

	if err := renderPostCreateCampusBand(w, postCreatePageWidth); err != nil {
		return err
	}
	if _, err := fmt.Fprintln(w); err != nil {
		return err
	}

	switch page.Stage {
	case domain.PostCreateStageChooseCategory:
		if err := renderPostCreateCategoryStep(w, page, postCreatePageWidth); err != nil {
			return err
		}
	case domain.PostCreateStageChooseSubcategory:
		if err := renderPostCreateSubcategoryStep(w, page, postCreatePageWidth); err != nil {
			return err
		}
	case domain.PostCreateStageForm:
		if err := renderPostCreateFormStep(w, page, postCreatePageWidth); err != nil {
			return err
		}
	default:
		if _, err := fmt.Fprintln(w, styleCentered("Unsupported post-create stage.", postCreatePageWidth, ansiGray)); err != nil {
			return err
		}
	}

	if _, err := fmt.Fprintln(w); err != nil {
		return err
	}
	return RenderPageFooter(w, PageFooterOptions{Width: postCreatePageWidth})
}

func postCreateBreadcrumb(page domain.PostCreatePage) *BreadcrumbOptions {
	if page.CategoryID <= 0 && page.SubcategoryID <= 0 {
		return nil
	}
	return &BreadcrumbOptions{
		CategoryID:    page.CategoryID,
		SubcategoryID: page.SubcategoryID,
	}
}

func renderPostCreateCampusBand(w io.Writer, width int) error {
	_, err := fmt.Fprintln(w, ansiHeader+renderSplitLine("", " (Stanford)", width)+ansiReset)
	return err
}

func renderPostCreateCategoryStep(w io.Writer, page domain.PostCreatePage, width int) error {
	if _, err := fmt.Fprintln(w, fitText("What type of post is this?", width)); err != nil {
		return err
	}
	if _, err := fmt.Fprintln(w); err != nil {
		return err
	}

	for _, category := range orderedPostCreateCategories(page.Categories) {
		label := postCreateCategoryLabel(category)
		if _, err := fmt.Fprintln(w, ansiBlue+fitText("  "+label, width)+ansiReset); err != nil {
			return err
		}
		if _, err := fmt.Fprintln(w); err != nil {
			return err
		}
	}
	return nil
}

func orderedPostCreateCategories(categories []domain.Category) []domain.Category {
	byID := make(map[int64]domain.Category, len(categories))
	for _, category := range categories {
		byID[category.ID] = category
	}

	ordered := make([]domain.Category, 0, len(postCreateCategoryMenuOrder))
	for _, id := range postCreateCategoryMenuOrder {
		category, ok := byID[id]
		if !ok {
			continue
		}
		ordered = append(ordered, category)
	}
	return ordered
}

func postCreateCategoryLabel(category domain.Category) string {
	if label, ok := postCreateCategoryLabels[category.ID]; ok {
		return label
	}
	if short := strings.TrimSpace(category.ShortName); short != "" {
		return short
	}
	return strings.TrimSpace(category.Name)
}

func renderPostCreateSubcategoryStep(w io.Writer, page domain.PostCreatePage, width int) error {
	if _, err := fmt.Fprintln(w, fitText("Please choose a category:", width)); err != nil {
		return err
	}
	if _, err := fmt.Fprintln(w); err != nil {
		return err
	}

	for _, subcategory := range page.Subcategories {
		label := strings.TrimSpace(subcategory.Name)
		if label == "" {
			continue
		}
		if _, err := fmt.Fprintln(w, ansiBlue+fitText("  "+label, width)+ansiReset); err != nil {
			return err
		}
		if _, err := fmt.Fprintln(w); err != nil {
			return err
		}
	}
	return nil
}

func renderPostCreateFormStep(w io.Writer, page domain.PostCreatePage, width int) error {
	if err := renderPostCreateFormField(w, width, "Post Title:", "[title]"); err != nil {
		return err
	}
	if domain.CategoryPriceAllowed(page.CategoryID) {
		if err := renderPostCreateFormField(w, width, "Price:", "[price]"); err != nil {
			return err
		}
	}
	if err := renderPostCreateFormField(w, width, "Post Description:", ""); err != nil {
		return err
	}

	for i := 0; i < 8; i++ {
		if _, err := fmt.Fprintln(w, fitText("  |", width)); err != nil {
			return err
		}
	}

	if err := renderPostCreateFormField(w, width, "Your Stanford Email:", "[you@stanford.edu]"); err != nil {
		return err
	}
	if _, err := fmt.Fprintln(w, fitText("  @stanford.edu, @alumni.stanford.edu, @stanfordalumni.org, @stanfordhealthcare.org", width)); err != nil {
		return err
	}
	if _, err := fmt.Fprintln(w, fitText("  We need your email to send you the self-publishing link.", width)); err != nil {
		return err
	}

	if err := renderPostCreateFormField(w, width, "Photos:", "[1] [2] [3] [4] (optional)"); err != nil {
		return err
	}
	if _, err := fmt.Fprintln(w); err != nil {
		return err
	}
	if _, err := fmt.Fprintln(w, fitText("[Preview]", width)); err != nil {
		return err
	}
	if _, err := fmt.Fprintln(w); err != nil {
		return err
	}
	if _, err := fmt.Fprintln(w, styleCentered(postHousingNotice, width, ansiGray)); err != nil {
		return err
	}
	return nil
}

func renderPostCreateFormField(w io.Writer, width int, label, value string) error {
	line := "  " + strings.TrimSpace(label)
	if strings.TrimSpace(value) != "" {
		line += " " + strings.TrimSpace(value)
	}
	_, err := fmt.Fprintln(w, fitText(line, width))
	return err
}
