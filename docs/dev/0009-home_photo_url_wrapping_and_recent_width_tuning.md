# Home Photo URL Wrapping + Recent Width Tuning

Date: 2026-02-27

## Summary
Improved home page readability by separating overall page width from recent-post text width and by changing the photo URL strip to wrap long values across lines instead of truncating with ellipsis.

## What Changed

### 1. Split width constants for layout control
- Updated `internal/adapters/home_output.go`:
  - replaced single `homeRowWidth` with:
    - `homePageWidth = 118` for header/photo strip and section header
    - `homeRecentWidth = 54` for wrapped post text rows
- `RenderHomePosts` now uses these widths by section.

### 2. Photo URL row now wraps instead of truncates
- `renderHomePhotoStrip` now renders image URL rows using `renderWrappedColumnRows(...)`.
- Added helpers:
  - `renderWrappedColumnRows`
  - `wrapColumnValue`
- Behavior:
  - long column values split into multiple fixed-width lines
  - preserves all characters (no ellipsis)
  - still honors column alignment.

### 3. Tests added for non-truncating wrap behavior
- Updated `internal/adapters/home_output_test.go` with:
  - `TestWrapColumnValue_DoesNotTruncate`
  - `TestRenderWrappedColumnRows_ContainsNoEllipsis`
- Also added `strings` import for these assertions.

## Why This Matters
- Long ticker URLs are fully visible in the photo strip for debugging/verification.
- Recent listing rows keep tighter width for better terminal scanning.
- Layout behavior is test-backed and less likely to regress.

## Files in This Increment
- `internal/adapters/home_output.go`
- `internal/adapters/home_output_test.go`
- `docs/dev/0009-home_photo_url_wrapping_and_recent_width_tuning.md`
