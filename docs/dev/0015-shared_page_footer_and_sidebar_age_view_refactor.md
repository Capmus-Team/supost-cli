# Shared Page Footer + Sidebar Age-View Refactor

Date: 2026-02-27

## Summary
Added a reusable page footer for home-style pages and refactored sidebar rendering to precompute compact/detailed age strings once per section.

## What Changed

### 1. New reusable page footer adapter
- Added `internal/adapters/page_footer.go` with:
  - `PageFooterOptions`
  - `RenderPageFooter(...)`
- Footer content includes:
  - nav links (`post a job`, `post housing`, `post a car`, `about`, `contact`, `privacy`, `terms`, `help`)
  - attribution and disclaimer lines.

### 2. Home output now renders footer
- Updated `internal/adapters/home_output.go`:
  - `RenderHomePosts` now adds a blank spacer line and invokes `RenderPageFooter(...)` after overview/recent content.

### 3. Sidebar age formatting refactor
- Updated `internal/adapters/home_sidebar.go`:
  - added `homeSectionAgeView` struct (`Section`, `CompactAge`, `DetailedAge`)
  - added `buildHomeSectionAgeViews(...)`
  - `renderHomeOverviewRows` and `renderHomeCategoryDetailsRows` now consume precomputed age views
- Result: avoids recomputing age strings in multiple render passes.

### 4. Tests updated and expanded
- Added `internal/adapters/page_footer_test.go` to validate expected footer strings.
- Updated `internal/adapters/home_output_test.go` to use new section-age-view flow in overview/detail tests.

## Why This Matters
- Footer rendering is now modular and reusable across page-like commands.
- Sidebar rendering path is cleaner and slightly more efficient by centralizing age formatting.
- Footer and refactored sidebar behavior are test-covered.

## Files in This Increment
- `internal/adapters/page_footer.go`
- `internal/adapters/page_footer_test.go`
- `internal/adapters/home_output.go`
- `internal/adapters/home_sidebar.go`
- `internal/adapters/home_output_test.go`
- `docs/dev/0015-shared_page_footer_and_sidebar_age_view_refactor.md`
