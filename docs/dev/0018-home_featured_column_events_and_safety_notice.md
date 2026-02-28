# Home Featured Column Events + Safety Notice

Date: 2026-02-28

## Summary
Extended the home right column to include an `events` section placeholder and two safety/disclaimer notices beneath featured job posts, with test coverage updated to validate the added content in wrapped terminal output.

## What Changed

### 1. Added events/safety/disclaimer content to featured column
- Updated `internal/adapters/home_output.go`:
  - added constants:
    - `homeEventsPlaceholder`
    - `homeSafetyNotice`
    - `homeAffiliationNotice`
  - updated `renderFeaturedJobPostRows(...)` to append:
    - blank spacer
    - `events` header
    - events placeholder text
    - blank spacer
    - safety notice text
    - affiliation disclaimer text

### 2. Added shared wrapped-text helper
- Added `appendWrappedTextRows(...)` in `internal/adapters/home_output.go`.
- Helper wraps colored text with existing styled-word logic and appends rendered lines to the output row slice.

### 3. Updated tests for wrapped/normalized assertions
- Updated `internal/adapters/home_output_test.go`:
  - existing featured-column tests now assert for:
    - `events`
    - `events data placeholder`
    - full safety warning
    - Stanford affiliation disclaimer
  - switched matching to whitespace-normalized text (`strings.Fields`) to keep assertions stable across terminal wrapping.

## Why This Matters
- Keeps key trust/safety messaging visible directly in the home view.
- Adds a clear events slot in the right column without changing service/repository contracts.
- Maintains robust tests despite dynamic line wrapping widths.

## Files in This Increment
- `internal/adapters/home_output.go`
- `internal/adapters/home_output_test.go`
- `docs/dev/0018-home_featured_column_events_and_safety_notice.md`
