# Post Page Notice Placement + Housing Policy Notice

Date: 2026-02-28

## Summary
Refined post-page content flow by separating date and price into distinct lines, moving policy/warning notices into the main content column, and adding a housing-specific sublicensing notice for housing categories.

## What Changed

### 1. Post header date/price line split
- Updated `internal/adapters/post_output.go`:
  - `renderPostTopBlock(...)` now always prints `Date:` on its own line.
  - `Price:` now prints on a separate line only when `post.HasPrice` is true.

### 2. Notices moved to post main content area
- Updated `renderPostMessagePosterRows(...)` in `internal/adapters/post_output.go`:
  - removed trailing commercial warning from the right-side message panel.
- Updated `renderPostMainRows(...)`:
  - commercial-services warning remains in the main (left) content column under the body.

### 3. Housing-specific notice added
- Added constant `postHousingNotice` in `internal/adapters/post_output.go`.
- Added `isHousingCategory(categoryID int64)` helper.
- `renderPostMainRows(...)` now conditionally includes the housing notice for category IDs `3` and `4`.

### 4. Tests updated and expanded
- Updated `internal/adapters/post_output_test.go`:
  - existing page render test now asserts housing notice text and URL for housing posts.
  - added `TestRenderPostPage_NonHousingSkipsHousingNotice` to verify non-housing posts do not include the housing policy line.

## Why This Matters
- Keeps legal/policy text near the post content where users read details.
- Avoids overloading the right-side message form panel with unrelated copy.
- Ensures housing compliance reminder appears only when relevant.

## Files in This Increment
- `internal/adapters/post_output.go`
- `internal/adapters/post_output_test.go`
- `docs/dev/0020-post_page_notice_placement_and_housing_policy_notice.md`
