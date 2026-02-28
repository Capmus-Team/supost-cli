# Featured Jobs Service/Repo Wiring + Home Render Input

Date: 2026-02-28

## Summary
Completed end-to-end wiring for featured job posts so the `home` command fetches a dedicated featured jobs slice from the service/repository layer and passes it into terminal rendering, instead of deriving featured entries only from the general recent feed.

## What Changed

### 1. `home` command now fetches featured jobs explicitly
- Updated `cmd/home.go`:
  - added `featuredJobPostLimit = 3`
  - added `featuredJobs` retrieval via `HomeService.ListRecentActiveByCategory(..., CategoryJobsOffCampus, 3)`
  - passes `featuredJobs` into render path: `renderHomeOutput(..., posts, featuredJobs, sections)`
  - preserves graceful fallback behavior when featured query fails (warn in verbose mode, continue rendering)
  - adjusted cached-sections fallback path to match updated render signature.

### 2. Service interface/use-case expanded for category-scoped recents
- Updated `internal/service/home.go`:
  - `HomeRepository` now includes `ListRecentActivePostsByCategory(...)`
  - added `ListRecentActiveByCategory(...)` with default-limit guard.
- Updated `internal/service/home_test.go`:
  - verifies category ID and limit forwarding
  - verifies default limit behavior for non-positive limits.

### 3. Repository contract and adapters now support category filter
- Updated `internal/repository/interfaces.go`:
  - added `ListRecentActivePostsByCategory(...)` to `HomePostStore`.
- Updated `internal/repository/inmemory.go`:
  - extracted shared filtering/sorting helper `listRecentActivePosts(limit, categoryID *int64)`
  - implemented category-filtered variant with same active/newest ordering behavior.
- Updated `internal/repository/postgres.go`:
  - added parameterized query implementation for active posts by category
  - uses `status = $1 AND category_id = $2` with `LIMIT $3`.

### 4. Home adapters now accept explicit featured input
- Updated `internal/adapters/home_output.go`:
  - `RenderHomePosts` signature now includes `featuredPosts []domain.Post`
  - `renderHomeRecentAndFeaturedRows(...)` now prefers explicit `featuredPosts` and falls back to selecting from `posts` when empty.
- Updated `internal/adapters/home_sidebar.go`:
  - threaded `featuredPosts` through `renderHomeOverviewAndRecent(...)`.
- Updated `internal/adapters/home_output_test.go`:
  - adapted existing call sites for new signature
  - added explicit-featured test coverage.

## Why This Matters
- Moves featured-job selection responsibility into service/repository access patterns, aligning with project layering rules.
- Keeps home rendering resilient by supporting both explicit featured data and derived fallback.
- Improves test coverage around service contract forwarding and explicit featured-render behavior.

## Files in This Increment
- `cmd/home.go`
- `internal/service/home.go`
- `internal/service/home_test.go`
- `internal/repository/interfaces.go`
- `internal/repository/inmemory.go`
- `internal/repository/postgres.go`
- `internal/adapters/home_output.go`
- `internal/adapters/home_sidebar.go`
- `internal/adapters/home_output_test.go`
- `docs/dev/0017-featured_jobs_service_repo_wiring_and_home_render_input.md`
