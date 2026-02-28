# Search Pagination Labels and Empty-State Rendering

Date: 2026-02-28

## Summary
Improved search result page UX by adding previous/next pagination labels and a no-results empty state, with test coverage in both adapter and service layers.

## What Changed

### 1. Search output now renders previous/next pagination labels
- Updated `internal/adapters/search_output.go`:
  - when `page > 1`, render centered `previous N posts`
  - when `has_more = true`, render centered `next N posts`
  - supports rendering both labels together on middle pages.

### 2. Search output now handles empty result pages explicitly
- Updated `renderSearchGroupedPosts(...)`:
  - when post list is empty, renders centered gray message:
    - `No posts found for this page.`

### 3. Adapter tests expanded for pagination and empty state
- Updated `internal/adapters/search_output_test.go`:
  - added test for simultaneous previous+next labels on page > 1
  - added test for empty-state message and previous label when page > 1 with no posts.

### 4. Service tests expanded for explicit paging forwarding
- Updated `internal/service/search_test.go`:
  - added `TestSearchService_Search_ForwardsPageAndPerPage`
  - verifies non-default `page/per_page` values are forwarded to repository unchanged.

## Why This Matters
- Makes search navigation clearer across first/middle/empty pages.
- Prevents blank output ambiguity when filters return no results.
- Increases confidence in paging behavior across service and rendering boundaries.

## Files in This Increment
- `internal/adapters/search_output.go`
- `internal/adapters/search_output_test.go`
- `internal/service/search_test.go`
- `docs/dev/0024-search_pagination_labels_and_empty_state.md`
