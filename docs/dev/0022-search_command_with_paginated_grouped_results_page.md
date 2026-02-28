# Search Command With Paginated, Date-Grouped Results Page

Date: 2026-02-28

## Summary
Added a new `search` command and supporting domain/service/repository/adapter layers to render paginated active post results in a date-grouped page format with shared header/footer and taxonomy-aware breadcrumb/title context.

## What Changed

### 1. New `search` CLI command
- Added `cmd/search.go`:
  - command: `search`
  - flags:
    - `--category` (int64)
    - `--subcategory` (int64)
    - `--page` (default `1`)
    - `--per-page` (default `100`, capped in service)
  - uses configured adapter (in-memory or Postgres)
  - calls `SearchService.Search(...)`
  - renders via page adapter for default/text/table formats.

### 2. New search result domain contract
- Added `internal/domain/search_result.go`:
  - `SearchResultPage` with filter metadata, paging fields, `has_more`, and result post slice.

### 3. New search service and tests
- Added `internal/service/search.go`:
  - defines `SearchRepository` interface where consumed
  - normalizes paging (`page>=1`, `per_page` default `100`, max `100`)
  - returns structured `SearchResultPage`.
- Added `internal/service/search_test.go`:
  - validates paging normalization and filter forwarding.

### 4. Repository implementations for search
- Added `internal/repository/inmemory_search.go`:
  - filters active posts
  - optional category/subcategory filtering
  - sorts newest-first
  - supports pagination with `hasMore`.
- Added `internal/repository/postgres_search.go`:
  - parameterized query against `public.post`
  - optional filter conditions for category/subcategory
  - `LIMIT perPage+1` strategy to compute `hasMore`.

### 5. Search page renderer and test
- Added `internal/adapters/search_output.go`:
  - renders shared header/footer
  - applies taxonomy-aware breadcrumb via `BreadcrumbOptions`
  - shows search title from taxonomy lookup
  - groups posts by posting date header
  - wraps each post line, including title/email and image indicator
  - renders `next N posts` prompt when `has_more` is true.
- Added `internal/adapters/search_output_test.go`:
  - verifies breadcrumb/title/date-group headers and next-page indicator.

## Why This Matters
- Introduces a web-like search page flow in CLI form while preserving clean layering.
- Keeps output human-readable (grouped by day) and machine-friendly (domain contract remains explicit).
- Supports consistent behavior across in-memory prototype data and real Postgres data.

## Files in This Increment
- `cmd/search.go`
- `internal/domain/search_result.go`
- `internal/service/search.go`
- `internal/service/search_test.go`
- `internal/repository/inmemory_search.go`
- `internal/repository/postgres_search.go`
- `internal/adapters/search_output.go`
- `internal/adapters/search_output_test.go`
- `docs/dev/0022-search_command_with_paginated_grouped_results_page.md`
