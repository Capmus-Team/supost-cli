# Search Default "All Posts" Title and In-Memory Search Tests

Date: 2026-03-01

## Summary
Updated search UX/documentation to treat `search` as an all-posts default view, simplified the page title to `all posts`, and added in-memory repository tests to lock active-post ordering and pagination behavior.

## What Changed

### 1. Search command messaging now reflects all-post default behavior
- Updated `cmd/search.go`:
  - command short description now says `Render all-post search results`
  - long description now says `Show paginated active posts grouped by posting date.`

### 2. Search page title standardized to `all posts`
- Updated `internal/adapters/search_output.go`:
  - removed category/subcategory-specific title derivation
  - replaced title helper with a stable `all posts` title.

### 3. Search renderer tests updated for new title contract
- Updated `internal/adapters/search_output_test.go`:
  - updated expectations from category/subcategory title labels to `all posts`
  - retained breadcrumb inference checks to ensure category/subcategory context still renders in header breadcrumbs.

### 4. Added in-memory search repository regression tests
- Added `internal/repository/inmemory_search_test.go`:
  - verifies `SearchActivePosts` returns only active posts
  - verifies ordering is newest-first with deterministic tie behavior
  - verifies pagination `hasMore` behavior across page boundaries.

### 5. README examples/documentation aligned to new defaults
- Updated `README.md`:
  - quick-start and command examples now show `supost search` as default all-post view
  - examples emphasize optional filters and pagination rather than category-only usage
  - command tree description now labels `search` as rendering all recent active posts.

## Why This Matters
- Clarifies that `search` works out of the box without required filters.
- Keeps page title semantics stable for parity with Craigslist-style all-post browsing.
- Adds repository-level regression coverage for core search ordering/pagination behavior.

## Files in This Increment
- `README.md`
- `cmd/search.go`
- `internal/adapters/search_output.go`
- `internal/adapters/search_output_test.go`
- `internal/repository/inmemory_search_test.go`
- `docs/dev/0032-search_default_all_posts_title_and_inmemory_search_tests.md`
