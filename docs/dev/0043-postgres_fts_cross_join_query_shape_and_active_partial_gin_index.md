# Postgres FTS Cross-Join Query Shape and Active Partial GIN Index

Date: 2026-03-01

## Summary
Refined Postgres search SQL generation to use a single `plainto_tsquery` binding via `CROSS JOIN` for keyword searches, updated SQL-shape tests accordingly, and added a partial active-post GIN index on `post.fts`.

## What Changed

### 1. Updated keyword search SQL to reuse tsquery value
- Updated `internal/repository/postgres_search.go`:
  - removed static `FROM public.post p` from base select constant
  - dynamic builder now composes `fromClause`
  - keyword path now uses:
    - `CROSS JOIN plainto_tsquery('english', $n) q`
    - predicate `p.fts @@ q`
    - ranking `ts_rank(p.fts, q) DESC`
  - default path still uses `FROM public.post p` and non-FTS ordering.

### 2. Updated repository tests for new SQL shape
- Updated `internal/repository/postgres_search_test.go`:
  - asserts default SQL includes `FROM public.post p`
  - FTS test now asserts:
    - `CROSS JOIN plainto_tsquery('english', $2) q`
    - `p.fts @@ q`
    - `ts_rank(p.fts, q)`
  - preserves arg-position and limit/offset checks.

### 3. Added partial active-post GIN index migration
- Added `supabase/migrations/20260301009000_post_active_fts_partial_index.sql`:
  - creates `idx_post_active_fts_idx`
  - `GIN (fts)` partial index with `WHERE status = 1`.

## Why This Matters
- Avoids repeating `plainto_tsquery(...)` calls in predicate and ranking expressions.
- Keeps generated SQL explicit and test-verified for both default and keyword branches.
- Improves targeted FTS performance for active-post search paths.

## Files in This Increment
- `internal/repository/postgres_search.go`
- `internal/repository/postgres_search_test.go`
- `supabase/migrations/20260301009000_post_active_fts_partial_index.sql`
- `docs/dev/0043-postgres_fts_cross_join_query_shape_and_active_partial_gin_index.md`
