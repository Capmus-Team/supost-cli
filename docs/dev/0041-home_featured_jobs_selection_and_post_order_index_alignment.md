# Home Featured Jobs Selection and Post-Order Index Alignment

Date: 2026-03-01

## Summary
Optimized home page featured-jobs loading to reuse already-fetched home posts first, updated Postgres ordering/query shapes for active post retrieval and category section recency, and added matching partial indexes for home ordering patterns.

## What Changed

### 1. Reused home feed posts for featured jobs before fallback query
- Updated `cmd/home.go`:
  - added `selectFeaturedJobsFromPosts(posts, limit)` helper
  - derives featured jobs from the already-loaded home posts first
  - only falls back to `ListRecentActiveByCategory(...jobs...)` when fewer than required results are found.

### 2. Standardized Postgres home ordering clauses
- Updated `internal/repository/postgres.go`:
  - changed active post ordering to `ORDER BY public.post.time_posted DESC, public.post.id DESC`
  - applies both in:
    - `ListRecentActivePosts`
    - `ListRecentActiveByCategory`.

### 3. Refined home category-section recency query
- Updated `internal/repository/postgres.go` `ListHomeCategorySections`:
  - switched from grouped aggregate over `public.post` to per-category `JOIN LATERAL` lookup
  - now reads latest active post timestamp per category using ordered `LIMIT 1`
  - keeps ordered output by category id.

### 4. Added indexes aligned to home query sort/filter patterns
- Added `supabase/migrations/20260301007000_home_post_order_indexes.sql`:
  - `idx_post_active_time_posted_id_desc` on `(time_posted desc, id desc)` where `status = 1`
  - `idx_post_active_category_time_posted_id_desc` on `(category_id, time_posted desc, id desc)` where `status = 1`.

## Why This Matters
- Reduces redundant DB work for featured jobs on home.
- Improves planner alignment for the most common active-post ordering paths.
- Supports fast “latest per category” lookups with query shapes that pair with the new indexes.

## Files in This Increment
- `cmd/home.go`
- `internal/repository/postgres.go`
- `supabase/migrations/20260301007000_home_post_order_indexes.sql`
- `docs/dev/0041-home_featured_jobs_selection_and_post_order_index_alignment.md`
