# Home Callout Column + Strip Layout Refinement

Date: 2026-02-27

## Summary
Enhanced the home photo strip layout by introducing a dedicated left callout column and a configurable right photo-grid region, improving visual structure and readability.

## What Changed

### 1. New layout constants for strip composition
- Updated `internal/adapters/home_output.go` with explicit constants:
  - `homeStripGap`
  - `homeCalloutWidth`
  - `homePhotoColumns`
  - `homePhotoColumnGap`
- Keeps strip geometry configurable and easier to tune.

### 2. Split-width strip rendering
- `renderHomePhotoStrip` now computes left/right widths via `calculateStripWidths(totalWidth)`.
- Right side continues to render wrapped image URLs + title/time rows.
- Left side now renders a dedicated callout block aligned by row.

### 3. Added callout copy renderer
- New helpers:
  - `renderHomeCalloutRows(width)`
  - `styleCentered(text, width, color)`
  - `centerText(text, width)`
- New callout content includes:
  - `post to classifieds`
  - `@stanford.edu required`
  - `post a job`
  - `post housing`
  - `post a car`
  - `open for all emails`

### 4. Test coverage for new callout and centering behavior
- Updated `internal/adapters/home_output_test.go` with:
  - `TestRenderHomeCalloutRows_ContainsRequestedCopy`
  - `TestCenterText_ProducesFixedWidth`
- Confirms expected copy is present and centering logic preserves fixed width.

## Why This Matters
- Home output now has a clearer two-column structure resembling traditional SUPost-style composition.
- The callout column improves discoverability of posting actions.
- Layout behavior is now test-backed for safer future UI tweaks.

## Files in This Increment
- `internal/adapters/home_output.go`
- `internal/adapters/home_output_test.go`
- `docs/dev/0010-home_callout_column_and_strip_layout_refinement.md`
