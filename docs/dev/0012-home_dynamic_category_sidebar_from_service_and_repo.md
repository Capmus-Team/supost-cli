# Home Dynamic Category Sidebar from Service + Repository

Date: 2026-02-27

## Summary
Extended the `home` flow so category/sidebar data now comes through service and repository layers (with graceful fallbacks), and updated home rendering to use dynamic category sections plus subcategory detail blocks.

## What Changed

### 1. Added home category domain model
- New file: `internal/domain/home_category.go`
- Introduced `HomeCategorySection` with `json` and `db` tags:
  - `category_id`, `category_name`, `subcategory_names`, `last_posted_at`
- Added category ID constants for known home sections.

### 2. Repository contract now supports home category sections
- Updated `internal/repository/interfaces.go`:
  - `HomePostStore` now includes `ListHomeCategorySections(ctx)`.
- In-memory adapter (`internal/repository/inmemory.go`) implements this as `nil, nil` so renderer fallback metadata is used when DB taxonomy is unavailable.
- Postgres adapter (`internal/repository/postgres.go`) adds `ListHomeCategorySections` query:
  - joins `category` + `subcategory`
  - computes latest active post time per category
  - groups rows into ordered `[]HomeCategorySection`.

### 3. Service layer exposes category sections
- Updated `internal/service/home.go`:
  - `HomeRepository` now includes `ListHomeCategorySections`
  - new `ListCategorySections(ctx)` method on `HomeService`.
- Updated tests in `internal/service/home_test.go` to cover section retrieval.

### 4. Command flow now loads posts + sections with fallback behavior
- Updated `cmd/home.go`:
  - keeps cache-based post fallback behavior more explicitly (`usedCache`, `cacheLoadErr`)
  - fetches sections through service via `ListCategorySections`
  - if section fetch fails, command continues (warns in verbose mode)
  - `renderHomeOutput` now accepts both `posts` and `sections`.

### 5. Home rendering moved to dynamic sidebar helpers
- Updated `internal/adapters/home_output.go`:
  - `RenderHomePosts` now accepts `sections []domain.HomeCategorySection`
  - delegates sidebar/overview layout to new helpers in `internal/adapters/home_sidebar.go`
  - widened sidebar width (`homeCalloutWidth` from 28 to 36)
- New file: `internal/adapters/home_sidebar.go`:
  - default category definitions
  - section normalization/fallback logic
  - overview rows with relative ages
  - category detail rows with subcategory formatting.

### 6. Adapter tests updated for dynamic sidebar behavior
- Updated `internal/adapters/home_output_test.go` to validate:
  - overview rendering with normalized sections
  - category detail rendering includes expected subcategories.

## Why This Matters
- Home sidebar data now follows the same domain/service/repository architecture as the post feed.
- DB-backed category taxonomy can drive rendering when available, while in-memory mode remains zero-dependency.
- Command behavior is more resilient: category/sidebar load failures no longer block home output.

## Files in This Increment
- `cmd/home.go`
- `internal/domain/home_category.go`
- `internal/repository/interfaces.go`
- `internal/repository/inmemory.go`
- `internal/repository/postgres.go`
- `internal/service/home.go`
- `internal/service/home_test.go`
- `internal/adapters/home_output.go`
- `internal/adapters/home_output_test.go`
- `internal/adapters/home_sidebar.go`
- `docs/dev/0012-home_dynamic_category_sidebar_from_service_and_repo.md`
