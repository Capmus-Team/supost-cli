# Home Header Bars + Timestamp Metadata

Date: 2026-02-27

## Summary
Enhanced `supost home` terminal output to better match a classic SUPost-style layout by adding top/meta bars, widening row layout, and rendering a right-aligned updated timestamp.

## What Changed

### 1. Added top-level terminal bars
- Updated `internal/adapters/home_output.go` to render two new bars before `recently posted`:
  - top bar (`ansiTopBar`) with SUPost/search/location styling
  - meta bar (`ansiMetaBar`) with location breadcrumb on the left and timestamp on the right

### 2. Increased render width for homepage layout
- `homeRowWidth` increased from `54` to `118` to better fit multi-part header content and improve post row readability.

### 3. Timestamp formatting for metadata row
- Added `formatHomeUpdatedTimestamp(now time.Time)` returning:
  - `Mon, Jan 2, 2006 03:04 PM - Updated`
- Added `renderHomeMetaBar(now, width)` to right-align the timestamp.

### 4. New line composition helpers
- Added:
  - `renderThreePartLine(left, center, right, width)`
  - `renderSplitLine(left, right, width)`
- These support fixed-width alignment for top/meta bars.

### 5. Tests expanded for new metadata behavior
- Updated `internal/adapters/home_output_test.go` with:
  - `TestFormatHomeUpdatedTimestamp`
  - `TestRenderHomeMetaBar_RightAlignsTimestamp`
- Existing wrapping/email-format tests remain in place.

## Why This Matters
- Home screen now better matches target CLI visual style.
- Timestamp display is deterministic and test-covered.
- Layout helper functions make future style adjustments lower risk.

## Files in This Increment
- `internal/adapters/home_output.go`
- `internal/adapters/home_output_test.go`
- `docs/dev/0007-home_header_bars_and_timestamp_metadata.md`
