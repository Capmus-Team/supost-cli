# Home Category Cache + Seed Taxonomy Alignment

Date: 2026-02-27

## Summary
Improved `home` startup performance and data consistency by caching category sections alongside posts, simplifying Postgres category timing queries, and adding seed taxonomy snapshots for categories/subcategories.

## What Changed

### 1. Home command now caches category sections too
- Updated `cmd/home.go`:
  - added cached section load path (`getCachedHomeCategorySections`)
  - returns immediately from cache when both posts and sections are cached
  - persists section cache after fresh load when DB mode + cache enabled
  - keeps verbose warnings for cache/section load issues without failing command.

### 2. Added category section cache read/write helpers
- Updated `internal/adapters/home_cache.go` with:
  - `homeCategorySectionsCache` payload model
  - `LoadHomeCategorySectionsCache(ttl)`
  - `SaveHomeCategorySectionsCache(sections)`
  - `homeCategorySectionsCachePath()`
- Cache file path:
  - `.../supost-cli/home_category_sections.json`

### 3. In-memory repository now computes category section times
- Updated `internal/repository/inmemory.go`:
  - `ListHomeCategorySections` now scans in-memory active posts
  - computes latest timestamp per category
  - returns ordered `[]domain.HomeCategorySection` by `category_id`.

### 4. Postgres category-section query simplified to timing-only
- Updated `internal/repository/postgres.go`:
  - `ListHomeCategorySections` now reads only from `public.post`
  - groups by `category_id`
  - computes `MAX(time_posted_at | to_timestamp(time_posted))`
- Taxonomy names/subcategory labels are handled from local/default metadata instead of DB joins.

### 5. Added taxonomy seed snapshots
- New files:
  - `testdata/seed/category_rows.json`
  - `testdata/seed/subcategory_rows.json`
- Purpose: provide reproducible local category/subcategory reference data aligned with production taxonomy snapshots.

## Why This Matters
- Reduces repeated DB calls for sidebar/category data on frequent `home` runs.
- Keeps in-memory mode closer to DB-backed behavior for category last-post timing.
- Decouples category labels from DB joins while preserving time freshness from live post data.
- Adds explicit taxonomy seed artifacts for local development and fallback mapping.

## Files in This Increment
- `cmd/home.go`
- `internal/adapters/home_cache.go`
- `internal/repository/inmemory.go`
- `internal/repository/postgres.go`
- `testdata/seed/category_rows.json`
- `testdata/seed/subcategory_rows.json`
- `docs/dev/0013-home_category_cache_and_seed_taxonomy_alignment.md`
