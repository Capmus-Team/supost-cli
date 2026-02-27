# Shared Page Header + Home Photo Strip

Date: 2026-02-27

## Summary
Refactored shared SUPost header rendering into a reusable adapter and updated home output to include a four-column image ticker strip before the recent posts list.

## What Changed

### 1. New shared page header adapter
- Added `internal/adapters/page_header.go` with:
  - `PageHeaderOptions`
  - `RenderPageHeader(...)`
  - shared line helpers (`renderThreePartLine`, `renderSplitLine`)
  - `formatUpdatedTimestamp(...)`
- Header now supports configurable width/location/right-label and deterministic timestamp injection for testing.

### 2. Home output now uses shared page header
- Updated `internal/adapters/home_output.go`:
  - removed duplicated top/meta bar rendering from home-specific code
  - calls `RenderPageHeader(...)` with home defaults

### 3. Added home photo strip section
- Added home-specific strip rendering before `recently posted`:
  - `renderHomePhotoStrip`
  - `selectRecentImagePosts`
  - `formatTickerImageURL`
  - `photoColumnWidth`
  - `renderColumnRow`
  - `fitText`
- Behavior:
  - picks up to 4 recent posts with images
  - renders URL row, title row (blue), and relative-time row (magenta)
  - uses fixed-width 4-column layout with truncation/ellipsis where needed

### 4. Tests updated and expanded
- Added `internal/adapters/page_header_test.go` for:
  - timestamp formatting
  - right-aligned metadata line behavior
- Updated `internal/adapters/home_output_test.go` for:
  - image post selection order and cap
  - ticker URL formatting
  - column row width behavior
- Removed home tests that are now covered in shared header tests.

## Why This Matters
- Eliminates duplicated header logic and makes future non-home pages easier to implement.
- Improves home visual fidelity with an explicit image ticker section.
- Keeps layout/format logic test-backed and easier to refactor safely.

## Files in This Increment
- `internal/adapters/page_header.go`
- `internal/adapters/page_header_test.go`
- `internal/adapters/home_output.go`
- `internal/adapters/home_output_test.go`
- `docs/dev/0008-shared_page_header_and_home_photo_strip.md`
