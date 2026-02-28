# Breadcrumb Category Inference From Subcategory

Date: 2026-02-28

## Summary
Improved adaptive breadcrumb generation so pages can infer and display parent category names when only `subcategory_id` is provided, then validated this behavior in both page-header and search-page tests.

## What Changed

### 1. Breadcrumb taxonomy cache now tracks subcategory → category mapping
- Updated `internal/adapters/page_header.go`:
  - added `subcategoryCategoryByID map[int64]int64`
  - populated this map during taxonomy seed loading.

### 2. Adaptive breadcrumb now infers category from subcategory
- Updated `buildAdaptiveBreadcrumbWithTitleLimit(...)`:
  - if `CategoryID` is missing (`<= 0`) and `SubcategoryID` is present, parent category ID is inferred via taxonomy mapping.
  - breadcrumb rendering then uses inferred category label instead of omitting category.
- Added helper `lookupSubcategoryCategoryID(...)`.

### 3. Tests expanded for inference behavior
- Updated `internal/adapters/page_header_test.go`:
  - added `TestBuildAdaptiveBreadcrumb_InfersCategoryFromSubcategory`.
- Updated `internal/adapters/search_output_test.go`:
  - added `TestRenderSearchResults_SubcategoryOnlyInfersParentCategoryInBreadcrumb`.

## Why This Matters
- Search/post pages that only have a subcategory context now render complete, user-friendly breadcrumbs.
- Keeps header semantics consistent and avoids partial taxonomy paths.

## Files in This Increment
- `internal/adapters/page_header.go`
- `internal/adapters/page_header_test.go`
- `internal/adapters/search_output_test.go`
- `docs/dev/0023-breadcrumb_category_inference_from_subcategory.md`
