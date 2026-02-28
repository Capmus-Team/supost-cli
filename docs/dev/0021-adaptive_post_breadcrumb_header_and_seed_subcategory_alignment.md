# Adaptive Post Breadcrumb Header + Seed Subcategory Alignment

Date: 2026-02-28

## Summary
Added adaptive taxonomy-aware breadcrumbs to the shared page header, wired post pages to pass breadcrumb context, and updated in-memory post seed subcategory IDs to match taxonomy seed data.

## What Changed

### 1. Shared page header now supports adaptive breadcrumbs
- Updated `internal/adapters/page_header.go`:
  - added `BreadcrumbOptions` and optional `PageHeaderOptions.Breadcrumb`
  - added adaptive breadcrumb generation using location/category/subcategory/post title
  - added width-aware truncation logic with ellipsis for long post titles
  - loads category/subcategory names from `testdata/seed/category_rows.json` and `subcategory_rows.json` (memoized with `sync.Once`) for readable breadcrumb labels.

### 2. Post page now passes breadcrumb context into header
- Updated `internal/adapters/post_output.go`:
  - `RenderPostPage` now passes category ID, subcategory ID, post ID, and post title via `BreadcrumbOptions`.

### 3. In-memory post seed taxonomy IDs aligned
- Updated `internal/repository/inmemory.go` seed posts:
  - housing sample posts moved to subcategory `59` (`rooms & shares`)
  - for-sale sample posts moved to subcategory `9` (`clothing & accessories`)
- This aligns sample post taxonomy with breadcrumb name lookups and category seed data.

### 4. Header tests expanded
- Updated `internal/adapters/page_header_test.go`:
  - added coverage for adaptive breadcrumb generation and title truncation
  - added rendering test verifying breadcrumb content + timestamp in header meta line.

## Why This Matters
- Post pages now show richer navigation context (`location » category » subcategory » post`) in the existing shared header.
- Breadcrumb labels are human-readable from real seed taxonomy data instead of numeric IDs.
- Seed alignment keeps local prototype output consistent with expected category/subcategory naming.

## Files in This Increment
- `internal/adapters/page_header.go`
- `internal/adapters/page_header_test.go`
- `internal/adapters/post_output.go`
- `internal/repository/inmemory.go`
- `docs/dev/0021-adaptive_post_breadcrumb_header_and_seed_subcategory_alignment.md`
