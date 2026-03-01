# Dynamic Postgres Search SQL and Subcategory Order Index

Date: 2026-03-01

## Summary
Refactored Postgres search SQL construction to build parameterized WHERE/ORDER clauses dynamically (without OR-guard patterns), updated repository tests for the new query shape/arg positions, and added a partial index for active-post subcategory ordering.

## What Changed

### 1. Reworked Postgres search statement builder to dynamic clause assembly
- Updated `internal/repository/postgres_search.go`:
  - replaced separate fixed query constants with `sqlQuerySearchSelect` base select
  - `buildSearchActivePostsStatement(...)` now:
    - appends only relevant filters (`category_id`, `subcategory_id`) when set
    - appends FTS predicate only when query text is non-empty
    - computes positional placeholders dynamically for query, limit, and offset
    - builds order clause as:
      - default: `p.time_posted DESC, p.id DESC`
      - keyword query: `ts_rank(...) DESC, p.time_posted DESC, p.id DESC`.

### 2. Updated repository tests for new SQL and argument contracts
- Updated `internal/repository/postgres_search_test.go`:
  - removed fixed-constant SQL equality assertions
  - now asserts SQL contains expected dynamic clauses/placeholder ordering
  - verifies OR-guard patterns and `NULLS LAST` are absent
  - updates expected argument positions/counts for keyword/default paths
  - adds coverage ensuring unset filters are omitted from SQL.

### 3. Added active-post subcategory ordering index migration
- Added `supabase/migrations/20260301008000_search_subcategory_order_index.sql`:
  - creates partial btree index:
    - `(subcategory_id, time_posted desc, id desc)`
    - `WHERE status = 1`.

## Why This Matters
- Produces tighter SQL for active search filters by removing always-true OR guards.
- Keeps placeholder numbering consistent with only the parameters actually used.
- Improves planner support for subcategory-filtered search ordering.

## Files in This Increment
- `internal/repository/postgres_search.go`
- `internal/repository/postgres_search_test.go`
- `supabase/migrations/20260301008000_search_subcategory_order_index.sql`
- `docs/dev/0042-dynamic_postgres_search_sql_and_subcategory_order_index.md`
