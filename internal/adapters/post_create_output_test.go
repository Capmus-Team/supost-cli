package adapters

import (
	"bytes"
	"strings"
	"testing"

	"github.com/Capmus-Team/supost-cli/internal/domain"
)

func TestRenderPostCreatePage_ChooseCategory(t *testing.T) {
	var out bytes.Buffer
	page := domain.PostCreatePage{
		Stage: domain.PostCreateStageChooseCategory,
		Categories: []domain.Category{
			{ID: 5, Name: "for sale/wanted", ShortName: "for sale"},
			{ID: 8, Name: "personals/dating", ShortName: "personals"},
			{ID: 3, Name: "housing (offering)", ShortName: "housing"},
		},
	}

	if err := RenderPostCreatePage(&out, page); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	plain := stripANSI(out.String())
	for _, needle := range []string{
		"SUPost » Stanford, California",
		"What type of post is this?",
		"for sale / wanted",
		"personals",
		"housing (offering)",
	} {
		if !strings.Contains(plain, needle) {
			t.Fatalf("missing %q in output", needle)
		}
	}
	for _, unwanted := range []string{"resumes", "events"} {
		if strings.Contains(plain, unwanted) {
			t.Fatalf("unexpected %q in category menu", unwanted)
		}
	}
}

func TestRenderPostCreatePage_ChooseSubcategory(t *testing.T) {
	var out bytes.Buffer
	page := domain.PostCreatePage{
		Stage:      domain.PostCreateStageChooseSubcategory,
		CategoryID: 8,
		Subcategories: []domain.Subcategory{
			{ID: 130, CategoryID: 8, Name: "friendship"},
			{ID: 135, CategoryID: 8, Name: "general romance"},
		},
	}

	if err := RenderPostCreatePage(&out, page); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	plain := stripANSI(out.String())
	for _, needle := range []string{
		"SUPost » Stanford, California » personals",
		"Please choose a category:",
		"friendship",
		"general romance",
	} {
		if !strings.Contains(plain, needle) {
			t.Fatalf("missing %q in output", needle)
		}
	}
}

func TestRenderPostCreatePage_FormStage(t *testing.T) {
	var out bytes.Buffer
	page := domain.PostCreatePage{
		Stage:         domain.PostCreateStageForm,
		CategoryID:    5,
		SubcategoryID: 14,
	}

	if err := RenderPostCreatePage(&out, page); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	plain := stripANSI(out.String())
	for _, needle := range []string{
		"SUPost » Stanford, California » for sale » furniture",
		"Post Title: [title]",
		"Price: [price]",
		"Post Description:",
		"Your Stanford Email: [you@stanford.edu]",
		"Photos: [1] [2] [3] [4] (optional)",
		"[Preview]",
	} {
		if !strings.Contains(plain, needle) {
			t.Fatalf("missing %q in output", needle)
		}
	}
}

func TestRenderPostCreatePage_FormStage_NoPriceForPersonals(t *testing.T) {
	var out bytes.Buffer
	page := domain.PostCreatePage{
		Stage:         domain.PostCreateStageForm,
		CategoryID:    8,
		SubcategoryID: 130,
	}

	if err := RenderPostCreatePage(&out, page); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	plain := stripANSI(out.String())
	if strings.Contains(plain, "Price: [price]") {
		t.Fatalf("did not expect price field for personals form")
	}
	if !strings.Contains(plain, "Post Title: [title]") {
		t.Fatalf("missing post title field")
	}
}
