# Home Recent + Featured Jobs Dual-Column Content

Date: 2026-02-28

## Summary
Extended the home content area so the right-side section can render two aligned columns: `recently posted` and a new `featured job posts` panel sourced from active off-campus jobs.

## What Changed

### 1. Added dual-column home content renderer
- Updated `internal/adapters/home_output.go` with:
  - `renderHomeRecentAndFeaturedRows(...)`
  - `splitHomeContentWidths(...)`
  - `combineHomeContentColumns(...)`
  - `padANSIVisibleWidth(...)`
  - `ansiVisibleRuneLen(...)`
- The content region now splits width into left/right columns with a fixed gap and keeps ANSI-colored output visually aligned.

### 2. Added featured jobs selection and rendering
- Added `selectFeaturedJobPosts(...)` to:
  - include only `active` posts
  - include only `CategoryJobsOffCampus`
  - sort newest-first by post timestamp (with ID tie-break)
  - cap output to `homeFeaturedLimit` (3)
- Added `renderFeaturedJobPostRows(...)` with its own section header and wrapped title/email lines.

### 3. Wired home sidebar/content integration
- Updated `internal/adapters/home_sidebar.go`:
  - replaced direct `renderRecentPostRows(...)` usage with `renderHomeRecentAndFeaturedRows(...)` in `renderHomeOverviewAndRecent(...)`.

### 4. Added tests for behavior and ordering
- Updated `internal/adapters/home_output_test.go` with:
  - `TestSelectFeaturedJobPosts_FiltersActiveJobsAndOrdersNewest`
  - `TestRenderHomeRecentAndFeaturedRows_ContainsFeaturedSection`
- Coverage verifies filtering rules, ordering, and that the featured section appears in rendered home content.

## Why This Matters
- Improves home-page scanability by surfacing current job opportunities without removing the general recent feed.
- Keeps terminal layout stable by accounting for ANSI escape sequences when padding columns.
- Preserves adapter-level test coverage for the new selection and render flow.

## Files in This Increment
- `internal/adapters/home_output.go`
- `internal/adapters/home_sidebar.go`
- `internal/adapters/home_output_test.go`
- `docs/dev/0016-home_recent_plus_featured_jobs_dual_column.md`
