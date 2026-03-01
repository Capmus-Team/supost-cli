# Postgres FTS Search Path and Supabase Schema Snapshot Artifacts

Date: 2026-03-01

## Summary
Refined Postgres search implementation to branch between default listing SQL and full-text search SQL, added focused tests for statement construction/command behavior, introduced a migration for generated `post.fts` index support, and added `supabase/schema/` snapshot + inventory artifacts.

## What Changed

### 1. Split Postgres search into explicit default vs FTS SQL paths
- Updated `internal/repository/postgres_search.go`:
  - extracted `sqlQuerySearchDefault` for non-keyword search
  - extracted `sqlQuerySearchFTS` for keyword search against `post.fts`
  - FTS query uses `plainto_tsquery('english', ...)`, `p.fts @@ q`, and rank ordering via `ts_rank`
  - added `buildSearchActivePostsStatement(...)` helper to normalize page/per-page and choose SQL/args.

### 2. Added repository unit tests for SQL selection + paging normalization
- Added `internal/repository/postgres_search_test.go`:
  - verifies empty query selects default SQL/args
  - verifies keyword query selects FTS SQL and trims query text
  - verifies paging defaults (`page`, `per_page`, `limit`, `offset`) are normalized.

### 3. Added command-level search keyword regression test
- Added `cmd/search_test.go`:
  - runs search command with keyword args
  - validates rendered output includes keyword title and matching post content.

### 4. Added migration to support indexed full-text search
- Added `supabase/migrations/20260301006000_post_fts_name_body.sql`:
  - adds generated `public.post.fts` tsvector column from weighted `name` + `body`
  - creates GIN index `post_fts_idx`.

### 5. Added Supabase schema snapshot + optimization inventory artifacts
- Added `supabase/schema/` artifacts:
  - `full_schema.sql` schema-only dump
  - `indexes.csv`
  - `index_usage_stats.csv`
  - `table_scan_stats.csv`
  - `table_sizes.csv`
  - `triggers.csv`
  - `README.md` with refresh guidance.

## Why This Matters
- Improves search relevance/performance for keyword queries by using dedicated full-text indexes and ranking.
- Preserves clear execution paths for filtered listing vs keyword search.
- Captures current remote schema/index usage state to support future tuning and migration planning.

## Files in This Increment
- `cmd/search_test.go`
- `internal/repository/postgres_search.go`
- `internal/repository/postgres_search_test.go`
- `supabase/migrations/20260301006000_post_fts_name_body.sql`
- `supabase/schema/README.md`
- `supabase/schema/full_schema.sql`
- `supabase/schema/indexes.csv`
- `supabase/schema/index_usage_stats.csv`
- `supabase/schema/table_scan_stats.csv`
- `supabase/schema/table_sizes.csv`
- `supabase/schema/triggers.csv`
- `docs/dev/0040-postgres_fts_search_path_and_supabase_schema_snapshot_artifacts.md`
