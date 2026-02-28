# Home Overview + Recent Two-Column Layout

Date: 2026-02-27

## Summary
Restructured the lower home section to render an `overview` block in the left column and `recently posted` feed in the right column, replacing the previous single full-width recent-post rendering.

## What Changed

### 1. New combined section renderer for lower home area
- Updated `RenderHomePosts` to call:
  - `renderHomeOverviewAndRecent(...)`
- Removed the direct full-width `recently posted` loop from the top-level function.

### 2. New overview panel rendering
- Added:
  - `renderHomeOverviewRows(width)`
  - `renderOverviewRow(label, age, width)`
- Overview panel now includes static rows:
  - housing, for sale, jobs, personals, campus job, community, services
- Each row uses blue labels and magenta age text.

### 3. New recent-post row rendering helper
- Added:
  - `renderRecentPostRows(posts, now, wrapWidth, sectionWidth)`
- Includes a right-column `recently posted` header and wrapped body rows.
- Wrap width is constrained by `minInt(homeRecentWidth, rightWidth)`.

### 4. Utility addition
- Added `minInt(a, b int)` helper for width calculation.

### 5. Tests expanded for overview/recent behavior
- Updated `internal/adapters/home_output_test.go` with:
  - `TestRenderHomeOverviewRows_ContainsRequestedCopy`
  - `TestRenderRecentPostRows_RespectsWrapWidth`
- Verifies expected overview content and right-column wrapping limits.

## Why This Matters
- Better visual parity with classic SUPost composition.
- Clearer information hierarchy: category snapshot on left, freshest posts on right.
- Safer evolution due to dedicated test coverage for both new sections.

## Files in This Increment
- `internal/adapters/home_output.go`
- `internal/adapters/home_output_test.go`
- `docs/dev/0011-home_overview_plus_recent_two_column_layout.md`
